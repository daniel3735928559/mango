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

