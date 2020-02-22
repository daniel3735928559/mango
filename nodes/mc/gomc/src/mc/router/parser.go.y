%{
	package router
	import (
		"strconv"
		//"fmt"
	)
	type Token struct {
		token   int
		literal string
		position Position
	}
%}
%union{
	token Token
	routes []*Route
	transforms *Route
	transform *Transform
	expression *Expression
	statement *Statement
	writeable *WriteableValue
	script []*Statement
	node *Node
}
%type<routes> route
%type<node> node
%type<transforms> transforms
%type<transform> transform
%type<script> script
%type<writeable> dstexpr
%type<expression> varexpr
%type<statement> stmt
%type<expression> expr
%type<expression> mapexprs
%type<expression> listexprs

%token<token> IDENT VAR NUMBER STRING THIS AND OR EQ NE LE GE PE ME TE DE RE XE SUB '?' '%' '=' '{' '}' '[' ']' '<' '>' ':' '+' '-' '*' '/' '&' '|', '^', '!', '~'

%left AND OR
%left GE LE EQ '<' '>'
%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left UNARY '!'
%%
route   : node '>' node
{
	// fmt.Println("C")
	$$ = nil
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = []*Route{&Route{Source: $1.Name, Dest: $3.Name}}
	}
}
| node '<' node
{
	// fmt.Println("B")
	$$ = nil
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = []*Route{&Route{Source: $3.Name, Dest: $1.Name}}
	}
}
| node '<' '>' node
{
	$$ = nil
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = []*Route{
			&Route{Source: $1.Name, Dest: $4.Name},
			&Route{Source: $4.Name, Dest: $1.Name}}
	}
}
| node '>' transforms
{
	$$ = nil
	// fmt.Println("A")
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = []*Route{
				    &Route{
				    Source: $1.Name,
				    Dest: $3.Dest,
				    Transforms: $3.Transforms}}
}
}
;
node    : IDENT
{
	$$ = &Node{Name: $1.literal}
}
| IDENT '/' IDENT
{
	$$ = &Node{Group: $1.literal, Name: $3.literal}
}
;
transforms : transform '>' node
{
	$$ = &Route{
		Dest: $3.Name,
		Transforms: []*Transform{$1}}
}
| transform '>' transforms
{
	$$ = &Route{
		Dest: $3.Dest,
		Transforms: append([]*Transform{$1}, $3.Transforms...)}
}
;
transform : '?' '{' expr '}'
{
	$$ = &Transform{
		Type: TR_FILTER,
		Condition: $3}
}
| '%' '{' script '}'
{
	$$ = &Transform{
		Type: TR_EDIT,
		Script: $3}
}
| '=' '{' mapexprs '}'
{
	$$ = &Transform{
		Type: TR_REPLACE,
		Replace: $3}
}
| '?' '{' expr '}' '%' '{' script '}'
{
	$$ = &Transform{
		Type: TR_COND_EDIT,
		Condition: $3,
		Script: $7}
}
| '?' '{' expr '}' '=' '{' mapexprs '}'
{
	$$ = &Transform{
		Type: TR_COND_REPLACE,
		Condition: $3,
		Replace: $7}
}
;
script
: stmt
{
	$$ = []*Statement{$1}
}
| stmt script
{
	$$ = append([]*Statement{$1}, $2...)
}

stmt : dstexpr '=' expr ';'
{
	$$ = MakeAssignment($1, $3)
}
| dstexpr PE expr ';'
{
	$$ = MakeAssignment($1, &Expression{
		Operation: OP_PLUS,
		Args: []*Expression{$1.ToExpression(), $3}})
}
| dstexpr ME expr ';'
{
	$$ = MakeAssignment($1, &Expression{
		Operation: OP_MINUS,
		Args: []*Expression{$1.ToExpression(), $3}})
}
| dstexpr TE expr ';'
{
	$$ = MakeAssignment($1, &Expression{
		Operation: OP_MUL,
		Args: []*Expression{$1.ToExpression(), $3}})
}
| dstexpr DE expr ';'
{
	$$ = MakeAssignment($1, &Expression{
		Operation: OP_DIV,
		Args: []*Expression{$1.ToExpression(), $3}})
}
| dstexpr RE expr ';'
{
	$$ = MakeAssignment($1, &Expression{
		Operation: OP_MOD,
		Args: []*Expression{$1.ToExpression(), $3}})
}

