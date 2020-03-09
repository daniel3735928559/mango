package route

import (
	"fmt"
	"strings"
	"mc/value"
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

func (rt *Route) Run(command string, args *value.Value) (string, *value.Value, error) {
	for _, t := range rt.Transforms {
		new_command, new_args, err := t.Execute(command, args)
		if err != nil {
			fmt.Println("Error executing transform",t.ToString(),"on",args.ToString(),"ERROR",err)
			return "", nil, err
		}
		if new_args == nil {
			return "", nil, nil
		}
		
		fmt.Println("TRANSFORMED",new_args.ToString())
		args = new_args
		command = new_command
	}
	return command, args, nil
}
