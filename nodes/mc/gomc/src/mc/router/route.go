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
	this := MakeValue(args)
	if this.Type != VAL_MAP {
		return nil, errors.New("args must be a map")
	}
	
	// Apply transforms
	var err error
	for _, t := range rt.Transforms {
		this, err = t.Execute(this)
		if err != nil {
			fmt.Println("Error executing transform",t.ToString(),"on",this,"ERROR",err)
			return nil, err
		}
	}
	
	output_prim := this.ToPrimitive()
	if output_args, ok := output_prim.(map[string]interface{}); ok {
		return output_args, nil
	}
	return nil, errors.New(fmt.Sprintf("Failed sending on %s",rt.ToString()))
}
