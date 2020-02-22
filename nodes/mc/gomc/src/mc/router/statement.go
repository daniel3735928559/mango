package router

import (
	"fmt"
	"errors"
)

type StatementType int
type PathStepType int

const (
	STMT_ASSIGN StatementType = iota + 1
	STMT_DECLARE
	STMT_DELETE
)

const (
	PATH_MAP PathStepType = iota + 1
	PATH_LIST
)

type Statement struct {
	Type StatementType
	Destination *WriteableValue
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

func (w *WriteableValue) ToExpression() *Expression {
	base := &Expression{
		Operation:OP_VAR,
		Value:&Value{
			Type:VAL_NAME,
			NameVal:w.Base}}
	current := base
	for _, pe := range w.Path {
		if pe.Type == PATH_MAP {
			current = &Expression{
				Operation:OP_MAPVAR,
				Args:[]*Expression{current, MakeNameExpression(pe.MapKey)}}
		} else if pe.Type  == PATH_LIST {
			current = &Expression{
				Operation:OP_LISTVAR,
				Args:[]*Expression{current, pe.ListIndex}}
		}
	}
	return current
}

func (w *WriteableValue) Write(this *Value, vars map[string]*Value, arg *Expression) (*Value, map[string]*Value, error) {
	content, err := arg.Evaluate(this, vars)
	if err != nil {
		return this, vars, err
	}
	var dest *Value
	dest_name := w.Base
	is_this_var := true
	is_local_var := true
	if dest_name == "this" {
		dest = this.Clone()
		is_this_var = false
		is_local_var = false
	} else if v, ok := vars[dest_name]; ok {
		dest = v.Clone()
		is_this_var = false
		is_local_var = true
	} else if v, ok := this.MapVal[dest_name]; ok && this.Type == VAL_MAP {
		dest = v.Clone()
		is_this_var = true
		is_local_var = false
	} else if len(w.Path) == 0 {
		dest = MakeEmptyValue()
		is_this_var = true
		is_local_var = false
	} else {
		return this, vars, errors.New(fmt.Sprintf("No such variable: %s", dest_name))
	}

	new_this := this.Clone()
	target_base := dest
	target := target_base
	if len(w.Path) == 0 {
		if is_this_var {
			// We're assigning to a value in this
			new_this.MapVal[dest_name] = content
			return new_this, vars, nil
		} else if is_local_var {
			// We're assigning an existing local variable
			vars[dest_name] = content
			return new_this, vars, nil
		} else {
			// We're returning a wholesale this object
			return content, vars, nil
		}
	}
	for i, e := range w.Path {
		if e.Type == PATH_MAP {
			if target.Type != VAL_MAP {
				return this, vars, errors.New(fmt.Sprintf("Attempted to access key %s in non-map %d != %d", e.MapKey,target.Type,VAL_MAP))
			}
			if i == len(w.Path) - 1 {
				target.MapVal[e.MapKey] = content
				if is_this_var {
					// We're assigning to a value in this
					new_this.MapVal[dest_name] = target_base
					return new_this, vars, nil
				} else if is_local_var {
					// We're assigning an existing local variable
					vars[dest_name] = target_base
					return new_this, vars, nil
				} else {
					// We're returning a wholesale this object
					return content, vars, nil
				}
			} else if new_target, ok := target.MapVal[e.MapKey]; ok {
				fmt.Println("Descending into map by key",e.MapKey,"new type",new_target.Type)
				target = new_target
			} else {
				return this, vars, errors.New(fmt.Sprintf("Attempted to access non-existent key %s", e.MapKey))
			}
		} else if e.Type == PATH_LIST {
			if target.Type != VAL_LIST {
				return this, vars, errors.New(fmt.Sprintf("Attempted to access index in non-list type %d != %d",target.Type,VAL_LIST))
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
				target.ListVal[list_index] = content
				if is_this_var {
					// We're assigning to a value in this
					new_this.MapVal[dest_name] = target_base
					return new_this, vars, nil
				} else if is_local_var {
					// We're assigning an existing local variable
					vars[dest_name] = target_base
					return new_this, vars, nil
				} else {
					// We're returning a wholesale this object
					return content, vars, nil
				}
			} else {
				fmt.Println("Descending into list by index",list_index,"new type",target.ListVal[list_index].Type)
				target = target.ListVal[list_index]
			}
		}
	}
	return this, vars, errors.New("Something went wrong?")
}

func MakeAssignmentStatement(dest *WriteableValue, val *Expression) *Statement {
	return &Statement {
		Type: STMT_ASSIGN,
		Destination: dest,
		Args: []*Expression{val}}
}
func MakeDeclarationStatement(name string) *Statement {
	return &Statement {
		Type: STMT_DECLARE,
		Destination: &WriteableValue{Base:name}}
}
func MakeDeletionStatement(name string) *Statement {
	return &Statement {
		Type: STMT_DELETE,
		Destination: &WriteableValue{Base:name}}
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
	} else if s.Type == STMT_DECLARE {
		name := s.Destination.Base
		if this.Type == VAL_MAP {
			if _, ok := this.MapVal[name]; ok {
				return this, vars, errors.New(fmt.Sprintf("Variable already exists in this: %s", name))
			}
		}
		if _, ok := vars[name]; ok {
			return this, vars, errors.New(fmt.Sprintf("Variable already exists as local variable: %s", name))
		}
		vars[name] = MakeEmptyValue()
		return this, vars, nil
	} else if s.Type == STMT_DELETE {
		name := s.Destination.Base
		if this.Type == VAL_MAP {
			if _, ok := this.MapVal[name]; ok {
				delete(this.MapVal, name)
				return this, vars, nil
			}
		}
		if _, ok := vars[name]; ok {
			delete(vars, name)
			return this, vars, nil
		}
		return this, vars, errors.New(fmt.Sprintf("Could not find variable: %s", name))
	}
	return nil, nil, errors.New("Unknown statement type")
}
