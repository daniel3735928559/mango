package router

import (
	"fmt"
	"strings"
)

type Route struct {
	Source string
	Transforms []*Transform
	Dest string
}

func (r *Route) ToString() string {
	if r.Transforms != nil && len(r.Transforms) > 0 {
		tforms_strings := make([]string, len(r.Transforms))
		for i, t := range r.Transforms {
			tforms_strings[i] = t.ToString()
		}
		tforms := strings.Join(tforms_strings, " > ")
		return fmt.Sprintf("%s > %s > %s", r.Source, tforms, r.Dest)
	}
	return fmt.Sprintf("%s > %s", r.Source, r.Dest)
}
