package nodetype

import (
	"fmt"
	"strings"
)

type NodeType struct {
	Name string
	Interface *NodeInterface
	Usage string
}

func Parse(spec string) (*NodeType, error) {
	mode := "none"
	interface_spec := ""
	ans := &NodeType{
		Name: "",
		Usage: ""}
	for lineno, line := range strings.Split(spec, "\n") {
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
				if len(fs) != 2 {
					return nil, fmt.Errorf("Error on line %d: config name line should be of the form `name <node_name>`", lineno)
				}
				ans.Name = fs[1]
			} else {
				return nil, fmt.Errorf(`Error on line %d: config lines supported: 
name <node_name>`, lineno)
			}
		} else if mode == "interface" {
			interface_spec += line + "\n"
		} else if mode == "usage" {
			ans.Usage += line + "\n"
		}
		new_if, err := ParseNodeInterface(interface_spec)
		if err != nil {
			return nil, err
		} else {
			ans.Interface = new_if
		}
	}
	return ans, nil
}

func (nt *NodeType) Run(args map[string]string) {
	
}

func (nt *NodeType) ToString() string {
	return fmt.Sprintf("%s", nt.Name)
}

