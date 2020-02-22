package router

import (
	"fmt"
	"errors"
)

func CallHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	// TODO: error if called function does not exist
	fmt.Println("CALL",args[0].NameVal)
	return this.Clone(), nil
}
func MapHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	// TODO: Error if odd number of args or if keys are not Type:"name"
	mapval := make(map[string]*Value)
	for i := 0; i < len(args); i += 2 {
		mapval[args[i].NameVal] = args[i+1]
	}
	return &Value{
		Type:VAL_MAP,
		MapVal: mapval}, nil
}
func ListHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return &Value{
		Type:VAL_LIST,
		ListVal: args}, nil
}
func MapGetHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	// TODO: error if key does not exist
	item := args[1].NameVal
	return args[0].MapVal[item], nil
}
func ListGetHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	// TODO: error if index out of bounds
	idx := uint(args[1].NumVal)
	if idx > uint(len(args[0].ListVal)) {
		return nil, errors.New(fmt.Sprintf("Index out of bounds: %d", idx))
	}
	return args[0].ListVal[idx], nil
}
// func ExprHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
// 	return nil, nil, nil
// }
// func ValHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
// 	return nil, nil, nil
// }
func VarHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	fmt.Println("VAR Handler",primitive)
	if primitive.NameVal == "this" {
		return this, nil
	}
	if this.Type == VAL_MAP {
		if v, ok := this.MapVal[primitive.NameVal]; ok {
			return v, nil
		}
	}
	if v, ok := local_vars[primitive.NameVal]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("No such variable: %s",primitive.NameVal))
}
func NumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(primitive.NumVal), nil
}
func StringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeStringValue(primitive.StringVal), nil
}
func UnaryMinusNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(-args[0].NumVal), nil
}
func AddNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(args[0].NumVal + args[1].NumVal), nil
}
func AddStringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeStringValue(args[0].StringVal + args[1].StringVal), nil
}
func SubNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(args[0].NumVal - args[1].NumVal), nil
}
func MulNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(args[0].NumVal * args[1].NumVal), nil
}
func MulStringNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	ans := args[0].StringVal
	// TODO: error if count < 0 or count is not int
	for i := 1; uint(i) < uint(args[1].NumVal); i++ {
		ans += args[0].StringVal
	}
	return MakeStringValue(ans), nil
}
func DivNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	// TODO: error arg[1] == 0
	return MakeFloatValue(args[0].NumVal / args[1].NumVal), nil
}
func ModNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) % int(args[1].NumVal))), nil
}
func XorNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) ^ int(args[1].NumVal))), nil
}
func AndNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) & int(args[1].NumVal))), nil
}
func OrNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) | int(args[1].NumVal))), nil
}
func LtNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal < args[1].NumVal), nil
}
func LtStringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal < args[1].StringVal), nil
}
func GtNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal > args[1].NumVal), nil
}
func GtStringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal > args[1].StringVal), nil
}
func EqNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal == args[1].NumVal), nil
}
func EqStringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal == args[1].StringVal), nil
}
func LeqNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal <= args[1].NumVal), nil
}
func LeqStringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal <= args[1].StringVal), nil
}
func GeqNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal >= args[1].NumVal), nil
}
func GeqStringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal >= args[1].StringVal), nil
}
func OrBoolHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].BoolVal || args[1].BoolVal), nil
}
func AndBoolHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].BoolVal && args[1].BoolVal), nil
}
func NotBoolHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(!args[0].BoolVal), nil
}
