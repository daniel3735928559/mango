package router

import (
	"fmt"
	"sort"
	"strings"
	"errors"
)

type ValueType int

const (
	VAL_MAP ValueType = iota + 1
	VAL_LIST
	VAL_NAME
	VAL_INT
	VAL_FLOAT
	VAL_NUM
	VAL_STRING
	VAL_BOOL
	VAL_ANY
	VAL_ANYANY
)

type Value struct {
	Type ValueType
	MapVal map[string]*Value
	ListVal []*Value
	NameVal string
	IntVal int
	FloatVal float64
	NumVal float64
	StringVal string
	BoolVal bool
}

func (v *Value) ToPrimitive() interface{} {
	if v.Type == VAL_MAP {
		ans := make(map[string]interface{})
		for key, val := range v.MapVal {
			ans[key] = val.ToPrimitive()
		}
		return ans
	} else if v.Type == VAL_LIST {
		ans := make([]interface{}, len(v.ListVal))
		for idx, val := range v.ListVal {
			ans[idx] = val.ToPrimitive()
		}
		return ans
	} else if v.Type == VAL_NAME {
		return v.NameVal
	} else if v.Type == VAL_INT {
		return v.IntVal
	} else if v.Type == VAL_FLOAT {
		return v.FloatVal
	} else if v.Type == VAL_NUM {
		return v.NumVal
	} else if v.Type == VAL_STRING {
		return v.StringVal
	} else if v.Type == VAL_BOOL {
		return v.BoolVal
	} else if v.Type == VAL_ANY {
		return "any"
	}
	return "[unknown value type]"
	
}

func (v *Value) Equals(x *Value) bool {
	fmt.Println("CMP",v,x)
	if v.Type != x.Type {
		return false
	} else if v.Type == VAL_NUM {
		return v.NumVal == x.NumVal
	} else if v.Type == VAL_STRING {
		return v.StringVal == x.StringVal
	} else if v.Type == VAL_LIST {
		if len(v.ListVal) == len(x.ListVal) {
			for i, item := range v.ListVal {
				if !item.Equals(x.ListVal[i]) {
					return false
				}
			}
			return true
		}
		return false
	} else if v.Type == VAL_MAP {
		if len(v.MapVal) == len(x.MapVal) {
			for key, val := range v.MapVal {
				if other_val, ok := x.MapVal[key]; ok {
					if !val.Equals(other_val) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
		return false
	} else if v.Type == VAL_BOOL {
		return v.BoolVal == x.BoolVal
	}
	return false
}

func (v *Value) Clone() *Value {
	ans := &Value{
		Type: v.Type,
		NameVal: v.NameVal,
		IntVal: v.IntVal,
		FloatVal: v.FloatVal,
		NumVal: v.NumVal,
		StringVal: v.StringVal,
		BoolVal: v.BoolVal}
	if v.MapVal != nil {
		ans.MapVal = make(map[string]*Value)
		for key, val := range v.MapVal {
			ans.MapVal[key] = val.Clone()
		}
	}
	if v.ListVal != nil {
		ans.ListVal = make([]*Value, len(v.ListVal))
		for i, val := range v.ListVal {
			ans.ListVal[i] = val.Clone()
		}
	}
	return ans
}

func MakeValue(args interface{}) (*Value, error) {
	var err error
	if mapvals, ok := args.(map[string]interface{}); ok {
		mapval := make(map[string]*Value)
		for k, v := range mapvals {
			mapval[k], err = MakeValue(v)
			if err != nil {
				return nil, err
			}
		}
		return &Value{
			Type: VAL_MAP,
			MapVal: mapval}, nil
	} else if listvals, ok := args.([]interface{}); ok {
		fmt.Println("MAKEVALUE list")
		listval := make([]*Value, len(listvals))
		for i, v := range listvals {
			listval[i], err = MakeValue(v)
			if err != nil {
				return nil, err
			}
		}
		return &Value{
			Type: VAL_LIST,
			ListVal: listval}, nil
	} else if boolval, ok := args.(bool); ok {
		return &Value{
			Type: VAL_BOOL,
			BoolVal: boolval}, nil
	} else if numval, ok := args.(float64); ok {
		return &Value{
			Type: VAL_NUM,
			NumVal: numval}, nil
	} else if numval, ok := args.(int); ok {
		return &Value{
			Type: VAL_NUM,
			NumVal: float64(numval)}, nil
	} else if strval, ok := args.(string); ok {
		return &Value{
			Type: VAL_STRING,
			StringVal: strval}, nil
	}
	return nil, errors.New(fmt.Sprintf("No conversion found for %v",args))
}

func MakeIntValue(x int) *Value {
	return &Value{
		Type: VAL_NUM,
		NumVal: float64(x)}
}

func MakeFloatValue(x float64) *Value {
	return &Value{
		Type: VAL_NUM,
		NumVal: float64(x)}
}

func MakeStringValue(x string) *Value {
	return &Value{
		Type: VAL_STRING,
		StringVal: x}
}

func MakeBoolValue(x bool) *Value {
	return &Value{
		Type: VAL_BOOL,
		BoolVal: x}
}

// func AssignValue(dst, src *Value, vars *map[string]*Value) {
// 	if src.Type == VAL_MAP {
// 		dst.MapVal = src.MapVal
// 	} else if src.Type == VAL_LIST {
// 		dst.ListVal = src.ListVal
// 	} else if src.Type == VAL_NAME {
// 		AssignValue(dst, (*vars)[src.NameVal], vars)
// 	} else if src.Type == VAL_INT {
// 		dst.IntVal = src.IntVal
// 	} else if src.Type == VAL_FLOAT {
// 		dst.FloatVal = src.FloatVal
// 	} else if src.Type == VAL_NUM {
// 		dst.NumVal = src.NumVal
// 	} else if src.Type == VAL_STRING {
// 		dst.StringVal = src.StringVal
// 	} else if src.Type == VAL_BOOL {
// 		dst.BoolVal = src.BoolVal
// 	}
// }

func (v *Value) ToString() string {
	if v.Type == VAL_MAP {
		map_keys := make([]string, 0)
		for k, _ := range v.MapVal {
			map_keys = append(map_keys, k)
		}
		sort.Strings(map_keys)
		val_strs := make([]string, len(map_keys))
		for i, k := range map_keys {
			val_strs[i] = fmt.Sprintf("%s:%s",k,v.MapVal[k].ToString())
		}
		return fmt.Sprintf("{%s}",strings.Join(val_strs ,","))
	} else if v.Type == VAL_LIST {
		val_strs := make([]string, len(v.ListVal))
		for i, val := range v.ListVal {
			val_strs[i] = val.ToString()
		}
		return fmt.Sprintf("[%s]",strings.Join(val_strs ,","))
	} else if v.Type == VAL_NAME {
		return fmt.Sprintf("VAR(%s)",v.NameVal)
	} else if v.Type == VAL_INT {
		return fmt.Sprintf("INT(%d)",v.IntVal)
	} else if v.Type == VAL_FLOAT {
		return fmt.Sprintf("FLOAT(%f)",v.FloatVal)
	} else if v.Type == VAL_NUM {
		return fmt.Sprintf("NUM(%f)",v.NumVal)
	} else if v.Type == VAL_STRING {
		return fmt.Sprintf("STRING(%s)",v.StringVal)
	} else if v.Type == VAL_BOOL {
		return fmt.Sprintf("BOOL(%v)",v.BoolVal)
	}
	return "[unknown type]"
}
