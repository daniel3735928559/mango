package value

import (
	"testing"
)

func TestFromObject(t *testing.T) {
	objs := []interface{}{
		map[string]interface{}{"key1":[]interface{}{1,2,3},"key2":map[string]interface{}{"key3":"hello"}},
		3,
		[]interface{}{"a","b","c"},
		map[string]interface{}{"key1":2}}
	expected := []string{
		`{"key1":[1,2,3],"key2":{"key3":"hello"}}`,
		`3`,
		`["a","b","c"]`,
		`{"key1":2}`}
	for i, o := range objs {
		v, err := FromObject(o)
		if err != nil {
			t.Errorf("Failed converting object %v: %v", o, err)
		}
		if v.ToString() != expected[i] {
			t.Errorf("Unexpected conversion--expected %s, got %s", expected[i], v.ToString())
		}
	}
}

func TestToObject(t *testing.T) {
	objs := []interface{}{
		map[string]interface{}{"key1":[]interface{}{1,2,3},"key2":map[string]interface{}{"key3":"hello"}},
		3,
		[]interface{}{"a","b","c"},
		map[string]interface{}{"key1":2}}
	expected := []string{
		`{"key1":[1,2,3],"key2":{"key3":"hello"}}`,
		`3`,
		`["a","b","c"]`,
		`{"key1":2}`}
	for i, o := range objs {
		v, _ := FromObject(o)
		o2 := v.ToObject()
		v2, _ := FromObject(o2)
		if v2.ToString() != expected[i] {
			t.Errorf("Unexpected conversion--expected %s, got %s", expected[i], v.ToString())
		}
	}
}

func TestClone(t *testing.T) {
	objs := []interface{}{
		map[string]interface{}{"key1":[]interface{}{1,2,3},"key2":map[string]interface{}{"key3":"hello"}},
		3,
		[]interface{}{"a","b","c"},
		map[string]interface{}{"key1":2}}
	for _, o := range objs {
		v, _ := FromObject(o)
		v2 := v.Clone()
		if !v.Equals(v2) {
			t.Errorf("Unexpected clone result--expected %s, got %s", v2.ToString(), v.ToString())
		}
	}
}


func TestMakeValues(t *testing.T) {
	observations := map[string]string{
		"true":MakeBoolValue(true).ToString(),
		"3":MakeIntValue(3).ToString(),
		"4":MakeFloatValue(4).ToString(),
		`"hello"`:MakeStringValue("hello").ToString()}
	for ans, obs := range observations {
		if ans != obs {
			t.Errorf("%s type expected: `%s`, got: `%s`", ans, ans, obs)
		}
	}
}

func TestParse(t *testing.T) {
	objs := []interface{}{
		map[string]interface{}{"key1":[]interface{}{1,2,3},"key2":map[string]interface{}{"key3":"hello"}},
		3,
		[]interface{}{"a","b","c"},
		map[string]interface{}{"key1":2},
		true,
		[]interface{}{},
		map[string]interface{}{}}
	ans := []string{
		`{key1:[1,2,3],key2:{key3:"hello"}}`,
		`3`,
		`["a","b","c"]`,
		`{key1:2}`,
                "true",
		"[]",
		"{}"}
	for i, o := range objs {
		v, _ := FromObject(o)
		v2, err := Parse(ans[i])
		if err != nil {
			t.Errorf("Failed parsing object %s: %v", ans[i], err)
		}
		if v2 == nil {
			t.Errorf("Failed to parse: %s", ans[i])
		} else if v.ToString() != v2.ToString() {
			t.Errorf("Unexpected conversion--expected %s, got %s", v2.ToString(), v.ToString())
		}
	}
}
