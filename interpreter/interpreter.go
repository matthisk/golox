package interpreter

import (
	"errors"
	"fmt"

	"github.com/matthisk/lox/lexer"
	"github.com/matthisk/lox/parser"
)

type ControlFlowStmtT string

const (
	BREAK    ControlFlowStmtT = "BREAK"
	CONTINUE                  = "CONTINUE"
	RETURN                    = "RETURN"
)

type ControlFlowStmt struct {
	t ControlFlowStmtT
	v interface{}
}

func Break() ControlFlowStmt {
	return ControlFlowStmt{
		t: BREAK,
	}
}

func Continue() ControlFlowStmt {
	return ControlFlowStmt{
		t: CONTINUE,
	}
}

func Return(v interface{}) ControlFlowStmt {
	return ControlFlowStmt{
		t: RETURN,
		v: v,
	}
}

type Interpreter struct {
	printer Printer
	env     *Environment
	globals *Environment
	locals  map[parser.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)

	globals.Define("clock", &Clock{})

	return &Interpreter{
		printer: DefaultPrinter{},
		env:     globals,
		globals: globals,
		locals:  make(map[parser.Expr]int),
	}
}

func NewInterpreterWithPrinter(printer Printer) *Interpreter {
	globals := NewEnvironment(nil)

	globals.Define("clock", &Clock{})

	return &Interpreter{
		printer: printer,
		env:     globals,
		globals: globals,
		locals:  make(map[parser.Expr]int),
	}
}

func (i *Interpreter) VisitThis(t *parser.This) (interface{}, error) {
	if d, ok := i.locals[t]; ok {
		return i.env.GetAt(d, "this")
	}
	return i.globals.Get("this")
}

func (i *Interpreter) VisitSet(s *parser.SetExpr) (interface{}, error) {
	name := s.Name.Lexeme.(string)
	object, err := i.evaluate(s.Object)
	if err != nil {
		return nil, err
	}

	if obj, ok := object.(*LoxInstance); ok {
		val, err := i.evaluate(s.Value)
		if err != nil {
			return nil, err
		}

		obj.Set(name, val)
		return nil, nil
	} else {
		return nil, fmt.Errorf("'%s' only instances have fields", name)
	}
}

func (i *Interpreter) VisitGet(g *parser.GetExpr) (interface{}, error) {
	instance, err := i.evaluate(g.From)
	if err != nil {
		return nil, err
	}

	if ins, ok := instance.(*LoxInstance); ok {
		return ins.Get(g.Property.Lexeme.(string))
	}

	return nil, errors.New("only instances have fields")
}

