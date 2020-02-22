// Code generated by goyacc -o src/mc/router/parser.go -v src/mc/router/parser.output src/mc/router/parser.go.y. DO NOT EDIT.

//line src/mc/router/parser.go.y:2
package router

import __yyfmt__ "fmt"

//line src/mc/router/parser.go.y:2
import (
	"strconv"
	//"fmt"
)

type Token struct {
	token    int
	literal  string
	position Position
}

//line src/mc/router/parser.go.y:13
type yySymType struct {
	yys        int
	token      Token
	routes     []*Route
	transforms *Route
	transform  *Transform
	expression *Expression
	statement  *Statement
	writeable  *WriteableValue
	script     []*Statement
	node       *Node
}

const IDENT = 57346
const VAR = 57347
const NUMBER = 57348
const STRING = 57349
const THIS = 57350
const AND = 57351
const OR = 57352
const EQ = 57353
const NE = 57354
const LE = 57355
const GE = 57356
const PE = 57357
const ME = 57358
const TE = 57359
const DE = 57360
const RE = 57361
const XE = 57362
const SUB = 57363
const UNARY = 57364
const IS = 57365

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"VAR",
	"NUMBER",
	"STRING",
	"THIS",
	"AND",
	"OR",
	"EQ",
	"NE",
	"LE",
	"GE",
	"PE",
	"ME",
	"TE",
	"DE",
	"RE",
	"XE",
	"SUB",
	"'?'",
	"'%'",
	"'='",
	"'{'",
	"'}'",
	"'['",
	"']'",
	"'<'",
	"'>'",
	"':'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'&'",
	"'|'",
	"'^'",
	"'!'",
	"'~'",
	"UNARY",
	"IS",
	"';'",
	"'('",
	"')'",
	"','",
	"'.'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line src/mc/router/parser.go.y:402

