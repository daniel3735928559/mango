package valuetype

import (
	"fmt"
)

type ValueTypeLexer struct {
	s *ValueTypeScanner
	lexerErrors *[]string
	result *ValueType
}

func (l *ValueTypeLexer) Lex(lval *ValueTypeParserSymType) int {
	tok, lit, pos := l.s.Scan()
	if tok == EOF {
		return 0
	}
	lval.token = Token{token: tok, literal: lit, position: pos}
	//fmt.Println("Lexed",tok,lit,pos)
	return tok
}

func (l *ValueTypeLexer) Error(e string) {
	le := *(l.lexerErrors)
	*(l.lexerErrors) = append(le, e)
	fmt.Println("ERROR",e)
}
