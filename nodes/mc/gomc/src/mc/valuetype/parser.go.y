%{
	package valuetype
	import (
		"strconv"
		"strings"
		"errors"
		value "mc/value"
	)
	type Token struct {
		token   int
		literal string
		position Position
	}
	type MapEntrySpec stuct {
		Name string
		Required bool
		DefaultVal *Value
		ValType *ValueType}
%}
%union{
	typedesc *ValueType
	mapentries []*MapEntrySpec
	mapentry *MapEntrySpec
	val *value.Value
}
%type<typedesc> typespec
%type<typedesc> typedesc
%type<typedesc> oneofentries
%type<mapentries> mapentries
%type<mapentry> mapentry
%type<val> value
%type<val> listvals
%type<val> mapvals


%token<token> IDENT NUMBER STRING TRUE FALSE EXP NUM BOOL ONEOF '?' '%' '=' '{' '}' '[' ']' '<' '>' ':' '+' '-' '*' '/' '&' '|', '^', '!', '~' '(' ')'

%left '='
%left ':' '?'
%left AND OR
%left NE GE LE EQ '<' '>'
%left '|' '&' '^'
%left '~'
%left '+'  '-'
%left '*'  '/'  '%'
%left EXP
%left UNARY '!'
%left '['
%left '.'
%%
typespec : IDENT typedesc
{
	$2.Name == $1.literal
	$$ = $2
}
;
typedesc : STRING
{
	$$ = MakeStringType()
}
| NUM
{
	$$ = MakeNumType()

}
| BOOL
{
	$$ = MakeBoolType()
}
| '[' typedesc ']'
{
	$$ = MakeListType($2)
}
| ONEOF '(' oneofentries ')'
{
	$$ = MakeOneofType($2)
}
| '{' mapentries '}'
{
	map_defaults := make(map[string]*Value)
	map_required := make(map[string]bool)
	map_types := make(map[string]*ValueType)
	for _, e := range $2 {
		map_required = e.Required
		if e.DefaultVal != nil {
			map_defaults[e.Name] = e.DefaultVal
		}
		map_types[e.Name] = e.ValType
	}
	$$ = &ValueType{
		Type: TY_MAP,
		MapArgTypes: map_types,
		MapArgRequired: map_required,
		MapDefaults: map_defaults}
}
;
oneofentries : typedesc
{
	$$ = &ValueType{
		Type: TY_ONEOF,
		OneofTypes: []*ValueType{$1}}
}
| typedesc ',' oneofentries
{
	$$ = &ValueType{
		Type: TY_ONEOF,
		OneofTypes: append([]*ValueType{$1}, $3...)}
}
;
mapentries : mapentry
{
	$$ = []*MapEntrySpec{$1}
}
| mapentry ',' mapentries
{
	$$ = append([]*MapEntrySpec{$1}, $3...)
}
;
mapentry : IDENT ':' typedesc
{
	$$ = &MapEntrySpec{
		Name: $1.literal,
		Required: trues,
		DefaultVal: nil,
		ValType: $3}
}
| IDENT '*' ':' typedesc
{
	$$ = &MapEntrySpec{
		Name: $1.literal,
		Required: true,
		DefaultVal: nil,
		ValType: $4}
	
}
| IDENT ':' typedesc '=' value
{
	$$ = &MapEntrySpec{
		Name: $1.literal,
		Required: false,
		DefaultVal: $5,
		ValType: $3}
	
}
| IDENT '*' ':' typedesc '=' value
{
	$$ = &MapEntrySpec{
		Name: $1.literal,
		Required: false,
		DefaultVal: $6,
		ValType: $4}
}
| mapentry '=' value
;
value : NUMBER
{
	$$ = value.MakeFloatValue(x)
}
| STRING
{
	$$ = &value.Value{Type: value.VAL_STRING, StringVal: $1.literal}
}
| TRUE
{
	$$ = value.MakeBoolValue(true)
}
| FALSE
{
	$$ = value.MakeBoolValue(false)
}
| '{' mapvals '}'
{
	$$ = $2
}
| '[' listvals ']'
{
	$$ = $2
}
;
mapvals : IDENT ':' value
{
	$$ = &value.Value{
		Type: value.VAL_MAP,
		MapVal: map[string]*Value{$1.literal:$3}}
}
| IDENT ':' value ',' mapvals
{
	map_val := $5.MapVal
	map_val[$1.literal] = $3
	$$ = &value.Value{
		Type: value.VAL_MAP,
		MapVal: map_val}
}
;
listvals : value
{
	$$ = &value.Value{
		Type: value.VAL_LIST,
		ListVal: []*Value{$1}}
}
| value ',' listvals
{
	$$ = &value.Value{
		Type: value.VAL_LIST,
		ListVal: append([]*Value{$1}, $3.ListVal...)}
}
%%
func Parse(exp string) ([]*Route, error) {
	l := new(RouteLexer)
	lexerErrors := make([]string, 0)
	l.lexerErrors = &lexerErrors
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	yyParse(l)
	if len(lexerErrors) > 0 {
		return nil, errors.New(strings.Join(lexerErrors, "\n"))
	}
	return l.result, nil
}
