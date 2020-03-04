package router

import (
	"fmt"
)

type RouteLexer struct {
	s *RouteScanner
	lexerErrors *[]string
	result []*Route
}

func (l *RouteLexer) Lex(lval *RouteParserSymType) int {
	tok, lit, pos := l.s.Scan()
	if tok == EOF {
		return 0
	}
	lval.token = Token{token: tok, literal: lit, position: pos}
	fmt.Println("Lexed",tok,lit,pos)
	return tok
}

func (l *RouteLexer) Error(e string) {
	le := *(l.lexerErrors)
	*(l.lexerErrors) = append(le, e)
	fmt.Println("ERROR",e)
}
