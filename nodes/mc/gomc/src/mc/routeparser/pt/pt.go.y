%{
	package routeparser
	import (
		"fmt"
		"text/scanner"
		"strings"
	)
	type Token struct {
		token   int
		literal string
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
%token<token> NAME NUMBER STRING NOT AND OR IS EQ '?' '%' '=' '{' '}' '<' '>' ':' '+' '-' '*' '/'

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
;
node    : NAME
{
	fmt.Println("NODE",$1.literal)
	$$ = &Node{Name: $1.literal}
}
;
%%
type RouteLexer struct {
	scanner.Scanner
	Nodes map[string]string
	result RouteList
}
func (l *RouteLexer) Lex(lval *yySymType) int {
	token := l.Scan()
	lit := l.TokenText()
	fmt.Println("lexint",lit)
	tok := int(token)
	switch tok {
	case scanner.Int:
		tok = NUMBER
	default:
		switch lit {
		case "IS":
			tok = IS
		case "NOT":
			tok = NOT
		case "AND":
			tok = AND
		case "OR":
			tok = OR
		case "<=":
			tok = LE
		case ">=":
			tok = GE
		default:
			fmt.Println("lit",lit)
			if v, ok := l.Nodes[lit]; ok {
				tok = NAME
				fmt.Println("Node",v)
				lit = v
			}
		}
	}
	lval.token = Token{token: tok, literal: lit}
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
	l.Init(strings.NewReader(exp))
	yyParse(l)
	return l.result
}
