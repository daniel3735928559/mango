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
	type Script struct {
		Statements []*Statement
	}
	type Statement struct {
		Name string
		Val *Expression 
	}
%}
%union{
	token Token
	routes []*Route
	transforms *Route
	transform *Transform
	expression *Expression
	script []*Expression
	node *Node
}
%type<routes> route
%type<node> node
%type<transforms> transforms
%type<transform> transform
%type<script> script
%type<expression> varexpr
%type<expression> stmt
%type<expression> expr
%type<expression> mapexprs
%type<expression> listexprs
%token<token> IDENT VAR NAME NUMBER STRING AND OR EQ LE GE PE ME TE DE RE XE SUB '?' '%' '=' '{' '}' '[' ']' '<' '>' ':' '+' '-' '*' '/' '&' '|', '^'

%right UNARY
%left IS
%left GE LE
%left AND
%left OR
%right NOT
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
	// fmt.Println("NODE",$1.literal)
	$$ = &Node{Name: $1.literal}
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
		Type: "filter",
		Source: $3.ToString()}
}
| '%' '{' script '}'
{
	$$ = &Transform{
		Type: "map"}
}
| '=' '{' IDENT ':' expr '}'
{
	$$ = &Transform{
		Type: "replace",
		Source: $5.ToString()}
}
;
script
: stmt
{
	$$ = []*Expression{$1}
}
| stmt script
{
	$$ = append([]*Expression{$1}, $2...)
}

stmt : expr ';'
{
	$$ = $1
}
| varexpr '=' expr ';'
{
	$$ = &Expression{
		Operation: OP_ASSIGN,
		Args: []*Expression{$1,$3}}
}
| varexpr PE expr ';'
{
	$$ = &Expression{
		Operation: OP_ASSIGN,
		Args: []*Expression{
			$1,
			&Expression{
				Operation: OP_PLUS,
				Args: []*Expression{$1, $3}}}}
}
| varexpr ME expr ';'
{
	$$ = &Expression{
		Operation: OP_ASSIGN,
		Args: []*Expression{
			$1,
			&Expression{
				Operation: OP_MINUS,
				Args: []*Expression{$1, $3}}}}
}
| varexpr TE expr ';'
{
	$$ = &Expression{
		Operation: OP_ASSIGN,
		Args: []*Expression{
			$1,
			&Expression{
				Operation: OP_MUL,
				Args: []*Expression{$1, $3}}}}
}
| varexpr DE expr ';'
{
	$$ = &Expression{
		Operation: OP_ASSIGN,
		Args: []*Expression{
			$1,
			&Expression{
				Operation: OP_DIV,
				Args: []*Expression{$1, $3}}}}
}
| varexpr RE expr ';'
{
	$$ = &Expression{
		Operation: OP_ASSIGN,
		Args: []*Expression{
			$1,
			&Expression{
				Operation: OP_MOD,
				Args: []*Expression{$1, $3}}}}
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
	$$ = &Expression{
		Operation: OP_VAR,
		Value: &Value{Type: VAL_NAME, NameVal: $1.literal}}
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
