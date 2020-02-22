


package router

import __yyfmt__ "fmt"


import (
	"strconv"
	//"fmt"
)

type Token struct {
	token    int
	literal  string
	position Position
}


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
const DEL = 57348
const NUMBER = 57349
const STRING = 57350
const THIS = 57351
const TRUE = 57352
const FALSE = 57353
const AND = 57354
const OR = 57355
const EQ = 57356
const NE = 57357
const LE = 57358
const GE = 57359
const PE = 57360
const ME = 57361
const TE = 57362
const DE = 57363
const RE = 57364
const AE = 57365
const OE = 57366
const XE = 57367
const SUB = 57368
const UNARY = 57369

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"VAR",
	"DEL",
	"NUMBER",
	"STRING",
	"THIS",
	"TRUE",
	"FALSE",
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
	"AE",
	"OE",
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
	"'.'",
	"';'",
	"'('",
	"')'",
	"','",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16



func Parse(exp string) []*Route {
	l := new(RouteLexer)
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	yyParse(l)
	return l.result
}


var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 718

var yyAct = [...]int{

	42, 65, 35, 64, 130, 66, 76, 77, 78, 79,
	80, 73, 74, 75, 124, 123, 62, 72, 52, 23,
	81, 6, 62, 86, 16, 5, 4, 108, 63, 47,
	48, 61, 3, 67, 68, 82, 69, 61, 146, 71,
	145, 107, 3, 85, 70, 127, 126, 89, 90, 91,
	92, 93, 94, 95, 96, 97, 98, 99, 100, 101,
	102, 103, 104, 14, 106, 10, 11, 12, 87, 88,
	110, 19, 18, 17, 112, 113, 114, 115, 116, 117,
	118, 119, 120, 121, 2, 8, 43, 122, 125, 7,
	13, 105, 57, 58, 53, 54, 56, 55, 84, 20,
	83, 21, 22, 3, 40, 38, 39, 15, 52, 41,
	36, 33, 62, 129, 60, 59, 37, 45, 46, 47,
	48, 49, 50, 51, 9, 1, 0, 61, 143, 142,
	0, 141, 57, 58, 53, 54, 56, 55, 0, 0,
	0, 0, 144, 0, 0, 0, 0, 0, 52, 0,
	0, 0, 62, 0, 60, 59, 0, 45, 46, 47,
	48, 49, 50, 51, 0, 0, 0, 61, 0, 0,
	0, 109, 57, 58, 53, 54, 56, 55, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 52, 0,
	0, 0, 62, 0, 60, 59, 0, 45, 46, 47,
	48, 49, 50, 51, 0, 0, 0, 61, 0, 0,
	111, 57, 58, 53, 54, 56, 55, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 52, 0, 0,
	0, 62, 0, 60, 59, 0, 45, 46, 47, 48,
	49, 50, 51, 0, 0, 0, 61, 139, 57, 58,
	53, 54, 56, 55, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 52, 0, 0, 0, 62, 0,
	60, 59, 0, 45, 46, 47, 48, 49, 50, 51,
	0, 0, 0, 61, 138, 57, 58, 53, 54, 56,
	55, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 52, 0, 0, 0, 62, 0, 60, 59, 0,
	45, 46, 47, 48, 49, 50, 51, 0, 0, 0,
	61, 137, 57, 58, 53, 54, 56, 55, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 52, 0,
	0, 0, 62, 0, 60, 59, 0, 45, 46, 47,
	48, 49, 50, 51, 0, 0, 0, 61, 136, 57,
	58, 53, 54, 56, 55, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 52, 0, 0, 0, 62,
	0, 60, 59, 0, 45, 46, 47, 48, 49, 50,
	51, 0, 0, 0, 61, 135, 57, 58, 53, 54,
	56, 55, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 52, 0, 0, 0, 62, 0, 60, 59,
	0, 45, 46, 47, 48, 49, 50, 51, 0, 0,
	0, 61, 134, 57, 58, 53, 54, 56, 55, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 52,
	0, 0, 0, 62, 0, 60, 59, 0, 45, 46,
	47, 48, 49, 50, 51, 0, 0, 0, 61, 133,
	57, 58, 53, 54, 56, 55, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 52, 0, 0, 0,
	62, 0, 60, 59, 0, 45, 46, 47, 48, 49,
	50, 51, 0, 0, 0, 61, 132, 57, 58, 53,
	54, 56, 55, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 52, 0, 0, 0, 62, 0, 60,
	59, 0, 45, 46, 47, 48, 49, 50, 51, 0,
	0, 0, 61, 131, 57, 58, 53, 54, 56, 55,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	52, 0, 0, 0, 62, 140, 60, 59, 0, 45,
	46, 47, 48, 49, 50, 51, 0, 0, 0, 61,
	57, 58, 53, 54, 56, 55, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 52, 0, 0, 0,
	62, 128, 60, 59, 0, 45, 46, 47, 48, 49,
	50, 51, 0, 0, 0, 61, 57, 58, 53, 54,
	56, 55, 29, 0, 0, 24, 30, 0, 25, 26,
	0, 0, 52, 0, 0, 44, 62, 0, 60, 59,
	0, 45, 46, 47, 48, 49, 50, 51, 27, 0,
	28, 61, 53, 54, 56, 55, 31, 0, 0, 0,
	0, 0, 34, 0, 0, 0, 52, 32, 0, 0,
	62, 0, 60, 59, 0, 45, 46, 47, 48, 49,
	50, 51, 52, 0, 0, 61, 62, 0, 0, 0,
	0, 45, 46, 47, 48, 49, 50, 51, 52, 0,
	0, 61, 62, 0, 0, 0, 0, 45, 46, 47,
	48, 0, 0, 0, 0, 0, 0, 61,
}
var yyPact = [...]int{

	99, -1000, -9, -19, 38, 28, 103, -1000, -1000, -11,
	43, 42, 41, -1000, 99, -1000, 38, 618, 100, 82,
	-1000, -1000, -1000, 604, -1000, -1000, -1000, 82, 618, -44,
	-1000, 618, 618, -1000, 618, 13, 100, -12, 96, 94,
	-1000, -1000, 12, -13, 40, 618, 618, 618, 618, 618,
	618, 618, 618, 618, 618, 618, 618, 618, 618, 618,
	618, 87, 618, 10, -6, 120, 618, -16, 160, -16,
	-1000, -1000, 618, 618, 618, 618, 618, 618, 618, 618,
	618, 618, 83, -33, -34, -1000, 618, 16, 15, -10,
	-10, -16, -16, 670, 670, 670, -16, 654, 654, 654,
	654, 638, 638, 654, 654, -1000, 568, -1000, -1000, 618,
	-46, -1000, 495, 458, 421, 384, 347, 310, 273, 236,
	199, 532, -1000, -1000, -1000, 80, 100, 82, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 82, 9, 7, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 125, 84, 85, 124, 2, 116, 111, 110, 1,
	0, 3,
}
var yyR1 = [...]int{

	0, 1, 1, 1, 1, 2, 2, 3, 3, 4,
	4, 4, 4, 4, 5, 5, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 10, 10, 11, 11, 7, 7,
	7, 6, 6, 6, 6,
}
var yyR2 = [...]int{

	0, 3, 3, 4, 3, 1, 3, 3, 3, 4,
	4, 4, 8, 8, 1, 2, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 3, 3, 1, 1, 1,
	3, 3, 4, 1, 2, 3, 1, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 2, 3, 3, 3, 5, 1, 3, 3, 4,
	1, 1, 1, 4, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 4, 35, 34, 40, -2, -3, -4,
	27, 28, 29, -2, 35, 4, 35, 30, 30, 30,
	-2, -2, -3, -9, 7, 10, 11, 30, 32, 4,
	8, 38, 49, -7, 44, -5, -8, -6, 5, 6,
	4, 9, -10, 4, 31, 37, 38, 39, 40, 41,
	42, 43, 28, 14, 15, 17, 16, 12, 13, 35,
	34, 47, 32, -10, -11, -9, 49, -9, -9, -9,
	31, -5, 29, 23, 24, 25, 18, 19, 20, 21,
	22, 32, 47, 4, 4, 31, 36, 28, 29, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, 4, -9, 31, 33, 51,
	-11, 50, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, 4, 48, 48, -9, 30, 30, 33, -11,
	50, 48, 48, 48, 48, 48, 48, 48, 48, 48,
	33, 51, -5, -10, -10, 31, 31,
}
var yyDef = [...]int{

	0, -2, 0, 5, 0, 0, 0, 1, 4, 0,
	0, 0, 0, 2, 0, 6, 0, 0, 0, 0,
	3, 7, 8, 0, 27, 28, 29, 0, 0, 60,
	33, 0, 0, 36, 0, 0, 14, 0, 0, 0,
	61, 62, 0, 0, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 56, 0, 34, 0, 51,
	10, 15, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 11, 0, 0, 0, 37,
	38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
	48, 49, 50, 52, 53, 58, 0, 30, 31, 0,
	0, 35, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 64, 25, 26, 54, 0, 0, 59, 57,
	32, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	63, 0, 0, 0, 55, 12, 13,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 44, 3, 3, 3, 28, 41, 3,
	49, 50, 39, 37, 51, 38, 47, 40, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 36, 48,
	34, 29, 35, 27, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 32, 3, 33, 43, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 30, 42, 31, 45,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 46,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}



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

		{
			// fmt.Println("C")
			yyVAL.routes = nil
			if l, ok := yylex.(*RouteLexer); ok {
				l.result = []*Route{&Route{Source: yyDollar[1].node.Name, Dest: yyDollar[3].node.Name}}
			}
		}
	case 2:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			// fmt.Println("B")
			yyVAL.routes = nil
			if l, ok := yylex.(*RouteLexer); ok {
				l.result = []*Route{&Route{Source: yyDollar[3].node.Name, Dest: yyDollar[1].node.Name}}
			}
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]

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

		{
			yyVAL.node = &Node{Name: yyDollar[1].token.literal}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.node = &Node{Group: yyDollar[1].token.literal, Name: yyDollar[3].token.literal}
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.transforms = &Route{
				Dest:       yyDollar[3].node.Name,
				Transforms: []*Transform{yyDollar[1].transform}}
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.transforms = &Route{
				Dest:       yyDollar[3].transforms.Dest,
				Transforms: append([]*Transform{yyDollar[1].transform}, yyDollar[3].transforms.Transforms...)}
		}
	case 9:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.transform = &Transform{
				Type:      TR_FILTER,
				Condition: yyDollar[3].expression}
		}
	case 10:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.transform = &Transform{
				Type:   TR_EDIT,
				Script: yyDollar[3].script}
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.transform = &Transform{
				Type:    TR_REPLACE,
				Replace: yyDollar[3].expression}
		}
	case 12:
		yyDollar = yyS[yypt-8 : yypt+1]

		{
			yyVAL.transform = &Transform{
				Type:      TR_COND_EDIT,
				Condition: yyDollar[3].expression,
				Script:    yyDollar[7].script}
		}
	case 13:
		yyDollar = yyS[yypt-8 : yypt+1]

		{
			yyVAL.transform = &Transform{
				Type:      TR_COND_REPLACE,
				Condition: yyDollar[3].expression,
				Replace:   yyDollar[7].expression}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.script = []*Statement{yyDollar[1].statement}
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]

		{
			yyVAL.script = append([]*Statement{yyDollar[1].statement}, yyDollar[2].script...)
		}
	case 16:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, yyDollar[3].expression)
		}
	case 17:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_BITWISEAND,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_BITWISEOR,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_BITWISEXOR,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_PLUS,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_MINUS,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 22:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_MUL,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_DIV,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 24:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.statement = MakeAssignmentStatement(yyDollar[1].writeable, &Expression{
				Operation: OP_MOD,
				Args:      []*Expression{yyDollar[1].writeable.ToExpression(), yyDollar[3].expression}})
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.statement = MakeDeclarationStatement(yyDollar[2].token.literal)
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.statement = MakeDeletionStatement(yyDollar[2].token.literal)
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			x, _ := strconv.ParseFloat(yyDollar[1].token.literal, 64)
			yyVAL.expression = &Expression{
				Operation: OP_NUM,
				Value:     MakeFloatValue(x)}
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_BOOL,
				Value:     MakeBoolValue(true)}
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_BOOL,
				Value:     MakeBoolValue(false)}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_CALL,
				Args: []*Expression{
					MakeNameExpression(yyDollar[1].token.literal),
					yyDollar[3].expression}}
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_STRING,
				Value:     &Value{Type: VAL_STRING, StringVal: yyDollar[1].token.literal}}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_UMINUS,
				Args:      []*Expression{yyDollar[2].expression}}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = yyDollar[2].expression
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_PLUS,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_MINUS,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_MUL,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_DIV,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_BITWISEAND,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_BITWISEOR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_BITWISEXOR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_MOD,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_EQ,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_NE,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_GE,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_LE,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_AND,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_OR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_NOT,
				Args:      []*Expression{yyDollar[2].expression}}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_GT,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_LT,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_MAP,
				Args: []*Expression{
					MakeNameExpression(yyDollar[1].token.literal),
					yyDollar[3].expression}}
		}
	case 55:
		yyDollar = yyS[yypt-5 : yypt+1]

		{
			args := []*Expression{
				MakeNameExpression(yyDollar[1].token.literal),
				yyDollar[3].expression}
			yyVAL.expression = &Expression{
				Operation: OP_MAP,
				Args:      append(args, yyDollar[5].expression.Args...)}
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_LIST,
				Args:      []*Expression{yyDollar[1].expression}}
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			args := []*Expression{yyDollar[1].expression}
			yyVAL.expression = &Expression{
				Operation: OP_LIST,
				Args:      append(args, yyDollar[3].expression.Args...)}
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_MAPVAR,
				Args:      []*Expression{yyDollar[1].expression, MakeNameExpression(yyDollar[3].token.literal)}}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyVAL.expression = &Expression{
				Operation: OP_LISTVAR,
				Args:      []*Expression{yyDollar[1].expression, yyDollar[3].expression}}
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.expression = MakeVarExpression(yyDollar[1].token.literal)
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.writeable = &WriteableValue{
				Base: yyDollar[1].token.literal,
				Path: []PathEntry{}}
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]

		{
			yyVAL.writeable = &WriteableValue{
				Base: "this",
				Path: []PathEntry{}}
		}
	case 63:
		yyDollar = yyS[yypt-4 : yypt+1]

		{
			yyDollar[1].writeable.Path = append(yyDollar[1].writeable.Path, PathEntry{Type: PATH_LIST, ListIndex: yyDollar[3].expression})
			yyVAL.writeable = yyDollar[1].writeable
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]

		{
			yyDollar[1].writeable.Path = append(yyDollar[1].writeable.Path, PathEntry{Type: PATH_MAP, MapKey: yyDollar[3].token.literal})
			yyVAL.writeable = yyDollar[1].writeable
		}
	}
	goto yystack /* stack new state and value */
}
