package router

import (
	"fmt"
	"testing"
	"encoding/json"
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
		`src > = {key1:"val1"} > dst`:`src > replace {NAME(key1):STRING(val1)} > dst`,
		`src > ? {key1 == "val1"} > dst`:`src > pass if {VAR(key1) == STRING(val1)} > dst`,
		`src > % {key1 = "val1";} > dst`:`src > edit {key1 = STRING(val1);} > dst`}
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

func TestParserSingleComplexTransforms(t *testing.T) {
	fmt.Println("ASD")
        examples := map[string]string{
		`src > ? {!(key2 >= 4 || key2 < 8)} = {key1:"val1",key2:-99,key3:1+1+1^6} > dst`: `src > replace {NAME(key1):STRING(val1),NAME(key2):-NUM(99.000000),NAME(key3):NUM(1.000000) + NUM(1.000000) + NUM(1.000000) ^ NUM(6.000000)} if {!VAR(key2) >= NUM(4.000000) || VAR(key2) < NUM(8.000000)} > dst`,
		`src > ? {key1 == "val1" && !(key2 == 4 || key2 >= 8)} > dst`: `src > pass if {VAR(key1) == STRING(val1) && !VAR(key2) == NUM(4.000000) || VAR(key2) >= NUM(8.000000)} > dst`,
		`src > ? {key1 == "val1" && !(key2 <= 4 || key3 < -100.9)} % {key1 = "val1"; key2[4] = key3.key4[1+1];} > dst`: `src > edit {key1 = STRING(val1);key2[NUM(4.000000)] = VAR(key3).NAME(key4)[NUM(1.000000) + NUM(1.000000)];} if {VAR(key1) == STRING(val1) && !VAR(key2) <= NUM(4.000000) || VAR(key3) < -NUM(0.000000)} > dst`}
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

func TestParserSingleComplexCollectionTransforms(t *testing.T) {
	fmt.Println("ASD")
        examples := map[string]string{
		`src > % {key2.key3.key4 += 5;} > dst`: `src > edit {key2.key3.key4 = VAR(key2).NAME(key3).NAME(key4) + NUM(5.000000);} > dst`,
		`src > % {key2[9].key4 += 5;} > dst`: `src > edit {key2[NUM(9.000000)].key4 = VAR(key2)[NUM(9.000000)].NAME(key4) + NUM(5.000000);} > dst`,
		`src > % {key2.key3[1] += 5;} > dst`: `src > edit {key2.key3[NUM(1.000000)] = VAR(key2).NAME(key3)[NUM(1.000000)] + NUM(5.000000);} > dst`}
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

func TestRouterSimple(t *testing.T) {
	r := MakeRouter()

	node0_msgs := []interface{}{}
	node1_msgs := []interface{}{}

	node0_handler := func(header map[string]string, args map[string]interface{}) {
		fmt.Println("node0 GOT",header,args)
		node0_msgs = append(node0_msgs, args)
	}
	node1_handler := func(header map[string]string, args map[string]interface{}) {
		fmt.Println("node1 GOT",header,args)
		node1_msgs = append(node1_msgs, args)
	}
	
	r.AddNode(&Node{
		Group:"group",
		Name:"node0",
		Handle:node0_handler})
	r.AddNode(&Node{
		Group:"group",
		Name:"node1",
		Handle:node1_handler})
	
	r.ParseAndAddRoutes("node0 > node1")
	
	r.Send("node0", map[string]string{}, map[string]interface{}{"key1":"val1"})
}

func RunMessagesThroughRoutes(t *testing.T, routes []string, messages map[string][]map[string]interface{}, expected map[string][]string) {
	// assumes nodes node0, node1, node2
	// args: 
	// - routes: list of routes to add to router
	// - messages: nodename -> message to send from that node
	// returns:
	// - nodename -> list of stringified messages received by that node
	r := MakeRouter()
	num_nodes := 10
	node_msgs := make([][]string, num_nodes)
	node_handlers := make([]func(header map[string]string, args map[string]interface{}), num_nodes)
	for i := 0; i < num_nodes; i++ {
		i := i
		node_msgs[i] = []string{}
		node_handlers[i] = func(header map[string]string, args map[string]interface{}) {
			fmt.Println("node",i,"GOT",header,args)
			data, _ := json.Marshal(args)
			node_msgs[i] = append(node_msgs[i], string(data))
		}
		r.AddNode(&Node{Group:"group",Name:fmt.Sprintf("node%d",i),Handle:node_handlers[i]})
	}

	for _, rt := range routes {
		r.ParseAndAddRoutes(rt)
	}
	for src_node, node_messages := range messages {
		for _, msg := range node_messages {
			fmt.Println("SEND",src_node,msg)
			r.Send(src_node, map[string]string{}, msg)
		}
	}
	results := make(map[string][]string)
	for i := 0; i < num_nodes; i++ {
		results[fmt.Sprintf("node%d",i)] = node_msgs[i]
	}
	for dst_node, recvd_messages := range results {
		if _, ok := expected[dst_node]; !ok {
			if len(recvd_messages) > 0 {
				t.Errorf("Did not have expected messages for node %s", dst_node)
			}
		} else {
			if len(recvd_messages) != len(expected[dst_node]) {
				t.Errorf("Expected %s to receive %d messages; got: %d", dst_node, len(expected[dst_node]), len(recvd_messages))
			} else {
				for i, msg := range recvd_messages {
					if expected[dst_node][i] != msg {
						t.Errorf("Expected %s received message %d to be %s; got: %s", dst_node, i, expected[dst_node][i], msg)
					}
				}
			}
		}
	}
}

func TestRouterEqFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1=="val1"} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":"val2"},
			map[string]interface{}{"key2":"val3"}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":"val1"}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":"val2"}`,`{"key2":"val3"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterNeFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 != "val1"} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":"val2"},
			map[string]interface{}{"key2":"val3"}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":"val2"}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":"val2"}`,`{"key2":"val3"}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterGtFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 > 0} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":10}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterLtFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 < 0} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":-1}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterGeFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 >= 0} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":10}`,`{"key1":0}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterLeFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 <= 0} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":-1}`,`{"key1":0}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterAndFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 >= 0 && key1 < 5} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":2},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node0":[]string{},
		"node1":[]string{`{"key1":2}`,`{"key1":0}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":2}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterOrFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1 < 0 || key1 >= 10} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":2},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":-1}`,`{"key1":10}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":2}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterNotFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {!(key1 <= 0)} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-1},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":10}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterListFilter(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > ? {key1[0] <= 0} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":[]interface{}{-1, 8}},
			map[string]interface{}{"key1":[]interface{}{10, -2}},
			map[string]interface{}{"key1":[]interface{}{0, 100}}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":[-1,8]}`,`{"key1":[0,100]}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":[-1,8]}`,`{"key1":[10,-2]}`,`{"key1":[0,100]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterArithmeticEdit(t *testing.T) {
	routes := []string{
		"node0 > node1",
		`node0 > % {key1 += 1;} > node2`,
		`node0 > % {key1 -= 1;} > node3`,
		`node0 > % {key1 *= 2;} > node4`,
		`node0 > % {key1 /= 2;} > node5`,
		`node0 > % {key1 %= 2;} > node6`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":-2},
			map[string]interface{}{"key1":10},
			map[string]interface{}{"key1":0}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":"val1"}`,`{"key1":-2}`,`{"key1":10}`,`{"key1":0}`},
		"node2":[]string{`{"key1":-1}`,`{"key1":11}`,`{"key1":1}`},
		"node3":[]string{`{"key1":-3}`,`{"key1":9}`,`{"key1":-1}`},
		"node4":[]string{`{"key1":"val1val1"}`,`{"key1":-4}`,`{"key1":20}`,`{"key1":0}`},
		"node5":[]string{`{"key1":-1}`,`{"key1":5}`,`{"key1":0}`},
		"node6":[]string{`{"key1":0}`,`{"key1":0}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}


func TestRouterListEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1[0] += 1;} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":[]interface{}{-1, 8}},
			map[string]interface{}{"key1":[]interface{}{10, -2}},
			map[string]interface{}{"key1":[]interface{}{0, 100}}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":[0,8]}`,`{"key1":[11,-2]}`,`{"key1":[1,100]}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":[-1,8]}`,`{"key1":[10,-2]}`,`{"key1":[0,100]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterListListEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1[1][0] += 1;} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":[]interface{}{[]interface{}{-1, -8}, []interface{}{8, 6}}},
			map[string]interface{}{"key1":[]interface{}{[]interface{}{0, 98}, []interface{}{-8, -6}}},
			map[string]interface{}{"key1":[]interface{}{[]interface{}{9, 18}}}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":[[-1,-8],[9,6]]}`,`{"key1":[[0,98],[-7,-6]]}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":[[-1,-8],[8,6]]}`,`{"key1":[[0,98],[-8,-6]]}`,`{"key1":[[9,18]]}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}

func TestRouterMapListEdit(t *testing.T) {
	routes := []string{
		"node0 > node2",
		`node0 > % {key1.foo[0] += 1;} > node1`}
	messages := map[string][]map[string]interface{}{
		"node0":[]map[string]interface{}{
			map[string]interface{}{"key1":"val1"},
			map[string]interface{}{"key1":map[string]interface{}{"foo":[]interface{}{-1, -8}, "bar":[]interface{}{8, 6}}},
			map[string]interface{}{"key1":map[string]interface{}{"foo":[]interface{}{100, -8}, "bar":[]interface{}{-8, -6}}},
			map[string]interface{}{"key1":map[string]interface{}{"bar":[]interface{}{8, 6}}}}}
	expected := map[string][]string{
		"node1":[]string{`{"key1":{"bar":[8,6],"foo":[0,-8]}}`,`{"key1":{"bar":[-8,-6],"foo":[101,-8]}}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":{"bar":[8,6],"foo":[-1,-8]}}`,`{"key1":{"bar":[-8,-6],"foo":[100,-8]}}`,`{"key1":{"bar":[8,6]}}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}
