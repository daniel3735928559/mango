package valuetype

import (
	// "fmt"
	// "strings"
	// "errors"
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
	MapArgTypes map[string]*ValueType
	MapArgRequired map[string]bool
	MapDefaults map[string]*value.Value
	ListArgType *ValueType
	OneofTypes []*ValueType
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
func (ty *ValueType) PossibleValueTypes() []value.ValueKind {
	possible := map[value.ValueKind]bool{
		value.VAL_MAP:false,
		value.VAL_LIST:false,
		value.VAL_NAME:false,
		value.VAL_NUM:false,
		value.VAL_STRING:false,
		value.VAL_BOOL:false}
	if ty.Type == TY_ONEOF {
		ans := []value.ValueKind{}
		for _, sty := range ty.OneofTypes {
			for _, tyty := range sty.PossibleValueTypes() {
				if !possible[tyty] {
					possible[tyty] = true
					ans = append(ans, tyty)
				}
			}
		}
		return ans
	}
	return []value.ValueKind{ValueTypeMapping[ty.Type]}
}

// Does v conform to ty?
func (ty *ValueType) Validate(v *value.Value) *value.Value {
	if ty.Type == TY_ONEOF {
		for _, sty := range ty.OneofTypes {
			if sv := sty.Validate(v); sv != nil {
				return sv
			}
		}
		return nil
	} else if ty.Type == TY_MAP && v.Type == value.VAL_MAP {
		// Check if all the keys in v are expected
		for k, _ := range v.MapVal {
			if _, ok := ty.MapArgTypes[k]; !ok {
				return nil
			}
		}
		// Check all required keys are present
		for k, r := range ty.MapArgRequired {
			if _, ok := v.MapVal[k]; r && !ok {
				return nil
			}
		}
		// Now check each is of the required type
		for k, sv := range v.MapVal {
			if svv := ty.MapArgTypes[k].Validate(sv); svv != nil {
				v.MapVal[k] = svv
			} else {
				return nil
			}
		}
		return v
	} else if ty.Type == TY_LIST && v.Type == value.VAL_LIST {
		for i, sv := range v.ListVal {
			if svv := ty.ListArgType.Validate(sv); svv != nil {
				v.ListVal[i] = svv
			} else {
				return nil
			}
			return v
		}
	} else if ty.Type == TY_NUM && v.Type == value.VAL_NUM {
		return v
	} else if ty.Type == TY_STRING && v.Type == value.VAL_STRING {
		return v
	} else if ty.Type == TY_BOOL && v.Type == value.VAL_BOOL {
		return v
	}
	return nil
}

