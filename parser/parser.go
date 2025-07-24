package parser

import (
	"errors"

	"github.com/matthisk/lox/lexer"
)

/*
program     → statement* EOF ;

statement   -> exprStmt | printStmt;

exprStmt    -> expression ";";
printStmt   -> "print" expression ";";

expression     → comma ;
comma          -> ternary ( "," ternary )*;
ternary        -> equality ("?" ternary ":" ternary)? ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
*/

type Parser struct {
	tokens []lexer.Token
	// State
	index int
}

func New(toks []lexer.Token) *Parser {
	return &Parser{
		tokens: toks,
		index:  0,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	return p.statements()
}

func (p *Parser) statements() ([]Stmt, error) {
	var result []Stmt

	for !p.atEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}

		result = append(result, stmt)
	}

	return result, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}

	return p.exprStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &PrintStmt{
		BaseNode: BaseNode{},
		expr:     expr,
	}, nil
}

func (p *Parser) exprStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ExprStmt{
		BaseNode: BaseNode{},
		expr:     expr,
	}, nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.atEnd() {
		if p.previous().Type == lexer.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case lexer.CLASS:
		case lexer.FUN:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		default:
		}

		p.advance()
	}
}

func (p *Parser) expression() (Expr, error) {
	return p.comma()
}

func (p *Parser) comma() (Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.COMMA) {
		right, err := p.ternary()
		if err != nil {
			return nil, err
		}

		expr = &Comma{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) ternary() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.QUESTION_MARK) {
		middle, err := p.ternary()
		if err != nil {
			return nil, err
		}

		if p.match(lexer.COLON) {
			right, err := p.ternary()
			if err != nil {
				return nil, err
			}

			expr = &Ternary{
				BaseNode: NewBaseNodeFromExprs(expr, right),
				left:     expr,
				middle:   middle,
				right:    right,
			}

		} else {
			return nil, errors.New("Expect colon.")
		}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			token:    operator.Type,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			token:    operator.Type,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.PLUS, lexer.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			token:    operator.Type,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.STAR, lexer.SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			token:    operator.Type,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right, _ := p.unary()
		return &Unary{
			BaseNode: BaseNode{
				startPos: operator.StartPos,
				endPos:   GetExprEndPos(right),
			},
			token: operator.Type,
			expr:  right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(lexer.NUMBER, lexer.STRING, lexer.FALSE, lexer.TRUE, lexer.NIL) {
		token := p.previous()
		return &Literal{
			BaseNode: BaseNode{
				startPos: token.StartPos,
				endPos:   token.EndPos,
			},
			token: token,
		}, nil
	}

	if p.match(lexer.LEFT_PAREN) {
		expr, _ := p.expression()
		err := p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		leftParen := p.tokens[p.index-1]
		rightParen := p.previous()
		return &Grouping{
			BaseNode: BaseNode{
				startPos: leftParen.StartPos,
				endPos:   rightParen.EndPos,
			},
			expr: expr,
		}, nil
	}

	return nil, errors.New("Expect expression.")
}

func (p *Parser) match(tokens ...lexer.TokenType) bool {
	for _, token := range tokens {
		if p.check(token) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tok lexer.TokenType, err string) error {
	if p.check(tok) {
		p.advance()
		return nil
	}

	return errors.New(err)
}

func (p *Parser) check(token lexer.TokenType) bool {
	if p.atEnd() {
		return false
	}
	return p.peek().Type == token
}

func (p *Parser) atEnd() bool {
	return p.peek().Type == lexer.EOF
}

func (p *Parser) previous() lexer.Token {
	if p.index == 0 {
		panic("Cannot peek at previous when at index 0")
	}
	return p.tokens[p.index-1]
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.index]
}

func (p *Parser) advance() {
	p.index++
}
