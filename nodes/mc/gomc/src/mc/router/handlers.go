package router

import (
	"fmt"
	"math"
	"regexp"
	"errors"
)

func CallHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	// TODO: error if called function does not exist
	name := args[0].NameVal
	fmt.Println("CALL",name)
	if name == "sub" {
		return this, nil
	}
	return nil, errors.New(fmt.Sprintf("No such function %s",name))
}
func TernaryHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	if args[0].BoolVal {
		return args[1], nil
	}
	return args[2], nil
}
func NameHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return primitive, nil
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
	item := args[1].NameVal
	if v, ok := args[0].MapVal[item]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("No such key %s found in map",item))
}
func ListGetHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	idx := uint(args[1].NumVal)
	if idx >= uint(len(args[0].ListVal)) {
		return nil, errors.New(fmt.Sprintf("Index out of bounds: %d", idx))
	}
	return args[0].ListVal[idx], nil
}
func StringGetHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	idx := int(args[1].NumVal)
	for i, rv := range args[0].StringVal {
		if i == idx {
			return MakeStringValue(string(rv)), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Index out of bounds: %d for %s", idx, args[0].StringVal))
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
func BoolHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(primitive.BoolVal), nil
}
func StringHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeStringValue(primitive.StringVal), nil
}
func MatchHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	matches, err := regexp.Match(args[1].StringVal, []byte(args[0].StringVal))
	if err != nil {
		return nil, err
	}
	return MakeBoolValue(matches), nil
}
func ExpNumHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	if args[0].NumVal == 0 && args[1].NumVal < 0 {
		return nil, errors.New("Divide by zero: zero to a negative power")
	}
	return MakeFloatValue(math.Pow(args[0].NumVal, args[1].NumVal)), nil
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
func AddListHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return &Value{
		Type:VAL_LIST,
		ListVal: append(args[0].ListVal, args[1].ListVal...)}, nil
}
func AddMapHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	ans := args[0].MapVal
	for k, v := range args[1].MapVal {
		ans[k] = v
	}
	return &Value{
		Type:VAL_MAP,
		MapVal: ans}, nil
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
	if args[1].NumVal == float64(0) {
		return nil, errors.New("Division by zero")
	}
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
func EqHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(args[0].Equals(args[1])), nil
}
func NeHandler(this *Value, local_vars map[string]*Value, args []*Value, primitive *Value) (*Value, error) {
	return MakeBoolValue(!args[0].Equals(args[1])), nil
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
