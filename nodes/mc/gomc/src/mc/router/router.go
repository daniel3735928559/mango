package router

import (
	"fmt"
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
	r.Nodes[n.Name] = n
}

func (r *Router) ParseAndAddRoutes(spec string) {
	rs := Parse(spec)
	r.Routes = append(r.Routes, rs...)
}

func (r *Router) AddRoute(rt *Route) {
	fmt.Println("adding route")
	r.Routes = append(r.Routes, rt)
}

func (r *Router) Send(src_node string, header map[string]string, args map[string]interface{}) {
	for _, rt := range r.Routes {
		if rt.Source == src_node {
			output, err := rt.Send(args)
			if err != nil {
				fmt.Println("ERROR SENDING ON",rt.ToString(), err)
			} else {
				fmt.Println("SENDING",output,"TO",r.Nodes[rt.Dest].Name)
				// Send the result to the destination
				r.Nodes[rt.Dest].Handle(header, output)
			}
		}
	}
}
