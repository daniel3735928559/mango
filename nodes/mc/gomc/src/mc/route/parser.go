


package route

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


type RouteParserSymType struct {
	yys        int
	token      Token
	routes     []*Route
	transforms *Route
	transform  *Transform
	expression *Expression
	statement  *Statement
	writeable  *WriteableValue
	script     []*Statement
	node       string
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
const EXP = 57369
const UNARY = 57370

var RouteParserToknames = [...]string{
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
	"EXP",
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
var RouteParserStatenames = [...]string{}

const RouteParserEofCode = 1
const RouteParserErrCode = 2
const RouteParserInitialStackSize = 16



func Parse(exp string) ([]*Route, error) {
	l := new(RouteLexer)
	lexerErrors := make([]string, 0)
	l.lexerErrors = &lexerErrors
	l.s = new(RouteScanner)
	l.s.Init(exp)
	//l.Init(strings.NewReader(exp))
	RouteParserParse(l)
	if len(lexerErrors) > 0 {
		return nil, errors.New(strings.Join(lexerErrors, "\n"))
	}
	return l.result, nil
}


var RouteParserExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const RouteParserPrivate = 57344

const RouteParserLast = 893

var RouteParserAct = [...]int{

	81, 52, 80, 44, 37, 153, 82, 32, 38, 146,
	33, 34, 145, 6, 104, 60, 3, 78, 16, 31,
	128, 78, 149, 93, 94, 95, 96, 97, 90, 91,
	92, 35, 77, 36, 147, 89, 77, 79, 98, 39,
	83, 84, 127, 85, 86, 42, 5, 4, 14, 88,
	40, 131, 103, 99, 87, 102, 105, 3, 54, 106,
	107, 108, 109, 110, 111, 112, 113, 114, 115, 116,
	117, 118, 119, 120, 121, 122, 123, 124, 51, 126,
	43, 10, 11, 12, 8, 130, 29, 30, 25, 22,
	134, 135, 136, 137, 138, 139, 140, 141, 142, 143,
	60, 28, 68, 53, 144, 148, 78, 73, 74, 69,
	70, 72, 71, 63, 64, 24, 21, 19, 23, 125,
	20, 77, 60, 59, 68, 101, 100, 3, 78, 15,
	76, 75, 152, 61, 62, 63, 64, 65, 66, 67,
	45, 58, 41, 77, 18, 46, 17, 164, 9, 56,
	55, 165, 73, 74, 69, 70, 72, 71, 49, 47,
	48, 1, 0, 50, 0, 0, 166, 60, 59, 68,
	0, 0, 0, 78, 0, 76, 75, 0, 61, 62,
	63, 64, 65, 66, 67, 2, 58, 0, 77, 0,
	7, 13, 129, 73, 74, 69, 70, 72, 71, 0,
	26, 0, 27, 0, 0, 0, 0, 0, 60, 59,
	68, 0, 0, 0, 78, 0, 76, 75, 0, 61,
	62, 63, 64, 65, 66, 67, 0, 58, 0, 77,
	0, 0, 132, 73, 74, 69, 70, 72, 71, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 60, 59,
	68, 0, 0, 0, 78, 0, 76, 75, 0, 61,
	62, 63, 64, 65, 66, 67, 0, 58, 0, 77,
	162, 73, 74, 69, 70, 72, 71, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 60, 59, 68, 0,
	0, 0, 78, 0, 76, 75, 0, 61, 62, 63,
	64, 65, 66, 67, 0, 58, 0, 77, 161, 73,
	74, 69, 70, 72, 71, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 60, 59, 68, 0, 0, 0,
	78, 0, 76, 75, 0, 61, 62, 63, 64, 65,
	66, 67, 0, 58, 0, 77, 160, 73, 74, 69,
	70, 72, 71, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 60, 59, 68, 0, 0, 0, 78, 0,
	76, 75, 0, 61, 62, 63, 64, 65, 66, 67,
	0, 58, 0, 77, 159, 73, 74, 69, 70, 72,
	71, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	60, 59, 68, 0, 0, 0, 78, 0, 76, 75,
	0, 61, 62, 63, 64, 65, 66, 67, 0, 58,
	0, 77, 158, 73, 74, 69, 70, 72, 71, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 60, 59,
	68, 0, 0, 0, 78, 0, 76, 75, 0, 61,
	62, 63, 64, 65, 66, 67, 0, 58, 0, 77,
	157, 73, 74, 69, 70, 72, 71, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 60, 59, 68, 0,
	0, 0, 78, 0, 76, 75, 0, 61, 62, 63,
	64, 65, 66, 67, 0, 58, 0, 77, 156, 73,
	74, 69, 70, 72, 71, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 60, 59, 68, 0, 0, 0,
	78, 0, 76, 75, 0, 61, 62, 63, 64, 65,
	66, 67, 0, 58, 0, 77, 155, 73, 74, 69,
	70, 72, 71, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 60, 59, 68, 0, 0, 0, 78, 0,
	76, 75, 0, 61, 62, 63, 64, 65, 66, 67,
	0, 58, 0, 77, 154, 73, 74, 69, 70, 72,
	71, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	60, 59, 68, 0, 0, 0, 78, 163, 76, 75,
	0, 61, 62, 63, 64, 65, 66, 67, 0, 58,
	0, 77, 73, 74, 69, 70, 72, 71, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 60, 59, 68,
	0, 0, 0, 78, 151, 76, 75, 0, 61, 62,
	63, 64, 65, 66, 67, 0, 58, 0, 77, 73,
	74, 69, 70, 72, 71, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 60, 59, 68, 0, 0, 0,
	78, 0, 76, 75, 150, 61, 62, 63, 64, 65,
	66, 67, 0, 58, 0, 77, 73, 74, 69, 70,
	72, 71, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 60, 59, 68, 0, 0, 133, 78, 0, 76,
	75, 0, 61, 62, 63, 64, 65, 66, 67, 0,
	58, 0, 77, 73, 74, 69, 70, 72, 71, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 60, 59,
	68, 0, 0, 57, 78, 0, 76, 75, 0, 61,
	62, 63, 64, 65, 66, 67, 0, 58, 0, 77,
	73, 74, 69, 70, 72, 71, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 60, 0, 68, 0, 0,
	0, 78, 0, 76, 75, 0, 61, 62, 63, 64,
	65, 66, 67, 0, 58, 0, 77, 69, 70, 72,
	71, 37, 0, 0, 32, 38, 0, 33, 34, 0,
	60, 0, 68, 0, 0, 0, 78, 0, 76, 75,
	0, 61, 62, 63, 64, 65, 66, 67, 35, 58,
	36, 77, 0, 0, 0, 0, 39, 60, 0, 68,
	0, 0, 42, 78, 0, 0, 0, 40, 61, 62,
	63, 64, 65, 66, 67, 60, 58, 68, 77, 0,
	0, 78, 0, 0, 0, 0, 61, 62, 63, 64,
	0, 60, 0, 68, 58, 0, 77, 78, 0, 0,
	0, 0, 61, 62, 63, 64, 0, 0, 0, 0,
	0, 0, 77,
}
var RouteParserPact = [...]int{

	123, -1000, 11, -28, 53, 12, 125, -1000, -1000, -18,
	113, 85, 84, -1000, 123, -1000, 53, 57, 797, 49,
	-1000, 154, 47, -1000, 99, 27, -1000, -1000, -1000, 85,
	84, 711, -1000, -1000, -1000, 99, 797, -44, -1000, 797,
	797, -1000, 797, 797, 22, 154, 5, 122, 121, -1000,
	-1000, 154, 20, -23, 99, -1000, -1000, -1000, 797, 797,
	797, 797, 797, 797, 797, 797, 797, 797, 797, 797,
	797, 797, 797, 797, 797, 797, 797, 115, 797, 10,
	-14, 140, 0, -16, 181, -16, 674, -1000, -1000, 797,
	797, 797, 797, 797, 797, 797, 797, 797, 797, 100,
	-37, -40, 2, -1000, 797, -10, 844, 637, -16, 73,
	73, -12, -12, 828, 828, 828, -12, 810, 810, 810,
	810, 783, 783, 810, 810, -1000, 600, -1000, -1000, 797,
	-46, -1000, -1000, -1000, 525, 487, 449, 411, 373, 335,
	297, 259, 221, 563, -1000, -1000, -1000, -1000, 95, -1000,
	797, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 99, 748, -1000,
}
var RouteParserPgo = [...]int{

	0, 161, 185, 84, 148, 146, 120, 118, 3, 145,
	142, 140, 0, 1, 2,
}
var RouteParserR1 = [...]int{

	0, 1, 1, 1, 1, 2, 2, 3, 3, 4,
	4, 4, 4, 4, 5, 5, 5, 7, 7, 7,
	6, 6, 8, 8, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 11, 11, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 13, 13, 14, 14,
	10, 10, 10, 9, 9, 9, 9,
}
var RouteParserR2 = [...]int{

	0, 3, 3, 4, 3, 1, 3, 3, 3, 2,
	2, 2, 4, 4, 3, 4, 1, 3, 4, 1,
	3, 4, 1, 2, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 3, 3, 1, 1, 1, 3, 3,
	4, 3, 1, 3, 5, 2, 3, 1, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 2, 3, 3, 3, 5, 1, 3,
	3, 4, 1, 1, 1, 4, 3,
}
var RouteParserChk = [...]int{

	-1000, -1, -2, 4, 36, 35, 41, -2, -3, -4,
	28, 29, 30, -2, 36, 4, 36, -5, 31, 4,
	-6, 31, 4, -7, 31, 4, -2, -2, -3, 29,
	30, -12, 7, 10, 11, 31, 33, 4, 8, 39,
	50, -10, 45, 31, -8, -11, -9, 5, 6, 4,
	9, 31, -13, 4, 31, -6, -7, 32, 46, 28,
	27, 38, 39, 40, 41, 42, 43, 44, 29, 14,
	15, 17, 16, 12, 13, 36, 35, 48, 33, -13,
	-14, -12, 50, -12, -12, -12, -12, 32, -8, 30,
	23, 24, 25, 18, 19, 20, 21, 22, 33, 48,
	4, 4, -8, 32, 37, -13, -12, -12, -12, -12,
	-12, -12, -12, -12, -12, -12, -12, -12, -12, -12,
	-12, -12, -12, -12, -12, 4, -12, 32, 34, 52,
	-14, 51, 51, 32, -12, -12, -12, -12, -12, -12,
	-12, -12, -12, -12, 4, 49, 49, 32, -12, 32,
	37, 34, -14, 51, 49, 49, 49, 49, 49, 49,
	49, 49, 49, 34, 52, -12, -13,
}
var RouteParserDef = [...]int{

	0, -2, 0, 5, 0, 0, 0, 1, 4, 0,
	0, 0, 0, 2, 0, 6, 0, 9, 0, 16,
	10, 0, 0, 11, 0, 19, 3, 7, 8, 0,
	0, 0, 35, 36, 37, 0, 0, 72, 42, 0,
	0, 47, 0, 0, 0, 22, 0, 0, 0, 73,
	74, 0, 0, 0, 0, 12, 13, 14, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 68, 0, 45, 0, 63, 0, 20, 23, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 17, 0, 0, 43, 0, 48, 49,
	50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	60, 61, 62, 64, 65, 70, 0, 38, 39, 0,
	0, 41, 46, 15, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 76, 33, 34, 21, 66, 18,
	0, 71, 69, 40, 24, 25, 26, 27, 28, 29,
	30, 31, 32, 75, 0, 44, 67,
}
var RouteParserTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 45, 3, 3, 3, 29, 42, 3,
	50, 51, 40, 38, 52, 39, 48, 41, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 37, 49,
	35, 30, 36, 28, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 33, 3, 34, 44, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 31, 43, 32, 46,
}
var RouteParserTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 47,
}
var RouteParserTok3 = [...]int{
	0,
}

var RouteParserErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}



/*	parser for yacc output	*/

var (
	RouteParserDebug        = 0
	RouteParserErrorVerbose = false
)

type RouteParserLexer interface {
	Lex(lval *RouteParserSymType) int
	Error(s string)
}

type RouteParserParser interface {
	Parse(RouteParserLexer) int
	Lookahead() int
}

type RouteParserParserImpl struct {
	lval  RouteParserSymType
	stack [RouteParserInitialStackSize]RouteParserSymType
	char  int
}

func (p *RouteParserParserImpl) Lookahead() int {
	return p.char
}

func RouteParserNewParser() RouteParserParser {
	return &RouteParserParserImpl{}
}

const RouteParserFlag = -1000

func RouteParserTokname(c int) string {
	if c >= 1 && c-1 < len(RouteParserToknames) {
		if RouteParserToknames[c-1] != "" {
			return RouteParserToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func RouteParserStatname(s int) string {
	if s >= 0 && s < len(RouteParserStatenames) {
		if RouteParserStatenames[s] != "" {
			return RouteParserStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func RouteParserErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !RouteParserErrorVerbose {
		return "syntax error"
	}

	for _, e := range RouteParserErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + RouteParserTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := RouteParserPact[state]
	for tok := TOKSTART; tok-1 < len(RouteParserToknames); tok++ {
		if n := base + tok; n >= 0 && n < RouteParserLast && RouteParserChk[RouteParserAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if RouteParserDef[state] == -2 {
		i := 0
		for RouteParserExca[i] != -1 || RouteParserExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; RouteParserExca[i] >= 0; i += 2 {
			tok := RouteParserExca[i]
			if tok < TOKSTART || RouteParserExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if RouteParserExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += RouteParserTokname(tok)
	}
	return res
}

func RouteParserlex1(lex RouteParserLexer, lval *RouteParserSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = RouteParserTok1[0]
		goto out
	}
	if char < len(RouteParserTok1) {
		token = RouteParserTok1[char]
		goto out
	}
	if char >= RouteParserPrivate {
		if char < RouteParserPrivate+len(RouteParserTok2) {
			token = RouteParserTok2[char-RouteParserPrivate]
			goto out
		}
	}
	for i := 0; i < len(RouteParserTok3); i += 2 {
		token = RouteParserTok3[i+0]
		if token == char {
			token = RouteParserTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = RouteParserTok2[1] /* unknown char */
	}
	if RouteParserDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", RouteParserTokname(token), uint(char))
	}
	return char, token
}

func RouteParserParse(RouteParserlex RouteParserLexer) int {
	return RouteParserNewParser().Parse(RouteParserlex)
}

func (RouteParserrcvr *RouteParserParserImpl) Parse(RouteParserlex RouteParserLexer) int {
	var RouteParsern int
	var RouteParserVAL RouteParserSymType
	var RouteParserDollar []RouteParserSymType
	_ = RouteParserDollar // silence set and not used
	RouteParserS := RouteParserrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	RouteParserstate := 0
	RouteParserrcvr.char = -1
	RouteParsertoken := -1 // RouteParserrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		RouteParserstate = -1
		RouteParserrcvr.char = -1
		RouteParsertoken = -1
	}()
	RouteParserp := -1
	goto RouteParserstack

ret0:
	return 0

ret1:
	return 1

RouteParserstack:
	/* put a state and value onto the stack */
	if RouteParserDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", RouteParserTokname(RouteParsertoken), RouteParserStatname(RouteParserstate))
	}

	RouteParserp++
	if RouteParserp >= len(RouteParserS) {
		nyys := make([]RouteParserSymType, len(RouteParserS)*2)
		copy(nyys, RouteParserS)
		RouteParserS = nyys
	}
	RouteParserS[RouteParserp] = RouteParserVAL
	RouteParserS[RouteParserp].yys = RouteParserstate

RouteParsernewstate:
	RouteParsern = RouteParserPact[RouteParserstate]
	if RouteParsern <= RouteParserFlag {
		goto RouteParserdefault /* simple state */
	}
	if RouteParserrcvr.char < 0 {
		RouteParserrcvr.char, RouteParsertoken = RouteParserlex1(RouteParserlex, &RouteParserrcvr.lval)
	}
	RouteParsern += RouteParsertoken
	if RouteParsern < 0 || RouteParsern >= RouteParserLast {
		goto RouteParserdefault
	}
	RouteParsern = RouteParserAct[RouteParsern]
	if RouteParserChk[RouteParsern] == RouteParsertoken { /* valid shift */
		RouteParserrcvr.char = -1
		RouteParsertoken = -1
		RouteParserVAL = RouteParserrcvr.lval
		RouteParserstate = RouteParsern
		if Errflag > 0 {
			Errflag--
		}
		goto RouteParserstack
	}

RouteParserdefault:
	/* default state action */
	RouteParsern = RouteParserDef[RouteParserstate]
	if RouteParsern == -2 {
		if RouteParserrcvr.char < 0 {
			RouteParserrcvr.char, RouteParsertoken = RouteParserlex1(RouteParserlex, &RouteParserrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if RouteParserExca[xi+0] == -1 && RouteParserExca[xi+1] == RouteParserstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			RouteParsern = RouteParserExca[xi+0]
			if RouteParsern < 0 || RouteParsern == RouteParsertoken {
				break
			}
		}
		RouteParsern = RouteParserExca[xi+1]
		if RouteParsern < 0 {
			goto ret0
		}
	}
	if RouteParsern == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			RouteParserlex.Error(RouteParserErrorMessage(RouteParserstate, RouteParsertoken))
			Nerrs++
			if RouteParserDebug >= 1 {
				__yyfmt__.Printf("%s", RouteParserStatname(RouteParserstate))
				__yyfmt__.Printf(" saw %s\n", RouteParserTokname(RouteParsertoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for RouteParserp >= 0 {
				RouteParsern = RouteParserPact[RouteParserS[RouteParserp].yys] + RouteParserErrCode
				if RouteParsern >= 0 && RouteParsern < RouteParserLast {
					RouteParserstate = RouteParserAct[RouteParsern] /* simulate a shift of "error" */
					if RouteParserChk[RouteParserstate] == RouteParserErrCode {
						goto RouteParserstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if RouteParserDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", RouteParserS[RouteParserp].yys)
				}
				RouteParserp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if RouteParserDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", RouteParserTokname(RouteParsertoken))
			}
			if RouteParsertoken == RouteParserEofCode {
				goto ret1
			}
			RouteParserrcvr.char = -1
			RouteParsertoken = -1
			goto RouteParsernewstate /* try again in the same state */
		}
	}

	/* reduction by production RouteParsern */
	if RouteParserDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", RouteParsern, RouteParserStatname(RouteParserstate))
	}

	RouteParsernt := RouteParsern
	RouteParserpt := RouteParserp
	_ = RouteParserpt // guard against "declared and not used"

	RouteParserp -= RouteParserR2[RouteParsern]
	// RouteParserp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if RouteParserp+1 >= len(RouteParserS) {
		nyys := make([]RouteParserSymType, len(RouteParserS)*2)
		copy(nyys, RouteParserS)
		RouteParserS = nyys
	}
	RouteParserVAL = RouteParserS[RouteParserp+1]

	/* consult goto table to find next state */
	RouteParsern = RouteParserR1[RouteParsern]
	RouteParserg := RouteParserPgo[RouteParsern]
	RouteParserj := RouteParserg + RouteParserS[RouteParserp].yys + 1

	if RouteParserj >= RouteParserLast {
		RouteParserstate = RouteParserAct[RouteParserg]
	} else {
		RouteParserstate = RouteParserAct[RouteParserj]
		if RouteParserChk[RouteParserstate] != -RouteParsern {
			RouteParserstate = RouteParserAct[RouteParserg]
		}
	}
	// dummy call; replaced with literal code
	switch RouteParsernt {

	case 1:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			// fmt.Println("C")
			RouteParserVAL.routes = nil
			if l, ok := RouteParserlex.(*RouteLexer); ok {
				l.result = []*Route{&Route{Source: RouteParserDollar[1].node, Dest: RouteParserDollar[3].node}}
			}
		}
	case 2:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			// fmt.Println("B")
			RouteParserVAL.routes = nil
			if l, ok := RouteParserlex.(*RouteLexer); ok {
				l.result = []*Route{&Route{Source: RouteParserDollar[3].node, Dest: RouteParserDollar[1].node}}
			}
		}
	case 3:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.routes = nil
			if l, ok := RouteParserlex.(*RouteLexer); ok {
				l.result = []*Route{
					&Route{Source: RouteParserDollar[1].node, Dest: RouteParserDollar[4].node},
					&Route{Source: RouteParserDollar[4].node, Dest: RouteParserDollar[1].node}}
			}
		}
	case 4:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.routes = nil
			// fmt.Println("A")
			if l, ok := RouteParserlex.(*RouteLexer); ok {
				l.result = []*Route{
					&Route{
						Source:     RouteParserDollar[1].node,
						Dest:       RouteParserDollar[3].transforms.Dest,
						Transforms: RouteParserDollar[3].transforms.Transforms}}
			}
		}
	case 5:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.node = RouteParserDollar[1].token.literal
		}
	case 6:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.node = fmt.Sprintf("%s/%s", RouteParserDollar[1].token.literal, RouteParserDollar[3].token.literal)
		}
	case 7:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.transforms = &Route{
				Dest:       RouteParserDollar[3].node,
				Transforms: []*Transform{RouteParserDollar[1].transform}}
		}
	case 8:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.transforms = &Route{
				Dest:       RouteParserDollar[3].transforms.Dest,
				Transforms: append([]*Transform{RouteParserDollar[1].transform}, RouteParserDollar[3].transforms.Transforms...)}
		}
	case 9:
		RouteParserDollar = RouteParserS[RouteParserpt-2 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:             TR_FILTER,
				CommandCondition: RouteParserDollar[2].transform.CommandCondition,
				Condition:        RouteParserDollar[2].transform.Condition}
		}
	case 10:
		RouteParserDollar = RouteParserS[RouteParserpt-2 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_EDIT,
				CommandReplace: RouteParserDollar[2].transform.CommandReplace,
				Script:         RouteParserDollar[2].transform.Script}
		}
	case 11:
		RouteParserDollar = RouteParserS[RouteParserpt-2 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_REPLACE,
				CommandReplace: RouteParserDollar[2].transform.CommandReplace,
				Replace:        RouteParserDollar[2].transform.Replace}
		}
	case 12:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:             TR_COND_EDIT,
				CommandCondition: RouteParserDollar[2].transform.CommandCondition,
				CommandReplace:   RouteParserDollar[4].transform.CommandReplace,
				Condition:        RouteParserDollar[2].transform.Condition,
				Script:           RouteParserDollar[4].transform.Script}
		}
	case 13:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:             TR_COND_REPLACE,
				CommandCondition: RouteParserDollar[2].transform.CommandCondition,
				Condition:        RouteParserDollar[2].transform.Condition,
				CommandReplace:   RouteParserDollar[4].transform.CommandReplace,
				Replace:          RouteParserDollar[4].transform.Replace}
		}
	case 14:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:             TR_FILTER,
				CommandCondition: "",
				Condition:        RouteParserDollar[2].expression}
		}
	case 15:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:             TR_EDIT,
				CommandCondition: RouteParserDollar[1].token.literal,
				Condition:        RouteParserDollar[3].expression}
		}
	case 16:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:             TR_EDIT,
				CommandCondition: RouteParserDollar[1].token.literal,
				Condition:        nil}
		}
	case 17:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_REPLACE,
				CommandReplace: "",
				Replace:        RouteParserDollar[2].expression}
		}
	case 18:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_REPLACE,
				CommandReplace: RouteParserDollar[1].token.literal,
				Replace:        RouteParserDollar[3].expression}
		}
	case 19:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_REPLACE,
				CommandReplace: RouteParserDollar[1].token.literal,
				Replace:        nil}
		}
	case 20:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_EDIT,
				CommandReplace: "",
				Script:         RouteParserDollar[2].script}
		}
	case 21:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.transform = &Transform{
				Type:           TR_EDIT,
				CommandReplace: RouteParserDollar[1].token.literal,
				Script:         RouteParserDollar[3].script}
		}
	case 22:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.script = []*Statement{RouteParserDollar[1].statement}
		}
	case 23:
		RouteParserDollar = RouteParserS[RouteParserpt-2 : RouteParserpt+1]

		{
			RouteParserVAL.script = append([]*Statement{RouteParserDollar[1].statement}, RouteParserDollar[2].script...)
		}
	case 24:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, RouteParserDollar[3].expression)
		}
	case 25:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_BITWISEAND,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 26:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_BITWISEOR,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 27:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_BITWISEXOR,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 28:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_PLUS,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 29:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_MINUS,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 30:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_MUL,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 31:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_DIV,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 32:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeAssignmentStatement(RouteParserDollar[1].writeable, &Expression{
				Operation: OP_MOD,
				Args:      []*Expression{RouteParserDollar[1].writeable.ToExpression(), RouteParserDollar[3].expression}})
		}
	case 33:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeDeclarationStatement(RouteParserDollar[2].token.literal)
		}
	case 34:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.statement = MakeDeletionStatement(RouteParserDollar[2].token.literal)
		}
	case 35:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			x, _ := strconv.ParseFloat(RouteParserDollar[1].token.literal, 64)
			RouteParserVAL.expression = &Expression{
				Operation: OP_NUM,
				Value:     value.MakeFloatValue(x)}
		}
	case 36:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_BOOL,
				Value:     value.MakeBoolValue(true)}
		}
	case 37:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_BOOL,
				Value:     value.MakeBoolValue(false)}
		}
	case 38:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = RouteParserDollar[2].expression
		}
	case 39:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = RouteParserDollar[2].expression
		}
	case 40:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_CALL,
				Args: []*Expression{
					MakeNameExpression(RouteParserDollar[1].token.literal),
					RouteParserDollar[3].expression}}
		}
	case 41:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_CALL,
				Args: []*Expression{
					MakeNameExpression(RouteParserDollar[1].token.literal),
					&Expression{
						Operation: OP_LIST,
						Args:      []*Expression{}}}}
		}
	case 42:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_STRING,
				Value:     &value.Value{Type: value.VAL_STRING, StringVal: RouteParserDollar[1].token.literal}}
		}
	case 43:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_MATCH,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 44:
		RouteParserDollar = RouteParserS[RouteParserpt-5 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_TERNARY,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression, RouteParserDollar[5].expression}}
		}
	case 45:
		RouteParserDollar = RouteParserS[RouteParserpt-2 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_UMINUS,
				Args:      []*Expression{RouteParserDollar[2].expression}}
		}
	case 46:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = RouteParserDollar[2].expression
		}
	case 47:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.expression = RouteParserDollar[1].expression
		}
	case 48:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_EXP,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 49:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_PLUS,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 50:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_MINUS,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 51:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_MUL,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 52:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_DIV,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 53:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_BITWISEAND,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 54:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_BITWISEOR,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 55:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_BITWISEXOR,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 56:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_MOD,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 57:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_EQ,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 58:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_NE,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 59:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_GE,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 60:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_LE,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 61:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_AND,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 62:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_OR,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 63:
		RouteParserDollar = RouteParserS[RouteParserpt-2 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_NOT,
				Args:      []*Expression{RouteParserDollar[2].expression}}
		}
	case 64:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_GT,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 65:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_LT,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 66:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_MAP,
				Args: []*Expression{
					MakeNameExpression(RouteParserDollar[1].token.literal),
					RouteParserDollar[3].expression}}
		}
	case 67:
		RouteParserDollar = RouteParserS[RouteParserpt-5 : RouteParserpt+1]

		{
			args := []*Expression{
				MakeNameExpression(RouteParserDollar[1].token.literal),
				RouteParserDollar[3].expression}
			RouteParserVAL.expression = &Expression{
				Operation: OP_MAP,
				Args:      append(args, RouteParserDollar[5].expression.Args...)}
		}
	case 68:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_LIST,
				Args:      []*Expression{RouteParserDollar[1].expression}}
		}
	case 69:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			args := []*Expression{RouteParserDollar[1].expression}
			RouteParserVAL.expression = &Expression{
				Operation: OP_LIST,
				Args:      append(args, RouteParserDollar[3].expression.Args...)}
		}
	case 70:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_MAPVAR,
				Args:      []*Expression{RouteParserDollar[1].expression, MakeNameExpression(RouteParserDollar[3].token.literal)}}
		}
	case 71:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserVAL.expression = &Expression{
				Operation: OP_LISTVAR,
				Args:      []*Expression{RouteParserDollar[1].expression, RouteParserDollar[3].expression}}
		}
	case 72:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.expression = MakeVarExpression(RouteParserDollar[1].token.literal)
		}
	case 73:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.writeable = &WriteableValue{
				Base: RouteParserDollar[1].token.literal,
				Path: []PathEntry{}}
		}
	case 74:
		RouteParserDollar = RouteParserS[RouteParserpt-1 : RouteParserpt+1]

		{
			RouteParserVAL.writeable = &WriteableValue{
				Base: "this",
				Path: []PathEntry{}}
		}
	case 75:
		RouteParserDollar = RouteParserS[RouteParserpt-4 : RouteParserpt+1]

		{
			RouteParserDollar[1].writeable.Path = append(RouteParserDollar[1].writeable.Path, PathEntry{Type: PATH_LIST, ListIndex: RouteParserDollar[3].expression})
			RouteParserVAL.writeable = RouteParserDollar[1].writeable
		}
	case 76:
		RouteParserDollar = RouteParserS[RouteParserpt-3 : RouteParserpt+1]

		{
			RouteParserDollar[1].writeable.Path = append(RouteParserDollar[1].writeable.Path, PathEntry{Type: PATH_MAP, MapKey: RouteParserDollar[3].token.literal})
			RouteParserVAL.writeable = RouteParserDollar[1].writeable
		}
	}
	goto RouteParserstack /* stack new state and value */
}
