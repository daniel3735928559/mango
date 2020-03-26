package valuetype

import (
	"fmt"
)

const (
	EOF     = -1
	UNKNOWN = 0
)

var (
	keywords = map[string]int{"string":STR, "num":NUM, "bool": BOOL, "true": TRUE, "false": FALSE, "oneof": ONEOF, "any": ANY}
	charsyms = map[rune]int {
		'=':'=',
		',':',',
		'{':'{',
		'}':'}',
		'[':'[',
		']':']',
		'(':'(',
		')':')',
		';':';',
		':':':',
		'*':'*'}
)

type Position struct {
	Line   int
	Column int
}

type ValueTypeScanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *ValueTypeScanner) Init(src string) {
	s.src = []rune(src)
}

func (s *ValueTypeScanner) Scan() (tok int, lit string, pos Position) {
	s.skipWhiteSpace()
	pos = s.position()
	switch ch := s.peek(); {
	case isLetter(ch):
		lit = s.scanIdentifier()
		//fmt.Println("scanning kwds",lit)
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

func (s *ValueTypeScanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *ValueTypeScanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *ValueTypeScanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *ValueTypeScanner) position() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *ValueTypeScanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *ValueTypeScanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}


func (s *ValueTypeScanner) scanSym() (int, string) {
	fc := s.peek()
	if fc == -1 {
		return EOF, ""
	} else if tok, ok := charsyms[fc]; ok {
		s.next()
		lit := string(fc)
		return tok, lit
	}
	return -1, "error"
}

func (s *ValueTypeScanner) scanNumber() string {
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

func (s *ValueTypeScanner) scanString() string {
	var ret []rune
	s.next()
	for s.peek() != '"' || (len(ret) > 0 && ret[len(ret)-1] == '\\') {
		ret = append(ret, s.peek())
		s.next()
	}
	s.next()
	return string(ret)
}
