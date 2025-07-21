package parser

import (
	"fmt"

	"github.com/matthisk/lox/lexer"
)

type Interpreter struct {
}

func (i Interpreter) VisitBinary(node *Binary) (interface{}, error) {
	l, err := i.evaluate(node.left)
	if err != nil {
		return nil, err
	}
	r, err := i.evaluate(node.right)
	if err != nil {
		return nil, err
	}

	switch node.token {
	case lexer.PLUS:
		if ls, ok := l.(string); ok {
			if rs, ok := r.(string); ok {
				return ls + rs, nil
			}
		}

		if ls, ok := l.(float64); ok {
			if rs, ok := r.(float64); ok {
				return ls + rs, nil
			}
		}
	case lexer.MINUS:
		err := isNumber(node.token, l, r)
		if err != nil {
			return nil, err
		}

		return l.(float64) - r.(float64), nil
	case lexer.SLASH:
		err := isNumber(node.token, l, r)
		if err != nil {
			return nil, err
		}

		return l.(float64) / r.(float64), nil
	case lexer.STAR:
		err := isNumber(node.token, l, r)
		if err != nil {
			return nil, err
		}

		return l.(float64) * r.(float64), nil
	case lexer.EQUAL_EQUAL:
		return isEqual(l, r), nil
	case lexer.BANG_EQUAL:
		return !isEqual(l, r), nil
	case lexer.GREATER:
		return l.(float64) > r.(float64), nil
	case lexer.GREATER_EQUAL:
		return l.(float64) >= r.(float64), nil
	case lexer.LESS:
		return l.(float64) < r.(float64), nil
	case lexer.LESS_EQUAL:
		return l.(float64) <= r.(float64), nil
	default:
		panic("unhandled default case")
	}

	return nil, fmt.Errorf("unsupported binary operation or operand types")
}

func (i Interpreter) VisitLiteral(node *Literal) (interface{}, error) {
	switch node.token.Type {
	case lexer.TRUE:
		return true, nil
	case lexer.FALSE:
		return false, nil
	case lexer.NIL:
		return nil, nil
	case lexer.NUMBER, lexer.STRING:
		return node.token.Lexeme, nil
	default:
		return node.token.Lexeme, nil
	}
}

func (i Interpreter) VisitUnary(node *Unary) (interface{}, error) {
	val, err := i.evaluate(node.expr)
	if err != nil {
		return nil, err
	}

	switch node.token {
	case lexer.BANG:
		return !isTruthy(val), nil
	case lexer.MINUS:
		return -val.(float64), nil
	default:
		return nil, fmt.Errorf("unhandled default case")
	}
}

func (i Interpreter) VisitComma(node *Comma) (interface{}, error) {
	_, err := i.evaluate(node.left)
	if err != nil {
		return nil, err
	}
	return i.evaluate(node.right)
}

func (i Interpreter) VisitGrouping(node *Grouping) (interface{}, error) {
	return i.evaluate(node.expr)
}

func (i Interpreter) VisitTernary(node *Ternary) (interface{}, error) {
	left, err := i.evaluate(node.left)
	if err != nil {
		return nil, err
	}
	if isTruthy(left) {
		return i.evaluate(node.middle)
	}
	return i.evaluate(node.right)
}

func isNumber(op lexer.TokenType, l, r interface{}) error {
	if _, ok := l.(float64); !ok {
		return fmt.Errorf("operand %v is not a number for %v", l, op)
	}
	if _, ok := r.(float64); !ok {
		return fmt.Errorf("operand %v is not a number for %v", r, op)
	}
	return nil
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func isEqual(l, r interface{}) bool {
	if l == nil && r == nil {
		return true
	}
	if l == nil || r == nil {
		return false
	}
	if b1, ok1 := l.(bool); ok1 {
		if b2, ok2 := r.(bool); ok2 {
			return b1 == b2
		}
	}
	if f1, ok1 := l.(float64); ok1 {
		if f2, ok2 := r.(float64); ok2 {
			return f1 == f2
		}
	}
	if s1, ok1 := l.(string); ok1 {
		if s2, ok2 := r.(string); ok2 {
			return s1 == s2
		}
	}
	return false
}

func (i Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i Interpreter) Run(expr Expr) (interface{}, error) {
	return i.evaluate(expr)
}
