package router

import (
	"fmt"
	"errors"
)

type StatementType int
type PathStepType int

const (
	STMT_ASSIGN StatementType = iota + 1
)

const (
	PATH_MAP PathStepType = iota + 1
	PATH_LIST
)

type Statement struct {
	Type StatementType
	Destination WriteableValue
	Args []*Expression
}

type PathEntry struct {
	Type PathStepType
	ListIndex *Expression
	MapKey string
}

type WriteableValue struct {
	Base string
	Path []PathEntry
}

func (w *WriteableValue) Write(this *Value, vars map[string]*Value, arg *Expression) (*Value, map[string]*Value, error) {
	content, err := arg.Evaluate(this, vars)
	if err != nil {
		return this, vars, err
	}
	var dest *Value
	dest_name := w.Base
	is_this := true
	if dest_name == "this" {
		dest = this.Clone()
	} else if v, ok := vars[dest_name]; ok {
		dest = v.Clone()
		is_this = false
	} else {
		return this, vars, errors.New(fmt.Sprintf("No such variable: %s", dest_name))
	}
	
	target_base := dest.Clone()
	target := target_base
	for i, e := range w.Path {
		if e.Type == PATH_MAP {
			if target.Type != VAL_MAP {
				return this, vars, errors.New(fmt.Sprintf("Attempted to access key %s in non-map", e.MapKey))
			}
			if i == len(w.Path) - 1 {
				target.MapVal[e.MapKey] = content
				if is_this {
					return target_base, vars, nil
				} else {
					vars[dest_name] = target_base
					return this, vars, nil
				}
			} else if new_target, ok := target.MapVal[e.MapKey]; ok {
				target = new_target
			} else {
				return this, vars, errors.New(fmt.Sprintf("Attempted to access non-existent key %s", e.MapKey))
			}
		} else if e.Type == PATH_LIST {
			if target.Type != VAL_LIST {
				return this, vars, errors.New("Attempted to access index in non-map")
			}
			idx, err := e.ListIndex.Evaluate(this, vars)
			if err != nil {
				return this, vars, err
			}
			if idx.Type != VAL_NUM {
				return this, vars, errors.New("List subscript must be integer")
			}
			list_index := int(idx.NumVal)
			if list_index >= len(target.ListVal) || list_index < 0 {
				return this, vars, errors.New(fmt.Sprintf("Attempted to access out-of-bounds index %d", list_index))
			}
			if i == len(w.Path) - 1 {
				// We're at the end--write this value
				target.ListVal[list_index] = content
				if is_this {
					return target_base, vars, nil
				} else {
					vars[dest_name] = target_base
					return this, vars, nil
				}
			} else {
				target = target.ListVal[list_index]
			}
		}
	}
	return this, vars, errors.New("Something went wrong?")
}

func MakeAssignment(dest WriteableValue, val *Expression) *Statement {
	return &Statement {
		Type: STMT_ASSIGN,
		Destination: dest,
		Args: []*Expression{val}}
}

func (s *Statement) ToString() string {
	if s.Type == STMT_ASSIGN {
		path_str := s.Destination.Base
		for _, e := range s.Destination.Path {
			if e.Type == PATH_MAP {
				path_str += fmt.Sprintf(".%s", e.MapKey)
			} else if e.Type == PATH_LIST {
				path_str += fmt.Sprintf("[%s]", e.ListIndex.ToString())
			}
		}
		return fmt.Sprintf("%s = %s", path_str, s.Args[0].ToString())
	}
	return "[unknown statement type]"
}

func (s *Statement) Execute(this *Value, vars map[string]*Value) (*Value, map[string]*Value, error) {
	fmt.Println("EXEC",s.ToString())
	if s.Type == STMT_ASSIGN {
		return s.Destination.Write(this, vars, s.Args[0])
	}
	return nil, nil, errors.New("Unknown statement type")
}
