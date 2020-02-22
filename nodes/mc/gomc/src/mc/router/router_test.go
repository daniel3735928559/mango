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
		`src > = {key1:"val1"} > dst`:`src > replace {VAR(key1):STRING(val1)} > dst`,
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
		`src > ? {!(key2 >= 4 || key2 < 8)} = {key1:"val1",key2:-99,key3:1+1+1^6} > dst`: `src > replace {VAR(key1):STRING(val1),VAR(key2):-NUM(99.000000),VAR(key3):NUM(1.000000) + NUM(1.000000) + NUM(1.000000) ^ NUM(6.000000)} if {!VAR(key2) >= NUM(4.000000) || VAR(key2) < NUM(8.000000)} > dst`,
		`src > ? {key1 == "val1" && !(key2 == 4 || key2 >= 8)} > dst`: `src > pass if {VAR(key1) == STRING(val1) && !VAR(key2) == NUM(4.000000) || VAR(key2) >= NUM(8.000000)} > dst`,
		`src > ? {key1 == "val1" && !(key2 <= 4 || key3 < -100.9)} % {key1 = "val1"; key2[4] = key3.key4[1+1];} > dst`: `src > edit {key1 = STRING(val1);key2[NUM(4.000000)] = VAR(key3).VAR(key4)[NUM(1.000000) + NUM(1.000000)];} if {VAR(key1) == STRING(val1) && !VAR(key2) <= NUM(4.000000) || VAR(key3) < -NUM(0.000000)} > dst`}
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
	node0_msgs,node1_msgs,node2_msgs := []string{}, []string{}, []string{}
	node0_handler := func(header map[string]string, args map[string]interface{}) {
		fmt.Println("node0 GOT",header,args)
		data, _ := json.Marshal(args)
		node0_msgs = append(node0_msgs, string(data))
	}
	node1_handler := func(header map[string]string, args map[string]interface{}) {
		fmt.Println("node1 GOT",header,args)
		data, _ := json.Marshal(args)
		node1_msgs = append(node1_msgs, string(data))
	}
	node2_handler := func(header map[string]string, args map[string]interface{}) {
		fmt.Println("node2 GOT",header,args)
		data, _ := json.Marshal(args)
		node2_msgs = append(node2_msgs, string(data))
	}
	r := MakeRouter()	
	r.AddNode(&Node{Group:"group",Name:"node0",Handle:node0_handler})
	r.AddNode(&Node{Group:"group",Name:"node1",Handle:node1_handler})
	r.AddNode(&Node{Group:"group",Name:"node2",Handle:node2_handler})

	for _, rt := range routes {
		r.ParseAndAddRoutes(rt)
	}
	for src_node, node_messages := range messages {
		for _, msg := range node_messages {
			fmt.Println("SEND",src_node,msg)
			r.Send(src_node, map[string]string{}, msg)
		}
	}
	results := map[string][]string {
		"node0":node0_msgs,
		"node1":node1_msgs,
		"node2":node2_msgs}
	for dst_node, recvd_messages := range results {
		if _, ok := expected[dst_node]; !ok {
			t.Errorf("Did not have expected messages for node %s", dst_node)
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
		"node0":[]string{},
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
		"node0":[]string{},
		"node1":[]string{`{"key1":10}`},
		"node2":[]string{`{"key1":"val1"}`,`{"key1":-1}`,`{"key1":10}`,`{"key1":0}`}}
	RunMessagesThroughRoutes(t, routes, messages, expected)
}
