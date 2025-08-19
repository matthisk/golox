package parser

import (
	"fmt"
	"strings"

	"github.com/matthisk/lox/lexer"
)

type AstPrinter struct{}

func (e AstPrinter) VisitThis(t *This) (interface{}, error) {
	return "this", nil
}

func (e AstPrinter) VisitGet(g *GetExpr) (interface{}, error) {
	from, err := g.From.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s.%s", from, g.Property.Lexeme.(string)), nil
}

func (e AstPrinter) VisitSet(s *SetExpr) (interface{}, error) {
	object, err := s.Object.Accept(e)
	if err != nil {
		return nil, err
	}

	value, err := s.Value.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s.%s = %s", object, s.Name.Lexeme.(string), value), nil
}

func (e AstPrinter) VisitReturnStmt(r *ReturnStmt) (interface{}, error) {
	if r.Expr == nil {
		return "return;", nil
	}

	expr, err := r.Expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("return %s;", expr), nil
}

func (e AstPrinter) VisitClass(s *ClassStatement) (interface{}, error) {
	var methods []string

	for _, method := range s.Methods {
		methodStr, err := method.Accept(e)
		if err != nil {
			return nil, err
		}
		methods = append(methods, methodStr.(string))
	}

	var body strings.Builder
	body.WriteString("{\n")
	for _, method := range methods {
		body.WriteString(fmt.Sprintf("  %s\n", method))
	}
	body.WriteString("}")

	return fmt.Sprintf("class %s %s", s.Name.Lexeme.(string), body.String()), nil
}

func (e AstPrinter) VisitCall(c *Call) (interface{}, error) {
	callee, err := c.Callee.Accept(e)
	if err != nil {
		return nil, err
	}

	var args []string
	for _, arg := range c.Arguments {
		argStr, err := arg.Accept(e)
		if err != nil {
			return nil, err
		}
		args = append(args, argStr.(string))
	}

	return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", ")), nil
}

func (e AstPrinter) VisitFunction(f *Function) (interface{}, error) {
	var params []string
	for _, param := range f.Params {
		params = append(params, fmt.Sprintf("%s", param.Lexeme))
	}

	var body strings.Builder
	body.WriteString("{\n")

	for _, stmt := range f.Body {
		stmtStr, err := stmt.Accept(e)
		if err != nil {
			return nil, err
		}
		body.WriteString(fmt.Sprintf("%s\n", stmtStr))
	}

	body.WriteString("}")

	return fmt.Sprintf("fun %s(%s) %s", f.Name.Lexeme, strings.Join(params, ", "), body.String()), nil
}

func (e AstPrinter) VisitBlock(b *Block) (interface{}, error) {
	var bd strings.Builder

	bd.WriteString("{\n")

	for i := range b.Stmts {
		stmt, err := b.Stmts[i].Accept(e)
		if err != nil {
			return nil, err
		}
		bd.WriteString(fmt.Sprintf("%s\n", stmt))
	}

	bd.WriteString("}")

	return bd.String(), nil
}

func (e AstPrinter) VisitIfStatement(s *IfStatement) (interface{}, error) {
	cond, err := s.Cond.Accept(e)
	if err != nil {
		return nil, err
	}

	ifBlock, err := s.IfBlock.Accept(e)
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("if (%s) %s", cond, ifBlock)

	if s.ElseBlock != nil {
		elseBlock, err := s.ElseBlock.Accept(e)
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf(" else %s", elseBlock)
	}

	return result, nil
}

func (e AstPrinter) VisitWhileStatement(s *WhileStatement) (interface{}, error) {
	cond, err := s.Cond.Accept(e)
	if err != nil {
		return nil, err
	}

	body, err := s.Body.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("while (%s) %s", cond, body), nil
}

func (e AstPrinter) VisitForStatement(s *ForStatement) (interface{}, error) {
	var init string
	var cond string
	var inc string

	if s.Initializer != nil {
		initResult, err := s.Initializer.Accept(e)
		if err != nil {
			return nil, err
		}
		init = initResult.(string)
		// Remove the semicolon from var declarations and expression statements
		if len(init) > 0 && init[len(init)-1] == ';' {
			init = init[:len(init)-1]
		}
	}

	if s.Condition != nil {
		condResult, err := s.Condition.Accept(e)
		if err != nil {
			return nil, err
		}
		cond = condResult.(string)
	}

	if s.Increment != nil {
		incResult, err := s.Increment.Accept(e)
		if err != nil {
			return nil, err
		}
		inc = incResult.(string)
	}

	body, err := s.Body.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("for (%s; %s; %s) %s", init, cond, inc, body), nil
}

func (e AstPrinter) VisitContinueStmt(c *ContinueStmt) (interface{}, error) {
	return "continue;", nil
}

func (e AstPrinter) VisitBreakStmt(b *BreakStmt) (interface{}, error) {
	return "break;", nil
}

func (e AstPrinter) VisitVarDecl(vd *VarDecl) (interface{}, error) {
	if vd.Initializer == nil {
		return fmt.Sprintf("var %s;", vd.Name), nil
	}

	initializer, err := vd.Initializer.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("var %s = %s;", vd.Name, initializer), nil
}

func (e AstPrinter) VisitPrintStmt(node *PrintStmt) (interface{}, error) {
	expr, err := node.Expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("print %s;", expr), nil
}

func (e AstPrinter) VisitExprStmt(node *ExprStmt) (interface{}, error) {
	expr, err := node.Expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s;", expr), nil
}

func (e AstPrinter) VisitLogical(b *Logical) (interface{}, error) {
	left, err := b.Left.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := b.Right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s %s)", b.Token, left, right), nil
}

func (e AstPrinter) VisitBinary(node *Binary) (interface{}, error) {
	left, err := node.Left.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := node.Right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s %s)", node.Token.String(), left.(string), right.(string)), nil
}

func (e AstPrinter) VisitLiteral(node *Literal) (interface{}, error) {
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

func (e AstPrinter) VisitVariable(b *Variable) (interface{}, error) {
	return b.Name, nil
}

func (e AstPrinter) VisitUnary(node *Unary) (interface{}, error) {
	expr, err := node.Expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s)", node.Token.String(), expr), nil
}

func (e AstPrinter) VisitComma(node *Comma) (interface{}, error) {
	left, err := node.Left.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := node.Right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(, %s %s)", left, right), nil
}

func (e AstPrinter) VisitTernary(node *Ternary) (interface{}, error) {
	left, err := node.Left.Accept(e)
	if err != nil {
		return nil, err
	}
	middle, err := node.Middle.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := node.Right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(TERNARY %s %s %s)", left.(string), middle.(string), right.(string)), nil
}

func (e AstPrinter) VisitAssign(b *Assign) (interface{}, error) {
	value, err := b.Value.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("(ASSIGN %s %s)", b.Name, value), nil
}

func (e AstPrinter) VisitGrouping(node *Grouping) (interface{}, error) {
	expr, err := node.Expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(group %s)", expr), nil
}

func Print(stmt Stmt) (string, error) {
	printer := AstPrinter{}
	res, err := stmt.Accept(printer)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}
