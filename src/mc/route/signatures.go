package route

import (
	"fmt"
	value "mc/value"
)

type Signature struct {
	Operation ExpressionOperationType
	ArgTypes []value.ValueKind
	ReturnType value.ValueKind
	Handler func(*value.Value, map[string]*value.Value, []*value.Value, *value.Value) (*value.Value, error) // this, local_vars, args, primitive -> result, err
}

func FindSignature(op ExpressionOperationType, args []*value.Value) *Signature {
	for i, a := range args {
		fmt.Printf("arg[%d]:%d\n",i,a.Type)
	}
	for _, s := range ExpressionSignatures {
		if s.Operation == op {
			if len(s.ArgTypes) == 1 && s.ArgTypes[0] == value.VAL_ANYANY {
				return s
			}
			if len(s.ArgTypes) == len(args) {
				ok := true
				fmt.Println("AT",s.ArgTypes)
				for i, a := range args {
					if a.Type != s.ArgTypes[i] && s.ArgTypes[i] != value.VAL_ANY {
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

// func (s *Signature) TypeCheck(args []*value.Value) error {
// 	if len(args) != len(s.ArgTypes) {
// 		return errors.New(fmt.Sprintf("Signature has arity %d, got %d args", len(s.ArgTypes), len(args)))
// 	}
// 	for i, a := range args {
// 		if s.ArgTypes[i] != a.Type {
// 			return errors.New(fmt.Sprintf("Argument %d has type mismatch: Wanted %d, got %d", i, s.ArgTypes[i], a.Type))
// 		}
// 	}
// 	return nil
// }

type StatementSignature struct {
	Operation StatementType
	ArgTypes []value.ValueKind
	Handler func(*value.Value, map[string]*value.Value, []*value.Value) (*value.Value, map[string]*value.Value, error) // this, local_vars, args -> updated_this, updated_local_vars, err
}

var (
	// StatementSignatures = []*StatementSignature{
	// 	&StatementSignature{
	// 		Operation: STMT_ASSIGN,
	// 		ArgTypes:[]value.ValueKind{value.VAL_ANY},
	// 		Handler: AssignHandler}}
	ExpressionSignatures = []*Signature{
		&Signature{
			Operation: OP_NAME,
			ArgTypes:[]value.ValueKind{},
			ReturnType: value.VAL_NAME,
			Handler: NameHandler},
		&Signature{
			Operation: OP_NUM,
			ArgTypes:[]value.ValueKind{},
			ReturnType: value.VAL_NUM,
			Handler: NumHandler},
		&Signature{
			Operation: OP_BOOL,
			ArgTypes:[]value.ValueKind{},
			ReturnType: value.VAL_NUM,
			Handler: BoolHandler},
		&Signature{
			Operation: OP_STRING,
			ArgTypes:[]value.ValueKind{},
			ReturnType: value.VAL_STRING,
			Handler: StringHandler},
		&Signature{
			Operation: OP_CALL,
			ArgTypes:[]value.ValueKind{value.VAL_NAME,value.VAL_LIST},
			ReturnType: value.VAL_ANY,
			Handler: CallHandler},
		&Signature{
			Operation: OP_MAP,
			ArgTypes:[]value.ValueKind{value.VAL_ANYANY},
			ReturnType: value.VAL_MAP,
			Handler: MapHandler},
		&Signature{
			Operation: OP_LIST,
			ArgTypes:[]value.ValueKind{value.VAL_ANYANY},
			ReturnType: value.VAL_LIST,
			Handler: ListHandler},
		&Signature{
			Operation: OP_MAPVAR,
			ArgTypes:[]value.ValueKind{value.VAL_MAP,value.VAL_NAME},
			ReturnType: value.VAL_ANY,
			Handler: MapGetHandler},
		&Signature{
			Operation: OP_LISTVAR,
			ArgTypes:[]value.ValueKind{value.VAL_LIST,value.VAL_NUM},
			ReturnType: value.VAL_ANY,
			Handler: ListGetHandler},
		&Signature{
			Operation: OP_LISTVAR,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_NUM},
			ReturnType: value.VAL_STRING,
			Handler: StringGetHandler},
		&Signature{
			Operation: OP_VAR,
			ArgTypes:[]value.ValueKind{},
			ReturnType: value.VAL_ANY,
			Handler: VarHandler},
		&Signature{
			Operation: OP_MATCH,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_STRING},
			ReturnType: value.VAL_BOOL,
			Handler: MatchHandler},
		&Signature{
			Operation: OP_TERNARY,
			ArgTypes:[]value.ValueKind{value.VAL_BOOL,value.VAL_ANY,value.VAL_ANY},
			ReturnType: value.VAL_ANY,
			Handler: TernaryHandler},
		&Signature{
			Operation: OP_EXP,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: ExpNumHandler},
		&Signature{
			Operation: OP_UMINUS,
			ArgTypes:[]value.ValueKind{value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: UnaryMinusNumHandler},
		&Signature{
			Operation: OP_PLUS,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: AddNumHandler},
		&Signature{
			Operation: OP_PLUS,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_STRING},
			ReturnType: value.VAL_STRING,
			Handler: AddStringHandler},
		&Signature{
			Operation: OP_PLUS,
			ArgTypes:[]value.ValueKind{value.VAL_LIST,value.VAL_LIST},
			ReturnType: value.VAL_LIST,
			Handler: AddListHandler},
		&Signature{
			Operation: OP_PLUS,
			ArgTypes:[]value.ValueKind{value.VAL_MAP,value.VAL_MAP},
			ReturnType: value.VAL_MAP,
			Handler: AddMapHandler},
		&Signature{
			Operation: OP_MINUS,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: SubNumHandler},
		&Signature{
			Operation: OP_MUL,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: MulNumHandler},
		&Signature{
			Operation: OP_MUL,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_NUM},
			ReturnType: value.VAL_STRING,
			Handler: MulStringNumHandler},
		&Signature{
			Operation: OP_DIV,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: DivNumHandler},
		&Signature{
			Operation: OP_MOD,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: ModNumHandler},
		&Signature{
			Operation: OP_BITWISEXOR,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: XorNumHandler},
		&Signature{
			Operation: OP_BITWISEAND,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: AndNumHandler},
		&Signature{
			Operation: OP_BITWISEOR,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_NUM,
			Handler: OrNumHandler},
		&Signature{
			Operation: OP_GT,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_BOOL,
			Handler: GtNumHandler},
		&Signature{
			Operation: OP_GT,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_STRING},
			ReturnType: value.VAL_BOOL,
			Handler: GtStringHandler},
		&Signature{
			Operation: OP_LT,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_BOOL,
			Handler: LtNumHandler},
		&Signature{
			Operation: OP_LT,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_STRING},
			ReturnType: value.VAL_BOOL,
			Handler: LtStringHandler},
		&Signature{
			Operation: OP_EQ,
			ArgTypes:[]value.ValueKind{value.VAL_ANY,value.VAL_ANY},
			ReturnType: value.VAL_BOOL,
			Handler: EqHandler},
		&Signature{
			Operation: OP_NE,
			ArgTypes:[]value.ValueKind{value.VAL_ANY,value.VAL_ANY},
			ReturnType: value.VAL_BOOL,
			Handler: NeHandler},
		&Signature{
			Operation: OP_LE,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_BOOL,
			Handler: LeqNumHandler},
		&Signature{
			Operation: OP_LE,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_STRING},
			ReturnType: value.VAL_BOOL,
			Handler: LeqStringHandler},
		&Signature{
			Operation: OP_GE,
			ArgTypes:[]value.ValueKind{value.VAL_NUM,value.VAL_NUM},
			ReturnType: value.VAL_BOOL,
			Handler: GeqNumHandler},
		&Signature{
			Operation: OP_GE,
			ArgTypes:[]value.ValueKind{value.VAL_STRING,value.VAL_STRING},
			ReturnType: value.VAL_BOOL,
			Handler: GeqStringHandler},
		&Signature{
			Operation: OP_OR,
			ArgTypes:[]value.ValueKind{value.VAL_BOOL,value.VAL_BOOL},
			ReturnType: value.VAL_BOOL,
			Handler: OrBoolHandler},
		&Signature{
			Operation: OP_AND,
			ArgTypes:[]value.ValueKind{value.VAL_BOOL,value.VAL_BOOL},
			ReturnType: value.VAL_BOOL,
			Handler: AndBoolHandler},
		&Signature{
			Operation: OP_NOT,
			ArgTypes:[]value.ValueKind{value.VAL_BOOL},
			ReturnType: value.VAL_BOOL,
			Handler: NotBoolHandler}}

)
