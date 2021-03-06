package value

import (
	//"fmt"
	"strconv"
)

const (
	EOF     = -1
	UNKNOWN = 0
)

var (
	keywords = map[string]int{"true": TRUE, "false": FALSE}
	charsyms = map[rune]int {
		',':',',
		'{':'{',
		'}':'}',
		'[':'[',
		']':']',
		'(':'(',
		')':')',
		':':':'}
)

type Position struct {
	Line   int
	Column int
}

type ValueScanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *ValueScanner) Init(src string) {
	s.src = []rune(src)
}

func (s *ValueScanner) Scan() (tok int, lit string, pos Position) {
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

func (s *ValueScanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *ValueScanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *ValueScanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *ValueScanner) position() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *ValueScanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *ValueScanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}


func (s *ValueScanner) scanSym() (int, string) {
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

func (s *ValueScanner) scanNumber() string {
	var ret []rune
	hasDecimal := false
	for isDigit(s.peek()) || (s.peek() == '.' && !hasDecimal) {
		ret = append(ret, s.peek())
		if s.peek() == '.' {
			hasDecimal = true
		}
		s.next()
	}
	return string(ret)
}

func (s *ValueScanner) scanString() string {
	var ret []rune
	ret = append(ret, s.peek())
	s.next()
	for s.peek() != '"' || (len(ret) > 0 && ret[len(ret)-1] == '\\') {
		ret = append(ret, s.peek())
		s.next()
	}
	ret = append(ret, s.peek())
	s.next()
	ans, err := strconv.Unquote(string(ret))
	if err != nil {
		return "FAILED UNQUOTE:"+string(ret)
	}
	return ans
}
