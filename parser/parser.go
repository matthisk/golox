package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/matthisk/lox/ds"
	"github.com/matthisk/lox/lexer"
)

/*
program     -> statement* EOF ;

// STATEMENT GRAMMAR
declaration -> funDecl | classDecl | varDecl | statement;
funDecl     -> "fun" function;
classDecl   -> "class" IDENTIFIER "{" function* "}";
function    -> IDENTIFIER "(" parameters* ")" block;
parameters  -> IDENTIFIER ( "," IDENTIFIER )*;
varDecl     -> "var" IDENTIFIER ( "=" expression )? ";" ;
statement   -> exprStmt | ifStmt | whileStmt | forStmt | printStmt | block;
exprStmt    -> expression ";";
ifStmt      -> "if" "(" expression ")" statement ( "else" statement )?;
whileStmt   -> "while" "(" expression ")" statement;
forStmt     -> "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement;
printStmt   -> "print" expression ";";
block       -> "{" declaration* "}";
controlStmt -> ( BREAK | CONTINUE ) ";"; // Can only happen while in for/while/if statement Body.

// EXPRESSION GRAMMAR
expression     -> comma;
comma          -> ternary ( "," ternary )* ;
ternary        -> equality ("?" ternary ":" ternary)? ;
assignment     -> ( call "." )?  IDENTIFIER "=" assignment | logic_or ;
logic_or       -> logic_and ( "or" logic_and )* ;
logic_and      -> equality ( "and" equality )* ;
equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           -> factor ( ( "-" | "+" ) factor )* ;
factor         -> unary ( ( "/" | "*" ) unary )* ;
unary          -> ( "!" | "-" ) unary ;
call           -> primary ( "(" arguments ")" | "." IDENTIFIER )*;
primary        -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
arguments      -> expression ( "," expression )* ;
*/

type ContextItem struct {
	Type     string
	CanBreak bool
}

type Parser struct {
	source string
	tokens []lexer.Token
	// State
	index        int
	contextStack *ds.Stack[*ContextItem]
	errors       *Errors
}

