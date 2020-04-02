package route

import (
	"fmt"
	"strings"
	"mc/value"
)

type Route struct {
	Id string
	Group string
	Source string
	Transforms []*Transform
	Dest string
}

func (rt *Route) GetSource() string {
	if strings.Contains(rt.Source, "/") {
		return rt.Source
	}
	return fmt.Sprintf("%s/%s", rt.Group, rt.Source)
}

func (rt *Route) GetDest() string {
	if strings.Contains(rt.Dest, "/") {
		return rt.Dest
	}
	return fmt.Sprintf("%s/%s", rt.Group, rt.Dest)
}

func (rt *Route) ToString() string {
	if rt.Transforms != nil && len(rt.Transforms) > 0 {
		tforms_strings := make([]string, len(rt.Transforms))
		for i, t := range rt.Transforms {
			tforms_strings[i] = t.ToString()
		}
		tforms := strings.Join(tforms_strings, " > ")
		return fmt.Sprintf("%s > %s > %s", rt.GetSource(), tforms, rt.GetDest())
	}
	return fmt.Sprintf("%s > %s", rt.GetSource(), rt.GetDest())
}

func (rt *Route) Run(command string, args *value.Value) (string, *value.Value, error) {
	for _, t := range rt.Transforms {
		fmt.Println("[MC] TRANSFORM",t.ToString())
		new_command, new_args, err := t.Execute(command, args)
		if err != nil {
			fmt.Println("Error executing transform",t.ToString(),"on",args.ToString(),"ERROR",err)
			return "", nil, err
		}
		if new_args == nil {
			return "", nil, nil
		}
		
		fmt.Println("[MC] TRANSFORMED",new_args.ToString())
		args = new_args
		command = new_command
	}
	return command, args, nil
}
