package router

type Transform struct {
	Type string
	Condition *Expression
	Script *Script
	Source string
}

func (t *Transform) ToString() string {
	return t.Type
}

func (t *Transform) Execute(input *Value) *Value {
	return nil
}