expr : NUMBER
{
	x, _ := strconv.Atoi($1.literal)
	$$ = &Expression{
		Operation: OP_NUM,
		Value: &Value{Type: VAL_NUM, NumVal: float64(x)}}
}
| '{' mapexprs '}'
{
	$$ = $2
}
| '[' listexprs ']'
{
	$$ = $2
}
| IDENT '(' listexprs ')'
{
	$$ = &Expression{
		Operation: OP_CALL,
		Args: []*Expression{
			MakeNameExpression($1.literal),
			$3}}
}
| STRING
{
	$$ = &Expression{
		Operation: OP_STRING,
		Value: &Value{Type: VAL_STRING, StringVal: $1.literal}}
}
| '-' expr      %prec UNARY
{
	$$ = &Expression{
		Operation: OP_UMINUS,
		Args: []*Expression{$2}}
}
| '(' expr ')'
{
	$$ = $2
}
| varexpr
{
	$$ = $1
}
| expr '+' expr
{
	$$ = &Expression{
		Operation: OP_PLUS,
		Args: []*Expression{$1, $3}}
}
| expr '-' expr
{
	$$ = &Expression{
		Operation: OP_MINUS,
		Args: []*Expression{$1, $3}}
}
| expr '*' expr
{
	$$ = &Expression{
		Operation: OP_MUL,
		Args: []*Expression{$1, $3}}
}
| expr '/' expr
{
	$$ = &Expression{
		Operation: OP_DIV,
		Args: []*Expression{$1, $3}}
}
| expr '&' expr
{
	$$ = &Expression{
		Operation: OP_BITWISEAND,
		Args: []*Expression{$1, $3}}
}
| expr '|' expr
{
	$$ = &Expression{
		Operation: OP_BITWISEOR,
		Args: []*Expression{$1, $3}}
}
| expr '^' expr
{
	$$ = &Expression{
		Operation: OP_BITWISEXOR,
		Args: []*Expression{$1, $3}}
}
| expr '%' expr
{
	$$ = &Expression{
		Operation: OP_MOD,
		Args: []*Expression{$1, $3}}
}
| expr EQ expr
{
	$$ = &Expression{
		Operation: OP_EQ,
		Args: []*Expression{$1, $3}}
}
| expr NE expr
{
	$$ = &Expression{
		Operation: OP_NE,
		Args: []*Expression{$1, $3}}
}
| expr GE expr
{
	$$ = &Expression{
		Operation: OP_GE,
		Args: []*Expression{$1, $3}}
}
| expr LE expr
{
	$$ = &Expression{
		Operation: OP_LE,
		Args: []*Expression{$1, $3}}
}
| expr AND expr
{
	$$ = &Expression{
		Operation: OP_AND,
		Args: []*Expression{$1, $3}}
}
| expr OR expr
{
	$$ = &Expression{
		Operation: OP_OR,
		Args: []*Expression{$1, $3}}
}
| '!' expr
{
	$$ = &Expression{
		Operation: OP_NOT,
		Args: []*Expression{$2}}
}
| expr '>' expr
{
	$$ = &Expression{
		Operation: OP_GT,
		Args: []*Expression{$1, $3}}
}
| expr '<' expr
{
	$$ = &Expression{
		Operation: OP_LT,
		Args: []*Expression{$1, $3}}
}
;
mapexprs : IDENT ':' expr
{
	$$ = &Expression{
		Operation: OP_MAP,
		Args: []*Expression{
			MakeNameExpression($1.literal),
			$3}}
}
| IDENT ':' expr ',' mapexprs
{
	args := []*Expression{
		MakeNameExpression($1.literal),
		$3}
	$$ = &Expression{
		Operation: OP_MAP,
		Args: append(args, $5.Args...)}
}
;
listexprs : expr
{
	$$ = &Expression{
		Operation: OP_LIST,
		Args: []*Expression{$1}}
}
| expr ',' listexprs
{
	args := []*Expression{$1}
	$$ = &Expression{
		Operation: OP_LIST,
		Args: append(args, $3.Args...)}
}
;
varexpr : expr '.' IDENT
{
	$$ = &Expression{
		Operation: OP_MAPVAR,
		Args: []*Expression{
			$1,
			MakeNameExpression($3.literal)}}
}
| expr '[' expr ']'
{
	$$ = &Expression{
		Operation: OP_LISTVAR,
		Args: []*Expression{$1, $3}}
}
| IDENT
{
	$$ = MakeNameExpression($1.literal)
}
;
dstexpr : IDENT
{
	$$ = &WriteableValue{
		Base: $1.literal,
		Path: []PathEntry{}}
}
| THIS
{
	$$ = &WriteableValue{
		Base: "this",
		Path: []PathEntry{}}
}
| dstexpr '[' expr ']'
{
	$1.Path = append($1.Path, PathEntry{Type:PATH_LIST,ListIndex:$3})
	$$ = $1
}
| dstexpr '.' IDENT
{
	$1.Path = append($1.Path, PathEntry{Type:PATH_MAP,MapKey:$3.literal})
	$$ = $1
}
;
%%
func Parse(exp string) []*Route {
	l := new(RouteLexer)
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	yyParse(l)
	return l.result
}
