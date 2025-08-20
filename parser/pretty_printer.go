package parser

import (
	"fmt"
	"strings"

	"github.com/matthisk/lox/lexer"
)

type PrettyPrinter struct {
	indent int
}

func (p *PrettyPrinter) VisitSuper(s *Super) (interface{}, error) {
	return fmt.Sprintf("super.%s", s.Method.Lexeme.(string)), nil
}

func (p *PrettyPrinter) VisitThis(t *This) (interface{}, error) {
	return "this", nil
}

func (p *PrettyPrinter) VisitGet(g *GetExpr) (interface{}, error) {
	from, err := g.From.Accept(p)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s.%s", from, g.Property.Lexeme.(string)), nil
}

func (p *PrettyPrinter) VisitSet(s *SetExpr) (interface{}, error) {
	object, err := s.Object.Accept(p)
	if err != nil {
		return nil, err
	}

	value, err := s.Value.Accept(p)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s.%s = %s", object, s.Name.Lexeme.(string), value), nil
}

func NewPrettyPrinter() *PrettyPrinter {
	return &PrettyPrinter{indent: 0}
}

func (p *PrettyPrinter) getIndent() string {
	return strings.Repeat("  ", p.indent)
}

func (p *PrettyPrinter) VisitClass(s *ClassStatement) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrettyPrinter) VisitReturnStmt(r *ReturnStmt) (interface{}, error) {
	if r.Expr == nil {
		return p.getIndent() + "return;", nil
	}
	expr, err := r.Expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("return %s;", expr), nil
}