func New(toks []lexer.Token, source string) *Parser {
	return &Parser{
		source:       source,
		tokens:       toks,
		index:        0,
		contextStack: ds.NewStack[*ContextItem](),
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	// In the rare case we want to call the parser twice on the same input,
	// we reset the state.
	p.index = 0
	p.errors = nil
	p.contextStack = ds.NewStack[*ContextItem]()

	return p.statements()
}

// ReportErrors returns a stream which
func (p *Parser) ReportErrors() string {
	if p.errors == nil {
		return ""
	}

	var w = bytes.NewBufferString("")

	for _, e := range p.errors.errs {
		fmt.Fprintf(w, "%s\n\n", e.Error())
		atLine := e.at.StartPos.Line
		sourceLines := strings.Split(p.source, "\n")
		if atLine > 2 {
			printSourceLine(w, atLine-2, sourceLines[atLine-3])
		}
		if atLine > 1 {
			printSourceLine(w, atLine-1, sourceLines[atLine-2])
		}
		printSourceLine(w, atLine, sourceLines[atLine-1])
		fmt.Fprintf(w, "%s  %s^\n", strings.Repeat(" ", len(strconv.Itoa(atLine))), strings.Repeat(" ", e.at.StartPos.Column-1))
		fmt.Fprintf(w, "\n")
	}

	return w.String()
}

func printSourceLine(w io.Writer, n int, line string) {
	fmt.Fprintf(w, "%d| %s\n", n, line)
}

func (p *Parser) Expression() (Expr, error) {
	return p.expression()
}

func (p *Parser) statements() ([]Stmt, error) {
	var result []Stmt
	var errors []*Error

	for !p.atEnd() {
		stmt, err := p.declStatement()

		if stmt != nil {
			result = append(result, stmt)
		}
		if err != nil {
			errors = append(errors, err.(*Error))
		}
	}

	if len(errors) > 0 {
		p.errors = NewErrors(errors)
		return result, p.errors
	}

	return result, nil
}

func (p *Parser) declStatement() (Stmt, error) {
	var stmt Stmt
	var err error

	if p.match(lexer.VAR) {
		stmt, err = p.varDeclStatement()
	} else if p.match(lexer.FUN) {
		stmt, err = p.function("function")
	} else if p.match(lexer.CLASS) {
		stmt, err = p.class()
	} else {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
	}

	return stmt, err
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.atEnd() {
		if p.previous().Type == lexer.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case lexer.CLASS:
			return
		case lexer.FUN:
			return
		case lexer.VAR:
			return
		case lexer.FOR:
			return
		case lexer.IF:
			return
		case lexer.WHILE:
			return
		case lexer.PRINT:
			return
		case lexer.RETURN:
			return
		default:
		}

		p.advance()
	}
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}

	if p.match(lexer.IF) {
		return p.ifStatement()
	}

	if p.match(lexer.WHILE) {
		return p.whileStatement()
	}

	if p.match(lexer.FOR) {
		return p.forStatement()
	}

	if p.match(lexer.LEFT_BRACE) {
		return p.blockStatement()
	}

	if p.match(lexer.RETURN) {
		return p.returnStmt()
	}

	if p.match(lexer.BREAK) {
		if cs, _ := p.contextStack.Peek(); cs == nil || !cs.CanBreak {
			return nil, p.error("Illegal jump target")
		}

		startPos := p.previous().StartPos

		err := p.consume(lexer.SEMICOLON, "Expect ';' after break statement.")
		if err != nil {
			return nil, err
		}

		endPos := p.previous().EndPos

		return &BreakStmt{BaseNode{StartPos: startPos, EndPos: endPos}}, nil
	}

	if p.match(lexer.CONTINUE) {
		if cs, _ := p.contextStack.Peek(); cs == nil || !cs.CanBreak {
			return nil, p.error("Illegal jump target")
		}

		startPos := p.previous().StartPos

		err := p.consume(lexer.SEMICOLON, "Expect ';' after continue statement.")
		if err != nil {
			return nil, err
		}

		endPos := p.previous().EndPos

		return &ContinueStmt{BaseNode{StartPos: startPos, EndPos: endPos}}, nil
	}

	return p.exprStatement()
}

func (p *Parser) class() (Stmt, error) {
	p.contextStack.Push(&ContextItem{"class", false})
	defer p.contextStack.Pop()

	startPos := p.previous().StartPos

	err := p.consume(lexer.IDENTIFIER, "I was expecting to see a class Name after the `class` keyword.")
	if err != nil {
		return nil, err
	}

	name := p.previous()

	err = p.consume(lexer.LEFT_BRACE, "I was expecting to see an opening brace `{` after the class Name.")
	if err != nil {
		return nil, err
	}

	var methods []*Function
	for p.peek().Type != lexer.RIGHT_BRACE && !p.atEnd() {
		fun, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, fun.(*Function))
	}

	err = p.consume(lexer.RIGHT_BRACE, "Expect '}' after class.")
	if err != nil {
		return nil, err
	}

	return &ClassStatement{
		BaseNode: BaseNode{
			StartPos: startPos,
			EndPos:   p.previous().EndPos,
		},
		Name:    name,
		Methods: methods,
	}, nil
}

