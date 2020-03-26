package nodetype

import (
	"fmt"
	"strings"
	"mc/value"
	"github.com/google/shlex"
	"github.com/docopt/docopt-go"
)

type NodeType struct {
	Name string
	Interface *NodeInterface
	Command string
	Usage string
	Environment map[string]string
	Validate bool
}

func Parse(spec string) (*NodeType, error) {
	mode := "none"
	interface_spec := ""
	ans := &NodeType{
		Name: "",
		Usage: "",
		Environment: make(map[string]string),
		Validate: true}
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
			config_usage := `Usage:
  config name <node_name>
  config command <cmd>
  config env <name> <val>
  config validate (yes|no)`
			config_args, err := shlex.Split(line)
			if err != nil {
				return nil, fmt.Errorf("Error on line %d: %v", lineno, err)
			}
			docopt.DefaultParser.HelpHandler = docopt.PrintHelpOnly
			args, err := docopt.ParseArgs(config_usage, config_args, "")
			if err != nil {
				return nil, fmt.Errorf("Error on line %d: %v", lineno, err)
			}
			if args["name"].(bool) {
				ans.Name = args["<node_name>"].(string)
			} else if args["command"].(bool) {
				ans.Command = args["<cmd>"].(string)
			} else if args["env"].(bool) {
				ans.Environment[args["<name>"].(string)] = args["<val>"].(string)
			} else if args["validate"].(bool) {
				ans.Validate = args["yes"].(bool)
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

