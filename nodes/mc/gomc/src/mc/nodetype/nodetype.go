package nodetype

import (
	"fmt"
)

type NodeType struct {
	Name string
	Interface *NodeInterface
	Usage string
}

func Parse(spec string) *NodeType {
	mode := "none"
	interface_spec := ""
	ans := &NodeType{
		Name: "",
		Usage: ""}
	for line := range strings.Split(spec, "\n") {
		line = strings.TrimSpace(line)
		if line == "[config]" {
			mode = "config"
		} else if line == "[interface]" {
			mode = "interface"
		} else if line == "[usage]" {
			mode = "usage"
		} else if mode == "config" {
			fs := strings.Fields(line)
			if fs[0] == "name" {
				ans.Name = fs[1]
			}
		} else if mode == "interface" {
			interface_spec += line + "\n"
		} else if mode == "usage" {
			ans.Usage += line + "\n"
		}
		ans.Interface = ParseNodeInterface(interface_spec)
	}
	return ans
}

func (nt *NodeType) Run(args map[string]string) {
	
}

func (nt *NodeType) ToString() string {
	return fmt.Sprintf("%s", nt.Name)
}

