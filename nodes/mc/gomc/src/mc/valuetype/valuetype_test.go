package valuetype

import (
	"fmt"
	"testing"
	value "mc/value"
)

func TestInterfaceParser(t *testing.T) {
	_, err := Parse(`{key:num}`)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestValidate(t *testing.T) {
	ty, e1 := Parse(`{key:num}`)
	fmt.Println("ty",ty,e1)
	v, e2 := value.FromObject(map[string]interface{}{"key":45})
	fmt.Println("v",v,e2)
	nv, err := ty.Validate(v, map[string]*ValueType{}, "")
	if err != nil {
		t.Errorf("%v", err)
	}
	if val, ok := nv.MapVal["key"]; ok {
		if val.NumVal != float64(45) {
			t.Errorf("Expected 45")
		}
	} else {
		t.Errorf("Expected `key`")
	}
}

