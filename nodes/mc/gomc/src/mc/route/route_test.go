package router

import (
	"fmt"
	"testing"
	value "mc/value"
)

func TestParserError(t *testing.T) {
        examples := map[string][]string{
		"node0 >":[]string{"node0 > node1"},
		"node0 <>":[]string{"node0 > node1", "node1 > node0"},
		"node0 < {}":[]string{"node1 > node0"}}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs, _ := Parse(s)
		fmt.Println("RS",rs,ans)
		if len(rs) > 0 {
			t.Errorf("expected no routes, got: %d", len(rs))
		}
        }
}

func TestParserWithoutTransform(t *testing.T) {
        examples := map[string][]string{
		"node0 > node1":[]string{"node0 > node1"},
		"node0 <> node1":[]string{"node0 > node1", "node1 > node0"},
		"node0 < node1":[]string{"node1 > node0"}}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs, _ := Parse(s)
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
		`src > = {key1:"val1"} > dst`:`src > replace {key1:"val1"} > dst`,
		`src > ? {key1 == "val1"} > dst`:`src > pass if {key1 == "val1"} > dst`,
		`src > % {key1 = "val1";} > dst`:`src > edit {key1 = "val1";} > dst`}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs, _ := Parse(s)
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

func TestParserSingleComplexTransforms(t *testing.T) {
	fmt.Println("ASD")
        examples := map[string]string{
		`src > ? {!(key2 >= 4 || key2 < 8)} = {key1:"val1",key2:-99,key3:1+1+1^6} > dst`: `src > replace {key1:"val1",key2:-99,key3:1 + 1 + 1 ^ 6} if {!key2 >= 4 || key2 < 8} > dst`,
		`src > ? {key1 == "val1" && !(key2 == 4 || key2 >= 8)} > dst`: `src > pass if {key1 == "val1" && !key2 == 4 || key2 >= 8} > dst`,
		`src > ? {key1 == "val1" && !(key2 <= 4 || key3 < -100.9)} % {key1 = "val1"; key2[4] = key3.key4[1+1];} > dst`: `src > edit {key1 = "val1";key2[4] = key3.key4[1 + 1];} if {key1 == "val1" && !key2 <= 4 || key3 < -100.9} > dst`}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs, _ := Parse(s)
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

func TestParserSingleComplexCollectionTransforms(t *testing.T) {
	fmt.Println("ASD")
        examples := map[string]string{
		`src > % {key2.key3.key4 += 5;} > dst`: `src > edit {key2.key3.key4 = key2.key3.key4 + 5;} > dst`,
		`src > % {key2[9].key4 += 5;} > dst`: `src > edit {key2[9].key4 = key2[9].key4 + 5;} > dst`,
		`src > % {key2.key3[1] += 5;} > dst`: `src > edit {key2.key3[1] = key2.key3[1] + 5;} > dst`}
        for s, ans := range examples {
		fmt.Println("PARSING",s)
		rs, _ := Parse(s)
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

func RunMessagesThroughRoutes(t *testing.T, routes []string, messages []map[string]interface{}, expected map[string][]string) {
	// assumes nodes node0, node1, node2
	// args: 
	// - routes: list of routes to add to router
	// - messages: nodename -> message to send from that node
	// returns:
	// - nodename -> list of stringified messages received by that node

	route_list := make([]*Route, 0)
	
	for _, routespec := range routes {
		rts, err := Parse(routespec)
		if err != nil {
			t.Errorf(fmt.Sprintf("Error parsing route: %s, Error: %v", routespec, err))
		} else {
			route_list = append(route_list, rts...)
		}
	}
	
	results := make(map[string][]string)
	for _, rt := range route_list {
		results[rt.Dest] = make([]string, 0)
		for _, msg := range messages {
			args, err := value.FromObject(msg)
			if err != nil {
				t.Errorf("Error converting %v",msg)
			}
			new_cmd, new_args, _ := rt.Run("test_cmd", args)
			if new_args != nil {
				fmt.Printf("%s recvd %s\n",rt.Dest,new_args.ToString())
				results[rt.Dest] = append(results[rt.Dest], fmt.Sprintf("%s %s", new_cmd, new_args.ToString()))
			}
		}
	}

	for dst_node, answers := range results {
		if len(answers) != len(expected[dst_node]) {
			t.Errorf("Expected %s to receive %d messages; got: %d", dst_node, len(expected[dst_node]), len(answers))
		} else {
			for i, answer := range answers {
				if expected[dst_node][i] != answer {
					t.Errorf("Expected %s received message %d to be %s; got: %s", dst_node, i, expected[dst_node][i], answer)
				}
			}
		}
	}
}


func TestRouterEqFilter(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > ? {key1 > "val0"} > node2`,
		`node0 > ? {key1==3} > node3`,
		`node0 > ? {key1==[1,2]} > node4`,
		`node0 > ? {key1=={b:2,a:1}} > node5`,
		`node0 > ? {key1 < "val2"} > node6`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":"val2"},
		map[string]interface{}{"key1":3},
		map[string]interface{}{"key1":4},
		map[string]interface{}{"key1":[]interface{}{1,2}},
		map[string]interface{}{"key1":map[string]interface{}{"a":1,"b":2}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":"val2"}`,`test_cmd {"key1":3}`,`test_cmd {"key1":4}`,`test_cmd {"key1":[1,2]}`,`test_cmd {"key1":{"a":1,"b":2}}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":"val2"}`},
		"node3":[]string{`test_cmd {"key1":3}`},
		"node4":[]string{`test_cmd {"key1":[1,2]}`},
		"node5":[]string{`test_cmd {"key1":{"a":1,"b":2}}`},
		"node6":[]string{`test_cmd {"key1":"val1"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterNeFilter(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > ? {key1 != "val1"} > node2`,
		`node0 > ? {key1 != 4} > node3`,
		`node0 > ? {key1 != [1,2]} > node4`,
		`node0 > ? {key1 != {a:1,b:2}} > node5`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":"val2"},
		map[string]interface{}{"key1":3},
		map[string]interface{}{"key1":4},
		map[string]interface{}{"key1":[]interface{}{1,2}},
		map[string]interface{}{"key1":map[string]interface{}{"a":1,"b":2}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":"val2"}`,`test_cmd {"key1":3}`,`test_cmd {"key1":4}`,`test_cmd {"key1":[1,2]}`,`test_cmd {"key1":{"a":1,"b":2}}`},
		"node2":[]string{`test_cmd {"key1":"val2"}`,`test_cmd {"key1":3}`,`test_cmd {"key1":4}`,`test_cmd {"key1":[1,2]}`,`test_cmd {"key1":{"a":1,"b":2}}`},
		"node3":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":"val2"}`,`test_cmd {"key1":3}`,`test_cmd {"key1":[1,2]}`,`test_cmd {"key1":{"a":1,"b":2}}`},
		"node4":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":"val2"}`,`test_cmd {"key1":3}`,`test_cmd {"key1":4}`,`test_cmd {"key1":{"a":1,"b":2}}`},
		"node5":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":"val2"}`,`test_cmd {"key1":3}`,`test_cmd {"key1":4}`,`test_cmd {"key1":[1,2]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterGtFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 > 0} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`test_cmd {"key1":10}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterLtFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 < 0} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`test_cmd {"key1":-1}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterGeFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 >= 0} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`test_cmd {"key1":10}`,`test_cmd {"key1":0}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterLeFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 <= 0} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`test_cmd {"key1":-1}`,`test_cmd {"key1":0}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterAndFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 >= 0 && key1 < 5} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":2},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`test_cmd {"key1":2}`,`test_cmd {"key1":0}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":2}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterOrFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 < 0 || key1 >= 10} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":2},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":2}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterNotFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {!(key1 <= 0)} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-1},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":10}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterListFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1[0] <= 0} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":[]interface{}{-1, 8}},
		map[string]interface{}{"key1":[]interface{}{10, -2}},
		map[string]interface{}{"key1":[]interface{}{0, 100}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":[-1,8]}`,`test_cmd {"key1":[0,100]}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":[-1,8]}`,`test_cmd {"key1":[10,-2]}`,`test_cmd {"key1":[0,100]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterArithmeticEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 += 1;} > node2`,
		`node0 > % {key1 -= 1;} > node3`,
		`node0 > % {key1 *= 2;} > node4`,
		`node0 > % {key1 /= 2;} > node5`,
		`node0 > % {key1 %= 2;} > node6`,
		`node0 > % {key1 = -key1;} > node7`,
		`node0 > % {key1 = 4-key1;} > node8`,
		`node0 > % {key1 = key1**3;} > node9`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-2},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-2}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`},
		"node2":[]string{`test_cmd {"key1":-1}`,`test_cmd {"key1":11}`,`test_cmd {"key1":1}`},
		"node3":[]string{`test_cmd {"key1":-3}`,`test_cmd {"key1":9}`,`test_cmd {"key1":-1}`},
		"node4":[]string{`test_cmd {"key1":"val1val1"}`,`test_cmd {"key1":-4}`,`test_cmd {"key1":20}`,`test_cmd {"key1":0}`},
		"node5":[]string{`test_cmd {"key1":-1}`,`test_cmd {"key1":5}`,`test_cmd {"key1":0}`},
		"node6":[]string{`test_cmd {"key1":0}`,`test_cmd {"key1":0}`,`test_cmd {"key1":0}`},
		"node7":[]string{`test_cmd {"key1":2}`,`test_cmd {"key1":-10}`,`test_cmd {"key1":-0}`},
		"node8":[]string{`test_cmd {"key1":6}`,`test_cmd {"key1":-6}`,`test_cmd {"key1":4}`},
		"node9":[]string{`test_cmd {"key1":-8}`,`test_cmd {"key1":1000}`,`test_cmd {"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterComplexArithmeticEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 = (key1+key1*3)/2;} > node2`,
		`node0 > % {key1 = (8/2-1)+key1%3;} > node3`,
		`node0 > % {key1 = (2*3+1+5*6)/(0);} > node4`,
		`node0 > % {key1 = (1.1-8.1)/(0.5*14);} > node5`,
		`node0 > % {key1 = 4/2/2;} > node6`,
		`node0 > % {key1 = 1--1-1;} > node7`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":-2},
		map[string]interface{}{"key1":10},
		map[string]interface{}{"key1":0}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":-2}`,`test_cmd {"key1":10}`,`test_cmd {"key1":0}`},
		"node2":[]string{`test_cmd {"key1":-4}`,`test_cmd {"key1":20}`,`test_cmd {"key1":0}`},
		"node3":[]string{`test_cmd {"key1":1}`,`test_cmd {"key1":4}`,`test_cmd {"key1":3}`},
		"node4":[]string{},
		"node5":[]string{`test_cmd {"key1":-1}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":-1}`,`test_cmd {"key1":-1}`},
		"node6":[]string{`test_cmd {"key1":1}`,`test_cmd {"key1":1}`,`test_cmd {"key1":1}`,`test_cmd {"key1":1}`},
		"node7":[]string{`test_cmd {"key1":1}`,`test_cmd {"key1":1}`,`test_cmd {"key1":1}`,`test_cmd {"key1":1}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterBitwiseEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1&=1;} > node2`,
		`node0 > % {key1|=1;} > node3`,
		`node0 > % {key1^=1;} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":0},
		map[string]interface{}{"key1":1},
		map[string]interface{}{"key1":2},
		map[string]interface{}{"key1":2.1}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":0}`,`test_cmd {"key1":1}`,`test_cmd {"key1":2}`,`test_cmd {"key1":2.1}`},
		"node2":[]string{`test_cmd {"key1":0}`,`test_cmd {"key1":1}`,`test_cmd {"key1":0}`,`test_cmd {"key1":0}`},
		"node3":[]string{`test_cmd {"key1":1}`,`test_cmd {"key1":1}`,`test_cmd {"key1":3}`,`test_cmd {"key1":3}`},
		"node4":[]string{`test_cmd {"key1":1}`,`test_cmd {"key1":0}`,`test_cmd {"key1":3}`,`test_cmd {"key1":3}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterBoolEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 = true;} > node2`,
		`node0 > % {key1 = false;} > node3`,
		`node0 > % {key1 = !key1;} > node4`,
		`node0 > % {key1 = key1 && false || !!true;} > node5`,
		`node0 > % {key1 = !key1 || !!false;} > node6`,
		`node0 > % {key1 = key1 || true && key1;} > node7`,
		`node0 > % {key1 = key1 || !key1;} > node8`,
		`node0 > % {key1 = key1 == true && !key1 == false;} > node9`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":true},
		map[string]interface{}{"key1":false}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":true}`,`test_cmd {"key1":false}`},
		"node2":[]string{`test_cmd {"key1":true}`,`test_cmd {"key1":true}`},
		"node3":[]string{`test_cmd {"key1":false}`,`test_cmd {"key1":false}`},
		"node4":[]string{`test_cmd {"key1":false}`,`test_cmd {"key1":true}`},
		"node5":[]string{`test_cmd {"key1":true}`,`test_cmd {"key1":true}`},
		"node6":[]string{`test_cmd {"key1":false}`,`test_cmd {"key1":true}`},
		"node7":[]string{`test_cmd {"key1":true}`,`test_cmd {"key1":false}`},
		"node8":[]string{`test_cmd {"key1":true}`,`test_cmd {"key1":true}`},
		"node9":[]string{`test_cmd {"key1":true}`,`test_cmd {"key1":false}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterListEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1[0] += 1;} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":[]interface{}{-1, 8}},
		map[string]interface{}{"key1":[]interface{}{10, -2}},
		map[string]interface{}{"key1":[]interface{}{0, 100}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":[0,8]}`,`test_cmd {"key1":[11,-2]}`,`test_cmd {"key1":[1,100]}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":[-1,8]}`,`test_cmd {"key1":[10,-2]}`,`test_cmd {"key1":[0,100]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterListListEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1[1][0] += 1;} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":[]interface{}{[]interface{}{-1, -8}, []interface{}{8, 6}}},
		map[string]interface{}{"key1":[]interface{}{[]interface{}{0, 98}, []interface{}{-8, -6}}},
		map[string]interface{}{"key1":[]interface{}{[]interface{}{9, 18}}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":[[-1,-8],[9,6]]}`,`test_cmd {"key1":[[0,98],[-7,-6]]}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":[[-1,-8],[8,6]]}`,`test_cmd {"key1":[[0,98],[-8,-6]]}`,`test_cmd {"key1":[[9,18]]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterMapMapEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1.foo1.foo2 += 1;} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":map[string]interface{}{"foo1":map[string]interface{}{"foo2":9, "bar2":-8}, "bar1":map[string]interface{}{"foo2":9, "bar2":-8}}},
		map[string]interface{}{"key1":map[string]interface{}{"foo1":map[string]interface{}{"bar2":-8}, "bar1":map[string]interface{}{"foo2":9, "bar2":-8}}},
		map[string]interface{}{"key1":map[string]interface{}{"bar1":map[string]interface{}{"foo2":9, "bar2":-8}}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":{"bar1":{"bar2":-8,"foo2":9},"foo1":{"bar2":-8,"foo2":10}}}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,
			`test_cmd {"key1":{"bar1":{"bar2":-8,"foo2":9},"foo1":{"bar2":-8,"foo2":9}}}`,
			`test_cmd {"key1":{"bar1":{"bar2":-8,"foo2":9},"foo1":{"bar2":-8}}}`,
			`test_cmd {"key1":{"bar1":{"bar2":-8,"foo2":9}}}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterMapListEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1.foo[0] += 1;} > node1`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"},
		map[string]interface{}{"key1":map[string]interface{}{"foo":[]interface{}{-1, -8}, "bar":[]interface{}{8, 6}}},
		map[string]interface{}{"key1":map[string]interface{}{"foo":[]interface{}{100, -8}, "bar":[]interface{}{-8, -6}}},
		map[string]interface{}{"key1":map[string]interface{}{"bar":[]interface{}{8, 6}}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":{"bar":[8,6],"foo":[0,-8]}}`,`test_cmd {"key1":{"bar":[-8,-6],"foo":[101,-8]}}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`,`test_cmd {"key1":{"bar":[8,6],"foo":[-1,-8]}}`,`test_cmd {"key1":{"bar":[-8,-6],"foo":[100,-8]}}`,`test_cmd {"key1":{"bar":[8,6]}}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterCollectionAssign(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 = {a:1,b:2};} > node2`,
		`node0 > % {key1 = ["a",2];} > node3`,
		`node0 > % {key1 = [[1,2],[3,4]];} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`},
		"node2":[]string{`test_cmd {"key1":{"a":1,"b":2}}`},
		"node3":[]string{`test_cmd {"key1":["a",2]}`},
		"node4":[]string{`test_cmd {"key1":[[1,2],[3,4]]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterCollectionEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 = key1[1]+key2.a;} > node2`,
		`node0 > % {key1[0] += key2.b;} > node3`,
		`node0 > % {key1[0] = [1,2];key2.a = {x:1};} > node4`,
		`node0 > % {key2=key1[0]+key2.b;} > node5`,
		`node0 > % {key1 += [3,4];} > node6`,
		`node0 > % {key2 += {c:9};} > node7`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":[]interface{}{"val1",2},"key2":map[string]interface{}{"a":1,"b":"val2"}}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":["val1",2],"key2":{"a":1,"b":"val2"}}`},
		"node2":[]string{`test_cmd {"key1":3,"key2":{"a":1,"b":"val2"}}`},
		"node3":[]string{`test_cmd {"key1":["val1val2",2],"key2":{"a":1,"b":"val2"}}`},
		"node4":[]string{`test_cmd {"key1":[[1,2],2],"key2":{"a":{"x":1},"b":"val2"}}`},
		"node5":[]string{`test_cmd {"key1":["val1",2],"key2":"val1val2"}`},
		"node6":[]string{`test_cmd {"key1":["val1",2,3,4],"key2":{"a":1,"b":"val2"}}`},
		"node7":[]string{`test_cmd {"key1":["val1",2],"key2":{"a":1,"b":"val2","c":9}}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterSimpleReplace(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > = {key1:"val1"} > node2`,
		`node0 > = {key1:1+1} > node3`,
		`node0 > = {key2:[1,2]} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val2"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val2"}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`},
		"node3":[]string{`test_cmd {"key1":2}`},
		"node4":[]string{`test_cmd {"key2":[1,2]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterSimpleCondReplace(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > ? {key1 >= "val"} = {key1:"val1"} > node2`,
		`node0 > ? {key1 == "val2"} = {key1:1+1} > node3`,
		`node0 > ? {key1 == "val3"} = {key2:[1,2]} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val2"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val2"}`},
		"node2":[]string{`test_cmd {"key1":"val1"}`},
		"node3":[]string{`test_cmd {"key1":2}`},
		"node4":[]string{`test_cmd {"key1":"val2"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterSimpleCondEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > ? {key1 >= "val"} % {key1="vallll";} > node2`,
		`node0 > ? {key1 == "val2"} % {key1=key1+key1;} > node3`,
		`node0 > ? {key1 <= "xal2"} % {key1=6.5*2;} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val2"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val2"}`},
		"node2":[]string{`test_cmd {"key1":"vallll"}`},
		"node3":[]string{`test_cmd {"key1":"val2val2"}`},
		"node4":[]string{`test_cmd {"key1":13}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}


func TestRouterStringEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 += "val2";} > node2`,
		`node0 > % {key1 = key1*2+key1;} > node3`,
		`node0 > % {key1 = key1[2];} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`},
		"node2":[]string{`test_cmd {"key1":"val1val2"}`},
		"node3":[]string{`test_cmd {"key1":"val1val1val1"}`},
		"node4":[]string{`test_cmd {"key1":"l"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterUnicodeEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 += "ç";} > node2`,
		`node0 > % {key1 = key1*2+key1;} > node3`,
		`node0 > % {key1 = key1[1];} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"héllö"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"héllö"}`},
		"node2":[]string{`test_cmd {"key1":"héllöç"}`},
		"node3":[]string{`test_cmd {"key1":"héllöhéllöhéllö"}`},
		"node4":[]string{`test_cmd {"key1":"é"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterDeletionEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {del key1;} > node2`,
		`node0 > % {del key2;} > node3`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1","key2":2}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1","key2":2}`},
		"node2":[]string{`test_cmd {"key2":2}`},
		"node3":[]string{`test_cmd {"key1":"val1"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterLocalVarsEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {var x;x=3;key2 *= x;} > node2`,
		`node0 > % {var x;x="val2";key1 += x;} > node3`,
		`node0 > % {var x;x="val2";x="val3";key1 += x;} > node4`,
		`node0 > % {var x;x=[1,2];x[1]=3;key2 += x[1];} > node5`,
		`node0 > % {var x;x={a:1};x.a=3;key2 += x.a;} > node6`,
		`node0 > % {var x;x=3;del x;key2 += x;} > node7`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1","key2":2}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1","key2":2}`},
		"node2":[]string{`test_cmd {"key1":"val1","key2":6}`},
		"node3":[]string{`test_cmd {"key1":"val1val2","key2":2}`},
		"node4":[]string{`test_cmd {"key1":"val1val3","key2":2}`},
		"node5":[]string{`test_cmd {"key1":"val1","key2":5}`},
		"node6":[]string{`test_cmd {"key1":"val1","key2":5}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterTernaryEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 = key1 > "a" ? "val2" : "val3";} > node2`,
		`node0 > % {key1 = key1 < "a" ? "val2" : "val3";} > node3`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1","key2":2}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1","key2":2}`},
		"node2":[]string{`test_cmd {"key1":"val2","key2":2}`},
		"node3":[]string{`test_cmd {"key1":"val3","key2":2}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterStringMatch(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key2 = key1 ~ "val[0-9]" ? "is val" : "non-val";} > node2`,
		`node0 > % {key2 = key1 ~ ".*v.*" ? "is val" : "non-val";} > node3`,
		`node0 > % {key2 = key1 ~ ".*b.*" ? "is val" : "non-val";} > node4`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`},
		"node2":[]string{`test_cmd {"key1":"val1","key2":"is val"}`},
		"node3":[]string{`test_cmd {"key1":"val1","key2":"is val"}`},
		"node4":[]string{`test_cmd {"key1":"val1","key2":"non-val"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterCommandEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > = test_cmd2 > node2`,
		`node0 > ? test_cmd = test_cmd3 > node3`,
		`node0 > ? test_cmd2 = test_cmd4 > node4`,
		`node0 > ? test_cmd > node5`,
		`node0 > ? test_cmd2 > node6`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`},
		"node2":[]string{`test_cmd2 {"key1":"val1"}`},
		"node3":[]string{`test_cmd3 {"key1":"val1"}`},
		"node4":[]string{`test_cmd {"key1":"val1"}`},
		"node5":[]string{`test_cmd {"key1":"val1"}`},
		"node6":[]string{}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterCommandEditComplexTransform(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > ? test_cmd % test_cmd2 {key1 += "2";} > node2`,
		`node0 > ? test_cmd = test_cmd3 {key2:4} > node3`,
		`node0 > ? test_cmd2 % {key1 += "2";} > node4`,
		`node0 > ? test_cmd2 = {key2:4} > node5`,
		`node0 > ? test_cmd {key1 ~ ".*[02468]"} = {key2:4} > node6`,
		`node0 > ? test_cmd2 {key1 == "val1"} = {key2:4} > node7`}
	messages := []map[string]interface{}{
		map[string]interface{}{"key1":"val1"}}
	expected := map[string][]string{
		"node1":[]string{`test_cmd {"key1":"val1"}`},
		"node2":[]string{`test_cmd2 {"key1":"val12"}`},
		"node3":[]string{`test_cmd3 {"key2":4}`},
		"node4":[]string{`test_cmd {"key1":"val1"}`},
		"node5":[]string{`test_cmd {"key1":"val1"}`},
		"node6":[]string{`test_cmd {"key1":"val1"}`},
		"node7":[]string{`test_cmd {"key1":"val1"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}
