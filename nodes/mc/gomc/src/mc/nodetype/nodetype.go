package nodetype

import (
	"fmt"
	"strings"
	"mc/value"
)

type NodeType struct {
	Name string
	Interface *NodeInterface
	Executable string
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
		if len(line) == 0 {
			continue
		}
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
			} else if fs[0] == "executable" {
				if len(fs) != 2 {
					return nil, fmt.Errorf("Error on line %d: config name line should be of the form `executable <path_to_executable>`", lineno)
				}
				ans.Executable = fs[1]
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
	if len(ans.Name) == 0 {
		return nil, fmt.Errorf("Must supply a type name")
	}
	return ans, nil
}

func (nt *NodeType) ValidateInput(name string, val *value.Value) (*value.Value, error) {
	return nt.Interface.ValidateInput(name, val)
}

func (nt *NodeType) ValidateOutput(name string, val *value.Value) (*value.Value, error) {
	return nt.Interface.ValidateOutput(name, val)
}

func (nt *NodeType) ToString() string {
	return fmt.Sprintf("%s", nt.Name)
}

