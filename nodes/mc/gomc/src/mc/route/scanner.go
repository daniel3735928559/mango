package router

import (
	"fmt"
)

const (
	EOF     = -1
	UNKNOWN = 0
)

var (
	keywords = map[string]int{"del": DEL, "var": VAR, "true":TRUE, "false":FALSE}
	syms = map[string]int{
		"==":EQ,
		"!=":NE,
		"<=":LE,
		">=":GE,
		"+=":PE,
		"-=":ME,
		"*=":TE,
		"/=":DE,
		"%=":RE,
		"&=":AE,
		"|=":OE,
		"^=":XE,
		"~=":SUB,
		"&&":AND,
		"||":OR,
		"**":EXP,
		"$_":THIS}
	charsyms = map[rune]int {
		'.':'.',
		',':',',
		'{':'{',
		'}':'}',
		'[':'[',
		']':']',
		'(':'(',
		')':')',
		';':';',
		':':':',
		'+':'+',
		'-':'-',
		'*':'*',
		'/':'/',
		'%':'%',
		'?':'?',
		'~':'~',
		'!':'!',
		'&':'&',
		'|':'|',
		'^':'^',
		'$':-1} // -1 means this is a prefix of a longer symbol
)

type Position struct {
	Line   int
	Column int
}

type RouteScanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *RouteScanner) Init(src string) {
	s.src = []rune(src)
}

func (s *RouteScanner) Scan() (tok int, lit string, pos Position) {
	s.skipWhiteSpace()
	pos = s.position()
	switch ch := s.peek(); {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if keyword, ok := keywords[lit]; ok {
			tok = keyword
		} else {
			tok = IDENT
		}
	case isDigit(ch):
		tok, lit = NUMBER, s.scanNumber()
	case ch == '"':
		tok, lit = STRING, s.scanString()
		//fmt.Println("STRING",lit)
	case ch == '=' || ch == '<' || ch == '>':
		tok, lit = s.scanTest()
		//fmt.Println("TEST",tok,lit)
	default:
		tok, lit = s.scanSym()
	}
	return
}

// ========================================

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (s *RouteScanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *RouteScanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *RouteScanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *RouteScanner) position() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *RouteScanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *RouteScanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}


func (s *RouteScanner) scanSym() (int, string) {
	fc := s.peek()
	if fc == -1 {
		return EOF, ""
	} else if tok, ok := charsyms[fc]; ok {
		s.next()
		lit := string(fc)
		long_lit := lit + string(s.peek())
		if long_tok, long_ok := syms[long_lit]; long_ok {
			s.next()
			return long_tok, long_lit
		}
		return tok, lit
	}
	return -1, "error"
}

func (s *RouteScanner) scanTest() (int, string) {
	fc := s.peek()
	if fc == '=' || fc == '<' || fc == '>' {
		s.next()
		//fmt.Printf("peek, fc:%s,p:%s\n",string(fc),string(s.peek()))
		if s.peek() == '=' {
			s.next()
			if fc == '=' {
				return EQ, "=="
			} else if fc == '<' {
				return LE, "<="
			} else if fc == '>' {
				return GE, ">="
			}
		}
	}
	return int(fc), string(fc)
}
	
func (s *RouteScanner) scanNumber() string {
	var ret []rune
	hasDecimal := false
	for isDigit(s.peek()) || (s.peek() == '.' && !hasDecimal) {
		ret = append(ret, s.peek())
		if s.peek() == '.' {
			hasDecimal = true
		}
		s.next()
	}
	fmt.Println("SN",string(ret))
	return string(ret)
}

func (s *RouteScanner) scanString() string {
	var ret []rune
	s.next()
	for s.peek() != '"' || (len(ret) > 0 && ret[len(ret)-1] == '\\') {
		ret = append(ret, s.peek())
		s.next()
	}
	s.next()
	return string(ret)
}
