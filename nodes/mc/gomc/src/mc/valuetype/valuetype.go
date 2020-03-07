package valuetype

import (
	"fmt"
	"sort"
	"strings"
	value "mc/value"
)


type ValueTypeKind int

const (
	TY_ONEOF ValueTypeKind = iota + 1
	TY_MAP
	TY_LIST
	TY_NUM
	TY_STRING
	TY_BOOL
	TY_EXT
)

var (
	ValueTypeMapping = map[ValueTypeKind]value.ValueKind{
		TY_ONEOF: value.VAL_ANY,
		TY_MAP: value.VAL_MAP,
		TY_LIST: value.VAL_LIST,
		TY_NUM: value.VAL_NUM,
		TY_STRING: value.VAL_STRING,
		TY_BOOL: value.VAL_BOOL}
)

type ValueType struct {
	Name string
	Type ValueTypeKind
	ExternalTypeName string
	MapArgTypes map[string]*ValueType
	MapArgRequired map[string]bool
	MapDefaults map[string]*value.Value
	ListArgType *ValueType
	OneofTypes []*ValueType
}

func (ty *ValueType) ToString() string {
	if ty.Type == TY_ONEOF {
		subtypes := make([]string, len(ty.OneofTypes))
		for i, sty := range ty.OneofTypes {
			subtypes[i] = sty.ToString()
		}
		return fmt.Sprintf("oneof(%s)", strings.Join(subtypes, ","))
	} else if ty.Type == TY_MAP {
		keys := make([]string, 0)
		for k, _ := range ty.MapArgTypes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		entries := make([]string, 0)
		// Check if all the keys in v are expected
		for _, k := range keys {
			sty := ty.MapArgTypes[k]
			req := "*"
			if ty.MapArgRequired[k] {
				req = ""
			}
			def := ""
			if defaultval, ok := ty.MapDefaults[k]; ok {
				def = "="+defaultval.ToString()
				req = ""
			}
			entries = append(entries, fmt.Sprintf("%s%s:%s%s", k, req, sty.ToString(), def))
		}
		return fmt.Sprintf("{%s}", strings.Join(entries, ","))
	} else if ty.Type == TY_LIST {
		return fmt.Sprintf("[%s]", ty.ListArgType.ToString())
	} else if ty.Type == TY_NUM {
		return "num"
	} else if ty.Type == TY_STRING {
		return "string"
	} else if ty.Type == TY_BOOL {
		return "bool"
	} else if ty.Type == TY_EXT {
		return ty.ExternalTypeName
	}
	return "[unknown]"
}

func MakeExtType(name string) *ValueType {
	return &ValueType{Type: TY_EXT, ExternalTypeName: name}
}

func MakeBoolType() *ValueType {
	return &ValueType{Type: TY_BOOL}
}
func MakeNumType() *ValueType {
	return &ValueType{Type: TY_NUM}
}
func MakeStringType() *ValueType {
	return &ValueType{Type: TY_STRING}
}
func MakeListType(subtype *ValueType) *ValueType {
	return &ValueType{
		Type: TY_LIST,
		ListArgType: subtype}
}
func MakeOneofType(subtypes []*ValueType) *ValueType {
	return &ValueType{
		Type: TY_ONEOF,
		OneofTypes: subtypes}
}

// Get the value types possibly conforming to ty
// func (ty *ValueType) PossibleValueTypes() []value.ValueKind {
// 	possible := map[value.ValueKind]bool{
// 		value.VAL_MAP:false,
// 		value.VAL_LIST:false,
// 		value.VAL_NAME:false,
// 		value.VAL_NUM:false,
// 		value.VAL_STRING:false,
// 		value.VAL_BOOL:false}
// 	if ty.Type == TY_ONEOF {
// 		ans := []value.ValueKind{}
// 		for _, sty := range ty.OneofTypes {
// 			for _, tyty := range sty.PossibleValueTypes() {
// 				if !possible[tyty] {
// 					possible[tyty] = true
// 					ans = append(ans, tyty)
// 				}
// 			}
// 		}
// 		return ans
// 	}
// 	return []value.ValueKind{ValueTypeMapping[ty.Type]}
// }

// Does v conform to ty?
func (ty *ValueType) Validate(v *value.Value, ext_types map[string]*ValueType, path string) (*value.Value, error) {
	fmt.Println("Validate",ty,v)
	if ty.Type == TY_ONEOF {
		for _, sty := range ty.OneofTypes {
			if sv, err := sty.Validate(v, ext_types, path); err == nil {
				return sv, nil
			}
		}
		return nil, fmt.Errorf("Error at %s: Subvalue does not validate as any of the given types", path)
	} else if ty.Type == TY_MAP && v.Type == value.VAL_MAP {
		fmt.Println("map")
		// Check if all the keys in v are expected
		for k, _ := range v.MapVal {
			if _, ok := ty.MapArgTypes[k]; !ok {
				return nil, fmt.Errorf("Error at %s: Unexpected map entry: %s", path, k)
			}
		}
		// Check all required keys are present
		for k, r := range ty.MapArgRequired {
			if _, ok := v.MapVal[k]; r && !ok {
				return nil, fmt.Errorf("Error at %s: Map entry `%s` required", path, k)
			}
		}
		// Populate default values of absent values
		for k, defval := range ty.MapDefaults {
			if _, ok := v.MapVal[k]; !ok {
				v.MapVal[k] = defval
			}
		}
		// Now check each is of the required type
		for k, sv := range v.MapVal {
			svv, err := ty.MapArgTypes[k].Validate(sv, ext_types, fmt.Sprintf("%s.%s",path,k))
			if err != nil {
				return nil, fmt.Errorf("Error at %s: %v", path, err)
			} else {
				v.MapVal[k] = svv
			}
		}
		return v, nil
	} else if ty.Type == TY_LIST && v.Type == value.VAL_LIST {
		for i, sv := range v.ListVal {
			svv, err := ty.ListArgType.Validate(sv, ext_types, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, fmt.Errorf("Error at %s: %v", path, err)
			} else {
				v.ListVal[i] = svv
			}
		}
		return v, nil
	} else if ty.Type == TY_NUM {
		if v.Type == value.VAL_NUM {
			return v, nil
		}
		return nil, fmt.Errorf("Error at %s: num expected", path)
	} else if ty.Type == TY_STRING {
		if v.Type == value.VAL_STRING {
			return v, nil
		}
		return nil, fmt.Errorf("Error at %s: string expected", path)
	} else if ty.Type == TY_BOOL {
		if v.Type == value.VAL_BOOL {
			return v, nil
		}
		return nil, fmt.Errorf("Error at %s: bool expected", path)
	} else if ty.Type == TY_EXT {
		if et, ok := ext_types[ty.ExternalTypeName]; ok {
			return et.Validate(v, ext_types, path)
		}
		return nil, fmt.Errorf("Error at %s: Unknown type `%s`", path, ty.ExternalTypeName)
	}
	return nil, fmt.Errorf("Invalid type: %s", path)
}

