package value

import (
	"fmt"
)

type ValueLexer struct {
	s *ValueScanner
	lexerErrors *[]string
	result *Value
}

func (l *ValueLexer) Lex(lval *ValueParserSymType) int {
	tok, lit, pos := l.s.Scan()
	if tok == EOF {
		return 0
	}
	lval.token = Token{token: tok, literal: lit, position: pos}
	fmt.Println("Lexed",tok,lit,pos)
	return tok
}

func (l *ValueLexer) Error(e string) {
	le := *(l.lexerErrors)
	*(l.lexerErrors) = append(le, e)
	fmt.Println("ERROR",e)
}
