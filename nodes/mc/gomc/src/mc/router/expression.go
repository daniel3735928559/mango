package router

import (
	"fmt"
	"errors"
	"strings"
)

type ExpressionOperationType int

const (
	OP_ASSIGN ExpressionOperationType = iota + 1
	OP_VAR
	OP_MAPVAR
	OP_LISTVAR
	OP_NUM
	OP_STRING
	OP_CALL
	OP_UMINUS
	OP_MAPGET
	OP_LISTGET
	OP_PLUS
	OP_MINUS
	OP_MUL
	OP_DIV
	OP_MOD
	OP_BITWISEAND
	OP_BITWISEOR
	OP_BITWISEXOR
	OP_MATCH
	OP_SUB
	OP_MAP
	OP_LIST
	OP_EQ
	OP_GT
	OP_LT
	OP_GE
	OP_LE
	OP_AND
	OP_OR
	OP_NOT
)

type Expression struct {
	Operation ExpressionOperationType
	Args []*Expression
	Value *Value
}			

func MakeNameExpression(name string) *Expression {
	return &Expression {
		Operation: OP_VAR,
		Value: &Value {
			Type:VAL_NAME,
			NameVal:name}}
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
	return fmt.Sprintf("%d(%s;%s)",e.Operation, subexprs, val)
}

func (e *Expression) TypeCheck() *Signature {
	arg_types := make([]ValueType, len(e.Args))
	for i, a := range e.Args {
		if a.Operation == OP_NUM || a.Operation == OP_STRING {
			arg_types[i] = a.Value.Type
		} else {
			sig := a.TypeCheck()
			if sig == nil {
			return nil
			}
			arg_types[i] = sig.ReturnType
		}
	}
		
	for _, sig := range ExpressionSignatures {
		if e.Operation == sig.Operation {
			ok := true
			for i, a := range arg_types {
				if i > len(sig.ArgTypes) || (sig.ArgTypes[i] != VAL_ANY && a != VAL_ANY && a != sig.ArgTypes[i]) {
					ok = false
				}
			}
			if ok {
				return sig
			}
		}
	}
	return nil
}

func (e *Expression) Evaluate(context *map[string]*Value) (*Value, error) {
	sig := e.TypeCheck()
	if sig == nil {
		return nil, errors.New("No valid type found for expression")
	}
	vals := make([]*Value, len(e.Args))
	for i, a := range e.Args {
		v, err := a.Evaluate(context)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	vars := make(map[string]*Value)
	return sig.Handler(vals, &vars)
}
