package nodetype

import (
	"fmt"
	"strings"
	"errors"
	valuetype "mc/valuetype"
	value "mc/value"
)

type NodeInterface struct {
	Imports []string
	Types map[string]*valuetype.ValueType
	Inputs map[string]*valuetype.ValueType
	Outputs map[string]*valuetype.ValueType
	ReturnTypes map[string]string // Input name -> Output name
}

func ParseNodeInterface(spec string) (*NodeInterface, error) {
	ans := &NodeInterface{
		Imports: make([]string, 0),
		Types: make(map[string]*valuetype.ValueType),
		Inputs: make(map[string]*valuetype.ValueType),
		Outputs: make(map[string]*valuetype.ValueType),
		ReturnTypes: make(map[string]string)}
	lines := strings.Split(spec, "\n")
	for lineno, line := range lines {
		fs := strings.Fields(line)
		if fs[0] == "import" {
			fn := strings.SplitN(line, " ", 2)[1]
			return nil, errors.New(fmt.Sprintf("Error at line %d: `%s`: import not yet implemented: %s", lineno, line, fn))
		} else if fs[0] == "type" || fs[0] == "input" || fs[0] == "output" {
			type_name := fs[1]
			fmt.Println("FS",line,fs)
			type_spec := strings.SplitN(line, " ", 3)[2]
			ty, err := valuetype.Parse(type_spec)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Error at line %d: `%s`: %v",lineno, line, err))
			}
			if fs[0] == "type" {
				if _, ok := ans.Types[type_name]; ok {
					return nil, errors.New(fmt.Sprintf("Error at line %d: `%s`: Type already defined: %s", lineno, line, type_name))
				}
				ans.Types[type_name] = ty
			} else if fs[0] == "input" {
				if _, ok := ans.Inputs[type_name]; ok {
					return nil, errors.New(fmt.Sprintf("Error at line %d: `%s`: Input already defined: %s", lineno, line, type_name))
				}
				ans.Inputs[type_name] = ty
			} else if fs[0] == "output" {
				if _, ok := ans.Outputs[type_name]; ok {
					return nil, errors.New(fmt.Sprintf("Error at line %d: `%s`: Output already defined: %s", lineno, line, type_name))
				}
				ans.Outputs[type_name] = ty
			}
		} else if fs[0] == "return" {
			if len(fs) != 3 {
				return nil, errors.New(fmt.Sprintf("Error at line %d: `%s`: Expected format `return <input_name> <output_name>`", lineno, line))
			}
			ans.ReturnTypes[fs[1]] = fs[2]
		}
	}
	return ans, nil
}

func (ni *NodeInterface) ToString() string {
	ans := ""
	for _, imp := range ni.Imports {
		ans += fmt.Sprintf("import %s\n",imp)
	}
	for name, ty := range ni.Types {
		ans += fmt.Sprintf("type %s %s\n",name, ty.ToString())
	}
	for name, ty := range ni.Inputs {
		ans += fmt.Sprintf("input %s %s\n",name, ty.ToString())
	}
	for name, ty := range ni.Outputs {
		ans += fmt.Sprintf("output %s %s\n",name, ty.ToString())
	}
	for name, retname := range ni.ReturnTypes {
		ans += fmt.Sprintf("return %s %s\n",name, retname)
	}
	return ans
}

func (ni *NodeInterface) ValidateInput(name string, val *value.Value) (*value.Value, error) {
	if ty, ok := ni.Inputs[name]; ok {
		return ty.Validate(val, ni.Types, "")
	}
	return nil, fmt.Errorf("No such input type found: %s", name)
}

func (ni *NodeInterface) ValidateOutput(name string, val *value.Value) (*value.Value, error) {
	if ty, ok := ni.Outputs[name]; ok {
		return ty.Validate(val, ni.Types, "")
	}
	return nil, fmt.Errorf("No such output type found: %s", name)
}