func (p *PrettyPrinter) VisitCall(c *Call) (interface{}, error) {
	callee, err := c.Callee.Accept(p)
	if err != nil {
		return nil, err
	}

	var args []string
	for _, arg := range c.Arguments {
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
	for _, param := range f.Params {
		params = append(params, fmt.Sprintf("%s", param.Lexeme))
	}

	result := p.getIndent() + fmt.Sprintf("fun %s(%s) {\n", f.Name.Lexeme, strings.Join(params, ", "))

	p.indent++
	for _, stmt := range f.Body {
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
	for _, stmt := range b.Stmts {
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
	cond, err := s.Cond.Accept(p)
	if err != nil {
		return nil, err
	}

	result := p.getIndent() + fmt.Sprintf("if (%s)", cond)

	// Handle if block formatting
	ifBlock, err := s.IfBlock.Accept(p)
	if err != nil {
		return nil, err
	}

	// Check if the if block is a Block or a single statement
	if _, ok := s.IfBlock.(*Block); ok {
		result += " " + strings.TrimPrefix(ifBlock.(string), p.getIndent())
	} else {
		result += " {\n"
		p.indent++
		result += fmt.Sprintf("%s\n", ifBlock)
		p.indent--
		result += p.getIndent() + "}"
	}

	// Handle else block if present
	if s.ElseBlock != nil {
		result += " else"

		elseBlock, err := s.ElseBlock.Accept(p)
		if err != nil {
			return nil, err
		}

		// Check if else block is an if statement (else if)
		if _, ok := s.ElseBlock.(*IfStatement); ok {
			result += " " + strings.TrimPrefix(elseBlock.(string), p.getIndent())
		} else if _, ok := s.ElseBlock.(*Block); ok {
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
	cond, err := s.Cond.Accept(p)
	if err != nil {
		return nil, err
	}

	result := p.getIndent() + fmt.Sprintf("while (%s)", cond)

	body, err := s.Body.Accept(p)
	if err != nil {
		return nil, err
	}

	// Check if Body is a Block or a single statement
	if _, ok := s.Body.(*Block); ok {
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

	if s.Initializer != nil {
		initResult, err := s.Initializer.Accept(&PrettyPrinter{indent: 0})
		if err != nil {
			return nil, err
		}
		init = strings.TrimSpace(initResult.(string))
		if strings.HasSuffix(init, ";") {
			init = init[:len(init)-1]
		}
	}

	if s.Condition != nil {
		condResult, err := s.Condition.Accept(&PrettyPrinter{indent: 0})
		if err != nil {
			return nil, err
		}
		cond = condResult.(string)
	}

	if s.Increment != nil {
		incResult, err := s.Increment.Accept(&PrettyPrinter{indent: 0})
		if err != nil {
			return nil, err
		}
		inc = incResult.(string)
	}

	result := p.getIndent() + fmt.Sprintf("for (%s; %s; %s)", init, cond, inc)

	body, err := s.Body.Accept(p)
	if err != nil {
		return nil, err
	}

	// Check if Body is a Block or a single statement
	if _, ok := s.Body.(*Block); ok {
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
	if vd.Initializer == nil {
		return p.getIndent() + fmt.Sprintf("var %s;", vd.Name), nil
	}
	initializer, err := vd.Initializer.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("var %s = %s;", vd.Name, initializer), nil
}

func (p *PrettyPrinter) VisitPrintStmt(node *PrintStmt) (interface{}, error) {
	expr, err := node.Expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("print %s;", expr), nil
}

func (p *PrettyPrinter) VisitExprStmt(node *ExprStmt) (interface{}, error) {
	expr, err := node.Expr.Accept(p)
	if err != nil {
		return nil, err
	}
	return p.getIndent() + fmt.Sprintf("%s;", expr), nil
}

func (p *PrettyPrinter) VisitLogical(b *Logical) (interface{}, error) {
	left, err := b.Left.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := b.Right.Accept(p)
	if err != nil {
		return nil, err
	}

	var op string
	switch b.Token {
	case lexer.AND:
		op = "and"
	case lexer.OR:
		op = "or"
	default:
		op = b.Token.String()
	}

	return fmt.Sprintf("%s %s %s", left, op, right), nil
}

func (p *PrettyPrinter) VisitBinary(node *Binary) (interface{}, error) {
	left, err := node.Left.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := node.Right.Accept(p)
	if err != nil {
		return nil, err
	}

	var op string
	switch node.Token {
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
		op = node.Token.String()
	}

	return fmt.Sprintf("%s %s %s", left, op, right), nil
}

func (p *PrettyPrinter) VisitLiteral(node *Literal) (interface{}, error) {
	switch node.Token.Type {
	case lexer.STRING:
		return fmt.Sprintf("\"%s\"", node.Token.Lexeme), nil
	case lexer.NUMBER:
		return fmt.Sprintf("%v", node.Token.Lexeme), nil
	case lexer.TRUE:
		return "true", nil
	case lexer.FALSE:
		return "false", nil
	case lexer.NIL:
		return "nil", nil
	default:
		return fmt.Sprintf("%v", node.Token.Lexeme), nil
	}
}

func (p *PrettyPrinter) VisitVariable(b *Variable) (interface{}, error) {
	return b.Name, nil
}

func (p *PrettyPrinter) VisitUnary(node *Unary) (interface{}, error) {
	expr, err := node.Expr.Accept(p)
	if err != nil {
		return nil, err
	}

	var op string
	switch node.Token {
	case lexer.BANG:
		op = "!"
	case lexer.MINUS:
		op = "-"
	default:
		op = node.Token.String()
	}

	return fmt.Sprintf("%s%s", op, expr), nil
}

func (p *PrettyPrinter) VisitComma(node *Comma) (interface{}, error) {
	left, err := node.Left.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := node.Right.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s, %s", left, right), nil
}

func (p *PrettyPrinter) VisitTernary(node *Ternary) (interface{}, error) {
	left, err := node.Left.Accept(p)
	if err != nil {
		return nil, err
	}
	middle, err := node.Middle.Accept(p)
	if err != nil {
		return nil, err
	}
	right, err := node.Right.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s ? %s : %s", left, middle, right), nil
}

func (p *PrettyPrinter) VisitAssign(b *Assign) (interface{}, error) {
	value, err := b.Value.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s = %s", b.Name, value), nil
}

func (p *PrettyPrinter) VisitGrouping(node *Grouping) (interface{}, error) {
	expr, err := node.Expr.Accept(p)
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
