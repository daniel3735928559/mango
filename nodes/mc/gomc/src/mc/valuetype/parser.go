


package valuetype

import __yyfmt__ "fmt"


import (
	"errors"
	"fmt"
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


var ValueTypeParserExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const ValueTypeParserPrivate = 57344

const ValueTypeParserLast = 51

var ValueTypeParserAct = [...]int{

	36, 38, 39, 16, 2, 30, 31, 32, 33, 15,
	9, 20, 4, 3, 5, 7, 34, 41, 35, 8,
	19, 6, 21, 24, 25, 10, 42, 14, 40, 28,
	29, 11, 26, 17, 27, 46, 43, 22, 18, 37,
	13, 12, 1, 0, 44, 45, 0, 47, 0, 0,
	23,
}
var ValueTypeParserPact = [...]int{

	3, -1000, -1000, -1000, -1000, -1000, 3, 5, 36, 8,
	3, 16, 24, -2, -1000, 1, 23, -1000, 36, 3,
	2, -1000, 3, -1000, 19, 3, -1000, 0, -1000, -1000,
	-1000, -1000, -1000, -1000, 35, 0, 11, -5, 7, 22,
	-1000, 0, -1000, 0, 21, -1000, 35, -1000,
}
var ValueTypeParserPgo = [...]int{

	0, 42, 3, 9, 31, 41, 2, 1, 0,
}
var ValueTypeParserR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 3, 3,
	4, 4, 5, 5, 5, 6, 6, 6, 6, 6,
	6, 8, 8, 7, 7,
}
var ValueTypeParserR2 = [...]int{

	0, 1, 1, 1, 1, 3, 4, 3, 1, 3,
	1, 3, 3, 4, 5, 1, 1, 1, 1, 3,
	3, 3, 5, 1, 3,
}
var ValueTypeParserChk = [...]int{

	-1000, -1, -2, 10, 9, 11, 18, 12, 16, -2,
	20, -4, -5, 4, 19, -3, -2, 17, 14, 22,
	13, 21, 14, -4, -2, 22, -3, 15, -2, -6,
	5, 6, 7, 8, 16, 18, -8, 4, -7, -6,
	17, 22, 19, 14, -6, -7, 14, -8,
}
var ValueTypeParserDef = [...]int{

	0, -2, 1, 2, 3, 4, 0, 0, 0, 0,
	0, 0, 10, 0, 5, 0, 8, 7, 0, 0,
	0, 6, 0, 11, 12, 0, 9, 0, 13, 14,
	15, 16, 17, 18, 0, 0, 0, 0, 0, 23,
	19, 0, 20, 0, 21, 24, 0, 22,
}
var ValueTypeParserTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	20, 21, 13, 3, 14, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 22, 3,
	3, 15, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 18, 3, 19, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 16, 3, 17,
}
var ValueTypeParserTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12,
}
var ValueTypeParserTok3 = [...]int{
	0,
}

var ValueTypeParserErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}



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

		{
			ValueTypeParserVAL.typedesc = nil
			if l, ok := ValueTypeParserlex.(*ValueTypeLexer); ok {
				l.result = ValueTypeParserDollar[1].typedesc
			}
		}
	case 2:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = MakeStringType()
		}
	case 3:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = MakeNumType()
		}
	case 4:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = MakeBoolType()
		}
	case 5:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = MakeListType(ValueTypeParserDollar[2].typedesc)
		}
	case 6:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-4 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = ValueTypeParserDollar[3].typedesc
		}
	case 7:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

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
	case 8:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = &ValueType{
				Type:       TY_ONEOF,
				OneofTypes: []*ValueType{ValueTypeParserDollar[1].typedesc}}
		}
	case 9:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.typedesc = &ValueType{
				Type:       TY_ONEOF,
				OneofTypes: append([]*ValueType{ValueTypeParserDollar[1].typedesc}, ValueTypeParserDollar[3].typedesc.OneofTypes...)}
		}
	case 10:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.mapentries = []*MapEntrySpec{ValueTypeParserDollar[1].mapentry}
		}
	case 11:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.mapentries = append([]*MapEntrySpec{ValueTypeParserDollar[1].mapentry}, ValueTypeParserDollar[3].mapentries...)
		}
	case 12:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			fmt.Println("sty", ValueTypeParserDollar[1].token.literal, ValueTypeParserDollar[3].typedesc)
			ValueTypeParserVAL.mapentry = &MapEntrySpec{
				Name:       ValueTypeParserDollar[1].token.literal,
				Required:   true,
				DefaultVal: nil,
				ValType:    ValueTypeParserDollar[3].typedesc}
		}
	case 13:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-4 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.mapentry = &MapEntrySpec{
				Name:       ValueTypeParserDollar[1].token.literal,
				Required:   false,
				DefaultVal: nil,
				ValType:    ValueTypeParserDollar[4].typedesc}

		}
	case 14:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-5 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.mapentry = &MapEntrySpec{
				Name:       ValueTypeParserDollar[1].token.literal,
				Required:   false,
				DefaultVal: ValueTypeParserDollar[5].val,
				ValType:    ValueTypeParserDollar[3].typedesc}
		}
	case 15:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			x, _ := strconv.ParseFloat(ValueTypeParserDollar[1].token.literal, 64)
			ValueTypeParserVAL.val = value.MakeFloatValue(x)
		}
	case 16:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = &value.Value{Type: value.VAL_STRING, StringVal: ValueTypeParserDollar[1].token.literal}
		}
	case 17:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = value.MakeBoolValue(true)
		}
	case 18:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = value.MakeBoolValue(false)
		}
	case 19:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = ValueTypeParserDollar[2].val
		}
	case 20:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = ValueTypeParserDollar[2].val
		}
	case 21:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = &value.Value{
				Type:   value.VAL_MAP,
				MapVal: map[string]*value.Value{ValueTypeParserDollar[1].token.literal: ValueTypeParserDollar[3].val}}
		}
	case 22:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-5 : ValueTypeParserpt+1]

		{
			map_val := ValueTypeParserDollar[5].val.MapVal
			map_val[ValueTypeParserDollar[1].token.literal] = ValueTypeParserDollar[3].val
			ValueTypeParserVAL.val = &value.Value{
				Type:   value.VAL_MAP,
				MapVal: map_val}
		}
	case 23:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-1 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = &value.Value{
				Type:    value.VAL_LIST,
				ListVal: []*value.Value{ValueTypeParserDollar[1].val}}
		}
	case 24:
		ValueTypeParserDollar = ValueTypeParserS[ValueTypeParserpt-3 : ValueTypeParserpt+1]

		{
			ValueTypeParserVAL.val = &value.Value{
				Type:    value.VAL_LIST,
				ListVal: append([]*value.Value{ValueTypeParserDollar[1].val}, ValueTypeParserDollar[3].val.ListVal...)}
		}
	}
	goto ValueTypeParserstack /* stack new state and value */
}
