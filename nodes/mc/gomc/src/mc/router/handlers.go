package router

import (
	"fmt"
	"errors"
)

func AssignHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	// if variable, ok := (*vars)[args[0].NameVal]; !ok {
	// 	return nil, errors.New(fmt.Sprintf("No such variable: %s", args[0].NameVal))
	// }
	// if args[0].PathVal != nil && len(args[0].PathVal) > 0 {
	// 	v := ResolvePathValue(variable, args[0].PathVal)
	// }
	return nil, nil
}
func CallHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	// TODO: error if called function does not exist
	fmt.Println("CALL",args[0].NameVal)
	return nil, nil
}
func MapHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	// TODO: Error if odd number of args or if keys are not Type:"name"
	mapval := make(map[string]*Value)
	for i := 0; i < len(args); i += 2 {
		mapval[args[i].NameVal] = args[i+1]
	}
	return &Value{
		Type:VAL_MAP,
		MapVal: mapval}, nil
}
func ListHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return &Value{
		Type:VAL_LIST,
		ListVal: args}, nil
}
func MapGetHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	// TODO: error if key does not exist
	item := args[1].NameVal
	return args[0].MapVal[item], nil
}
func ListGetHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	// TODO: error if index out of bounds
	idx := uint(args[1].NumVal)
	if idx > uint(len(args[0].ListVal)) {
		return nil, errors.New(fmt.Sprintf("Index out of bounds: %d", idx))
	}
	return args[0].ListVal[idx], nil
}
func ExprHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return nil, nil
}
func ValHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return nil, nil
}
func UnaryMinusNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(-args[0].NumVal), nil
}
func AddNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(args[0].NumVal + args[1].NumVal), nil
}
func AddStringHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeStringValue(args[0].StringVal + args[1].StringVal), nil
}
func SubNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(args[0].NumVal - args[1].NumVal), nil
}
func MulNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(args[0].NumVal * args[1].NumVal), nil
}
func MulStringNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	ans := args[0].StringVal
	// TODO: error if count < 0 or count is not int
	for i := 0; uint(i) < uint(args[1].NumVal); i++ {
		ans += args[0].StringVal
	}
	return MakeStringValue(ans), nil
}
func DivNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	// TODO: error arg[1] == 0
	return MakeFloatValue(args[0].NumVal / args[1].NumVal), nil
}
func ModNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) % int(args[1].NumVal))), nil
}
func XorNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) ^ int(args[1].NumVal))), nil
}
func AndNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) & int(args[1].NumVal))), nil
}
func OrNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeFloatValue(float64(int(args[0].NumVal) | int(args[1].NumVal))), nil
}
func LtNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal < args[1].NumVal), nil
}
func LtStringHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal < args[1].StringVal), nil
}
func GtNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal > args[1].NumVal), nil
}
func GtStringHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal > args[1].StringVal), nil
}
func EqNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal == args[1].NumVal), nil
}
func EqStringHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal == args[1].StringVal), nil
}
func LeqNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal <= args[1].NumVal), nil
}
func LeqStringHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal <= args[1].StringVal), nil
}
func GeqNumHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].NumVal >= args[1].NumVal), nil
}
func GeqStringHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].StringVal >= args[1].StringVal), nil
}
func OrBoolHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].BoolVal || args[1].BoolVal), nil
}
func AndBoolHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(args[0].BoolVal && args[1].BoolVal), nil
}
func NotBoolHandler(args []*Value, vars *map[string]*Value) (*Value, error) {
	return MakeBoolValue(!args[0].BoolVal), nil
}
