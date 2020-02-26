package router

import (
	"fmt"
	"strings"
)

type Router struct {
	Routes []*Route
	Nodes map[string]*Node
}

func MakeRouter() *Router {
	return &Router {
		Routes: make([]*Route, 0),
		Nodes: make(map[string]*Node)}
}

func (r *Router) AddNode(n *Node) {
	r.Nodes[n.ToString()] = n
}

func (r *Router) FindNode(node_id string) *Node {
	parts := strings.SplitN(node_id, "/" , 2)
	var group, name string
	if len(parts) == 1 {
		group = "root"
		name = parts[0]
	} else if len(parts) == 2 {
		group = parts[0]
		name = parts[1]
	}
	fmt.Println("FIND",group,name)
	fmt.Println(r.Nodes)
	if node, ok := r.Nodes[fmt.Sprintf("%s/%s",group,name)]; ok {
		fmt.Println("FOUND")
		return node
	}
	fmt.Println("NOT FOUND")
	return nil
}

func (r *Router) ParseAndAddRoutes(spec string) {
	rs := Parse(spec)
	for _, rt := range rs {
		r.addRoute(rt)
	}
}

func (r *Router) addRoute(rt *Route) {
	fmt.Println("adding route")
	src := r.FindNode(rt.Source.ToString())
	dst := r.FindNode(rt.Dest.ToString())
	if src != nil && dst != nil {
		rt.Source = src
		rt.Dest = dst
		r.Routes = append(r.Routes, rt)
	} else if src == nil {
		fmt.Printf("Error adding route: No such source node: %s/%s\n",rt.Source.Group,rt.Source.Name)
	} else if dst == nil {
		fmt.Printf("Error adding route: No such destination node: %s/%s\n",rt.Dest.Group,rt.Dest.Name)
	}
}

func (r *Router) Send(src_node *Node, message_id string, command string, args map[string]interface{}) {
	for _, rt := range r.Routes {
		if rt.Source == src_node {
			output_command, output, err := rt.Send(command, args)
			if err != nil {
				fmt.Println("ERROR SENDING ON",rt.ToString(), err)
			} else if output != nil {
				fmt.Println("SENDING",output,"TO",rt.Dest.ToString())
				// Send the result to the destination
				rt.Dest.Handler(src_node.ToString(), message_id, output_command, output, rt.Dest.ToString())
			}
		}
	}
}