func Parse(exp string) []*Route {
	l := new(RouteLexer)
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	yyParse(l)
	return l.result
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 651

var yyAct = [...]int{

	38, 61, 33, 60, 69, 70, 71, 72, 73, 116,
	62, 6, 77, 68, 3, 16, 74, 5, 4, 23,
	99, 129, 128, 98, 76, 66, 59, 113, 112, 78,
	79, 63, 64, 19, 65, 39, 75, 67, 18, 17,
	14, 110, 8, 80, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 90, 91, 92, 93, 94, 95, 22,
	97, 36, 96, 2, 3, 37, 101, 15, 7, 13,
	103, 104, 105, 106, 107, 108, 109, 34, 20, 111,
	21, 31, 53, 54, 49, 50, 52, 51, 35, 9,
	1, 0, 0, 0, 0, 0, 48, 0, 3, 0,
	58, 0, 56, 55, 115, 41, 42, 43, 44, 45,
	46, 47, 0, 0, 126, 125, 10, 11, 12, 124,
	57, 0, 0, 0, 0, 127, 53, 54, 49, 50,
	52, 51, 0, 0, 0, 0, 0, 0, 0, 0,
	48, 0, 0, 0, 58, 123, 56, 55, 0, 41,
	42, 43, 44, 45, 46, 47, 53, 54, 49, 50,
	52, 51, 0, 0, 57, 0, 0, 0, 0, 0,
	48, 0, 0, 0, 58, 0, 56, 55, 0, 41,
	42, 43, 44, 45, 46, 47, 0, 0, 0, 0,
	122, 0, 0, 0, 57, 53, 54, 49, 50, 52,
	51, 0, 0, 0, 0, 0, 0, 0, 0, 48,
	0, 0, 0, 58, 0, 56, 55, 0, 41, 42,
	43, 44, 45, 46, 47, 0, 0, 0, 0, 121,
	0, 0, 0, 57, 53, 54, 49, 50, 52, 51,
	0, 0, 0, 0, 0, 0, 0, 0, 48, 0,
	0, 0, 58, 0, 56, 55, 0, 41, 42, 43,
	44, 45, 46, 47, 0, 0, 0, 0, 120, 0,
	0, 0, 57, 53, 54, 49, 50, 52, 51, 0,
	0, 0, 0, 0, 0, 0, 0, 48, 0, 0,
	0, 58, 0, 56, 55, 0, 41, 42, 43, 44,
	45, 46, 47, 0, 0, 0, 0, 119, 0, 0,
	0, 57, 53, 54, 49, 50, 52, 51, 0, 0,
	0, 0, 0, 0, 0, 0, 48, 0, 0, 0,
	58, 0, 56, 55, 0, 41, 42, 43, 44, 45,
	46, 47, 0, 0, 0, 0, 118, 0, 0, 0,
	57, 53, 54, 49, 50, 52, 51, 0, 0, 0,
	0, 0, 0, 0, 0, 48, 0, 0, 0, 58,
	0, 56, 55, 0, 41, 42, 43, 44, 45, 46,
	47, 0, 0, 0, 0, 117, 0, 0, 0, 57,
	53, 54, 49, 50, 52, 51, 0, 0, 0, 0,
	0, 0, 0, 0, 48, 0, 0, 0, 58, 114,
	56, 55, 0, 41, 42, 43, 44, 45, 46, 47,
	53, 54, 49, 50, 52, 51, 0, 0, 57, 0,
	0, 0, 0, 0, 48, 0, 0, 0, 58, 0,
	56, 55, 0, 41, 42, 43, 44, 45, 46, 47,
	53, 54, 49, 50, 52, 51, 102, 0, 57, 0,
	0, 0, 0, 0, 48, 0, 0, 0, 58, 0,
	56, 55, 0, 41, 42, 43, 44, 45, 46, 47,
	53, 54, 49, 50, 52, 51, 0, 100, 57, 0,
	0, 0, 0, 0, 48, 0, 0, 40, 58, 0,
	56, 55, 0, 41, 42, 43, 44, 45, 46, 47,
	53, 54, 49, 50, 52, 51, 0, 0, 57, 0,
	0, 0, 0, 0, 48, 0, 0, 0, 58, 0,
	56, 55, 0, 41, 42, 43, 44, 45, 46, 47,
	53, 54, 49, 50, 0, 0, 0, 0, 57, 0,
	0, 0, 0, 0, 48, 0, 0, 0, 58, 0,
	56, 55, 0, 41, 42, 43, 44, 45, 46, 47,
	54, 49, 50, 0, 0, 0, 0, 0, 57, 0,
	0, 0, 0, 48, 0, 0, 0, 58, 0, 56,
	55, 0, 41, 42, 43, 44, 45, 46, 47, 49,
	50, 0, 0, 0, 0, 0, 0, 57, 0, 0,
	27, 48, 24, 28, 0, 58, 0, 56, 55, 0,
	41, 42, 43, 44, 45, 46, 47, 0, 0, 0,
	0, 25, 0, 26, 0, 57, 0, 0, 0, 29,
	0, 0, 0, 0, 0, 32, 0, 0, 0, 0,
	30,
}
var yyPact = [...]int{

	60, -1000, -12, -24, 94, 10, 63, -1000, -1000, -15,
	14, 13, 8, -1000, 60, -1000, 94, 606, 57, 31,
	-1000, -1000, -1000, 471, -1000, 31, 606, -34, -1000, 606,
	606, -1000, 606, -1, 57, -11, -1000, -1000, -2, -19,
	6, 606, 606, 606, 606, 606, 606, 606, 606, 606,
	606, 606, 606, 606, 606, 606, 606, 58, 606, -3,
	-8, 441, 606, 501, 411, 588, -1000, -1000, 606, 606,
	606, 606, 606, 606, 606, 37, -1000, 606, 3, 2,
	501, 501, 501, 501, 501, 501, 501, 501, 501, 501,
	531, 531, 560, 588, 501, 501, -1000, 381, -1000, -1000,
	606, -36, -1000, 342, 303, 264, 225, 186, 147, 117,
	-1000, 73, 57, 31, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 31, -4, -5, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 90, 63, 42, 89, 2, 88, 81, 77, 1,
	0, 3,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 1, 2, 2, 3, 3, 4,
	4, 4, 4, 4, 5, 5, 8, 8, 8, 8,
	8, 8, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 10, 10, 11,
	11, 7, 7, 7, 6, 6, 6, 6,
}
var yyR2 = [...]int{

	0, 3, 3, 4, 3, 1, 3, 3, 3, 4,
	4, 4, 8, 8, 1, 2, 4, 4, 4, 4,
	4, 4, 1, 3, 3, 4, 1, 2, 3, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 2, 3, 3, 3, 5, 1,
	3, 3, 4, 1, 1, 1, 4, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 4, 30, 29, 35, -2, -3, -4,
	22, 23, 24, -2, 30, 4, 30, 25, 25, 25,
	-2, -2, -3, -9, 6, 25, 27, 4, 7, 33,
	44, -7, 39, -5, -8, -6, 4, 8, -10, 4,
	26, 32, 33, 34, 35, 36, 37, 38, 23, 11,
	12, 14, 13, 9, 10, 30, 29, 47, 27, -10,
	-11, -9, 44, -9, -9, -9, 26, -5, 24, 15,
	16, 17, 18, 19, 27, 47, 26, 31, 23, 24,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, 4, -9, 26, 28,
	46, -11, 45, -9, -9, -9, -9, -9, -9, -9,
	4, -9, 25, 25, 28, -11, 45, 43, 43, 43,
	43, 43, 43, 28, 46, -5, -10, -10, 26, 26,
}
var yyDef = [...]int{

	0, -2, 0, 5, 0, 0, 0, 1, 4, 0,
	0, 0, 0, 2, 0, 6, 0, 0, 0, 0,
	3, 7, 8, 0, 22, 0, 0, 53, 26, 0,
	0, 29, 0, 0, 14, 0, 54, 55, 0, 0,
	9, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 49, 0, 27, 0, 44, 10, 15, 0, 0,
	0, 0, 0, 0, 0, 0, 11, 0, 0, 0,
	30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 45, 46, 51, 0, 23, 24,
	0, 0, 28, 0, 0, 0, 0, 0, 0, 0,
	57, 47, 0, 0, 52, 50, 25, 16, 17, 18,
	19, 20, 21, 56, 0, 0, 0, 48, 12, 13,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 39, 3, 3, 3, 23, 36, 3,
	44, 45, 34, 32, 46, 33, 47, 35, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 31, 43,
	29, 24, 30, 22, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 27, 3, 28, 38, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 25, 37, 26, 40,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	41, 42,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:45
		{
			// fmt.Println("C")
			yyVAL.routes = nil
			if l, ok := yylex.(*RouteLexer); ok {
				l.result = []*Route{&Route{Source: yyDollar[1].node.Name, Dest: yyDollar[3].node.Name}}
			}
		}
	case 2:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:53
		{
			// fmt.Println("B")
			yyVAL.routes = nil
			if l, ok := yylex.(*RouteLexer); ok {
				l.result = []*Route{&Route{Source: yyDollar[3].node.Name, Dest: yyDollar[1].node.Name}}
			}
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:61
		{
			yyVAL.routes = nil
			if l, ok := yylex.(*RouteLexer); ok {
				l.result = []*Route{
					&Route{Source: yyDollar[1].node.Name, Dest: yyDollar[4].node.Name},
					&Route{Source: yyDollar[4].node.Name, Dest: yyDollar[1].node.Name}}
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:70
		{
			yyVAL.routes = nil
			// fmt.Println("A")
			if l, ok := yylex.(*RouteLexer); ok {
				l.result = []*Route{
					&Route{
						Source:     yyDollar[1].node.Name,
						Dest:       yyDollar[3].transforms.Dest,
						Transforms: yyDollar[3].transforms.Transforms}}
			}
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:83
		{
			yyVAL.node = &Node{Name: yyDollar[1].token.literal}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:87
		{
			yyVAL.node = &Node{Group: yyDollar[1].token.literal, Name: yyDollar[3].token.literal}
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:92
		{
			yyVAL.transforms = &Route{
				Dest:       yyDollar[3].node.Name,
				Transforms: []*Transform{yyDollar[1].transform}}
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:98
		{
			yyVAL.transforms = &Route{
				Dest:       yyDollar[3].transforms.Dest,
				Transforms: append([]*Transform{yyDollar[1].transform}, yyDollar[3].transforms.Transforms...)}
		}
	case 9:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:105
		{
			yyVAL.transform = &Transform{
				Type:      TR_FILTER,
				Condition: yyDollar[3].expression}
		}
	case 10:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:111
		{
			yyVAL.transform = &Transform{
				Type:   TR_EDIT,
				Script: yyDollar[3].script}
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:117
		{
			yyVAL.transform = &Transform{
				Type:    TR_REPLACE,
				Replace: yyDollar[3].expression}
		}
	case 12:
		yyDollar = yyS[yypt-8 : yypt+1]
//line src/mc/router/parser.go.y:123
		{
			yyVAL.transform = &Transform{
				Type:      TR_COND_EDIT,
				Condition: yyDollar[3].expression,
				Script:    yyDollar[7].script}
		}
	case 13:
		yyDollar = yyS[yypt-8 : yypt+1]
//line src/mc/router/parser.go.y:130
		{
			yyVAL.transform = &Transform{
				Type:      TR_COND_REPLACE,
				Condition: yyDollar[3].expression,
				Replace:   yyDollar[7].expression}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:139
		{
			yyVAL.script = []*Statement{yyDollar[1].statement}
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
//line src/mc/router/parser.go.y:143
		{
			yyVAL.script = append([]*Statement{yyDollar[1].statement}, yyDollar[2].script...)
		}
	case 16:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:148
		{
			yyVAL.statement = MakeAssignment(yyDollar[1].writeable, yyDollar[3].expression)
		}
	case 17:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:152
		{
			yyVAL.statement = MakeAssignment(yyDollar[1].writeable, &Expression{
				Operation: OP_PLUS,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:158
		{
			yyVAL.statement = MakeAssignment(yyDollar[1].writeable, &Expression{
				Operation: OP_MINUS,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:164
		{
			yyVAL.statement = MakeAssignment(yyDollar[1].writeable, &Expression{
				Operation: OP_MUL,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:170
		{
			yyVAL.statement = MakeAssignment(yyDollar[1].writeable, &Expression{
				Operation: OP_DIV,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:176
		{
			yyVAL.statement = MakeAssignment(yyDollar[1].writeable, &Expression{
				Operation: OP_MOD,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:183
		{
			x, _ := strconv.Atoi(yyDollar[1].token.literal)
			yyVAL.expression = &Expression{
				Operation: OP_NUM,
				Value:     &Value{Type: VAL_NUM, NumVal: float64(x)}}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:190
		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:194
		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 25:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:198
		{
			yyVAL.expression = &Expression{
				Operation: OP_CALL,
				Args: []*Expression{
					MakeNameExpression(yyDollar[1].token.literal),
					yyDollar[3].expression}}
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:206
		{
			yyVAL.expression = &Expression{
				Operation: OP_STRING,
				Value:     &Value{Type: VAL_STRING, StringVal: yyDollar[1].token.literal}}
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
//line src/mc/router/parser.go.y:212
		{
			yyVAL.expression = &Expression{
				Operation: OP_UMINUS,
				Args:      []*Expression{yyDollar[2].expression}}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:218
		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:222
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:226
		{
			yyVAL.expression = &Expression{
				Operation: OP_PLUS,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:232
		{
			yyVAL.expression = &Expression{
				Operation: OP_MINUS,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:238
		{
			yyVAL.expression = &Expression{
				Operation: OP_MUL,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:244
		{
			yyVAL.expression = &Expression{
				Operation: OP_DIV,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:250
		{
			yyVAL.expression = &Expression{
				Operation: OP_BITWISEAND,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:256
		{
			yyVAL.expression = &Expression{
				Operation: OP_BITWISEOR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:262
		{
			yyVAL.expression = &Expression{
				Operation: OP_BITWISEXOR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:268
		{
			yyVAL.expression = &Expression{
				Operation: OP_MOD,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:274
		{
			yyVAL.expression = &Expression{
				Operation: OP_EQ,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:280
		{
			yyVAL.expression = &Expression{
				Operation: OP_NE,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:286
		{
			yyVAL.expression = &Expression{
				Operation: OP_GE,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:292
		{
			yyVAL.expression = &Expression{
				Operation: OP_LE,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:298
		{
			yyVAL.expression = &Expression{
				Operation: OP_AND,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:304
		{
			yyVAL.expression = &Expression{
				Operation: OP_OR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
//line src/mc/router/parser.go.y:310
		{
			yyVAL.expression = &Expression{
				Operation: OP_NOT,
				Args:      []*Expression{yyDollar[2].expression}}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:316
		{
			yyVAL.expression = &Expression{
				Operation: OP_GT,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:322
		{
			yyVAL.expression = &Expression{
				Operation: OP_LT,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:329
		{
			yyVAL.expression = &Expression{
				Operation: OP_MAP,
				Args: []*Expression{
					MakeNameExpression(yyDollar[1].token.literal),
					yyDollar[3].expression}}
		}
	case 48:
		yyDollar = yyS[yypt-5 : yypt+1]
//line src/mc/router/parser.go.y:337
		{
			args := []*Expression{
				MakeNameExpression(yyDollar[1].token.literal),
				yyDollar[3].expression}
			yyVAL.expression = &Expression{
				Operation: OP_MAP,
				Args:      append(args, yyDollar[5].expression.Args...)}
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:347
		{
			yyVAL.expression = &Expression{
				Operation: OP_LIST,
				Args:      []*Expression{yyDollar[1].expression}}
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:353
		{
			args := []*Expression{yyDollar[1].expression}
			yyVAL.expression = &Expression{
				Operation: OP_LIST,
				Args:      append(args, yyDollar[3].expression.Args...)}
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:361
		{
			yyVAL.expression = &Expression{
				Operation: OP_MAPVAR,
				Args: []*Expression{
					yyDollar[1].expression,
					MakeNameExpression(yyDollar[3].token.literal)}}
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:369
		{
			yyVAL.expression = &Expression{
				Operation: OP_LISTVAR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:375
		{
			yyVAL.expression = MakeNameExpression(yyDollar[1].token.literal)
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:380
		{
			yyVAL.writeable = &WriteableValue{
				Base: yyDollar[1].token.literal,
				Path: []PathEntry{}}
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
//line src/mc/router/parser.go.y:386
		{
			yyVAL.writeable = &WriteableValue{
				Base: "this",
				Path: []PathEntry{}}
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
//line src/mc/router/parser.go.y:392
		{
			yyDollar[1].writeable.Path = append(yyDollar[1].writeable.Path, PathEntry{Type: PATH_LIST, ListIndex: yyDollar[3].expression})
			yyVAL.writeable = yyDollar[1].writeable
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
//line src/mc/router/parser.go.y:397
		{
			yyDollar[1].writeable.Path = append(yyDollar[1].writeable.Path, PathEntry{Type: PATH_MAP, MapKey: yyDollar[3].token.literal})
			yyVAL.writeable = yyDollar[1].writeable
		}
	}
	goto yystack /* stack new state and value */
}
