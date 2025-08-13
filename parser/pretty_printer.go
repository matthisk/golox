package parser

import (
	"fmt"
	"strings"

	"github.com/matthisk/lox/lexer"
)

type PrettyPrinter struct {
	indent int
}

func NewPrettyPrinter() *PrettyPrinter {
	return &PrettyPrinter{indent: 0}
}

func (p *PrettyPrinter) getIndent() string {
	return strings.Repeat("  ", p.indent)
}

func (p *PrettyPrinter) VisitReturnStmt(r *ReturnStmt) (interface{}, error) {
	if r.expr == nil {
		return p.getIndent() + "return;", nil
	}
	expr, err := r.expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("return %s;", expr), nil
}

func (p *PrettyPrinter) VisitCall(c *Call) (interface{}, error) {
	callee, err := c.callee.Accept(p)
	if err != nil {
		return nil, err
	}

	var args []string
	for _, arg := range c.arguments {
		argStr, err := arg.Accept(p)
		if err != nil {
			return nil, err
		}
		args = append(args, argStr.(string))
	}

	return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", ")), nil
}

func (p *PrettyPrinter) VisitFunction(f *Function) (interface{}, error) {
	var params []string
	for _, param := range f.params {
		params = append(params, fmt.Sprintf("%s", param.Lexeme))
	}

	result := p.getIndent() + fmt.Sprintf("fun %s(%s) {\n", f.name.Lexeme, strings.Join(params, ", "))

	p.indent++
	for _, stmt := range f.body {
		stmtStr, err := stmt.Accept(p)
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf("%s\n", stmtStr)
	}
	p.indent--

	result += p.getIndent() + "}"

	return result, nil
}

func (p *PrettyPrinter) VisitBlock(b *Block) (interface{}, error) {
	result := p.getIndent() + "{\n"

	p.indent++
	for _, stmt := range b.stmts {
		stmtStr, err := stmt.Accept(p)
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf("%s\n", stmtStr)
	}
	p.indent--

	result += p.getIndent() + "}"

	return result, nil
}

func (p *PrettyPrinter) VisitIfStatement(s *IfStatement) (interface{}, error) {
	cond, err := s.cond.Accept(p)
	if err != nil {
		return nil, err
	}

	result := p.getIndent() + fmt.Sprintf("if (%s)", cond)

	// Handle if block formatting
	ifBlock, err := s.ifBlock.Accept(p)
	if err != nil {
		return nil, err
	}

	// Check if the if block is a Block or a single statement
	if _, ok := s.ifBlock.(*Block); ok {
		result += " " + strings.TrimPrefix(ifBlock.(string), p.getIndent())
	} else {
		result += " {\n"
		p.indent++
		result += fmt.Sprintf("%s\n", ifBlock)
		p.indent--
		result += p.getIndent() + "}"
	}

	// Handle else block if present
	if s.elseBlock != nil {
		result += " else"

		elseBlock, err := s.elseBlock.Accept(p)
		if err != nil {
			return nil, err
		}

		// Check if else block is an if statement (else if)
		if _, ok := s.elseBlock.(*IfStatement); ok {
			result += " " + strings.TrimPrefix(elseBlock.(string), p.getIndent())
		} else if _, ok := s.elseBlock.(*Block); ok {
			result += " " + strings.TrimPrefix(elseBlock.(string), p.getIndent())
		} else {
			result += " {\n"
			p.indent++
			result += fmt.Sprintf("%s\n", elseBlock)
			p.indent--
			result += p.getIndent() + "}"
		}
	}

	return result, nil
}

func (p *PrettyPrinter) VisitWhileStatement(s *WhileStatement) (interface{}, error) {
	cond, err := s.cond.Accept(p)
	if err != nil {
		return nil, err
	}

	result := p.getIndent() + fmt.Sprintf("while (%s)", cond)

	body, err := s.body.Accept(p)
	if err != nil {
		return nil, err
	}

	// Check if body is a Block or a single statement
	if _, ok := s.body.(*Block); ok {
		result += " " + strings.TrimPrefix(body.(string), p.getIndent())
	} else {
		result += " {\n"
		p.indent++
		result += fmt.Sprintf("%s\n", body)
		p.indent--
		result += p.getIndent() + "}"
	}

	return result, nil
}

func (p *PrettyPrinter) VisitForStatement(s *ForStatement) (interface{}, error) {
	var init, cond, inc string

	if s.initializer != nil {
		initResult, err := s.initializer.Accept(&PrettyPrinter{indent: 0})
		if err != nil {
			return nil, err
		}
		init = strings.TrimSpace(initResult.(string))
		if strings.HasSuffix(init, ";") {
			init = init[:len(init)-1]
		}
	}

	if s.condition != nil {
		condResult, err := s.condition.Accept(&PrettyPrinter{indent: 0})
		if err != nil {
			return nil, err
		}
		cond = condResult.(string)
	}

	if s.increment != nil {
		incResult, err := s.increment.Accept(&PrettyPrinter{indent: 0})
		if err != nil {
			return nil, err
		}
		inc = incResult.(string)
	}

	result := p.getIndent() + fmt.Sprintf("for (%s; %s; %s)", init, cond, inc)

	body, err := s.body.Accept(p)
	if err != nil {
		return nil, err
	}

	// Check if body is a Block or a single statement
	if _, ok := s.body.(*Block); ok {
		result += " " + strings.TrimPrefix(body.(string), p.getIndent())
	} else {
		result += " {\n"
		p.indent++
		result += fmt.Sprintf("%s\n", body)
		p.indent--
		result += p.getIndent() + "}"
	}

	return result, nil
}

