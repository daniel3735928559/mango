package route

import (
	"fmt"
	"math"
	"regexp"
	"errors"
	value "mc/value"
)

func CallHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	// TODO: error if called function does not exist
	name := args[0].NameVal
	fmt.Println("CALL",name)
	if name == "sub" {
		return this, nil
	}
	return nil, errors.New(fmt.Sprintf("No such function %s",name))
}
func TernaryHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	if args[0].BoolVal {
		return args[1], nil
	}
	return args[2], nil
}
func NameHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return primitive, nil
}
func MapHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	// TODO: Error if odd number of args or if keys are not Type:"name"
	mapval := make(map[string]*value.Value)
	for i := 0; i < len(args); i += 2 {
		mapval[args[i].NameVal] = args[i+1]
	}
	return &value.Value{
		Type:value.VAL_MAP,
		MapVal: mapval}, nil
}
func ListHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return &value.Value{
		Type:value.VAL_LIST,
		ListVal: args}, nil
}
func MapGetHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	item := args[1].NameVal
	if v, ok := args[0].MapVal[item]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("No such key %s found in map",item))
}
func ListGetHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	idx := uint(args[1].NumVal)
	if idx >= uint(len(args[0].ListVal)) {
		return nil, errors.New(fmt.Sprintf("Index out of bounds: %d", idx))
	}
	return args[0].ListVal[idx], nil
}
func StringGetHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	idx := int(args[1].NumVal)
	for i, rv := range args[0].StringVal {
		if i == idx {
			return value.MakeStringValue(string(rv)), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Index out of bounds: %d for %s", idx, args[0].StringVal))
}
// func ExprHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
// 	return nil, nil, nil
// }
// func ValHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
// 	return nil, nil, nil
// }
func VarHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	fmt.Println("VAR Handler",primitive)
	if primitive.NameVal == "this" {
		return this, nil
	}
	if this.Type == value.VAL_MAP {
		if v, ok := this.MapVal[primitive.NameVal]; ok {
			return v, nil
		}
	}
	if v, ok := local_vars[primitive.NameVal]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("No such variable: %s",primitive.NameVal))
}
func NumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(primitive.NumVal), nil
}
func BoolHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(primitive.BoolVal), nil
}
func StringHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeStringValue(primitive.StringVal), nil
}
func MatchHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	matches, err := regexp.Match(args[1].StringVal, []byte(args[0].StringVal))
	if err != nil {
		return nil, err
	}
	return value.MakeBoolValue(matches), nil
}
func ExpNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	if args[0].NumVal == 0 && args[1].NumVal < 0 {
		return nil, errors.New("Divide by zero: zero to a negative power")
	}
	return value.MakeFloatValue(math.Pow(args[0].NumVal, args[1].NumVal)), nil
}
func UnaryMinusNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(-args[0].NumVal), nil
}
func AddNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(args[0].NumVal + args[1].NumVal), nil
}
func AddStringHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeStringValue(args[0].StringVal + args[1].StringVal), nil
}
func AddListHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return &value.Value{
		Type:value.VAL_LIST,
		ListVal: append(args[0].ListVal, args[1].ListVal...)}, nil
}
func AddMapHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	ans := args[0].MapVal
	for k, v := range args[1].MapVal {
		ans[k] = v
	}
	return &value.Value{
		Type:value.VAL_MAP,
		MapVal: ans}, nil
}
func SubNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(args[0].NumVal - args[1].NumVal), nil
}
func MulNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(args[0].NumVal * args[1].NumVal), nil
}
func MulStringNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	ans := args[0].StringVal
	// TODO: error if count < 0 or count is not int
	for i := 1; uint(i) < uint(args[1].NumVal); i++ {
		ans += args[0].StringVal
	}
	return value.MakeStringValue(ans), nil
}
func DivNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	if args[1].NumVal == float64(0) {
		return nil, errors.New("Division by zero")
	}
	return value.MakeFloatValue(args[0].NumVal / args[1].NumVal), nil
}
func ModNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(float64(int(args[0].NumVal) % int(args[1].NumVal))), nil
}
func XorNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(float64(int(args[0].NumVal) ^ int(args[1].NumVal))), nil
}
func AndNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(float64(int(args[0].NumVal) & int(args[1].NumVal))), nil
}
func OrNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeFloatValue(float64(int(args[0].NumVal) | int(args[1].NumVal))), nil
}
func LtNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].NumVal < args[1].NumVal), nil
}
func LtStringHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].StringVal < args[1].StringVal), nil
}
func GtNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].NumVal > args[1].NumVal), nil
}
func GtStringHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].StringVal > args[1].StringVal), nil
}
func EqHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].Equals(args[1])), nil
}
func NeHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(!args[0].Equals(args[1])), nil
}
func LeqNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].NumVal <= args[1].NumVal), nil
}
func LeqStringHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].StringVal <= args[1].StringVal), nil
}
func GeqNumHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].NumVal >= args[1].NumVal), nil
}
func GeqStringHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].StringVal >= args[1].StringVal), nil
}
func OrBoolHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].BoolVal || args[1].BoolVal), nil
}
func AndBoolHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(args[0].BoolVal && args[1].BoolVal), nil
}
func NotBoolHandler(this *value.Value, local_vars map[string]*value.Value, args []*value.Value, primitive *value.Value) (*value.Value, error) {
	return value.MakeBoolValue(!args[0].BoolVal), nil
}
