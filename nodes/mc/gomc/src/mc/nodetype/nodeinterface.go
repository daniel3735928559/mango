package router

import (
	// "fmt"
)

type NodeInterface struct {
	Inputs map[string]*TypeDesc
	Outputs map[string]*TypeDesc
	ReturnTypes map[string]string // Input name -> Output name
}

func ParseNodeInterface(ifdesc string) *NodeInterface {
	return nil
}

func (ni *NodeInterface) ValidateInput(name string, val *Value) *Value {
	if ty, ok := ni.Inputs[name]; ok {
		return ty.Validate(val)
	}
	return nil
}

func (ni *NodeInterface) ValidateOutput(name string, val *Value) *Value {
	if ty, ok := ni.Outputs[name]; ok {
		return ty.Validate(val)
	}
	return nil
}
