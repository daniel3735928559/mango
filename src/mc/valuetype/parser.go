// Code generated by goyacc -o src/mc/valuetype/parser.go -p ValueTypeParser -v src/mc/valuetype/parser.output src/mc/valuetype/parser.go.y. DO NOT EDIT.

//line src/mc/valuetype/parser.go.y:2
package valuetype

import __yyfmt__ "fmt"

//line src/mc/valuetype/parser.go.y:2
import (
	//"fmt"
	"errors"
	value "mc/value"
	"strconv"
	"strings"
)

type Token struct {
	token    int
	literal  string
	position Position
}
type MapEntrySpec struct {
	Name       string
	Required   bool
	DefaultVal *value.Value
	ValType    *ValueType
}

//line src/mc/valuetype/parser.go.y:22
type ValueTypeParserSymType struct {
	yys        int
	token      Token
	typedesc   *ValueType
	mapentries []*MapEntrySpec
	mapentry   *MapEntrySpec
	val        *value.Value
}

const IDENT = 57346
const NUMBER = 57347
const STRING = 57348
const TRUE = 57349
const FALSE = 57350
const NUM = 57351
const STR = 57352
const BOOL = 57353
const ONEOF = 57354
const ANY = 57355

var ValueTypeParserToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"NUMBER",
	"STRING",
	"TRUE",
	"FALSE",
	"NUM",
	"STR",
	"BOOL",
	"ONEOF",
	"ANY",
	"'*'",
	"','",
	"'='",
	"'{'",
	"'}'",
	"'['",
	"']'",
	"'('",
	"')'",
	"':'",
}
var ValueTypeParserStatenames = [...]string{}

const ValueTypeParserEofCode = 1
const ValueTypeParserErrCode = 2
const ValueTypeParserInitialStackSize = 16

//line src/mc/valuetype/parser.go.y:224

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

//line yacctab:1
var ValueTypeParserExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const ValueTypeParserPrivate = 57344

const ValueTypeParserLast = 60

var ValueTypeParserAct = [...]int{

	39, 42, 44, 19, 2, 33, 34, 35, 36, 18,
	46, 28, 11, 24, 13, 12, 5, 37, 23, 38,
	43, 4, 3, 6, 9, 7, 27, 22, 47, 10,
	17, 8, 31, 32, 41, 29, 26, 45, 20, 33,
	34, 35, 36, 16, 30, 15, 51, 48, 40, 49,
	50, 37, 52, 38, 41, 25, 21, 14, 16, 1,
}
var ValueTypeParserPact = [...]int{

	12, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 12, -6,
	39, 10, 12, 20, -1000, 41, 4, -1000, -9, 40,
	-1000, 54, 12, -12, -1000, 12, -1000, 28, 12, -1000,
	34, -1000, -1000, -1000, -1000, -1000, -1000, 30, 0, 19,
	-1000, -13, 8, -1000, 32, -1000, 34, -1000, 34, 31,
	-1000, 50, -1000,
}
var ValueTypeParserPgo = [...]int{

	0, 59, 3, 9, 14, 45, 2, 1, 0,
}
var ValueTypeParserR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 4, 4, 5, 5, 5, 6, 6,
	6, 6, 6, 6, 6, 6, 8, 8, 7, 7,
}
var ValueTypeParserR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 3, 4, 3,
	2, 1, 3, 1, 3, 3, 4, 5, 1, 1,
	1, 1, 3, 2, 3, 2, 3, 5, 1, 3,
}
var ValueTypeParserChk = [...]int{

	-1000, -1, -2, 10, 9, 4, 11, 13, 19, 12,
	17, -2, 21, -4, 18, -5, 4, 20, -3, -2,
	18, 15, 23, 14, 22, 15, -4, -2, 23, -3,
	16, -2, -6, 5, 6, 7, 8, 17, 19, -8,
	18, 4, -7, 20, -6, 18, 23, 20, 15, -6,
	-7, 15, -8,
}
var ValueTypeParserDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 6, 0, 0,
	0, 0, 0, 0, 10, 13, 0, 7, 0, 11,
	9, 0, 0, 0, 8, 0, 14, 15, 0, 12,
	0, 16, 17, 18, 19, 20, 21, 0, 0, 0,
	23, 0, 0, 25, 28, 22, 0, 24, 0, 26,
	29, 0, 27,
}
var ValueTypeParserTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	21, 22, 14, 3, 15, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 23, 3,
	3, 16, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 19, 3, 20, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 17, 3, 18,
}
var ValueTypeParserTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13,
}
var ValueTypeParserTok3 = [...]int{
	0,
}

var ValueTypeParserErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	ValueTypeParserDebug        = 0
	ValueTypeParserErrorVerbose = false
)

type ValueTypeParserLexer interface {
	Lex(lval *ValueTypeParserSymType) int
	Error(s string)
}

type ValueTypeParserParser interface {
	Parse(ValueTypeParserLexer) int
	Lookahead() int
}

type ValueTypeParserParserImpl struct {
	lval  ValueTypeParserSymType
	stack [ValueTypeParserInitialStackSize]ValueTypeParserSymType
	char  int
}

func (p *ValueTypeParserParserImpl) Lookahead() int {
	return p.char
}

func ValueTypeParserNewParser() ValueTypeParserParser {
	return &ValueTypeParserParserImpl{}
}

const ValueTypeParserFlag = -1000

func ValueTypeParserTokname(c int) string {
	if c >= 1 && c-1 < len(ValueTypeParserToknames) {
		if ValueTypeParserToknames[c-1] != "" {
			return ValueTypeParserToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func ValueTypeParserStatname(s int) string {
	if s >= 0 && s < len(ValueTypeParserStatenames) {
		if ValueTypeParserStatenames[s] != "" {
			return ValueTypeParserStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func ValueTypeParserErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !ValueTypeParserErrorVerbose {
		return "syntax error"
	}

	for _, e := range ValueTypeParserErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + ValueTypeParserTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := ValueTypeParserPact[state]
	for tok := TOKSTART; tok-1 < len(ValueTypeParserToknames); tok++ {
		if n := base + tok; n >= 0 && n < ValueTypeParserLast && ValueTypeParserChk[ValueTypeParserAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if ValueTypeParserDef[state] == -2 {
		i := 0
		for ValueTypeParserExca[i] != -1 || ValueTypeParserExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; ValueTypeParserExca[i] >= 0; i += 2 {
			tok := ValueTypeParserExca[i]
			if tok < TOKSTART || ValueTypeParserExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if ValueTypeParserExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += ValueTypeParserTokname(tok)
	}
	return res
}

func ValueTypeParserlex1(lex ValueTypeParserLexer, lval *ValueTypeParserSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = ValueTypeParserTok1[0]
		goto out
	}
	if char < len(ValueTypeParserTok1) {
		token = ValueTypeParserTok1[char]
		goto out
	}
	if char >= ValueTypeParserPrivate {
		if char < ValueTypeParserPrivate+len(ValueTypeParserTok2) {
			token = ValueTypeParserTok2[char-ValueTypeParserPrivate]
			goto out
		}
	}
	for i := 0; i < len(ValueTypeParserTok3); i += 2 {
		token = ValueTypeParserTok3[i+0]
		if token == char {
			token = ValueTypeParserTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = ValueTypeParserTok2[1] /* unknown char */
	}
	if ValueTypeParserDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", ValueTypeParserTokname(token), uint(char))
	}
	return char, token
}

func ValueTypeParserParse(ValueTypeParserlex ValueTypeParserLexer) int {
	return ValueTypeParserNewParser().Parse(ValueTypeParserlex)
}

func (ValueTypeParserrcvr *ValueTypeParserParserImpl) Parse(ValueTypeParserlex ValueTypeParserLexer) int {
	var ValueTypeParsern int
	var ValueTypeParserVAL ValueTypeParserSymType
	var ValueTypeParserDollar []ValueTypeParserSymType
	_ = ValueTypeParserDollar // silence set and not used
	ValueTypeParserS := ValueTypeParserrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	ValueTypeParserstate := 0
	ValueTypeParserrcvr.char = -1
	ValueTypeParsertoken := -1 // ValueTypeParserrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		ValueTypeParserstate = -1
		ValueTypeParserrcvr.char = -1
		ValueTypeParsertoken = -1
	}()
	ValueTypeParserp := -1
	goto ValueTypeParserstack

ret0:
	return 0

ret1:
	return 1

ValueTypeParserstack:
	/* put a state and value onto the stack */
	if ValueTypeParserDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", ValueTypeParserTokname(ValueTypeParsertoken), ValueTypeParserStatname(ValueTypeParserstate))
	}

	ValueTypeParserp++
	if ValueTypeParserp >= len(ValueTypeParserS) {
		nyys := make([]ValueTypeParserSymType, len(ValueTypeParserS)*2)
		copy(nyys, ValueTypeParserS)
		ValueTypeParserS = nyys
	}
	ValueTypeParserS[ValueTypeParserp] = ValueTypeParserVAL
	ValueTypeParserS[ValueTypeParserp].yys = ValueTypeParserstate

ValueTypeParsernewstate:
	ValueTypeParsern = ValueTypeParserPact[ValueTypeParserstate]
	if ValueTypeParsern <= ValueTypeParserFlag {
		goto ValueTypeParserdefault /* simple state */
	}
	if ValueTypeParserrcvr.char < 0 {
		ValueTypeParserrcvr.char, ValueTypeParsertoken = ValueTypeParserlex1(ValueTypeParserlex, &ValueTypeParserrcvr.lval)
	}
	ValueTypeParsern += ValueTypeParsertoken
	if ValueTypeParsern < 0 || ValueTypeParsern >= ValueTypeParserLast {
		goto ValueTypeParserdefault
	}
	ValueTypeParsern = ValueTypeParserAct[ValueTypeParsern]
	if ValueTypeParserChk[ValueTypeParsern] == ValueTypeParsertoken { /* valid shift */
		ValueTypeParserrcvr.char = -1
		ValueTypeParsertoken = -1
		ValueTypeParserVAL = ValueTypeParserrcvr.lval
		ValueTypeParserstate = ValueTypeParsern
		if Errflag > 0 {
			Errflag--
		}
		goto ValueTypeParserstack
	}

ValueTypeParserdefault:
	/* default state action */
	ValueTypeParsern = ValueTypeParserDef[ValueTypeParserstate]
	if ValueTypeParsern == -2 {
		if ValueTypeParserrcvr.char < 0 {
			ValueTypeParserrcvr.char, ValueTypeParsertoken = ValueTypeParserlex1(ValueTypeParserlex, &ValueTypeParserrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if ValueTypeParserExca[xi+0] == -1 && ValueTypeParserExca[xi+1] == ValueTypeParserstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			ValueTypeParsern = ValueTypeParserExca[xi+0]
			if ValueTypeParsern < 0 || ValueTypeParsern == ValueTypeParsertoken {
				break
			}
		}
		ValueTypeParsern = ValueTypeParserExca[xi+1]
		if ValueTypeParsern < 0 {
			goto ret0
		}
	}
	if ValueTypeParsern == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			ValueTypeParserlex.Error(ValueTypeParserErrorMessage(ValueTypeParserstate, ValueTypeParsertoken))
			Nerrs++
			if ValueTypeParserDebug >= 1 {
				__yyfmt__.Printf("%s", ValueTypeParserStatname(ValueTypeParserstate))
				__yyfmt__.Printf(" saw %s\n", ValueTypeParserTokname(ValueTypeParsertoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for ValueTypeParserp >= 0 {
				ValueTypeParsern = ValueTypeParserPact[ValueTypeParserS[ValueTypeParserp].yys] + ValueTypeParserErrCode
				if ValueTypeParsern >= 0 && ValueTypeParsern < ValueTypeParserLast {
					ValueTypeParserstate = ValueTypeParserAct[ValueTypeParsern] /* simulate a shift of "error" */
					if ValueTypeParserChk[ValueTypeParserstate] == ValueTypeParserErrCode {
						goto ValueTypeParserstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if ValueTypeParserDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", ValueTypeParserS[ValueTypeParserp].yys)
				}
				ValueTypeParserp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if ValueTypeParserDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", ValueTypeParserTokname(ValueTypeParsertoken))
			}
			if ValueTypeParsertoken == ValueTypeParserEofCode {
				goto ret1
			}
			ValueTypeParserrcvr.char = -1
			ValueTypeParsertoken = -1
			goto ValueTypeParsernewstate /* try again in the same state */
		}
	}

	/* reduction by production ValueTypeParsern */
	if ValueTypeParserDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", ValueTypeParsern, ValueTypeParserStatname(ValueTypeParserstate))
	}

	ValueTypeParsernt := ValueTypeParsern
	ValueTypeParserpt := ValueTypeParserp
	_ = ValueTypeParserpt // guard against "declared and not used"

	ValueTypeParserp -= ValueTypeParserR2[ValueTypeParsern]
	// ValueTypeParserp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if ValueTypeParserp+1 >= len(ValueTypeParserS) {
		nyys := make([]ValueTypeParserSymType, len(ValueTypeParserS)*2)
		copy(nyys, ValueTypeParserS)
		ValueTypeParserS = nyys
	}
	ValueTypeParserVAL = ValueTypeParserS[ValueTypeParserp+1]

	/* consult goto table to find next state */
	ValueTypeParsern = ValueTypeParserR1[ValueTypeParsern]
	ValueTypeParserg := ValueTypeParserPgo[ValueTypeParsern]
	ValueTypeParserj := ValueTypeParserg + ValueTypeParserS[ValueTypeParserp].yys + 1

	if ValueTypeParserj >= ValueTypeParserLast {
		ValueTypeParserstate = ValueTypeParserAct[ValueTypeParserg]
	} else {
		ValueTypeParserstate = ValueTypeParserAct[ValueTypeParserj]
		if ValueTypeParserChk[ValueTypeParserstate] != -ValueTypeParsern {
			ValueTypeParserstate = ValueTypeParserAct[ValueTypeParserg]
		}
	}
	// dummy call; replaced with literal code
	switch ValueTypeParsernt {

	case 1:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:45
		{
			ValueTypeParserVAL.typedesc = nil
			if l, ok := ValueTypeParserlex.(*ValueTypeLexer); ok {
				l.result = ValueTypeParserDollar[1].typedesc
			}
		}
	case 2:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:53
		{
			ValueTypeParserVAL.typedesc = MakeStringType()
		}
	case 3:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:57
		{
			ValueTypeParserVAL.typedesc = MakeNumType()
		}
	case 4:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:61
		{
			ValueTypeParserVAL.typedesc = MakeExtType(ValueTypeParserDollar[1].token.literal)
		}
	case 5:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:65
		{
			ValueTypeParserVAL.typedesc = MakeBoolType()
		}
	case 6:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:69
		{
			ValueTypeParserVAL.typedesc = MakeAnyType()
		}
	case 7:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:73
		{
			ValueTypeParserVAL.typedesc = MakeListType(ValueTypeParserDollar[2].typedesc)
		}
	case 8:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-4 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:77
		{
			ValueTypeParserVAL.typedesc = ValueTypeParserDollar[3].typedesc
		}
	case 9:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:81
		{
			map_defaults := make(map[string]*value.Value)
			map_required := make(map[string]bool)
			map_types := make(map[string]*ValueType)
			for _, e := range ValueTypeParserDollar[2].mapentries {
				map_required[e.Name] = e.Required
				if e.DefaultVal != nil {
					map_defaults[e.Name] = e.DefaultVal
				}
				map_types[e.Name] = e.ValType
			}
			ValueTypeParserVAL.typedesc = &ValueType{
				Type:           TY_MAP,
				MapArgTypes:    map_types,
				MapArgRequired: map_required,
				MapDefaults:    map_defaults}
		}
	case 10:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-2 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:99
		{
			map_defaults := make(map[string]*value.Value)
			map_required := make(map[string]bool)
			map_types := make(map[string]*ValueType)
			ValueTypeParserVAL.typedesc = &ValueType{
				Type:           TY_MAP,
				MapArgTypes:    map_types,
				MapArgRequired: map_required,
				MapDefaults:    map_defaults}
		}
	case 11:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:111
		{
			ValueTypeParserVAL.typedesc = &ValueType{
				Type:       TY_ONEOF,
				OneofTypes: []*ValueType{ValueTypeParserDollar[1].typedesc}}
		}
	case 12:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:117
		{
			ValueTypeParserVAL.typedesc = &ValueType{
				Type:       TY_ONEOF,
				OneofTypes: append([]*ValueType{ValueTypeParserDollar[1].typedesc}, ValueTypeParserDollar[3].typedesc.OneofTypes...)}
		}
	case 13:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:124
		{
			ValueTypeParserVAL.mapentries = []*MapEntrySpec{ValueTypeParserDollar[1].mapentry}
		}
	case 14:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:128
		{
			ValueTypeParserVAL.mapentries = append([]*MapEntrySpec{ValueTypeParserDollar[1].mapentry}, ValueTypeParserDollar[3].mapentries...)
		}
	case 15:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:133
		{
			//fmt.Println("sty",$1.literal,$3)
			ValueTypeParserVAL.mapentry = &MapEntrySpec{
				Name:       ValueTypeParserDollar[1].token.literal,
				Required:   true,
				DefaultVal: nil,
				ValType:    ValueTypeParserDollar[3].typedesc}
		}
	case 16:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-4 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:142
		{
			ValueTypeParserVAL.mapentry = &MapEntrySpec{
				Name:       ValueTypeParserDollar[1].token.literal,
				Required:   false,
				DefaultVal: nil,
				ValType:    ValueTypeParserDollar[4].typedesc}

		}
	case 17:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-5 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:151
		{
			ValueTypeParserVAL.mapentry = &MapEntrySpec{
				Name:       ValueTypeParserDollar[1].token.literal,
				Required:   false,
				DefaultVal: ValueTypeParserDollar[5].val,
				ValType:    ValueTypeParserDollar[3].typedesc}
		}
	case 18:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:160
		{
			x, _ := strconv.ParseFloat(ValueTypeParserDollar[1].token.literal, 64)
			ValueTypeParserVAL.val = value.MakeFloatValue(x)
		}
	case 19:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:165
		{
			ValueTypeParserVAL.val = &value.Value{Type: value.VAL_STRING, StringVal: ValueTypeParserDollar[1].token.literal}
		}
	case 20:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:169
		{
			ValueTypeParserVAL.val = value.MakeBoolValue(true)
		}
	case 21:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:173
		{
			ValueTypeParserVAL.val = value.MakeBoolValue(false)
		}
	case 22:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:177
		{
			ValueTypeParserVAL.val = ValueTypeParserDollar[2].val
		}
	case 23:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-2 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:181
		{
			ValueTypeParserVAL.val = &value.Value{
				Type:   value.VAL_MAP,
				MapVal: map[string]*value.Value{}}
		}
	case 24:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:187
		{
			ValueTypeParserVAL.val = ValueTypeParserDollar[2].val
		}
	case 25:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-2 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:191
		{
			ValueTypeParserVAL.val = &value.Value{
				Type:    value.VAL_LIST,
				ListVal: []*value.Value{}}
		}
	case 26:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:198
		{
			ValueTypeParserVAL.val = &value.Value{
				Type:   value.VAL_MAP,
				MapVal: map[string]*value.Value{ValueTypeParserDollar[1].token.literal: ValueTypeParserDollar[3].val}}
		}
	case 27:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-5 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:204
		{
			map_val := ValueTypeParserDollar[5].val.MapVal
			map_val[ValueTypeParserDollar[1].token.literal] = ValueTypeParserDollar[3].val
			ValueTypeParserVAL.val = &value.Value{
				Type:   value.VAL_MAP,
				MapVal: map_val}
		}
	case 28:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:213
		{
			ValueTypeParserVAL.val = &value.Value{
				Type:    value.VAL_LIST,
				ListVal: []*value.Value{ValueTypeParserDollar[1].val}}
		}
	case 29:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]
//line src/mc/valuetype/parser.go.y:219
		{
			ValueTypeParserVAL.val = &value.Value{
				Type:    value.VAL_LIST,
				ListVal: append([]*value.Value{ValueTypeParserDollar[1].val}, ValueTypeParserDollar[3].val.ListVal...)}
		}
	}
	goto ValueTypeParserstack /* stack new state and value */
}
