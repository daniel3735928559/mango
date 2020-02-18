package router

import (
	"fmt"
	"testing"
)

func TestParserWithoutTransform(t *testing.T) {
	fmt.Println("ASD")
        examples := map[string][]string{
		"node0 > node1":[]string{"node0 > node1"},
		"node0 <> node1":[]string{"node0 > node1", "node1 > node0"},
		"node0 < node1":[]string{"node1 > node0"}}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs := Parse(s)
		fmt.Println("RS",rs)
		for i, r := range rs {
			fmt.Println(r.ToString())
			if r.ToString() != ans[i] {
				t.Errorf("expected: %s, got: %s", ans[i], r.ToString())
			}
		}
        }
}

func TestParserSingleSimpleTransforms(t *testing.T) {
	fmt.Println("ASD")
        examples := map[string]string{
		`src > = {key1:"val1"} > dst`:`src > REPLACE(key1=STRING(val1)) > dst`,
		`src > ? {key1 == "val1"} > dst`:`src > IF(EQ(key1,STRING(val1))) > dst`,
		`src > % {key1 = "val1";} > dst`:`src > EDIT([SET(key1,STRING(val1))]) > dst`}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs := Parse(s)
		fmt.Println(rs)
		if len(rs) != 1 {
			t.Errorf("Expected 1 route in %s; got %d", s, len(rs))
		} else {
			r := rs[0]
			fmt.Println(r.ToString())
			if r.ToString() != ans {
				t.Errorf("expected: %s, got: %s", ans, r.ToString())
			}
		}
        }
}


func TestRouter(t *testing.T) {
	r := MakeRouter()

	for i := 0; i < 3; i++ {
		r.AddNode(&Node{Name:fmt.Sprintf("node%d",i)})

	}
	r.ParseAndAddRoutes("node0 > node1")
	r.ParseAndAddRoutes(`node0 > ? {key1=="val1"} > node1`)
}
