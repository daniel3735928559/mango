package router


type Signature struct {
	Operation ExpressionOperationType
	ArgTypes []ValueType
	ReturnType ValueType
	Handler func([]*Value, *map[string]*Value) (*Value, error)
}

var (
	ExpressionSignatures = []*Signature{
		&Signature{
			Operation: OP_ASSIGN,
			ArgTypes:[]ValueType{VAL_NAME,VAL_ANY},
			ReturnType: VAL_ANY,
			Handler: AssignHandler},
		&Signature{
			Operation: OP_CALL,
			ArgTypes:[]ValueType{VAL_NAME,VAL_LIST},
			ReturnType: VAL_ANY,
			Handler: CallHandler},
		&Signature{
			Operation: OP_MAP,
			ArgTypes:[]ValueType{VAL_ANY},
			ReturnType: VAL_ANY,
			Handler: MapHandler},
		&Signature{
			Operation: OP_LIST,
			ArgTypes:[]ValueType{VAL_ANY},
			ReturnType: VAL_ANY,
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
			ArgTypes:[]ValueType{VAL_NAME},
			ReturnType: VAL_ANY,
			Handler: ValHandler},
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
			ArgTypes:[]ValueType{VAL_NUM,VAL_NUM},
			ReturnType: VAL_BOOL,
			Handler: EqNumHandler},
		&Signature{
			Operation: OP_EQ,
			ArgTypes:[]ValueType{VAL_STRING,VAL_STRING},
			ReturnType: VAL_BOOL,
			Handler: EqStringHandler},
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
