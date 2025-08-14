package parser

import "github.com/matthisk/lox/lexer"

type Error struct {
	pos lexer.Pos
	at  lexer.Token
	msg string
}

func NewError(tok lexer.Token, err string) *Error {
	return &Error{
		pos: tok.StartPos,
		at:  tok,
		msg: err,
	}
}

func (p *Error) Error() string {
	return p.msg
}
