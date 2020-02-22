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
	OP_NAME
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
	OP_NE
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
		Operation: OP_NAME,
		Value: &Value {
			Type:VAL_NAME,
			NameVal:name}}
}

func MakeVarExpression(name string) *Expression {
	return &Expression {
		Operation: OP_VAR,
		Value: &Value {
			Type:VAL_NAME,
			NameVal:name}}
}

func (e *Expression) ToString() string {
	if e.Operation == OP_VAR {
		return fmt.Sprintf("%s", e.Value.ToString())
	} else if e.Operation == OP_NAME {
		return fmt.Sprintf("NAME(%s)", e.Value.NameVal)
	} else if e.Operation == OP_MAPVAR {
		return fmt.Sprintf("%s.%s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_LISTVAR {
		return fmt.Sprintf("%s[%s]", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_NUM || e.Operation == OP_STRING {
		return e.Value.ToString()
	} else if e.Operation == OP_CALL {
		return fmt.Sprintf("%s(%s)", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_UMINUS {
		return fmt.Sprintf("-%s", e.Args[0].ToString())
	} else if e.Operation == OP_MAPGET {
		return fmt.Sprintf("%s[%s]", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_LISTGET {
		return fmt.Sprintf("%s[%s]", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_MAP {
		fmt.Println("MA",e.Args,len(e.Args))
		mapargs := make([]string, len(e.Args)/2)
		for i := 0; i < len(e.Args)/2; i++ {
			mapargs[i] = fmt.Sprintf("%s:%s", e.Args[2*i].ToString(), e.Args[2*i+1].ToString())
			fmt.Println("ma",mapargs[i])
		}
		return fmt.Sprintf("{%s}",strings.Join(mapargs, ","))
	} else if e.Operation == OP_LIST {
		listargs := make([]string, len(e.Args))
		for i := 0; i < len(e.Args); i++ {
			listargs[i] = e.Args[i].ToString()
		}
		return fmt.Sprintf("[%s]",strings.Join(listargs, ","))
	} else if e.Operation == OP_PLUS {
		return fmt.Sprintf("%s + %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_MINUS {
		return fmt.Sprintf("%s - %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_MUL {
		return fmt.Sprintf("%s * %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_DIV {
		return fmt.Sprintf("%s / %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_MOD {
		return fmt.Sprintf("%s %% %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_BITWISEAND {
		return fmt.Sprintf("%s & %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_BITWISEOR {
		return fmt.Sprintf("%s | %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_BITWISEXOR {
		return fmt.Sprintf("%s ^ %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_MATCH {
		return fmt.Sprintf("%s ~ %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_SUB {
		return fmt.Sprintf("%s ~~ %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_EQ {
		return fmt.Sprintf("%s == %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_NE {
		return fmt.Sprintf("%s != %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_GT {
		return fmt.Sprintf("%s > %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_LT {
		return fmt.Sprintf("%s < %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_GE {
		return fmt.Sprintf("%s >= %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_LE {
		return fmt.Sprintf("%s <= %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_AND {
		return fmt.Sprintf("%s && %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_OR {
		return fmt.Sprintf("%s || %s", e.Args[0].ToString(), e.Args[1].ToString())
	} else if e.Operation == OP_NOT {
		return fmt.Sprintf("!%s", e.Args[0].ToString())
	}
	return fmt.Sprintf("[unknown expression type %d]",e.Operation)
	// subexprs := ""
	// val := ""
	// if e.Args != nil {
	// 	subexprs_strings := make([]string, len(e.Args))
	// 	for i, es := range e.Args {
	// 		subexprs_strings[i] = es.ToString()
	// 	}
	// 	subexprs = strings.Join(subexprs_strings, ",")
	// }
	// if e.Value != nil {
	// 	val = e.Value.ToString()
	// }
	// return fmt.Sprintf("%d(%s;%s)",e.Operation, subexprs, val)
}

// func (e *Expression) TypeCheck() *Signature {
// 	fmt.Println("TypeCheck",e.ToString(),"op",e.Operation)
// 	arg_types := make([]ValueType, len(e.Args))
// 	for i, a := range e.Args {
// 		fmt.Println("checking arg",a.ToString())
// 		sig := a.TypeCheck()
// 		if sig == nil {
// 			return nil
// 		}
// 		arg_types[i] = sig.ReturnType
// 		fmt.Println("arg type",i,sig.ReturnType)
// 	}
	
// 	for _, sig := range ExpressionSignatures {
// 		if e.Operation == sig.Operation {
// 			ok := true
// 			for i, a := range arg_types {
// 				if i > len(sig.ArgTypes) || (sig.ArgTypes[i] != VAL_ANY && a != VAL_ANY && a != sig.ArgTypes[i]) {
// 					ok = false
// 				}
// 			}
// 			if ok {
// 				fmt.Println("FOUND SIG",arg_types,sig.ArgTypes)
// 				return sig
// 			}
// 		}
// 	}
// 	fmt.Println("No type for ",e.ToString())
// 	return nil
// }

func (e *Expression) Evaluate(this *Value, vars map[string]*Value) (*Value, error) {
	if e == nil {
		return nil, errors.New("Invalid expression")
	}
	fmt.Println("EVAL",e.ToString())
	args := make([]*Value, len(e.Args))
	local_vars := vars
	var err error
	var arg *Value
	for i, a := range e.Args {
		arg, err = a.Evaluate(this, local_vars)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}
	fmt.Println("Searching signature",e.ToString())
	sig := FindSignature(e.Operation, args)
	if sig == nil {
		return nil, errors.New("No valid type found for expression")
	}
	ans, err := sig.Handler(this, local_vars, args, e.Value)
	if ans != nil {
		fmt.Println("EVALed",e.ToString(),"=",ans.ToString())
	} else {
		fmt.Println("EVALed",e.ToString(),"= nil")
	}
	//fmt.Println("EVALed",e.ToString(),"=",ans.ToString())
	return ans, err
}
