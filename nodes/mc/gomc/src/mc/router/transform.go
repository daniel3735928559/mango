package router

import (
	"fmt"
	"errors"
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
	Condition *Expression
	Replace *Expression
	Script []*Statement
}

func (t *Transform) ToString() string {
	if t.Type == TR_FILTER {
		return fmt.Sprintf("pass if {%s}", t.Condition.ToString())
	} else if t.Type == TR_EDIT {
		script := ""
		for _, e := range t.Script {
			script += fmt.Sprintf("%s;", e.ToString())
		}
		return fmt.Sprintf("edit {%s}", script)
	} else if t.Type == TR_REPLACE {
		return fmt.Sprintf("replace %s", t.Replace.ToString())
	} else if t.Type == TR_COND_EDIT {
		script := ""
		for _, e := range t.Script {
			script += fmt.Sprintf("%s;", e.ToString())
		}
		return fmt.Sprintf("edit {%s} if {%s}", script, t.Condition.ToString())
	} else if t.Type == TR_COND_REPLACE {
		return fmt.Sprintf("replace %s if {%s}", t.Replace.ToString(), t.Condition.ToString())
	} 
	return "[unknown transform type]"
}

func (t *Transform) EvaluateCondition(this *Value) (bool, error) {
	fmt.Println("Condition",t.Condition.ToString(),"on",this.ToString())
	if t.Condition == nil {
		return false, errors.New("Tried to evaluate nil condition")
	}
	res, err := t.Condition.Evaluate(this, make(map[string]*Value))
	if err != nil {
		return false, err
	}
	if res.Type != VAL_BOOL {
		return false, errors.New("condition does not evaluate to a boolean")
	}
	fmt.Println("Condition result",res.BoolVal)
	return res.BoolVal, nil
}

func (t *Transform) EvaluateScript(this *Value) (*Value, error) {
	vars := make(map[string]*Value)
	var err error
	for _, e := range t.Script {
		fmt.Println("Script",e.ToString(),"on",this.ToString())
		this, vars, err = e.Execute(this, vars)
		if err != nil {
			return nil, err
		}
	}
	return this, nil
}

func (t *Transform) EvaluateReplacement(this *Value) (*Value, error) {
	fmt.Println("Replacement",t.Replace.ToString(),"on",this.ToString())
	replacement, err := t.Replace.Evaluate(this, make(map[string]*Value))
	if err != nil {
		return nil, err
	}
	return replacement, nil
}



func (t *Transform) Execute(args *Value) (*Value, error) {
	this := args.Clone()
	if t.Type == TR_FILTER {
		to_pass, err := t.EvaluateCondition(this)
		if err != nil {
			return nil, err
		}
		if to_pass {
			return this, nil
		}
		return nil, nil
	} else if t.Type == TR_EDIT {
		return t.EvaluateScript(this)
	} else if t.Type == TR_REPLACE {
		return t.EvaluateReplacement(this)
	} else if t.Type == TR_COND_EDIT {
		to_pass, err := t.EvaluateCondition(this)
		if err != nil {
			return nil, err
		}
		if to_pass {
			return t.EvaluateScript(this)
		}
		return this, nil
	} else if t.Type == TR_COND_REPLACE {
		to_pass, err := t.EvaluateCondition(this)
		if err != nil {
			return nil, err
		}
		if to_pass {
			return t.EvaluateReplacement(this)
		}
		return this, nil
	}	
	return nil, errors.New("Transform of unknown type")
}
