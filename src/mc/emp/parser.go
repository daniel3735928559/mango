package emp

import (
	"fmt"
	"strings"
	"github.com/google/shlex"
	"github.com/docopt/docopt-go"
)

func (emp *EMP) ParseConfig(config_str string) error {
	config_def := `Usage: 
  config name <name>
  config param <name>
`
	config_args, err := shlex.Split(config_str)
	if err != nil {
		return err
	}
	docopt.DefaultParser.HelpHandler = docopt.PrintHelpOnly
	args, err:= docopt.ParseArgs(config_def, config_args, "")
	if err != nil {
		return err
	}
	if args["name"].(bool) {
		name, _ := args["<name>"].(string)
		emp.Name = name
	} else if args["param"].(bool) {
		name, _ := args["<name>"].(string)
		emp.ParamNames = append(emp.ParamNames, name)
	}
	return nil
}

func (emp *EMP) ParseNode(node_str string) error {
	node_def := `Usage: 
  node instance <type> <name> [<args>...]
  node dummy <name>
  node merge <name> <mergepoints>...
  node gen <name> <values>...
`
	node_args, err := shlex.Split(node_str)
	if err != nil {
		return err
	}
	docopt.DefaultParser.HelpHandler = docopt.PrintHelpOnly
	args, err:= docopt.ParseArgs(node_def, node_args, "")
	if err != nil {
		return err
	}
	fmt.Println("EA",args)
	
	if args["instance"].(bool) {
		node_type := args["<type>"].(string)
		node_name := args["<name>"].(string)
		emp.Nodes = append(emp.Nodes, EMPNode{TypeName: node_type, Name: node_name, Args: args["<args>"].([]string)})
	} else if args["dummy"].(bool) {
		node_name := args["<name>"].(string)
		emp.Nodes = append(emp.Nodes, EMPNode{TypeName: "dummy", Name: node_name})
	} else if args["merge"].(bool) {
		node_name := args["<name>"].(string)
		emp.Nodes = append(emp.Nodes, EMPNode{TypeName: "merge", Name: node_name, Args: args["<mergepoints>"].([]string)})
	} else if args["gen"].(bool) {
		return fmt.Errorf("Not yet implemented: gen")
	}
	return nil
}

func (emp *EMP) ParseRoute(routes_str string) error {
	if len(routes_str) > 0 {
		emp.Routes = append(emp.Routes, routes_str)
	}
	return nil
}

func Parse(emp_data string) (*EMP, error) {
	mode := "none"
	ans := &EMP{
		ParamNames: make([]string, 0)}
	for _, line := range strings.Split(emp_data, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
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