func (p *Parser) function(funType string) (Stmt, error) {
	p.contextStack.Push(&ContextItem{funType, false})
	defer p.contextStack.Pop()

	pos := p.previous().StartPos
	err := p.consume(lexer.IDENTIFIER, fmt.Sprintf("Expect %s Name.", funType))
	if err != nil {
		return nil, err
	}

	funName := p.previous()

	err = p.consume(lexer.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s Name.", funType))
	if err != nil {
		return nil, err
	}

	parameters, err := p.parameters()
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.RIGHT_PAREN, fmt.Sprintf("Expect ')' after %s parameters.", funType))
	if err != nil {
		return nil, err
	}

	err = p.consume(lexer.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s Body.", funType))
	if err != nil {
		return nil, err
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &Function{
		BaseNode: BaseNode{
			StartPos: pos,
			EndPos:   p.previous().EndPos,
		},
		Name:   funName,
		Params: parameters,
		Body:   body.(*Block).Stmts,
	}, nil
}

func (p *Parser) parameters() ([]lexer.Token, error) {
	var result []lexer.Token

	// The function has no parameters.
	if p.check(lexer.RIGHT_PAREN) {
		return nil, nil
	}

	for {
		err := p.consume(lexer.IDENTIFIER, "Expect 'IDENTIFIER'.")
		if err != nil {
			return nil, err
		}

		result = append(result, p.previous())

		if len(result) > 255 {
			return nil, p.error("Can't have more than 255 parameters.")
		}

		if !p.match(lexer.COMMA) {
			break
		}
	}

	return result, nil
}

func (p *Parser) varDeclStatement() (Stmt, error) {
	startPos := p.previous().StartPos
	var initializer Expr
	err := p.consume(lexer.IDENTIFIER, "Expect variable Name.")
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
			StartPos: startPos,
			EndPos:   p.previous().EndPos,
		},
		Name:        fmt.Sprintf("%s", name.Lexeme),
		Initializer: initializer,
	}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	startPos := p.previous().StartPos
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
			StartPos: startPos,
			EndPos:   GetExprEndPos(expr),
		},
		Expr: expr,
	}, nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	startPos := p.previous().StartPos
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
			StartPos: startPos,
			EndPos:   p.previous().EndPos,
		},
		Stmts: result,
	}, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	p.contextStack.Push(&ContextItem{"if", true})
	defer p.contextStack.Pop()

	pos := p.previous().StartPos

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
			StartPos: pos,
			EndPos:   p.previous().EndPos,
		},
		Cond:      cond,
		IfBlock:   ifBlock,
		ElseBlock: elseBlock,
	}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	p.contextStack.Push(&ContextItem{"while", true})
	defer p.contextStack.Pop()

	pos := p.previous().StartPos

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
			StartPos: pos,
			EndPos:   p.previous().EndPos,
		},
		Cond: cond,
		Body: body,
	}, nil
}

