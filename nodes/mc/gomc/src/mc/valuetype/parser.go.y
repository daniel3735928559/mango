%{
	package valuetype
	import (
		"fmt"
		"strings"
		"strconv"
		"errors"
		value "mc/value"
	)
	type Token struct {
		token   int
		literal string
		position Position
	}
	type MapEntrySpec struct {
		Name string
		Required bool
		DefaultVal *value.Value
		ValType *ValueType
	}
%}
%union{
	token Token
	typedesc *ValueType
	mapentries []*MapEntrySpec
	mapentry *MapEntrySpec
	val *value.Value
}
%type<typedesc> parseroot
%type<typedesc> typedesc
%type<typedesc> oneofentries
%type<mapentries> mapentries
%type<mapentry> mapentry
%type<val> value
%type<val> listvals
%type<val> mapvals

%token<token> IDENT NUMBER STRING TRUE FALSE NUM STR BOOL ONEOF '*' ',' '=' '{' '}' '[' ']' '(' ')'

%left '='
%left '*'
%left '[' '{' '('
%%
parseroot : typedesc
{
	$$ = nil
	if l, ok := ValueTypeParserlex.(*ValueTypeLexer); ok {
		l.result = $1
	}
}
;
typedesc : STR
{
	$$ = MakeStringType()
}
| NUM
{
	$$ = MakeNumType()
}
| IDENT
{
	$$ = MakeExtType($1.literal)
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
	$$ = $3
}
| '{' mapentries '}'
{
	map_defaults := make(map[string]*value.Value)
	map_required := make(map[string]bool)
	map_types := make(map[string]*ValueType)
	for _, e := range $2 {
		map_required[e.Name] = e.Required
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
		OneofTypes: append([]*ValueType{$1}, $3.OneofTypes...)}
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
	fmt.Println("sty",$1.literal,$3)
	$$ = &MapEntrySpec{
		Name: $1.literal,
		Required: true,
		DefaultVal: nil,
		ValType: $3}
}
| IDENT '*' ':' typedesc
{
	$$ = &MapEntrySpec{
		Name: $1.literal,
		Required: false,
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
;
value : NUMBER
{
	x, _ := strconv.ParseFloat($1.literal, 64)
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
| '{' '}'
{
	$$ = &value.Value{
		Type: value.VAL_MAP,
		MapVal: map[string]*value.Value{}}
}
| '[' listvals ']'
{
	$$ = $2
}
| '[' ']'
{
	$$ = &value.Value{
		Type: value.VAL_LIST,
		ListVal: []*value.Value{}}
}
;
mapvals : IDENT ':' value
{
	$$ = &value.Value{
		Type: value.VAL_MAP,
		MapVal: map[string]*value.Value{$1.literal:$3}}
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
		ListVal: []*value.Value{$1}}
}
| value ',' listvals
{
	$$ = &value.Value{
		Type: value.VAL_LIST,
		ListVal: append([]*value.Value{$1}, $3.ListVal...)}
}
%%
func Parse(exp string) (*ValueType, error) {
	l := new(ValueTypeLexer)
	lexerErrors := make([]string, 0)
	l.lexerErrors = &lexerErrors
	l.s = new(ValueTypeScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	ValueTypeParserParse(l)
	if len(lexerErrors) > 0 {
		return nil, errors.New(strings.Join(lexerErrors, "\n"))
	}
	return l.result, nil
}
