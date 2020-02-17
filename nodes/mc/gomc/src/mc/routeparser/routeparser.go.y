%{
	package routeparser
	import (
		"strconv"
		"fmt"
	)
	type Token struct {
		token   int
		literal string
		position Position
	}
	type Node struct {
		Name string
	}
	type Route struct {
		Source string
		Transforms []*Transform
		Dest string
	}
	type RouteList struct {
		Routes []*Route
	}
	type Transform struct {
		Type string
		Condition *Expression
		Script *Script
		Source string
	}
	type Expression struct {
		Operation string
		Args []*Expression
		Value *Value
	}
	type Script struct {
		Statements []*Statement
	}
	type Statement struct {
		Operation string
		Args []*Value
	}
	type Value struct {
		Type string
		NameVal string
		IntVal int
		FloatVal float64
		StringVal string
		ExprVal *Expression
	}

	// type Transform struct {
	// 	Type FILTER|EDIT|REPLACE|FILTEREDIT|FILTERREPLACE
	// }
	// type Operation struct {
	// 	Op string
	// 	Args []Operation
	// 	Kwargs map[string]Operation
	// 	Execute Script
	// }
	// type Script struct {
	// 	Ops []Op
	// }
%}
%union{
	token Token
	routes []*Route
	transforms []*Transform
	transform *Transform
	expression *Expression
	script []*Statement
	statement *Statement
	node *Node
}
%type<routes> route
%type<node> node
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
		l.result = RouteList{Routes:[]*Route{&Route{Source: $1.Name, Dest: $3.Name}}}
	}
}
| node '<' node
{
	fmt.Println("B")
	$$ = nil
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = RouteList{Routes:[]*Route{&Route{Source: $3.Name, Dest: $1.Name}}}
	}
}
| node '<' '>' node
{
	$$ = nil
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = RouteList{Routes: []*Route{
			&Route{Source: $1.Name, Dest: $4.Name},
			&Route{Source: $4.Name, Dest: $1.Name}}}
	}
}
| node '>' transform '>' node
{
	$$ = nil
	fmt.Println("A")
	if l, ok := yylex.(*RouteLexer); ok {
		l.result = RouteList{Routes:[]*Route{&Route{Source: $1.Name, Dest: $5.Name, Transforms: []*Transform{$3}}}}
	}
}
;
node    : IDENT
{
	fmt.Println("NODE",$1.literal)
	$$ = &Node{Name: $1.literal}
}
;
transform : '?' '{' expr '}'
{
	$$ = &Transform{Source: $1.literal}
}
| '%' '{' script '}'
{
	$$ = &Transform{Source: $1.literal}
}
| '=' '{' IDENT ':' expr '}'
{
	$$ = &Transform{Source: $1.literal}
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
type RouteLexer struct {
	s *RouteScanner
	Nodes map[string]string
	result RouteList
}
func (l *RouteLexer) Lex(lval *yySymType) int {
	tok, lit, pos := l.s.Scan()
	if tok == EOF {
		return 0
	}
	lval.token = Token{token: tok, literal: lit, position: pos}
	fmt.Println("Lexed",tok,lit,pos)
	return tok
}
func (l *RouteLexer) Error(e string) {
	fmt.Println("ERROR",e)
}
func (r *Route) ToString() string {
	return fmt.Sprintf("%s > %s", r.Source, r.Dest)
}
func Parse(exp string, nodes map[string]string) RouteList {
	l := new(RouteLexer)
	l.Nodes = nodes
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	yyParse(l)
	return l.result
}
