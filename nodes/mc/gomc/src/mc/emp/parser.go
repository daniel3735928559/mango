package emp

import (
	"fmt"
	"strings"
	"mc/nodetype"
	"github.com/google/shlex"
	"github.com/docopt/docopt-go"
)

func (emp *EMP) ParseConfig(config_str string) error {
	config_def := `config

Usage: 
  config name <name>
`
	config_args, err := shlex.Split(config_str)
	if err != nil {
		return err
	}
	args, err:= docopt.ParseArgs(config_def, config_args, "")
	if err != nil {
		return err
	}
	if name, ok := args["<name>"].(string); ok {
		emp.Name = name
	}
	return nil
}

func (emp *EMP) ParseNode(node_str string) error {
	if len(node_str) > 0 {
		fs := strings.SplitN(node_str," ",3)
		if len(fs) < 2 {
			return fmt.Errorf("Need <type> <name> <args>...: %s", node_str)
		}
		node_type := fs[0]
		node_name := fs[1]
		if _, ok := emp.NodeTypes[node_type]; ok {
			emp.Nodes = append(emp.Nodes, fs[3])
		} else {
			return fmt.Errorf("Unknown node type: %s", node_name)
		}
	}
	return nil
}

func (emp *EMP) ParseRoute(routes_str string) error {
	if len(routes_str) > 0 {
		emp.Routes = append(emp.Routes, routes_str)
	}
	return nil
}

func Parse(emp_data string, types map[string]*nodetype.NodeType) (*EMP, error) {
	mode := "none"
	ans := &EMP{
		NodeTypes: types}
	for _, line := range strings.Split(emp_data, "\n") {
		line = strings.TrimSpace(line)
		if line == "[config]" {
			mode = "config"
		} else if line == "[nodes]" {
			mode = "nodes"
		} else if line == "[routes]" {
			mode = "routes"
		} else if mode == "config" {
			err := ans.ParseConfig(line)
			if err != nil {
				return nil, err
			}
		} else if mode == "nodes" {
			err := ans.ParseNode(line)
			if err != nil {
				return nil, err
			}
		} else if mode == "routes" {
			err := ans.ParseRoute(line)
			if err != nil {
				return nil, err
			}
		}
	}
	return ans, nil
}
