package router

import (
	"fmt"
	"strings"
	"errors"
)

type Route struct {
	Source string
	Transforms []*Transform
	Dest string
}

func (rt *Route) ToString() string {
	if rt.Transforms != nil && len(rt.Transforms) > 0 {
		tforms_strings := make([]string, len(rt.Transforms))
		for i, t := range rt.Transforms {
			tforms_strings[i] = t.ToString()
		}
		tforms := strings.Join(tforms_strings, " > ")
		return fmt.Sprintf("%s > %s > %s", rt.Source, tforms, rt.Dest)
	}
	return fmt.Sprintf("%s > %s", rt.Source, rt.Dest)
}

func (rt *Route) Send(args map[string]interface{}) (map[string]interface{}, error) {
	// Convert args to a *Value to be the "this" object
	this, err := MakeValue(args)
	if err != nil {
		return nil, err
	}
	if this.Type != VAL_MAP {
		return nil, errors.New("args must be a map")
	}
	fmt.Println("MV",this)
	
	// Apply transforms
	
	for _, t := range rt.Transforms {
		new_this, err := t.Execute(this)
		if err != nil {
			fmt.Println("Error executing transform",t.ToString(),"on",this.ToString(),"ERROR",err)
			return nil, err
		}
		if new_this == nil {
			return nil, nil
		}
		
		fmt.Println("TRANSFORMED",new_this.ToString())
		this = new_this
	}
	
	output_prim := this.ToPrimitive()
	if output_args, ok := output_prim.(map[string]interface{}); ok {
		return output_args, nil
	}
	return nil, errors.New(fmt.Sprintf("Failed sending on %s",rt.ToString()))
}
