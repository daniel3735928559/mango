package routeparser

import (
	"fmt"
	"testing"
)

func TestPTParserWithoutTransform(t *testing.T) {
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
		if len(rs.Routes) != 1 {
			t.Errorf("Expected 1 route in %s; got %d", s, len(rs.Routes))
		}
		for i, r := range rs.Routes {
			fmt.Println(r.ToString())
			if r.ToString() != ans[i] {
				t.Errorf("expected: %s, got: %s", ans[i], r.ToString())
			}
		}
        }
}
