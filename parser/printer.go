package parser

import (
	"fmt"
	"github.com/matthisk/lox/lexer"
	"strings"
)

type AstPrinter struct{}

func (e AstPrinter) VisitBlock(b *Block) (interface{}, error) {
	var bd strings.Builder

	bd.WriteString("{\n")

	for i := range b.stmts {
		stmt, err := b.stmts[i].Accept(e)
		if err != nil {
			return nil, err
		}
		bd.WriteString(fmt.Sprintf("%s\n", stmt))
	}

	bd.WriteString("}")

	return bd.String(), nil
}

func (e AstPrinter) VisitIfStatement(s *IfStatement) (interface{}, error) {
	cond, err := s.cond.Accept(e)
	if err != nil {
		return nil, err
	}

	ifBlock, err := s.ifBlock.Accept(e)
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("if (%s) %s", cond, ifBlock)

	if s.elseBlock != nil {
		elseBlock, err := s.elseBlock.Accept(e)
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf(" else %s", elseBlock)
	}

	return result, nil
}

func (e AstPrinter) VisitWhileStatement(s *WhileStatement) (interface{}, error) {
	cond, err := s.cond.Accept(e)
	if err != nil {
		return nil, err
	}

	body, err := s.body.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("while (%s) %s", cond, body), nil
}

func (e AstPrinter) VisitVarDecl(vd *VarDecl) (interface{}, error) {
	initializer, err := vd.initializer.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("var %s = %s;", vd.name, initializer), nil
}

func (e AstPrinter) VisitPrintStmt(node *PrintStmt) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("print %s;", expr), nil
}

func (e AstPrinter) VisitExprStmt(node *ExprStmt) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s;", expr), nil
}

func (e AstPrinter) VisitLogical(b *Logical) (interface{}, error) {
	left, err := b.left.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := b.right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s %s)", b.token, left, right), nil
}

func (e AstPrinter) VisitBinary(node *Binary) (interface{}, error) {
	left, err := node.left.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := node.right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s %s)", node.token.String(), left.(string), right.(string)), nil
}

func (e AstPrinter) VisitLiteral(node *Literal) (interface{}, error) {
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

func (e AstPrinter) VisitVariable(b *Variable) (interface{}, error) {
	return b.name, nil
}

func (e AstPrinter) VisitUnary(node *Unary) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s)", node.token.String(), expr), nil
}

func (e AstPrinter) VisitComma(node *Comma) (interface{}, error) {
	left, err := node.left.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := node.right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(, %s %s)", left, right), nil
}

func (e AstPrinter) VisitTernary(node *Ternary) (interface{}, error) {
	left, err := node.left.Accept(e)
	if err != nil {
		return nil, err
	}
	middle, err := node.middle.Accept(e)
	if err != nil {
		return nil, err
	}
	right, err := node.right.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(TERNARY %s %s %s)", left.(string), middle.(string), right.(string)), nil
}

func (e AstPrinter) VisitAssign(b *Assign) (interface{}, error) {
	value, err := b.value.Accept(e)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("(ASSIGN %s %s)", b.name, value), nil
}

func (e AstPrinter) VisitGrouping(node *Grouping) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(group %s)", expr), nil
}

func Print(expr Expr) (string, error) {
	printer := AstPrinter{}
	res, err := expr.Accept(printer)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}
