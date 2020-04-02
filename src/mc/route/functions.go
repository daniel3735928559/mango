package route

import (
	"fmt"
	"time"
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
			Handler: NowFunc}}
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
