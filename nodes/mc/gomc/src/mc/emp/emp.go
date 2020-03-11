package emp

import (
	"fmt"
)

type EMPNode struct {
	Name string
	TypeName string
	Args string
}

type EMP struct {
	Name string
	Nodes []EMPNode
	Routes []string
}



func (emp *EMP) Run() {
	for n := range emp.Nodes {
		fmt.Println(n)
	}
}
