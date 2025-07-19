package parser

import (
	"github.com/matthisk/lox/lexer"
)

/**
Grammar:

expression → literal
		   | unary
		   | binary
		   | grouping ;

literal        → NUMBER | STRING | "true" | "false" | "nil" ;
grouping       → "(" expression ")" ;
unary          → ( "-" | "!" ) expression ;
binary         → expression operator expression ;
operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
               | "+"  | "-"  | "*" | "/" ;
*/

type Visitor interface {
	VisitBinary(node *Binary) (interface{}, error)
	VisitLiteral(node *Literal) (interface{}, error)
	VisitUnary(node *Unary) (interface{}, error)
	VisitComma(node *Comma) (interface{}, error)
	VisitGrouping(node *Grouping) (interface{}, error)
	VisitTernary(b *Ternary) (interface{}, error)
}

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type BaseNode struct {
	pos lexer.Pos
}

type Binary struct {
	BaseNode
	left  Expr
	token lexer.TokenType
	right Expr
}

func (b *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinary(b)
}

type Unary struct {
	BaseNode
	token lexer.TokenType
	expr  Expr
}

func (b *Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnary(b)
}

type Comma struct {
	BaseNode
	left  Expr
	right Expr
}

func (b *Comma) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitComma(b)
}

type Ternary struct {
	BaseNode
	left   Expr
	middle Expr
	right  Expr
}

func (b *Ternary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTernary(b)
}

type Grouping struct {
	BaseNode
	expr Expr
}

func (b *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGrouping(b)
}

type Literal struct {
	BaseNode
	token lexer.Token
}

func (b *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteral(b)
}
