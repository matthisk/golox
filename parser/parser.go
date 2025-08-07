package parser

import (
	"errors"
	"fmt"
	"github.com/matthisk/lox/lexer"
)

/*
program     → statement* EOF ;

declaration -> varDecl | statement;

varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;

statement   -> exprStmt | ifStmt | whileStmt | forStmt | printStmt | block;

exprStmt    -> expression ";";
ifStmt      -> "if" "(" expression ")" statement ( "else" statement )?;
whileStmt   -> "while" "(" expression ")" statement;
forStmt     -> "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement;
printStmt   -> "print" expression ";";
block       -> "{" declaration* "}";

expression     → comma ;
comma          -> ternary ( "," ternary )*;
ternary        -> equality ("?" ternary ":" ternary)? ;
assignment     -> IDENTIFIER "=" assignment | logic_or;
logic_or       -> logic_and ( "or" logic_and )*;
logic_and      -> equality ( "and" equality )*;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" | IDENTIFIER ;
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

func (p *Parser) Expression() (Expr, error) {
	return p.expression()
}

func (p *Parser) statements() ([]Stmt, error) {
	var result []Stmt

	for !p.atEnd() {
		stmt, err := p.declStatement()
		if err != nil {
			return nil, err
		}

		result = append(result, stmt)
	}

	return result, nil
}

func (p *Parser) declStatement() (Stmt, error) {
	var stmt Stmt
	var err error

	if p.match(lexer.VAR) {
		stmt, err = p.varDeclStatement(p.previous().StartPos)
	}

	if stmt == nil {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
	}

	return stmt, err
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(lexer.PRINT) {
		return p.printStatement(p.previous().StartPos)
	}

	if p.match(lexer.IF) {
		return p.ifStatement(p.previous().StartPos)
	}

	if p.match(lexer.WHILE) {
		return p.whileStatement(p.previous().StartPos)
	}

	if p.match(lexer.FOR) {
		return p.forStatement(p.previous().StartPos)
	}

	if p.match(lexer.LEFT_BRACE) {
		return p.blockStatement(p.previous().StartPos)
	}

	return p.exprStatement()
}

func (p *Parser) varDeclStatement(startPos lexer.Pos) (Stmt, error) {
	var initializer Expr
	err := p.consume(lexer.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	name := p.previous()

	if p.match(lexer.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &VarDecl{
		BaseNode: BaseNode{
			startPos: startPos,
			endPos:   p.previous().EndPos,
		},
		name:        fmt.Sprintf("%s", name.Lexeme),
		initializer: initializer,
	}, nil
}

func (p *Parser) printStatement(startPos lexer.Pos) (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &PrintStmt{
		BaseNode: BaseNode{
			startPos: startPos,
			endPos:   GetExprEndPos(expr),
		},
		expr: expr,
	}, nil
}

func (p *Parser) blockStatement(startPos lexer.Pos) (Stmt, error) {
	var result []Stmt
	for !p.check(lexer.RIGHT_BRACE) && !p.atEnd() {
		stmt, err := p.declStatement()
		if err != nil {
			return nil, err
		}

		result = append(result, stmt)
	}

	err := p.consume(lexer.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return &Block{
		BaseNode: BaseNode{
			startPos: startPos,
			endPos:   p.previous().EndPos,
		},
		stmts: result,
	}, nil
}

func (p *Parser) ifStatement(pos lexer.Pos) (Stmt, error) {
	err := p.consume(lexer.LEFT_PAREN, "Expect '(' after if statement.")
	if err != nil {
		return nil, err
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.RIGHT_PAREN, "Expect ')' after if statement.")
	if err != nil {
		return nil, err
	}

	ifBlock, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBlock Stmt
	if p.match(lexer.ELSE) {
		elseBlock, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStatement{
		BaseNode: BaseNode{
			startPos: pos,
			endPos:   p.previous().EndPos,
		},
		cond:      cond,
		ifBlock:   ifBlock,
		elseBlock: elseBlock,
	}, nil
}

func (p *Parser) whileStatement(pos lexer.Pos) (Stmt, error) {
	err := p.consume(lexer.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &WhileStatement{
		BaseNode: BaseNode{
			startPos: pos,
			endPos:   p.previous().EndPos,
		},
		cond: cond,
		body: body,
	}, nil
}

func (p *Parser) forStatement(pos lexer.Pos) (Stmt, error) {
	err := p.consume(lexer.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// Parse initializer: varDecl | exprStmt | ";"
	var initializer Stmt
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer, err = p.varDeclStatement(p.previous().StartPos)
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.exprStatement()
		if err != nil {
			return nil, err
		}
	}

	// Parse condition (optional)
	var condition Expr
	if !p.check(lexer.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	err = p.consume(lexer.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// Parse increment (optional)
	var increment Expr
	if !p.check(lexer.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	err = p.consume(lexer.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	// Parse body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ForStatement{
		BaseNode: BaseNode{
			startPos: pos,
			endPos:   p.previous().EndPos,
		},
		initializer: initializer,
		condition:   condition,
		increment:   increment,
		body:        body,
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
		BaseNode: BaseNode{
			startPos: GetExprStartPos(expr),
			endPos:   GetExprEndPos(expr),
		},
		expr: expr,
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
	expr, err := p.assignment()
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

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.logicOr()
	if err != nil {
		return nil, err
	}

	if p.match(lexer.EQUAL) {
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if e, ok := expr.(*Variable); ok {
			name := e.name
			return &Assign{
				BaseNode: BaseNode{
					startPos: GetExprStartPos(expr),
					endPos:   GetExprEndPos(e),
				},
				name:  name,
				value: value,
			}, nil
		}

		return nil, errors.New("Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) logicOr() (Expr, error) {
	expr, err := p.logicAnd()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.OR) {
		operator := p.previous()
		right, err := p.logicAnd()
		if err != nil {
			return nil, err
		}

		expr = &Logical{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			token:    operator.Type,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) logicAnd() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &Logical{
			BaseNode: NewBaseNodeFromExprs(expr, right),
			left:     expr,
			token:    operator.Type,
			right:    right,
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

	if p.match(lexer.IDENTIFIER) {
		token := p.previous()
		return &Variable{
			BaseNode: BaseNode{
				startPos: token.StartPos,
				endPos:   token.EndPos,
			},
			name: fmt.Sprintf("%s", token.Lexeme),
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
	if p.index >= len(p.tokens) {
		return true
	}

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