func (p *Parser) forStatement() (Stmt, error) {
	p.contextStack.Push(&ContextItem{"for", true})
	defer p.contextStack.Pop()

	pos := p.previous().StartPos

	err := p.consume(lexer.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// Parse initializer: varDecl | exprStmt | ";"
	var initializer Stmt
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer, err = p.varDeclStatement()
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

	// Parse Body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Desugar the for loop into a while loop
	if increment != nil {
		body = &Block{
			BaseNode: BaseNode{},
			Stmts:    []Stmt{body, &ExprStmt{Expr: increment}},
		}
	}

	if condition == nil {
		condition = &Literal{Token: lexer.Token{Type: lexer.TRUE}}
	}

	body = &WhileStatement{
		BaseNode: BaseNode{
			StartPos: pos,
			EndPos:   p.previous().EndPos,
		},
		Cond: condition,
		Body: body,
	}

	if initializer != nil {
		body = &Block{
			Stmts: []Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) returnStmt() (Stmt, error) {
	pos := p.previous().StartPos
	var expression Expr
	var err error

	if !p.check(lexer.SEMICOLON) {
		expression, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(lexer.SEMICOLON, "Expect ';' after statement.")
	if err != nil {
		return nil, err
	}

	return &ReturnStmt{
		BaseNode: BaseNode{
			StartPos: pos,
			EndPos:   p.previous().EndPos,
		},
		Expr: expression,
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
			StartPos: GetExprStartPos(expr),
			EndPos:   GetExprEndPos(expr),
		},
		Expr: expr,
	}, nil
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
			Left:     expr,
			Right:    right,
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
				Left:     expr,
				Middle:   middle,
				Right:    right,
			}

		} else {
			return nil, p.error("Expect colon.")
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
			name := e.Name
			return &Assign{
				BaseNode: BaseNode{
					StartPos: GetExprStartPos(expr),
					EndPos:   GetExprEndPos(value),
				},
				Name:  name,
				Value: value,
			}, nil
		} else if e, ok := expr.(*GetExpr); ok {
			return &SetExpr{
				BaseNode: BaseNode{
					StartPos: GetExprStartPos(expr),
					EndPos:   GetExprEndPos(value),
				},
				Object: e.From,
				Name:   e.Property,
				Value:  value,
			}, nil

		}

		return nil, p.error("invalid assignment target")
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
			Left:     expr,
			Token:    operator.Type,
			Right:    right,
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
			Left:     expr,
			Token:    operator.Type,
			Right:    right,
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
			Left:     expr,
			Token:    operator.Type,
			Right:    right,
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
			Left:     expr,
			Token:    operator.Type,
			Right:    right,
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
			Left:     expr,
			Token:    operator.Type,
			Right:    right,
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
			Left:     expr,
			Token:    operator.Type,
			Right:    right,
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
				StartPos: operator.StartPos,
				EndPos:   GetExprEndPos(right),
			},
			Token: operator.Type,
			Expr:  right,
		}, nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(lexer.LEFT_PAREN) {
			arguments, err := p.arguments()
			if err != nil {
				return nil, err
			}

			err = p.consume(lexer.RIGHT_PAREN, "Expect ')' closing paren after call.")
			if err != nil {
				return nil, err
			}

			expr = &Call{
				BaseNode: BaseNode{
					StartPos: GetExprStartPos(expr),
					EndPos:   p.previous().EndPos,
				},
				Callee:    expr,
				Paren:     p.previous(),
				Arguments: arguments,
			}
		} else if p.match(lexer.DOT) {
			err := p.consume(lexer.IDENTIFIER, "Expect identifier after '.'")
			if err != nil {
				return nil, err
			}

			expr = &GetExpr{
				BaseNode: BaseNode{
					StartPos: GetExprStartPos(expr),
					EndPos:   p.previous().EndPos,
				},
				From:     expr,
				Property: p.previous(),
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) arguments() ([]Expr, error) {
	// in case there are no arguments return immediately (i.e. fn()).
	if p.check(lexer.RIGHT_PAREN) {
		return nil, nil
	}

	var result []Expr

	for {
		// We don't parse an expression here, as the comma expression would consume the comma's separating arguments.
		expression, err := p.ternary()
		if err != nil {
			return nil, err
		}

		result = append(result, expression)

		if len(result) > 255 {
			return nil, p.error("Can't have more than 255 arguments.")
		}

		if !p.match(lexer.COMMA) {
			break
		}
	}

	return result, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(lexer.NUMBER, lexer.STRING, lexer.FALSE, lexer.TRUE, lexer.NIL) {
		token := p.previous()
		return &Literal{
			BaseNode: BaseNode{
				StartPos: token.StartPos,
				EndPos:   token.EndPos,
			},
			Token: token,
		}, nil
	}

	if p.match(lexer.IDENTIFIER) {
		token := p.previous()
		return &Variable{
			BaseNode: BaseNode{
				StartPos: token.StartPos,
				EndPos:   token.EndPos,
			},
			Name: fmt.Sprintf("%s", token.Lexeme),
		}, nil
	}

	if p.match(lexer.THIS) {
		token := p.previous()
		return &This{
			BaseNode: BaseNode{
				StartPos: token.StartPos,
				EndPos:   token.EndPos,
			},
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
				StartPos: leftParen.StartPos,
				EndPos:   rightParen.EndPos,
			},
			Expr: expr,
		}, nil
	}

	return nil, p.error("Expect expression.")
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

	return p.error(err)
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

func (p *Parser) error(msg string) *Error {
	return &Error{
		pos: p.previous().StartPos,
		at:  p.peek(),
		msg: msg,
	}
}
