package router

import (
	"fmt"
)

type Node struct {
	Name string
	Group string
	Handle func(map[string]string, map[string]interface{})
}

func (n *Node) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}

