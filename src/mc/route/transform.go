package route

import (
	"fmt"
	"errors"
	value "mc/value"
)

type TransformType int

const (
	TR_FILTER TransformType = iota + 1
	TR_EDIT
	TR_REPLACE
	TR_COND_EDIT
	TR_COND_REPLACE
)

type Transform struct {
	Type TransformType
	CommandCondition string
	CommandReplace string
	Condition *Expression
	Replace *Expression
	Script []*Statement
}

func (t *Transform) ToString() string {
	if t.Type == TR_FILTER {
		fmt.Println(t.Condition)
		cond := ""
		if t.Condition != nil {
			cond = t.Condition.ToString()
		}
		return fmt.Sprintf("pass if %s{%s}", t.CommandCondition, cond)
	} else if t.Type == TR_EDIT {
		script := ""
		if t.Script != nil {
			for _, e := range t.Script {
				script += fmt.Sprintf("%s;", e.ToString())
			}
		}
		return fmt.Sprintf("edit %s{%s}", t.CommandReplace, script)
	} else if t.Type == TR_REPLACE {
		repl := ""
		if t.Replace != nil {
			repl = t.Replace.ToString()
		}
		return fmt.Sprintf("replace %s%s", t.CommandReplace, repl)
	} else if t.Type == TR_COND_EDIT {
		script := ""
		if t.Script != nil {
			for _, e := range t.Script {
				script += fmt.Sprintf("%s;", e.ToString())
			}
		}
		cond := ""
		if t.Condition != nil {
			cond = t.Condition.ToString()
		}
		return fmt.Sprintf("edit %s{%s} if %s{%s}", t.CommandReplace, script, t.CommandCondition, cond)
	} else if t.Type == TR_COND_REPLACE {
		cond := ""
		repl := ""
		if t.Replace != nil {
			repl = t.Replace.ToString()
		}
		if t.Condition != nil {
			cond = t.Condition.ToString()
		}
		return fmt.Sprintf("replace %s%s if %s{%s}", t.CommandReplace, repl, t.CommandCondition, cond)
	}
	return "[unknown transform type]"
}

func (t *Transform) EvaluateCondition(command string, this *value.Value) (bool, error) {
	fmt.Println("Filter",t.ToString(),"on",command, this.ToString())
	if len(t.CommandCondition) > 0 {
		if t.CommandCondition != command {
			// The command has not matched--no need to check the condition
			return false, nil
		} else if t.Condition == nil {
			// We are only checking the command type, and so need not evaluate the condition
			return true, nil
		}
	}
	// The command has matched or was not checked
	if t.Condition == nil {
		return false, errors.New("Tried to evaluate nil condition")
	}
	res, err := t.Condition.Evaluate(this, make(map[string]*value.Value))
	if err != nil {
		return false, err
	}
	if res.Type != value.VAL_BOOL {
		return false, errors.New("condition does not evaluate to a boolean")
	}
	fmt.Println("Condition result",res.BoolVal)
	return res.BoolVal, nil
}

func (t *Transform) EvaluateScript(command string, this *value.Value) (string, *value.Value, error) {
	vars := make(map[string]*value.Value)
	var err error
	if t.Script != nil && len(t.Script) > 0 {
		for _, s := range t.Script {
			fmt.Println("Statement",s.ToString(),"on",this.ToString())
			this, vars, err = s.Execute(this, vars)
			if err != nil {
				return "", nil, err
			}
		}
	}
	// If there was no error (or no script), replace the command if needed
	if len(t.CommandReplace) > 0 {
		command = t.CommandReplace
	}
	return command, this, nil
}

func (t *Transform) EvaluateReplacement(command string, this *value.Value) (string, *value.Value, error) {
	fmt.Println("Replacement",t.ToString(),"on",command,this.ToString())
	if len(t.CommandReplace) > 0 {
		command = t.CommandReplace
	}
	if t.Replace != nil {
		replacement, err := t.Replace.Evaluate(this, make(map[string]*value.Value))
		if err != nil {
			return "", nil, err
		}
		return command, replacement, nil
	} else {
		return command, this, nil
	}
}

func (t *Transform) Execute(command string, args *value.Value) (string, *value.Value, error) {
	this := args.Clone()
	if t.Type == TR_FILTER {
		to_pass, err := t.EvaluateCondition(command, this)
		if err != nil {
			return "", nil, err
		}
		if to_pass {
			return command, this, nil
		}
		return "", nil, nil
	} else if t.Type == TR_EDIT {
		return t.EvaluateScript(command, this)
	} else if t.Type == TR_REPLACE {
		return t.EvaluateReplacement(command, this)
	} else if t.Type == TR_COND_EDIT {
		to_edit, err := t.EvaluateCondition(command, this)
		if err != nil {
			return "", nil, err
		}
		if to_edit {
			return t.EvaluateScript(command, this)
		}
		// The filter did not match--just pass through unedited
		return command, this, nil
	} else if t.Type == TR_COND_REPLACE {
		to_replace, err := t.EvaluateCondition(command, this)
		if err != nil {
			return "", nil, err
		}
		if to_replace {
			return t.EvaluateReplacement(command, this)
		}
		// The filter did not match--just pass through unreplaced
		return command, this, nil
	}	
	return "", nil, errors.New("Transform of unknown type")
}


// func (t *Transform) TransformType(ty *TypeDesc) []*TypeDesc {
// 	this := args.Clone()
// 	if t.Type == TR_FILTER {
// 		// Assume (generously) any command of the correct kind
// 		// can pass the expression check
// 		if t.CommandCondition == ty.Name {
// 			return []*TypeDesc{ty}
// 		}
// 		return []*TypeDesc{}
// 	} else if t.Type == TR_EDIT {
// 		return t.EvaluateScript(command, this)
// 	} else if t.Type == TR_REPLACE {
// 		return t.EvaluateReplacement(command, this)
// 	} else if t.Type == TR_COND_EDIT {
// 		to_edit, err := t.EvaluateCondition(command, this)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 		if to_edit {
// 			return t.EvaluateScript(command, this)
// 		}
// 		// The filter did not match--just pass through unedited
// 		return command, this, nil
// 	} else if t.Type == TR_COND_REPLACE {
// 		to_replace, err := t.EvaluateCondition(command, this)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 		if to_replace {
// 			return t.EvaluateReplacement(command, this)
// 		}
// 		// The filter did not match--just pass through unreplaced
// 		return command, this, nil
// 	}	
// 	return "", nil, errors.New("Transform of unknown type")
// }
