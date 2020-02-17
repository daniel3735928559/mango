package routeparser

import (
	"fmt"
	"testing"
)

func TestParserWithoutTransform(t *testing.T) {
	fmt.Println("ASD")
	nodes := map[string]string{
		"node0": "a",
		"node1": "b",
	}
        examples := map[string][]string{
		"node0 > node1":[]string{"a > b"},
		"node0 <> node1":[]string{"a > b", "b > a"},
		"node0 < node1":[]string{"b > a"}}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs := Parse(s, nodes)
		fmt.Println("RS",rs)
		for i, r := range rs.Routes {
			fmt.Println(r.ToString())
			if r.ToString() != ans[i] {
				t.Errorf("expected: %s, got: %s", ans[i], r.ToString())
			}
		}
        }
}


func TestParserSingleSimpleTransforms(t *testing.T) {
	fmt.Println("ASD")
	nodes := map[string]string{
		"src": "src",
		"dst": "dst",
	}
        examples := map[string]string{
		`src > = {key1:"val1"} > dst`:`src > REPLACE(key1=STRING(val1)) > dst`,
		`src > ? {key1 == "val1"} > dst`:`src > IF(EQ(key1,STRING(val1))) > dst`,
		`src > % {key1 = "val1";} > dst`:`src > EDIT([SET(key1,STRING(val1))]) > dst`}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs := Parse(s, nodes)
		fmt.Println(rs)
		if len(rs.Routes) != 1 {
			t.Errorf("Expected 1 route in %s; got %d", s, len(rs.Routes))
		} else {
			r := rs.Routes[0]
			fmt.Println(r.ToString())
			if r.ToString() != ans {
				t.Errorf("expected: %s, got: %s", ans, r.ToString())
			}
		}
        }
}
