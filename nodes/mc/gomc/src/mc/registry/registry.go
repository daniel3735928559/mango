package registry

import (
	"fmt"
	"strings"
	"mc/node"
	"mc/nodetype"
	"mc/route"
)

type Registry struct {
	Nodes map[string]*node.Node
	NodeTypes map[string]*nodetype.NodeType
	Routes map[string]*route.Route
}

func MakeRegistry() *Registry {
	return &Registry{
		Nodes: make(map[string]*node.Node),
		NodeTypes: make(map[string]*nodetype.NodeType),
		Routes: make(map[string]*route.Route)}
}

func (reg *Registry) AddNode(n *node.Node) error {
	if _, ok := reg.Nodes[n.Id]; ok {
		return fmt.Errorf("Node with id %s already exists", n.Id)
	}
	reg.Nodes[n.Id] = n
	return nil
}

func (reg *Registry) AddNodeType(n *nodetype.NodeType) error {
	if _, ok := reg.NodeTypes[n.Name]; ok {
		return fmt.Errorf("NodeType with name %s already exists", n.Name)
	}
	reg.NodeTypes[n.Name] = n
	return nil
}

func (reg *Registry) AddRoute(r *route.Route) error {
	if _, ok := reg.Routes[r.Id]; ok {
		return fmt.Errorf("Route with id %s already exists", r.Id)
	}
	reg.Routes[r.Id] = r
	return nil
}

func (reg *Registry) FindRoutesBySrc(src string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.Source == src {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) FindRoutesByDst(dst string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.Dest == dst {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) FindRoutesBySrcDst(src string, dst string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.Source == src && rt.Dest == dst {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) FindNodeById(node_id string) *node.Node {
	if n, ok := reg.Nodes[node_id]; ok {
		return n
	}
	return nil
}

func (reg *Registry) FindNodeByName(node_name string) *node.Node {
	fs := strings.Split(node_name, "/")
	if len(fs) != 2 {
		return nil
	}
	for _, n := range reg.Nodes {
		if n.Group == fs[0] && n.Name == fs[1] {
			return n
		}
	}
	return nil
}

func (reg *Registry) FindNodeType(nodetype string) *nodetype.NodeType {
	if nt, ok := reg.NodeTypes[nodetype]; ok {
		return nt
	}
	return nil
}
