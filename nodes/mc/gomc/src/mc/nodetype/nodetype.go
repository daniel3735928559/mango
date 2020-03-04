package nodetype

import (
	"fmt"
)

type NodeType struct {
	Name string
	Interface *NodeInterface
	Args []string
	Command string
}

func Parse(spec string) *NodeType {
	return nil
}

func (nt *NodeType) Run(args map[string]string) {
	
}

func (nt *NodeType) ToString() string {
	return fmt.Sprintf("%s", nt.Name)
}

