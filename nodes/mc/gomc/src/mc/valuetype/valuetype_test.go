package valuetype

import (
	"fmt"
	"testing"
	value "mc/value"
)

func TestInterfaceParser(t *testing.T) {
	ty, err := Parse(`{key:num}`)
	if err != nil {
		t.Errorf("%v", err)
	}
	tystr := ty.ToString()
	expected := `{key:num}`
	if tystr != expected {
		t.Errorf("Parsing problem--expected: %s, got: %s", expected, tystr)
	}
}

func TestInterfaceParserComplex(t *testing.T) {
	descs := []string{
		`{key1:[num]=[1,2,3],key2:{key3:oneof(string,num,[num])="hello"}}`,
		`{key1:[num]=[1,2],key2:{key3*:[oneof(string,num,[num])]}}`,
		`num`}
	for _, desc := range descs {
		ty, err := Parse(desc)
		if err != nil {
			t.Errorf("%v", err)
		}
		tystr := ty.ToString()
		if tystr != desc {
			t.Errorf("Parsing problem--expected: %s, got: %s", desc, tystr)
		}
	}
}

func TestValidate(t *testing.T) {
	ty, e1 := Parse(`{key:num}`)
	fmt.Println("ty",ty.ToString(), e1)
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


func TestValidateComplex(t *testing.T) {
	desc := `{key1:[num]=[1,2,3],key2:{key3:oneof(string,num,[num])="hello"}}`
	ty, _ := Parse(desc)
	fmt.Println(ty.ToString())
	v, _ := value.FromObject(map[string]interface{}{"key2":map[string]interface{}{}})
	ans, _ := value.FromObject(map[string]interface{}{"key1":[]interface{}{1,2,3},"key2":map[string]interface{}{"key3":"hello"}})
	nv, err := ty.Validate(v, map[string]*ValueType{}, "")
	if err != nil {
		t.Errorf("%v",err)
	}
	if nv == nil {
		t.Errorf("%s improperly failed validation",v.ToString())
	} else if !nv.Equals(ans) {
		t.Errorf("Expected %s, got %s", ans.ToString(), nv.ToString())
	}
}