func (p *PrettyPrinter) VisitContinueStmt(c *ContinueStmt) (interface{}, error) {
	return p.getIndent() + "continue;", nil
}

func (p *PrettyPrinter) VisitBreakStmt(b *BreakStmt) (interface{}, error) {
	return p.getIndent() + "break;", nil
}

func (p *PrettyPrinter) VisitVarDecl(vd *VarDecl) (interface{}, error) {
	if vd.initializer == nil {
		return p.getIndent() + fmt.Sprintf("var %s;", vd.name), nil
	}
	initializer, err := vd.initializer.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("var %s = %s;", vd.name, initializer), nil
}

func (p *PrettyPrinter) VisitPrintStmt(node *PrintStmt) (interface{}, error) {
	expr, err := node.expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("print %s;", expr), nil
}

func (p *PrettyPrinter) VisitExprStmt(node *ExprStmt) (interface{}, error) {
	expr, err := node.expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("%s;", expr), nil
}

func (p *PrettyPrinter) VisitLogical(b *Logical) (interface{}, error) {
	left, err := b.left.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := b.right.Accept(p)
	if err != nil {
		return nil, err
	}

	var op string
	switch b.token {
	case lexer.AND:
		op = "and"
	case lexer.OR:
		op = "or"
	default:
		op = b.token.String()
	}

	return fmt.Sprintf("%s %s %s", left, op, right), nil
}

func (p *PrettyPrinter) VisitBinary(node *Binary) (interface{}, error) {
	left, err := node.left.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := node.right.Accept(p)
	if err != nil {
		return nil, err
	}

	var op string
	switch node.token {
	case lexer.EQUAL_EQUAL:
		op = "=="
	case lexer.BANG_EQUAL:
		op = "!="
	case lexer.GREATER:
		op = ">"
	case lexer.GREATER_EQUAL:
		op = ">="
	case lexer.LESS:
		op = "<"
	case lexer.LESS_EQUAL:
		op = "<="
	case lexer.PLUS:
		op = "+"
	case lexer.MINUS:
		op = "-"
	case lexer.STAR:
		op = "*"
	case lexer.SLASH:
		op = "/"
	default:
		op = node.token.String()
	}

	return fmt.Sprintf("%s %s %s", left, op, right), nil
}

func (p *PrettyPrinter) VisitLiteral(node *Literal) (interface{}, error) {
	switch node.token.Type {
	case lexer.STRING:
		return fmt.Sprintf("\"%s\"", node.token.Lexeme), nil
	case lexer.NUMBER:
		return fmt.Sprintf("%v", node.token.Lexeme), nil
	case lexer.TRUE:
		return "true", nil
	case lexer.FALSE:
		return "false", nil
	case lexer.NIL:
		return "nil", nil
	default:
		return fmt.Sprintf("%v", node.token.Lexeme), nil
	}
}

func (p *PrettyPrinter) VisitVariable(b *Variable) (interface{}, error) {
	return b.name, nil
}

func (p *PrettyPrinter) VisitUnary(node *Unary) (interface{}, error) {
	expr, err := node.expr.Accept(p)
	if err != nil {
		return nil, err
	}

	var op string
	switch node.token {
	case lexer.BANG:
		op = "!"
	case lexer.MINUS:
		op = "-"
	default:
		op = node.token.String()
	}

	return fmt.Sprintf("%s%s", op, expr), nil
}

func (p *PrettyPrinter) VisitComma(node *Comma) (interface{}, error) {
	left, err := node.left.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := node.right.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s, %s", left, right), nil
}

func (p *PrettyPrinter) VisitTernary(node *Ternary) (interface{}, error) {
	left, err := node.left.Accept(p)
	if err != nil {
		return nil, err
	}
	middle, err := node.middle.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := node.right.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s ? %s : %s", left, middle, right), nil
}

func (p *PrettyPrinter) VisitAssign(b *Assign) (interface{}, error) {
	value, err := b.value.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s = %s", b.name, value), nil
}

func (p *PrettyPrinter) VisitGrouping(node *Grouping) (interface{}, error) {
	expr, err := node.expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s)", expr), nil
}

func PrettyPrint(stmts []Stmt) (string, error) {
	printer := NewPrettyPrinter()
	var result strings.Builder

	for i, stmt := range stmts {
		stmtStr, err := stmt.Accept(printer)
		if err != nil {
			return "", err
		}
		result.WriteString(stmtStr.(string))
		if i < len(stmts)-1 {
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}