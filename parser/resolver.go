package parser

import (
	"errors"
	"fmt"
)

type FunctionType = string

const (
	NONE     FunctionType = "NONE"
	FUNCTION FunctionType = "FUNCTION"
)

type Resolver struct {
	interpreter     *Interpreter
	currentFunction FunctionType
	scopes          Stack[map[string]bool]
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, currentFunction: NONE, scopes: Stack[map[string]bool]{}}
}

func (r *Resolver) VisitPrintStmt(node *PrintStmt) (interface{}, error) {
	return r.resolveExpr(node.expr)
}

func (r *Resolver) VisitExprStmt(node *ExprStmt) (interface{}, error) {
	return r.resolveExpr(node.expr)
}

func (r *Resolver) VisitBinary(node *Binary) (interface{}, error) {
	_, err := r.resolveExpr(node.left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(node.right)
}

func (r *Resolver) VisitLiteral(node *Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitUnary(node *Unary) (interface{}, error) {
	return r.resolveExpr(node.expr)
}

func (r *Resolver) VisitComma(node *Comma) (interface{}, error) {
	_, err := r.resolveExpr(node.left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(node.right)
}

func (r *Resolver) VisitGrouping(node *Grouping) (interface{}, error) {
	return r.resolveExpr(node.expr)
}

func (r *Resolver) VisitTernary(b *Ternary) (interface{}, error) {
	_, err := r.resolveExpr(b.left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(b.middle)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(b.right)
}

func (r *Resolver) VisitVarDecl(vd *VarDecl) (interface{}, error) {
	err := r.declare(vd.name)
	if err != nil {
		return nil, err
	}
	defer r.define(vd.name)

	if vd.initializer != nil {
		return r.resolveExpr(vd.initializer)
	}

	return nil, nil
}

func (r *Resolver) VisitVariable(b *Variable) (interface{}, error) {
	if scope, ok := r.scopes.Peek(); ok {
		if defined, ok := scope[b.name]; ok && !defined {
			return nil, errors.New("can't access variable '" + b.name + "' in initializer")
		}
	}

	r.resolveLocal(b, b.name)

	return nil, nil
}

func (r *Resolver) resolveLocal(expr Expr, name string) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope := r.scopes.Get(i)
		if _, ok := scope[name]; ok {
			r.interpreter.resolve(expr, r.scopes.Size()-i-1)
			break
		}
	}
}

func (r *Resolver) VisitAssign(b *Assign) (interface{}, error) {
	_, err := r.resolveExpr(b.value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(b, b.name)
	return nil, nil
}

func (r *Resolver) VisitBlock(b *Block) (interface{}, error) {
	r.beginScope()
	_, err := r.resolveStmts(b.stmts)
	if err != nil {
		return nil, err
	}
	r.endScope()

	return nil, nil
}

func (r *Resolver) VisitIfStatement(s *IfStatement) (interface{}, error) {
	_, err := r.resolveExpr(s.cond)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveStmt(s.ifBlock)
	if err != nil {
		return nil, err
	}

	if s.elseBlock != nil {
		_, err := r.resolveStmt(s.elseBlock)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitWhileStatement(s *WhileStatement) (interface{}, error) {
	_, err := r.resolveExpr(s.cond)
	if err != nil {
		return nil, err
	}

	return r.resolveStmt(s.body)
}

func (r *Resolver) VisitForStatement(s *ForStatement) (interface{}, error) {
	panic("for statement should have been desugared to while loop")
}

func (r *Resolver) VisitLogical(b *Logical) (interface{}, error) {
	_, err := r.resolveExpr(b.left)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(b.right)
}

func (r *Resolver) VisitContinueStmt(c *ContinueStmt) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitBreakStmt(b *BreakStmt) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitCall(c *Call) (interface{}, error) {
	_, err := r.resolveExpr(c.callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range c.arguments {
		_, err := r.resolveExpr(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitFunction(f *Function) (interface{}, error) {
	err := r.declare(f.name.Lexeme.(string))
	if err != nil {
		return nil, err
	}
	r.define(f.name.Lexeme.(string))

	err = r.resolveFunction(f)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) resolveFunction(f *Function) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION
	defer func() { r.currentFunction = enclosingFunction }()

	r.beginScope()
	defer r.endScope()

	for _, param := range f.params {
		err := r.declare(param.Lexeme.(string))
		if err != nil {
			return err
		}
		r.define(param.Lexeme.(string))
	}

	_, err := r.resolveStmts(f.body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitReturnStmt(s *ReturnStmt) (interface{}, error) {
	if r.currentFunction == NONE {
		return nil, errors.New("can't return from top-level code")
	}

	if s.expr != nil {
		return r.resolveExpr(s.expr)
	}

	return nil, nil
}

func (r *Resolver) resolveExpr(expr Expr) (interface{}, error) {
	_, err := expr.Accept(r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) resolveStmt(stmt Stmt) (interface{}, error) {
	_, err := stmt.Accept(r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) resolveStmts(stmts []Stmt) (interface{}, error) {
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
			return fmt.Errorf("Already a variable with this name in this scope.")
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
