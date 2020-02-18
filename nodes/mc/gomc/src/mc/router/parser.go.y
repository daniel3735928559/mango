%{
	package router
	import (
		"strconv"
		"fmt"
	)
	type Token struct {
		token   int
		literal string
		position Position
	}
	type RouteList struct {
		Routes []*Route
	}
	type Script struct {
		Statements []*Statement
	}
	type Statement struct {
		Operation string
		Args []*Value
	}
%}
%union{
	token Token
	routes []*Route
	transforms *Route
	transform *Transform
	expression *Expression
	script []*Statement
	statement *Statement
	node *Node
}
%type<routes> route
%type<node> node
%type<transforms> transforms
%type<transform> transform
%type<script> script
%type<statement> stmt
%type<expression> expr
%token<token> IDENT VAR NAME NUMBER STRING NOT AND OR IS EQ '?' '%' '=' '{' '}' '<' '>' ':' '+' '-' '*' '/'

%right UNARY
%left IS
%left GE LE
%left AND
%left OR
%right NOT
%%
route   : node '>' node
{
	fmt.Println("C")
	$$ = nil
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = []*Route{&Route{Source: $1.Name, Dest: $3.Name}}
	}
}
| node '<' node
{
	fmt.Println("B")
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
	fmt.Println("A")
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
	fmt.Println("NODE",$1.literal)
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
	$$ = []*Statement{$1}
}
| stmt script
{
	$$ = append([]*Statement{$1}, $2...)
}

stmt : expr ';'
{
	$$ = &Statement{
		Operation: "expr",
		Args: []*Value{
			&Value{Type:"expr", ExprVal: $1}}}
}
| IDENT '=' expr ';'
{
	$$ = &Statement{
		Operation: "assign",
		Args: []*Value{
			&Value{Type: "var", NameVal: $1.literal},
			&Value{Type: "expr", ExprVal: $3}}}
}
expr : NUMBER
{
	x, _ := strconv.Atoi($1.literal)
	$$ = &Expression{
		Operation: "val",
		Value: &Value{Type: "int", IntVal: x}}
}
| IDENT
{
	$$ = &Expression{
		Operation: "val",
		Value: &Value{Type: "var", NameVal: $1.literal}}
}
| STRING
{
	$$ = &Expression{
		Operation: "string",
		Value: &Value{Type: "var", StringVal: $1.literal}}
}
| '-' expr      %prec UNARY
{
	$$ = &Expression{
		Operation: "-",
		Args: []*Expression{$2}}
}
| '(' expr ')'
{
	$$ = $2
}
| expr '+' expr
{
	$$ = &Expression{
		Operation: "+",
		Args: []*Expression{$1, $3}}
}
| expr '-' expr
{
	$$ = &Expression{
		Operation: "-",
		Args: []*Expression{$1, $3}}
}
| expr '*' expr
{
	$$ = &Expression{
		Operation: "*",
		Args: []*Expression{$1, $3}}
}
| expr '/' expr
{
	$$ = &Expression{
		Operation: "/",
		Args: []*Expression{$1, $3}}
}
| expr '%' expr
{
	$$ = &Expression{
		Operation: "%",
		Args: []*Expression{$1, $3}}
}
| expr EQ expr
{
	$$ = &Expression{
		Operation: "==",
		Args: []*Expression{$1, $3}}
}
| expr GE expr
{
	$$ = &Expression{
		Operation: ">=",
		Args: []*Expression{$1, $3}}
}
| expr LE expr
{
	$$ = &Expression{
		Operation: "<=",
		Args: []*Expression{$1, $3}}
}
| expr '>' expr
{
	$$ = &Expression{
		Operation: ">",
		Args: []*Expression{$1, $3}}
}
| expr '<' expr
{
	$$ = &Expression{
		Operation: "<",
		Args: []*Expression{$1, $3}}
}
;
%%
func (v *Value) ToString() string {
	return v.Type
}
func Parse(exp string) []*Route {
	l := new(RouteLexer)
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	yyParse(l)
	return l.result
}
