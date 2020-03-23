


package value

import __yyfmt__ "fmt"


import (
	//"fmt"
	"errors"
	"strconv"
	"strings"
)

type Token struct {
	token    int
	literal  string
	position Position
}


type ValueParserSymType struct {
	yys   int
	token Token
	val   *Value
}

const IDENT = 57346
const NUMBER = 57347
const STRING = 57348
const TRUE = 57349
const FALSE = 57350

var ValueParserToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"NUMBER",
	"STRING",
	"TRUE",
	"FALSE",
	"','",
	"'{'",
	"'}'",
	"'['",
	"']'",
	"'('",
	"')'",
	"':'",
}
var ValueParserStatenames = [...]string{}

const ValueParserEofCode = 1
const ValueParserErrCode = 2
const ValueParserInitialStackSize = 16



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


var ValueParserExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const ValueParserPrivate = 57344

const ValueParserLast = 29

var ValueParserAct = [...]int{

	9, 3, 4, 5, 6, 12, 7, 16, 8, 13,
	14, 2, 3, 4, 5, 6, 17, 7, 11, 8,
	15, 21, 22, 11, 20, 10, 18, 19, 1,
}
var ValueParserPact = [...]int{

	7, -1000, -1000, -1000, -1000, -1000, -1000, 14, -4, 9,
	-1000, -9, 3, -1000, 17, -1000, 7, -1000, 7, 12,
	-1000, 19, -1000,
}
var ValueParserPgo = [...]int{

	0, 28, 10, 5, 0,
}
var ValueParserR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	4, 4, 3, 3,
}
var ValueParserR2 = [...]int{

	0, 1, 1, 1, 1, 1, 3, 2, 3, 2,
	3, 5, 1, 3,
}
var ValueParserChk = [...]int{

	-1000, -1, -2, 5, 6, 7, 8, 10, 12, -4,
	11, 4, -3, 13, -2, 11, 16, 13, 9, -2,
	-3, 9, -4,
}
var ValueParserDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 0, 0, 0,
	7, 0, 0, 9, 12, 6, 0, 8, 0, 10,
	13, 0, 11,
}
var ValueParserTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	14, 15, 3, 3, 9, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 16, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 12, 3, 13, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 10, 3, 11,
}
var ValueParserTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8,
}
var ValueParserTok3 = [...]int{
	0,
}

var ValueParserErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}



/*	parser for yacc output	*/

var (
	ValueParserDebug        = 0
	ValueParserErrorVerbose = false
)

type ValueParserLexer interface {
	Lex(lval *ValueParserSymType) int
	Error(s string)
}

type ValueParserParser interface {
	Parse(ValueParserLexer) int
	Lookahead() int
}

type ValueParserParserImpl struct {
	lval  ValueParserSymType
	stack [ValueParserInitialStackSize]ValueParserSymType
	char  int
}

func (p *ValueParserParserImpl) Lookahead() int {
	return p.char
}

func ValueParserNewParser() ValueParserParser {
	return &ValueParserParserImpl{}
}

const ValueParserFlag = -1000

