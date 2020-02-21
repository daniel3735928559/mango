package router

import (
	"fmt"
)

type RouteLexer struct {
	s *RouteScanner
	result []*Route
}

func (l *RouteLexer) Lex(lval *yySymType) int {
	tok, lit, pos := l.s.Scan()
	if tok == EOF {
		return 0
	}
	lval.token = Token{token: tok, literal: lit, position: pos}
	fmt.Println("Lexed",tok,lit,pos)
	return tok
}

func (l *RouteLexer) Error(e string) {
	fmt.Println("ERROR",e)
}
