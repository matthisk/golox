package parser

import (
	"fmt"

	"github.com/matthisk/lox/lexer"
)

type Errors struct {
	errs []*Error
}

func NewErrors(errs []*Error) *Errors {
	return &Errors{errs: errs}
}

func (p *Errors) Error() string {
	return fmt.Sprintf("found %d errors while parsing", len(p.errs))
}

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