func ValueParserTokname(c int) string {
	if c >= 1 && c-1 < len(ValueParserToknames) {
		if ValueParserToknames[c-1] != "" {
			return ValueParserToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func ValueParserStatname(s int) string {
	if s >= 0 && s < len(ValueParserStatenames) {
		if ValueParserStatenames[s] != "" {
			return ValueParserStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func ValueParserErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !ValueParserErrorVerbose {
		return "syntax error"
	}

	for _, e := range ValueParserErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + ValueParserTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := ValueParserPact[state]
	for tok := TOKSTART; tok-1 < len(ValueParserToknames); tok++ {
		if n := base + tok; n >= 0 && n < ValueParserLast && ValueParserChk[ValueParserAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if ValueParserDef[state] == -2 {
		i := 0
		for ValueParserExca[i] != -1 || ValueParserExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; ValueParserExca[i] >= 0; i += 2 {
			tok := ValueParserExca[i]
			if tok < TOKSTART || ValueParserExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if ValueParserExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += ValueParserTokname(tok)
	}
	return res
}

func ValueParserlex1(lex ValueParserLexer, lval *ValueParserSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = ValueParserTok1[0]
		goto out
	}
	if char < len(ValueParserTok1) {
		token = ValueParserTok1[char]
		goto out
	}
	if char >= ValueParserPrivate {
		if char < ValueParserPrivate+len(ValueParserTok2) {
			token = ValueParserTok2[char-ValueParserPrivate]
			goto out
		}
	}
	for i := 0; i < len(ValueParserTok3); i += 2 {
		token = ValueParserTok3[i+0]
		if token == char {
			token = ValueParserTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = ValueParserTok2[1] /* unknown char */
	}
	if ValueParserDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", ValueParserTokname(token), uint(char))
	}
	return char, token
}

func ValueParserParse(ValueParserlex ValueParserLexer) int {
	return ValueParserNewParser().Parse(ValueParserlex)
}

func (ValueParserrcvr *ValueParserParserImpl) Parse(ValueParserlex ValueParserLexer) int {
	var ValueParsern int
	var ValueParserVAL ValueParserSymType
	var ValueParserDollar []ValueParserSymType
	_ = ValueParserDollar // silence set and not used
	ValueParserS := ValueParserrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	ValueParserstate := 0
	ValueParserrcvr.char = -1
	ValueParsertoken := -1 // ValueParserrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		ValueParserstate = -1
		ValueParserrcvr.char = -1
		ValueParsertoken = -1
	}()
	ValueParserp := -1
	goto ValueParserstack

ret0:
	return 0

ret1:
	return 1

ValueParserstack:
	/* put a state and value onto the stack */
	if ValueParserDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", ValueParserTokname(ValueParsertoken), ValueParserStatname(ValueParserstate))
	}

	ValueParserp++
	if ValueParserp >= len(ValueParserS) {
		nyys := make([]ValueParserSymType, len(ValueParserS)*2)
		copy(nyys, ValueParserS)
		ValueParserS = nyys
	}
	ValueParserS[ValueParserp] = ValueParserVAL
	ValueParserS[ValueParserp].yys = ValueParserstate

ValueParsernewstate:
	ValueParsern = ValueParserPact[ValueParserstate]
	if ValueParsern <= ValueParserFlag {
		goto ValueParserdefault /* simple state */
	}
	if ValueParserrcvr.char < 0 {
		ValueParserrcvr.char, ValueParsertoken = ValueParserlex1(ValueParserlex, &ValueParserrcvr.lval)
	}
	ValueParsern += ValueParsertoken
	if ValueParsern < 0 || ValueParsern >= ValueParserLast {
		goto ValueParserdefault
	}
	ValueParsern = ValueParserAct[ValueParsern]
	if ValueParserChk[ValueParsern] == ValueParsertoken { /* valid shift */
		ValueParserrcvr.char = -1
		ValueParsertoken = -1
		ValueParserVAL = ValueParserrcvr.lval
		ValueParserstate = ValueParsern
		if Errflag > 0 {
			Errflag--
		}
		goto ValueParserstack
	}

ValueParserdefault:
	/* default state action */
	ValueParsern = ValueParserDef[ValueParserstate]
	if ValueParsern == -2 {
		if ValueParserrcvr.char < 0 {
			ValueParserrcvr.char, ValueParsertoken = ValueParserlex1(ValueParserlex, &ValueParserrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if ValueParserExca[xi+0] == -1 && ValueParserExca[xi+1] == ValueParserstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			ValueParsern = ValueParserExca[xi+0]
			if ValueParsern < 0 || ValueParsern == ValueParsertoken {
				break
			}
		}
		ValueParsern = ValueParserExca[xi+1]
		if ValueParsern < 0 {
			goto ret0
		}
	}
	if ValueParsern == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			ValueParserlex.Error(ValueParserErrorMessage(ValueParserstate, ValueParsertoken))
			Nerrs++
			if ValueParserDebug >= 1 {
				__yyfmt__.Printf("%s", ValueParserStatname(ValueParserstate))
				__yyfmt__.Printf(" saw %s\n", ValueParserTokname(ValueParsertoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for ValueParserp >= 0 {
				ValueParsern = ValueParserPact[ValueParserS[ValueParserp].yys] + ValueParserErrCode
				if ValueParsern >= 0 && ValueParsern < ValueParserLast {
					ValueParserstate = ValueParserAct[ValueParsern] /* simulate a shift of "error" */
					if ValueParserChk[ValueParserstate] == ValueParserErrCode {
						goto ValueParserstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if ValueParserDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", ValueParserS[ValueParserp].yys)
				}
				ValueParserp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if ValueParserDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", ValueParserTokname(ValueParsertoken))
			}
			if ValueParsertoken == ValueParserEofCode {
				goto ret1
			}
			ValueParserrcvr.char = -1
			ValueParsertoken = -1
			goto ValueParsernewstate /* try again in the same state */
		}
	}

	/* reduction by production ValueParsern */
	if ValueParserDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", ValueParsern, ValueParserStatname(ValueParserstate))
	}

	ValueParsernt := ValueParsern
	ValueParserpt := ValueParserp
	_ = ValueParserpt // guard against "declared and not used"

	ValueParserp -= ValueParserR2[ValueParsern]
	// ValueParserp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if ValueParserp+1 >= len(ValueParserS) {
		nyys := make([]ValueParserSymType, len(ValueParserS)*2)
		copy(nyys, ValueParserS)
		ValueParserS = nyys
	}
	ValueParserVAL = ValueParserS[ValueParserp+1]

	/* consult goto table to find next state */
	ValueParsern = ValueParserR1[ValueParsern]
	ValueParserg := ValueParserPgo[ValueParsern]
	ValueParserj := ValueParserg + ValueParserS[ValueParserp].yys + 1

	if ValueParserj >= ValueParserLast {
		ValueParserstate = ValueParserAct[ValueParserg]
	} else {
		ValueParserstate = ValueParserAct[ValueParserj]
		if ValueParserChk[ValueParserstate] != -ValueParsern {
			ValueParserstate = ValueParserAct[ValueParserg]
		}
	}
	// dummy call; replaced with literal code
	switch ValueParsernt {

	case 1:
		ValueParserDollar = ValueParserS[ValueParserpt-1 : ValueParserpt+1]

		{
			ValueParserVAL.val = nil
			if l, ok := ValueParserlex.(*ValueLexer); ok {
				l.result = ValueParserDollar[1].val
			}
		}
	case 2:
		ValueParserDollar = ValueParserS[ValueParserpt-1 : ValueParserpt+1]

		{
			x, _ := strconv.ParseFloat(ValueParserDollar[1].token.literal, 64)
			ValueParserVAL.val = MakeFloatValue(x)
		}
	case 3:
		ValueParserDollar = ValueParserS[ValueParserpt-1 : ValueParserpt+1]

		{
			ValueParserVAL.val = &Value{Type: VAL_STRING, StringVal: ValueParserDollar[1].token.literal}
		}
	case 4:
		ValueParserDollar = ValueParserS[ValueParserpt-1 : ValueParserpt+1]

		{
			ValueParserVAL.val = MakeBoolValue(true)
		}
	case 5:
		ValueParserDollar = ValueParserS[ValueParserpt-1 : ValueParserpt+1]

		{
			ValueParserVAL.val = MakeBoolValue(false)
		}
	case 6:
		ValueParserDollar = ValueParserS[ValueParserpt-3 : ValueParserpt+1]

		{
			ValueParserVAL.val = ValueParserDollar[2].val
		}
	case 7:
		ValueParserDollar = ValueParserS[ValueParserpt-2 : ValueParserpt+1]

		{
			ValueParserVAL.val = &Value{
				Type:   VAL_MAP,
				MapVal: map[string]*Value{}}
		}
	case 8:
		ValueParserDollar = ValueParserS[ValueParserpt-3 : ValueParserpt+1]

		{
			ValueParserVAL.val = ValueParserDollar[2].val
		}
	case 9:
		ValueParserDollar = ValueParserS[ValueParserpt-2 : ValueParserpt+1]

		{
			ValueParserVAL.val = &Value{
				Type:    VAL_LIST,
				ListVal: []*Value{}}
		}
	case 10:
		ValueParserDollar = ValueParserS[ValueParserpt-3 : ValueParserpt+1]

		{
			ValueParserVAL.val = &Value{
				Type:   VAL_MAP,
				MapVal: map[string]*Value{ValueParserDollar[1].token.literal: ValueParserDollar[3].val}}
		}
	case 11:
		ValueParserDollar = ValueParserS[ValueParserpt-5 : ValueParserpt+1]

		{
			map_val := ValueParserDollar[5].val.MapVal
			map_val[ValueParserDollar[1].token.literal] = ValueParserDollar[3].val
			ValueParserVAL.val = &Value{
				Type:   VAL_MAP,
				MapVal: map_val}
		}
	case 12:
		ValueParserDollar = ValueParserS[ValueParserpt-1 : ValueParserpt+1]

		{
			ValueParserVAL.val = &Value{
				Type:    VAL_LIST,
				ListVal: []*Value{ValueParserDollar[1].val}}
		}
	case 13:
		ValueParserDollar = ValueParserS[ValueParserpt-3 : ValueParserpt+1]

		{
			ValueParserVAL.val = &Value{
				Type:    VAL_LIST,
				ListVal: append([]*Value{ValueParserDollar[1].val}, ValueParserDollar[3].val.ListVal...)}
		}
	}
	goto ValueParserstack /* stack new state and value */
}
