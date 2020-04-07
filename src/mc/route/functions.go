package route

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"regexp"
	"math/rand"
	"encoding/json"
	value "mc/value"
)

type FuncHandler func(*value.Value, []*value.Value) (*value.Value, error) // this, args -> result, err

type FunctionSignature struct {
	Name string
	ArgTypes []value.ValueKind
	ReturnType value.ValueKind
	Handler FuncHandler
}

var (
	ExpressionFunctions = []*FunctionSignature{
		&FunctionSignature{
			Name: "raw",
			ArgTypes: []value.ValueKind{},
			ReturnType: value.VAL_STRING,
			Handler: RawFunc},
		&FunctionSignature{
			Name: "now",
			ArgTypes: []value.ValueKind{},
			ReturnType: value.VAL_NUM,
			Handler: NowFunc},
		&FunctionSignature{
			Name: "split",
			ArgTypes: []value.ValueKind{value.VAL_STRING, value.VAL_STRING},
			ReturnType: value.VAL_LIST,
			Handler: SplitFunc},
		&FunctionSignature{
			Name: "match",
			ArgTypes: []value.ValueKind{value.VAL_STRING, value.VAL_STRING},
			ReturnType: value.VAL_STRING,
			Handler: MatchFunc},
		&FunctionSignature{
			Name: "int",
			ArgTypes: []value.ValueKind{value.VAL_STRING},
			ReturnType: value.VAL_NUM,
			Handler: IntFunc},
		&FunctionSignature{
			Name: "randint",
			ArgTypes: []value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: RandIntFunc}}
)

func GetFunction(name string, args []*value.Value) (FuncHandler, error) {
	for _, s := range ExpressionFunctions {
		if s.Name != name {
			continue
		}
		if len(s.ArgTypes) != len(args) {
			continue
		}
		ok := true
		for i, ty := range s.ArgTypes {
			if ty != value.VAL_ANY && args[i].Type != ty {
				ok = false
				break
			}
		}
		if ok {
			return s.Handler, nil
		}
	}
	return nil, fmt.Errorf("No matching function signature found for: %s", name)
}

func RawFunc(this *value.Value, args []*value.Value) (*value.Value, error) {
	raw_data, _ := json.Marshal(this.ToObject())
	raw_val, _ := value.FromObject(string(raw_data))
	return raw_val, nil
}

func NowFunc(this *value.Value, args []*value.Value) (*value.Value, error) {
	time_val, _ := value.FromObject(float64(time.Now().UnixNano())/1000000.0)
	return time_val, nil
}

func SplitFunc(this *value.Value, args []*value.Value) (*value.Value, error) {
	split := make([]interface{}, 0)
	fmt.Println("SPLIT",args[0].StringVal,args[1].StringVal)
	for _, s := range strings.Split(args[0].StringVal, args[1].StringVal) {
		fmt.Println("SS",s)
		split = append(split, s)
	}
	ans, err := value.FromObject(split)
	if err != nil {
return nil, err
	}
	return ans, nil
}

func MatchFunc(this *value.Value, args []*value.Value) (*value.Value, error) {
	re, err := regexp.Compile(args[1].StringVal)
	if err != nil {
		return nil, err
	}
	ans := value.MakeStringValue(re.FindString(args[0].StringVal))
	return ans, nil
}

func IntFunc(this *value.Value, args []*value.Value) (*value.Value, error) {
	ans, err := strconv.Atoi(args[0].StringVal)
	if err != nil {
		return value.MakeIntValue(0), nil
	}
	return value.MakeIntValue(ans), nil
}

func RandIntFunc(this *value.Value, args []*value.Value) (*value.Value, error) {
	x := int(args[0].NumVal)
	y := int(args[1].NumVal)
	if x > y {
		return nil, fmt.Errorf("Invalid arguments: %d > %d", x, y)
	}
	ans := rand.Intn(y-x+1)+x
	return value.MakeIntValue(ans), nil
}
