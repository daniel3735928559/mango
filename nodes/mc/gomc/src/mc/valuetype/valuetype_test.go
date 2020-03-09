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


func TestMakeTypes(t *testing.T) {
	observations := map[string]string{
		"bool":MakeBoolType().ToString(),
		"num":MakeNumType().ToString(),
		"string":MakeStringType().ToString(),
		"foo":MakeExtType("foo").ToString()}
	for ans, obs := range observations {
		if ans != obs {
			t.Errorf("%s type expected: `%s`, got: `%s`", ans, ans, obs)
		}
	}
}

func TestValidateComplex(t *testing.T) {
	desc := `{key1:[num]=[1,2,3],key2:{key3:oneof(bool,string,[num])=false}}`
	ty, _ := Parse(desc)
	fmt.Println(ty.ToString())
	v, _ := value.FromObject(map[string]interface{}{"key2":map[string]interface{}{}})
	ans, _ := value.FromObject(map[string]interface{}{"key1":[]interface{}{1,2,3},"key2":map[string]interface{}{"key3":false}})
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

type testcase struct {
	expected string
	inp interface{}
}

func TestValidateFailures(t *testing.T) {
	desc := `{key1:[num]=[1,2,3],key2:{key3:oneof(bool,string,[num])=false}}`
	testcases := []testcase{
		testcase{
			expected:`Error at root0: Invalid type: root0.key1`,
			inp:map[string]interface{}{"key1":"hello","key2":map[string]interface{}{}}},
		testcase{
			expected:`Error at root1: Invalid type: root1.key1`,
			inp:map[string]interface{}{"key1":"hello","key2":map[string]interface{}{}}},
		testcase{
			expected:`Error at root2: Invalid type: root2.key1`,
			inp:map[string]interface{}{"key1":"hello","key2":map[string]interface{}{}}},
		testcase{
			expected:`Error at root3: Invalid type: root3.key1`,
			inp:map[string]interface{}{"key1":"hello","key2":map[string]interface{}{}}}}
	ty, _ := Parse(desc)
	fmt.Println(ty.ToString())
	idx := 0
	for _, tc := range testcases {
		errmsg := tc.expected
		o := tc.inp
		v, _ := value.FromObject(o)
		_, err := ty.Validate(v, map[string]*ValueType{}, fmt.Sprintf("root%d",idx))
		if err == nil {
			t.Errorf("value was not supposed to validate: %s", v.ToString())
		} else if fmt.Sprintf("%v",err) != errmsg {
			t.Errorf("Expected error when processing `%s`: `%s`, got: `%v`", v.ToString(), errmsg, err)
		}
		idx++
	}
}

func TestValidateExtTypes(t *testing.T) {
	foodesc := `{key1:[num]=[1,2,3],key2:{key3:oneof(bool,string,[num])=false}}`
	footy, _ := Parse(foodesc)
	desc := `{key1*:[foo],key2:oneof(bool,foo)}`
	ty, _ := Parse(desc)
	fmt.Println(ty.ToString())
	testcases := []string{
		`{key2:true}`,
		`{key1:[{key1:[1],key2:{}}],key2:true}`,
		`{key1:[{key2:{}}],key2:{key2:{}}}`,
		`{key1:[{key2:{}}],key2:{key2:{}}}`,
		`{key1:[{key2:{}}],key2:{key1:[3,2,1],key2:{key3:"hello"}}}`}
	answers := []string{
		`{key2:true}`,
		`{key1:[{key1:[1],key2:{key3:false}}],key2:true}`,
		`{key1:[{key1:[1,2,3],key2:{key3:false}}],key2:{key1:[1,2,3],key2:{key3:false}}}`,
		`{key1:[{key1:[1,2,3],key2:{key3:false}}],key2:{key1:[1,2,3],key2:{key3:false}}}`,
		`{key1:[{key1:[1,2,3],key2:{key3:false}}],key2:{key1:[3,2,1],key2:{key3:"hello"}}}`}
	for i, tc := range testcases {
		fmt.Println("PARSING",tc)
		v, _ := value.Parse(tc)
		ans, _ := value.Parse(answers[i])
		nv, err := ty.Validate(v, map[string]*ValueType{"foo":footy}, "")
		if err != nil {
			t.Errorf("Failed to validate: %v", err)
		} else if nv == nil {
			t.Errorf("Failed to validate: nv is null; %s", v.ToString())
		} else {
			if !nv.Equals(ans) {
				t.Errorf("Testcase %d: Expected: %s, Got: %s", i, ans.ToString(), nv.ToString())
			}
		}
	}
}

