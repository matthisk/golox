package parser

import "fmt"

type ExprPrinter struct{}

func (e ExprPrinter) VisitVarDecl(vd *VarDecl) (interface{}, error) {
	initializer, err := vd.initializer.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("var %v = %s;", vd.name.Lexeme, initializer), nil
}

func (e ExprPrinter) VisitPrintStmt(node *PrintStmt) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("print %s;", expr), nil
}

func (e ExprPrinter) VisitExprStmt(node *ExprStmt) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s;", expr), nil
}

func (e ExprPrinter) VisitBinary(node *Binary) (interface{}, error) {
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

func (e ExprPrinter) VisitLiteral(node *Literal) (interface{}, error) {
	switch node.token.Lexeme.(type) {
	case float64:
	case int:
		return fmt.Sprintf("%d", node.token.Lexeme), nil
	case string:
		return node.token.Lexeme, nil
	}

	return nil, fmt.Errorf("Unexpected lexeme type.")
}

func (e ExprPrinter) VisitVariable(b *Variable) (interface{}, error) {
	return b.name, nil
}

func (e ExprPrinter) VisitUnary(node *Unary) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(%s %s)", node.token.String(), expr), nil
}

func (e ExprPrinter) VisitComma(node *Comma) (interface{}, error) {
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

func (e ExprPrinter) VisitTernary(node *Ternary) (interface{}, error) {
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

func (e ExprPrinter) VisitGrouping(node *Grouping) (interface{}, error) {
	expr, err := node.expr.Accept(e)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("(group %s)", expr), nil
}

func Print(expr Expr) (string, error) {
	printer := ExprPrinter{}
	res, err := expr.Accept(printer)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}
