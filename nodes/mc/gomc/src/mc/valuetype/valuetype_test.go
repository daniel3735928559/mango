package valuetype

import (
	//"fmt"
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
	ty, _ := Parse(`{key:num}`)
	v, _ := value.FromObject(map[string]int{"key":45})
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

