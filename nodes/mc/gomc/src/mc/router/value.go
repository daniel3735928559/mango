package router

import (
	"fmt"
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

func AssignValue(dst, src *Value, vars *map[string]*Value) {
	if src.Type == VAL_MAP {
		dst.MapVal = src.MapVal
	} else if src.Type == VAL_LIST {
		dst.ListVal = src.ListVal
	} else if src.Type == VAL_NAME {
		AssignValue(dst, (*vars)[src.NameVal], vars)
	} else if src.Type == VAL_INT {
		dst.IntVal = src.IntVal
	} else if src.Type == VAL_FLOAT {
		dst.FloatVal = src.FloatVal
	} else if src.Type == VAL_NUM {
		dst.NumVal = src.NumVal
	} else if src.Type == VAL_STRING {
		dst.StringVal = src.StringVal
	} else if src.Type == VAL_BOOL {
		dst.BoolVal = src.BoolVal
	}
}

func (v *Value) ToString() string {
	if v.Type == VAL_MAP {
		return fmt.Sprintf("%v",v.MapVal)
	} else if v.Type == VAL_LIST {
		return fmt.Sprintf("%v",v.ListVal)
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