func (i *Interpreter) VisitClass(s *parser.ClassStatement) (interface{}, error) {
	name := s.Name.Lexeme.(string)

	i.env.Define(name, nil)
	class := &LoxClass{
		name:    name,
		methods: make(map[string]LoxCallable),
	}

	for _, method := range s.Methods {
		class.methods[method.Name.Lexeme.(string)] = &LoxFunction{
			declaration:   method,
			closure:       i.env,
			isInitializer: method.Name.Lexeme.(string) == "init",
		}
	}

	err := i.env.Assign(name, class)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitFunction(f *parser.Function) (interface{}, error) {
	fun := &LoxFunction{
		declaration:   f,
		closure:       i.env,
		isInitializer: false,
	}
	i.env.Define(f.Name.Lexeme.(string), fun)
	return nil, nil
}

func (i *Interpreter) VisitBlock(b *parser.Block) (interface{}, error) {
	i.env = NewEnvironment(i.env)
	defer func() { i.env = i.env.enclosing }()

	for _, stmt := range b.Stmts {
		result, err := stmt.Accept(i)
		if err != nil {
			return nil, err
		}

		// in case we encounter a return, break or continue statement we return back to
		// the control flow statement.
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitIfStatement(s *parser.IfStatement) (interface{}, error) {
	cond, err := s.Cond.Accept(i)
	if err != nil {
		return nil, err
	}

	if isTruthy(cond) {
		return s.IfBlock.Accept(i)
	} else {
		if s.ElseBlock != nil {
			return s.ElseBlock.Accept(i)
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitWhileStatement(s *parser.WhileStatement) (interface{}, error) {
	cond, err := s.Cond.Accept(i)
	if err != nil {
		return nil, err
	}

	for isTruthy(cond) {
		result, err := s.Body.Accept(i)
		if err != nil {
			return nil, err
		}

		if c, ok := result.(ControlFlowStmt); ok {
			if c.t == BREAK {
				break
			}

			if c.t == RETURN {
				return result, nil
			}
		}

		cond, err = s.Cond.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitForStatement(s *parser.ForStatement) (interface{}, error) {
	// We don't implement for statements, they are desugared to a while loop by the parser
	return nil, nil
}

func (i *Interpreter) VisitBreakStmt(b *parser.BreakStmt) (interface{}, error) {
	return Break(), nil
}

func (i *Interpreter) VisitContinueStmt(c *parser.ContinueStmt) (interface{}, error) {
	return Continue(), nil
}

func (i *Interpreter) VisitReturnStmt(r *parser.ReturnStmt) (interface{}, error) {
	var val interface{}
	var err error

	if r.Expr != nil {
		val, err = r.Expr.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	return Return(val), nil
}

func (i *Interpreter) VisitPrintStmt(node *parser.PrintStmt) (interface{}, error) {
	expr, err := node.Expr.Accept(i)
	if err != nil {
		return nil, err
	}

	i.printer.Print(expr)
	return nil, nil
}

func (i *Interpreter) VisitCall(c *parser.Call) (interface{}, error) {
	callee, err := i.evaluate(c.Callee)
	if err != nil {
		return nil, err
	}

	var args []interface{}
	for _, argument := range c.Arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	if c, ok := callee.(LoxCallable); ok {
		if len(args) != c.Arity() {
			return nil, fmt.Errorf("Expected %d arguments but got %d", c.Arity(), len(args))
		}

		return c.Call(i, args)
	} else {
		return nil, errors.New("Can only call functions and classes.")
	}
}

func (i *Interpreter) VisitExprStmt(node *parser.ExprStmt) (interface{}, error) {
	_, err := node.Expr.Accept(i)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitVarDecl(vd *parser.VarDecl) (interface{}, error) {
	if vd.Initializer != nil {
		val, err := vd.Initializer.Accept(i)
		if err != nil {
			return nil, err
		}

		i.env.Define(vd.Name, val)
	} else {
		i.env.Define(vd.Name, nil)
	}

	return nil, nil
}

func (i *Interpreter) VisitAssign(b *parser.Assign) (interface{}, error) {
	value, err := i.evaluate(b.Value)
	if err != nil {
		return nil, err
	}

	err = i.env.Assign(b.Name, value)

	return value, err
}

func (i *Interpreter) VisitLogical(b *parser.Logical) (interface{}, error) {
	left, err := b.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	if b.Token == lexer.OR {
		if isTruthy(left) {
			return left, nil
		}
	} else if b.Token == lexer.AND {
		if !isTruthy(left) {
			return left, nil
		}
	} else {
		panic("Illegal AST node Logical with token type")
	}

	return i.evaluate(b.Right)
}

func (i *Interpreter) VisitBinary(node *parser.Binary) (interface{}, error) {
	l, err := i.evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	r, err := i.evaluate(node.Right)
	if err != nil {
		return nil, err
	}

	switch node.Token {
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
		err := isNumber(node.Token, l, r)
		if err != nil {
			return nil, err
		}

		return l.(float64) - r.(float64), nil
	case lexer.SLASH:
		err := isNumber(node.Token, l, r)
		if err != nil {
			return nil, err
		}

		return l.(float64) / r.(float64), nil
	case lexer.STAR:
		err := isNumber(node.Token, l, r)
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

func (i *Interpreter) VisitLiteral(node *parser.Literal) (interface{}, error) {
	switch node.Token.Type {
	case lexer.TRUE:
		return true, nil
	case lexer.FALSE:
		return false, nil
	case lexer.NIL:
		return nil, nil
	case lexer.NUMBER, lexer.STRING:
		return node.Token.Lexeme, nil
	default:
		return node.Token.Lexeme, nil
	}
}

func (i *Interpreter) VisitUnary(node *parser.Unary) (interface{}, error) {
	val, err := i.evaluate(node.Expr)
	if err != nil {
		return nil, err
	}

	switch node.Token {
	case lexer.BANG:
		return !isTruthy(val), nil
	case lexer.MINUS:
		return -val.(float64), nil
	default:
		return nil, fmt.Errorf("unhandled default case")
	}
}

func (i *Interpreter) VisitComma(node *parser.Comma) (interface{}, error) {
	_, err := i.evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	return i.evaluate(node.Right)
}

func (i *Interpreter) VisitGrouping(node *parser.Grouping) (interface{}, error) {
	return i.evaluate(node.Expr)
}

func (i *Interpreter) VisitTernary(node *parser.Ternary) (interface{}, error) {
	left, err := i.evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	if isTruthy(left) {
		return i.evaluate(node.Middle)
	}
	return i.evaluate(node.Right)
}

func (i *Interpreter) VisitVariable(b *parser.Variable) (interface{}, error) {
	if depth, ok := i.locals[b]; ok {
		return i.env.GetAt(depth, b.Name)
	}
	return i.globals.Get(b.Name)
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

func (i *Interpreter) evaluate(expr parser.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) Run(stmts []parser.Stmt) error {
	for _, stmt := range stmts {
		_, err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) EvaluateExpression(expr parser.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) executeBlock(body []parser.Stmt, env *Environment) (interface{}, error) {
	oldEnv := i.env
	i.env = env
	defer func() { i.env = oldEnv }()

	for j := range body {
		res, err := body[j].Accept(i)
		if err != nil {
			return nil, err
		}

		// Return statements break the execution of the function body.
		if c, ok := res.(ControlFlowStmt); ok && c.t == RETURN {
			return c.v, nil
		}
	}

	return nil, nil
}

func (i *Interpreter) resolve(expr parser.Expr, depth int) {
	i.locals[expr] = depth
}
