package router

type MCType interface {
	Value() []byte
	
}

type MCListType struct {
	ElementTypeName string
}

type MCDictType struct {
	ElementTypeNames map[string]string
}

type MCStringType struct {
	Encoding string
}

type MCIntType struct {
	
}

type MCFloatType struct {
	
}

type MCValue struct {
	Type *MCType // list|dict|string|int|float
	Value interface{}
}

func (v *MCValue) IntValue() int {
	return v.Value.(int)
}

func (v *MCValue) GetListElement(i int) *MCValue {
	return v.Value.([]*MCValue)[i]
}

func (v *MCValue) GetDictElement(key string) *MCValue {
	return v.Value.(map[string]*MCValue)[key]
}

func (v *MCValue) GetDictElement(key string) *MCValue {
	return v.Value.(map[string]*MCValue)[key]
}
