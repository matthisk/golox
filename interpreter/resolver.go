package interpreter

import (
	"errors"
	"fmt"

	"github.com/matthisk/lox/ds"
	"github.com/matthisk/lox/parser"
)

type FunctionType = string

const (
	NONE     FunctionType = "NONE"
	FUNCTION FunctionType = "FUNCTION"
	METHOD   FunctionType = "METHOD"
)

type Resolver struct {
	interpreter     *Interpreter
	currentFunction FunctionType
	scopes          ds.Stack[map[string]bool]
}

func (r *Resolver) VisitSet(s *parser.SetExpr) (interface{}, error) {
	_, err := r.resolveExpr(s.Object)
	if err != nil {
		return nil, err
	}
	return r.resolveExpr(s.Value)
}

func (r *Resolver) VisitClass(s *parser.ClassStatement) (interface{}, error) {
	err := r.declare(s.Name.Lexeme.(string))
	if err != nil {
		return nil, err
	}
	r.define(s.Name.Lexeme.(string))

	for _, method := range s.Methods {
		err := r.resolveFunction(method, METHOD)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, currentFunction: NONE, scopes: ds.Stack[map[string]bool]{}}
}

func (r *Resolver) VisitPrintStmt(node *parser.PrintStmt) (interface{}, error) {
	return r.resolveExpr(node.Expr)
}

func (r *Resolver) VisitExprStmt(node *parser.ExprStmt) (interface{}, error) {
	return r.resolveExpr(node.Expr)
}

func (r *Resolver) VisitBinary(node *parser.Binary) (interface{}, error) {
	_, err := r.resolveExpr(node.Left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(node.Right)
}

func (r *Resolver) VisitLiteral(node *parser.Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitUnary(node *parser.Unary) (interface{}, error) {
	return r.resolveExpr(node.Expr)
}

func (r *Resolver) VisitComma(node *parser.Comma) (interface{}, error) {
	_, err := r.resolveExpr(node.Left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(node.Right)
}

func (r *Resolver) VisitGrouping(node *parser.Grouping) (interface{}, error) {
	return r.resolveExpr(node.Expr)
}

func (r *Resolver) VisitGet(g *parser.GetExpr) (interface{}, error) {
	return r.resolveExpr(g.From)
}

func (r *Resolver) VisitTernary(b *parser.Ternary) (interface{}, error) {
	_, err := r.resolveExpr(b.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(b.Middle)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(b.Right)
}

func (r *Resolver) VisitVarDecl(vd *parser.VarDecl) (interface{}, error) {
	err := r.declare(vd.Name)
	if err != nil {
		return nil, err
	}
	defer r.define(vd.Name)

	if vd.Initializer != nil {
		return r.resolveExpr(vd.Initializer)
	}

	return nil, nil
}

func (r *Resolver) VisitVariable(b *parser.Variable) (interface{}, error) {
	if scope, ok := r.scopes.Peek(); ok {
		if defined, ok := scope[b.Name]; ok && !defined {
			return nil, errors.New("can't access variable '" + b.Name + "' in initializer")
		}
	}

	r.resolveLocal(b, b.Name)

	return nil, nil
}

func (r *Resolver) resolveLocal(expr parser.Expr, name string) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope := r.scopes.Get(i)
		if _, ok := scope[name]; ok {
			r.interpreter.resolve(expr, r.scopes.Size()-i-1)
			break
		}
	}
}

func (r *Resolver) VisitAssign(b *parser.Assign) (interface{}, error) {
	_, err := r.resolveExpr(b.Value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(b, b.Name)
	return nil, nil
}

func (r *Resolver) VisitBlock(b *parser.Block) (interface{}, error) {
	r.beginScope()
	_, err := r.resolveStmts(b.Stmts)
	if err != nil {
		return nil, err
	}
	r.endScope()

	return nil, nil
}

func (r *Resolver) VisitIfStatement(s *parser.IfStatement) (interface{}, error) {
	_, err := r.resolveExpr(s.Cond)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveStmt(s.IfBlock)
	if err != nil {
		return nil, err
	}

	if s.ElseBlock != nil {
		_, err := r.resolveStmt(s.ElseBlock)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitWhileStatement(s *parser.WhileStatement) (interface{}, error) {
	_, err := r.resolveExpr(s.Cond)
	if err != nil {
		return nil, err
	}

	return r.resolveStmt(s.Body)
}

func (r *Resolver) VisitForStatement(s *parser.ForStatement) (interface{}, error) {
	panic("for statement should have been desugared to while loop")
}

func (r *Resolver) VisitLogical(b *parser.Logical) (interface{}, error) {
	_, err := r.resolveExpr(b.Left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(b.Right)
}

func (r *Resolver) VisitContinueStmt(c *parser.ContinueStmt) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitBreakStmt(b *parser.BreakStmt) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitCall(c *parser.Call) (interface{}, error) {
	_, err := r.resolveExpr(c.Callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range c.Arguments {
		_, err := r.resolveExpr(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitFunction(f *parser.Function) (interface{}, error) {
	err := r.declare(f.Name.Lexeme.(string))
	if err != nil {
		return nil, err
	}
	r.define(f.Name.Lexeme.(string))

	err = r.resolveFunction(f, FUNCTION)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) resolveFunction(f *parser.Function, t FunctionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = t
	defer func() { r.currentFunction = enclosingFunction }()

	r.beginScope()
	defer r.endScope()

	for _, param := range f.Params {
		err := r.declare(param.Lexeme.(string))
		if err != nil {
			return err
		}
		r.define(param.Lexeme.(string))
	}

	_, err := r.resolveStmts(f.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitReturnStmt(s *parser.ReturnStmt) (interface{}, error) {
	if r.currentFunction == NONE {
		return nil, errors.New("can't return from top-level code")
	}

	if s.Expr != nil {
		return r.resolveExpr(s.Expr)
	}

	return nil, nil
}

func (r *Resolver) resolveExpr(expr parser.Expr) (interface{}, error) {
	_, err := expr.Accept(r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) resolveStmt(stmt parser.Stmt) (interface{}, error) {
	_, err := stmt.Accept(r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) resolveStmts(stmts []parser.Stmt) (interface{}, error) {
	for i := range stmts {
		_, err := stmts[i].Accept(r)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) beginScope() {
	scope := make(map[string]bool)
	r.scopes.Push(scope)
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name string) error {
	if scope, ok := r.scopes.Peek(); ok {
		if _, ok := scope[name]; ok {
			return fmt.Errorf("Already a variable with this Name in this scope.")
		}
		scope[name] = false
	}
	return nil
}

func (r *Resolver) define(name string) {
	if scope, ok := r.scopes.Peek(); ok {
		scope[name] = true
	}
}
