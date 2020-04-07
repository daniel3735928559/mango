package registry

import (
	"fmt"
	"strings"
	"mc/node"
	"mc/nodetype"
	"mc/route"
)

type Registry struct {
	Groups map[string]bool
	Nodes map[string]node.Node
	NodeTypes map[string]*nodetype.NodeType
	Routes map[string]*route.Route
}

func MakeRegistry() *Registry {
	return &Registry{
		Nodes: make(map[string]node.Node),
		NodeTypes: make(map[string]*nodetype.NodeType),
		Routes: make(map[string]*route.Route),
		Groups: make(map[string]bool)}
}

func (reg *Registry) AddNode(n node.Node) error {
	if _, ok := reg.Nodes[n.GetId()]; ok {
		return fmt.Errorf("Node with id %s already exists", n.GetId())
	}
	reg.Nodes[n.GetId()] = n
	reg.Groups[n.GetGroup()] = true
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

func (reg *Registry) GetRoutes() []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		ans = append(ans, rt)
	}
	return ans
}

func (reg *Registry) GetGroups() []string {
	ans := make([]string, 0)
	for g, _ := range reg.Groups {
		ans = append(ans, g)
	}
	return ans
}

func (reg *Registry) FindRoutesBySrc(src string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.GetSource() == src {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) FindRoutesByDst(dst string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.GetDest() == dst {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) FindRoutesBySrcDst(src string, dst string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.GetSource() == src && rt.GetDest() == dst {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) FindRoutesByGroup(group_name string) []*route.Route {
	ans := make([]*route.Route, 0)
	for _, rt := range reg.Routes {
		if rt.Group == group_name {
			ans = append(ans, rt)
		}
	}
	return ans
}

func (reg *Registry) GetNodes() []node.Node {
	ans := make([]node.Node, 0)
	for _, n := range reg.Nodes {
		ans = append(ans, n)
	}
	return ans
}

func (reg *Registry) FindNodeById(node_id string) node.Node {
	if n, ok := reg.Nodes[node_id]; ok {
		return n
	}
	return nil
}

func (reg *Registry) FindNodeByName(node_name string) node.Node {
	fs := strings.Split(node_name, "/")
	if len(fs) != 2 {
		return nil
	}
	for _, n := range reg.Nodes {
		if n.GetGroup() == fs[0] && n.GetName() == fs[1] {
			return n
		}
	}
	return nil
}

func (reg *Registry) FindNodesByGroup(group_name string) []node.Node {
	ans := make([]node.Node, 0)
	for _, n := range reg.Nodes {
		if n.GetGroup() == group_name {
			ans = append(ans, n)
		}
	}
	return ans
}

func (reg *Registry) FindNodeType(nodetype string) *nodetype.NodeType {
	if nt, ok := reg.NodeTypes[nodetype]; ok {
		return nt
	}
	return nil
}

func (reg *Registry) DelGroup(name string) {
	nodes := reg.FindNodesByGroup(name)
	for _, n := range nodes {
		reg.DelNode(n.GetId())
	}
	rts := reg.FindRoutesByGroup(name)
	for _, rt := range rts {
		reg.DelRoute(rt.Id)
	}
	delete(reg.Groups, name)
}

func (reg *Registry) DelNode(node_id string) bool {
	if n, ok := reg.Nodes[node_id]; ok {
		node_routes := make([]*route.Route, 0)
		src_routes := reg.FindRoutesBySrc(n.ToString())
		dst_routes := reg.FindRoutesByDst(n.ToString())
		node_routes = append(node_routes, src_routes...)
		node_routes = append(node_routes, dst_routes...)
		for _, rt := range node_routes {
			reg.DelRoute(rt.Id)
		}
		delete(reg.Nodes, node_id)
		return true
	}
	return false
}

func (reg *Registry) DelRoute(rt_id string) bool {
	if _, ok := reg.Routes[rt_id]; ok {
		delete(reg.Routes, rt_id)
		return true
	}
	return false
}

