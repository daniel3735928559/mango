package emp

import (
	"fmt"
	"mc/nodetype"
)

type EMP struct {
	Name string
	NodeTypes map[string]*nodetype.NodeType
	Nodes []string
	Routes []string
}

func (emp *EMP) Run() {
	for n := range emp.Nodes {
		fmt.Println(n)
	}
}
