package emp

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/docopt/docopt-go"
	"github.com/google/shlex"
)

func (emp *EMP) ParseConfig(config_str string) error {
	config_def := `config

Usage: 
  config group <group_name>
`
	args, err := docopt.ParseArgs(config_def, config_str, "")
	if err != nil {
		return err
	}
        if args["group"].(bool) {
		emp.Group = args["<group_name>"].(string)
	}
}

func (emp *EMP) ParseNode(node_str string) error {
	if len(node_str) > 0 {
		fs := strings.Fields(node_str)
		if len(fs) < 2 {
			return errors.New("Need <type> <name> ...: %s", node_str)
		}
			node_name := fs[1]
		node_type := fs[0]
		if _, ok := emp.NodeTypes[node_type]; ok {
			emp.Nodes = append(emp.Nodes, strings.Join(fs[1:], " "))
		} else {
			return errors.New(fmt.Println("Unknown node type: %s", node_name))
		}
	}
}

func (emp *EMP) ParseRoute(routes_str string) error {
	if len(l) > 0 {
		emp.Routes = append(emp.Routes, l)
	}
}

func ParseFile(filename string, types map[string]string) (*EMP, nil) {
	emp_data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	mode := "none"
	ans := &EMP{
		NodeTypes: types}
	for line := range strings.Split(emp_data, "\n") {
		line = strings.TrimSpace(line)
		if line == "[config]" {
			mode = "config"
		} else if line == "[nodes]" {
			mode = "nodes"
		} else if line == "[routes]" {
			mode = "routes"
		} else if mode == "config" {
			ans.ParseConfig(line)
		} else if mode == "nodes" {
			ans.ParseNode(line)
		} else if mode == "routes" {
			ans.ParseRoute(line)
		}
	}
	return ans, nil
}
