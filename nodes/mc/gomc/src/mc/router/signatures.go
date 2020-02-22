package router

import (
	"fmt"
	"errors"
)

type Signature struct {
	Operation ExpressionOperationType
	ArgTypes []ValueType
	ReturnType ValueType
	Handler func(*Value, map[string]*Value, []*Value, *Value) (*Value, error) // this, local_vars, args, primitive -> result, err
}

func FindSignature(op ExpressionOperationType, args []*Value) *Signature {
	for i, a := range args {
		fmt.Printf("arg[%d]:%d\n",i,a.Type)
	}
	for _, s := range ExpressionSignatures {
		if s.Operation == op {
			if len(s.ArgTypes) == 1 && s.ArgTypes[0] == VAL_ANYANY {
				return s
			}
			if len(s.ArgTypes) == len(args) {
				ok := true
				fmt.Println("AT",s.ArgTypes)
				for i, a := range args {
					if a.Type != s.ArgTypes[i] && s.ArgTypes[i] != VAL_ANY {
						ok = false
					}
				}
				if ok {
					return s
				}
			}
		}
	}
	return nil
}

func (s *Signature) TypeCheck(args []*Value) error {
	if len(args) != len(s.ArgTypes) {
		return errors.New(fmt.Sprintf("Signature has arity %d, got %d args", len(s.ArgTypes), len(args)))
	}
	for i, a := range args {
		if s.ArgTypes[i] != a.Type {
			return errors.New(fmt.Sprintf("Argument %d has type mismatch: Wanted %d, got %d", i, s.ArgTypes[i], a.Type))
		}
	}
	return nil
}

type StatementSignature struct {
	Operation StatementType
	ArgTypes []ValueType
	Handler func(*Value, map[string]*Value, []*Value) (*Value, map[string]*Value, error) // this, local_vars, args -> updated_this, updated_local_vars, err
}

var (
	// StatementSignatures = []*StatementSignature{
	// 	&StatementSignature{
	// 		Operation: STMT_ASSIGN,
	// 		ArgTypes:[]ValueType{VAL_ANY},
	// 		Handler: AssignHandler}}
	ExpressionSignatures = []*Signature{
		&Signature{
			Operation: OP_NAME,
			ArgTypes:[]ValueType{},
			ReturnType: VAL_NAME,
			Handler: NameHandler},
		&Signature{
			Operation: OP_NUM,
			ArgTypes:[]ValueType{},
			ReturnType: VAL_NUM,
			Handler: NumHandler},
		&Signature{
			Operation: OP_STRING,
			ArgTypes:[]ValueType{},
			ReturnType: VAL_STRING,
			Handler: StringHandler},
		&Signature{
			Operation: OP_CALL,
			ArgTypes:[]ValueType{VAL_NAME,VAL_LIST},
			ReturnType: VAL_ANY,
			Handler: CallHandler},
		&Signature{
			Operation: OP_MAP,
			ArgTypes:[]ValueType{VAL_ANYANY},
			ReturnType: VAL_MAP,
			Handler: MapHandler},
		&Signature{
			Operation: OP_LIST,
			ArgTypes:[]ValueType{VAL_ANYANY},
			ReturnType: VAL_LIST,
			Handler: ListHandler},
		&Signature{
			Operation: OP_MAPVAR,
			ArgTypes:[]ValueType{VAL_MAP,VAL_NAME},
			ReturnType: VAL_ANY,
			Handler: MapGetHandler},
		&Signature{
			Operation: OP_LISTVAR,
			ArgTypes:[]ValueType{VAL_LIST,VAL_NUM},
			ReturnType: VAL_ANY,
			Handler: ListGetHandler},
		&Signature{
			Operation: OP_VAR,
			ArgTypes:[]ValueType{},
			ReturnType: VAL_ANY,
			Handler: VarHandler},
		&Signature{
			Operation: OP_UMINUS,
			ArgTypes:[]ValueType{VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: UnaryMinusNumHandler},
		&Signature{
			Operation: OP_PLUS,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: AddNumHandler},
		&Signature{
			Operation: OP_PLUS,
			ArgTypes:[]ValueType{VAL_STRING,VAL_STRING},
			ReturnType: VAL_STRING,
			Handler: AddStringHandler},
		&Signature{
			Operation: OP_MINUS,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: SubNumHandler},
		&Signature{
			Operation: OP_MUL,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: MulNumHandler},
		&Signature{
			Operation: OP_MUL,
			ArgTypes:[]ValueType{VAL_STRING,VAL_NUM},
			ReturnType: VAL_STRING,
			Handler: MulStringNumHandler},
		&Signature{
			Operation: OP_DIV,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: DivNumHandler},
		&Signature{
			Operation: OP_MOD,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: ModNumHandler},
		&Signature{
			Operation: OP_BITWISEXOR,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: XorNumHandler},
		&Signature{
			Operation: OP_BITWISEAND,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: AndNumHandler},
		&Signature{
			Operation: OP_BITWISEOR,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_NUM,
			Handler: OrNumHandler},
		&Signature{
			Operation: OP_GT,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_BOOL,
			Handler: GtNumHandler},
		&Signature{
			Operation: OP_GT,
			ArgTypes:[]ValueType{VAL_STRING,VAL_STRING},
			ReturnType: VAL_BOOL,
			Handler: GtStringHandler},
		&Signature{
			Operation: OP_LT,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_BOOL,
			Handler: LtNumHandler},
		&Signature{
			Operation: OP_LT,
			ArgTypes:[]ValueType{VAL_STRING,VAL_STRING},
			ReturnType: VAL_BOOL,
			Handler: LtStringHandler},
		&Signature{
			Operation: OP_EQ,
			ArgTypes:[]ValueType{VAL_ANY,VAL_ANY},
			ReturnType: VAL_BOOL,
			Handler: EqHandler},
		&Signature{
			Operation: OP_NE,
			ArgTypes:[]ValueType{VAL_ANY,VAL_ANY},
			ReturnType: VAL_BOOL,
			Handler: NeHandler},
		&Signature{
			Operation: OP_LE,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_BOOL,
			Handler: LeqNumHandler},
		&Signature{
			Operation: OP_LE,
			ArgTypes:[]ValueType{VAL_STRING,VAL_STRING},
			ReturnType: VAL_BOOL,
			Handler: LeqStringHandler},
		&Signature{
			Operation: OP_GE,
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_BOOL,
			Handler: GeqNumHandler},
		&Signature{
			Operation: OP_GE,
			ArgTypes:[]ValueType{VAL_STRING,VAL_STRING},
			ReturnType: VAL_BOOL,
			Handler: GeqStringHandler},
		&Signature{
			Operation: OP_OR,
			ArgTypes:[]ValueType{VAL_BOOL,VAL_BOOL},
			ReturnType: VAL_BOOL,
			Handler: OrBoolHandler},
		&Signature{
			Operation: OP_AND,
			ArgTypes:[]ValueType{VAL_BOOL,VAL_BOOL},
			ReturnType: VAL_BOOL,
			Handler: AndBoolHandler},
		&Signature{
			Operation: OP_NOT,
			ArgTypes:[]ValueType{VAL_BOOL},
			ReturnType: VAL_BOOL,
			Handler: NotBoolHandler}}

)
