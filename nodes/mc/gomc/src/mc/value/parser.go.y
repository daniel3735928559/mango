%{
	package value
	import (
		//"fmt"
		"strings"
		"strconv"
		"errors"
	)
	type Token struct {
		token   int
		literal string
		position Position
	}
%}
%union{
	token Token
	val *Value
}

%type<val> parseroot
%type<val> value
%type<val> listvals
%type<val> mapvals

%token<token> IDENT NUMBER STRING TRUE FALSE ',' '{' '}' '[' ']' '(' ')'

%left '[' '{' '('
%%
parseroot : value
{
	$$ = nil
	if l, ok := ValueParserlex.(*ValueLexer); ok {
		l.result = $1
	}
}
;
value : NUMBER
{
	x, _ := strconv.ParseFloat($1.literal, 64)
	$$ = MakeFloatValue(x)
}
| STRING
{
	$$ = &Value{Type: VAL_STRING, StringVal: $1.literal}
}
| TRUE
{
	$$ = MakeBoolValue(true)
}
| FALSE
{
	$$ = MakeBoolValue(false)
}
| '{' mapvals '}'
{
	$$ = $2
}
| '{' '}'
{
	$$ = &Value{
		Type: VAL_MAP,
		MapVal: map[string]*Value{}}
}
| '[' listvals ']'
{
	$$ = $2
}
| '[' ']'
{
	$$ = &Value{
		Type: VAL_LIST,
		ListVal: []*Value{}}
}
;
mapvals : IDENT ':' value
{
	$$ = &Value{
		Type: VAL_MAP,
		MapVal: map[string]*Value{$1.literal:$3}}
}
| IDENT ':' value ',' mapvals
{
	map_val := $5.MapVal
	map_val[$1.literal] = $3
	$$ = &Value{
		Type: VAL_MAP,
		MapVal: map_val}
}
;
listvals : value
{
	$$ = &Value{
		Type: VAL_LIST,
		ListVal: []*Value{$1}}
}
| value ',' listvals
{
	$$ = &Value{
		Type: VAL_LIST,
		ListVal: append([]*Value{$1}, $3.ListVal...)}
}
%%
func Parse(exp string) (*Value, error) {
	l := new(ValueLexer)
	lexerErrors := make([]string, 0)
	l.lexerErrors = &lexerErrors
	l.s = new(ValueScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	ValueParserParse(l)
	if len(lexerErrors) > 0 {
		return nil, errors.New(strings.Join(lexerErrors, "\n"))
	}
	return l.result, nil
}
