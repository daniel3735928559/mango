package router

import (
	"fmt"
	"strings"
)

type Expression struct {
	Operation string
	Args []*Expression
	Value *Value
}			

func (e *Expression) ToString() string {
	subexprs := ""
	val := ""
	if e.Args != nil {
		subexprs_strings := make([]string, len(e.Args))
		for i, es := range e.Args {
			subexprs_strings[i] = es.ToString()
		}
		subexprs = strings.Join(subexprs_strings, ",")
	}
	if e.Value != nil {
		val = e.Value.ToString()
	}
	return fmt.Sprintf("%s(%s;%s)",e.Operation, val, subexprs)
}
