package parser

import (
	"github.com/matthisk/lox/lexer"
)

type Visitor interface {
	VisitPrintStmt(node *PrintStmt) (interface{}, error)
	VisitExprStmt(node *ExprStmt) (interface{}, error)
	VisitBinary(node *Binary) (interface{}, error)
	VisitLiteral(node *Literal) (interface{}, error)
	VisitUnary(node *Unary) (interface{}, error)
	VisitComma(node *Comma) (interface{}, error)
	VisitGrouping(node *Grouping) (interface{}, error)
	VisitTernary(b *Ternary) (interface{}, error)
	VisitVarDecl(vd *VarDecl) (interface{}, error)
	VisitVariable(b *Variable) (interface{}, error)
	VisitAssign(b *Assign) (interface{}, error)
	VisitBlock(b *Block) (interface{}, error)
	VisitIfStatement(s *IfStatement) (interface{}, error)
	VisitWhileStatement(s *WhileStatement) (interface{}, error)
	VisitForStatement(s *ForStatement) (interface{}, error)
	VisitLogical(b *Logical) (interface{}, error)
	VisitContinueStmt(c *ContinueStmt) (interface{}, error)
	VisitBreakStmt(b *BreakStmt) (interface{}, error)
	VisitCall(c *Call) (interface{}, error)
	VisitFunction(f *Function) (interface{}, error)
	VisitReturnStmt(r *ReturnStmt) (interface{}, error)
	VisitClass(s *ClassStatement) (interface{}, error)
	VisitGet(g *GetExpr) (interface{}, error)
	VisitSet(s *SetExpr) (interface{}, error)
	VisitThis(t *This) (interface{}, error)
}

type Stmt interface {
	Accept(v Visitor) (interface{}, error)
}

type Block struct {
	BaseNode
	Stmts []Stmt
}

func (b *Block) Accept(v Visitor) (interface{}, error) {
	return v.VisitBlock(b)
}

type VarDecl struct {
	BaseNode
	Name        string
	Initializer Expr
}

func (vd *VarDecl) Accept(v Visitor) (interface{}, error) {
	return v.VisitVarDecl(vd)
}

type ExprStmt struct {
	BaseNode
	Expr Expr
}

func (e *ExprStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitExprStmt(e)
}

type PrintStmt struct {
	BaseNode
	Expr Expr
}

func (p *PrintStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitPrintStmt(p)
}

type ReturnStmt struct {
	BaseNode
	Expr Expr
}

func (r *ReturnStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitReturnStmt(r)
}

type BreakStmt struct {
	BaseNode
}

func (b *BreakStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitBreakStmt(b)
}

type ContinueStmt struct {
	BaseNode
}

func (c *ContinueStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitContinueStmt(c)
}

type IfStatement struct {
	BaseNode
	Cond      Expr
	IfBlock   Stmt
	ElseBlock Stmt
}

func (s *IfStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitIfStatement(s)
}

type WhileStatement struct {
	BaseNode
	Cond Expr
	Body Stmt
}

func (s *WhileStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitWhileStatement(s)
}

type ClassStatement struct {
	BaseNode
	Name    lexer.Token
	Methods []*Function
}

func (s *ClassStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitClass(s)
}

type Function struct {
	BaseNode
	Name   lexer.Token
	Params []lexer.Token
	Body   []Stmt
}

func (f *Function) Accept(v Visitor) (interface{}, error) {
	return v.VisitFunction(f)
}

type ForStatement struct {
	BaseNode
	Initializer Stmt // Can be VarDecl, ExprStmt, or nil
	Condition   Expr // Can be nil
	Increment   Expr // Can be nil
	Body        Stmt
}

func (s *ForStatement) Accept(v Visitor) (interface{}, error) {
	return v.VisitForStatement(s)
}

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type BaseNode struct {
	StartPos lexer.Pos
	EndPos   lexer.Pos
}

// GetStartPos returns the start position of the BaseNode
func (b *BaseNode) GetStartPos() lexer.Pos {
	return b.StartPos
}

// GetEndPos returns the end position of the BaseNode
func (b *BaseNode) GetEndPos() lexer.Pos {
	return b.EndPos
}

// SetPos sets both start and end positions of the BaseNode
func (b *BaseNode) SetPos(start, end lexer.Pos) {
	b.StartPos = start
	b.EndPos = end
}

type Call struct {
	BaseNode
	Callee    Expr
	Paren     lexer.Token
	Arguments []Expr
}

func (c *Call) Accept(v Visitor) (interface{}, error) {
	return v.VisitCall(c)
}

type Binary struct {
	BaseNode
	Left  Expr
	Token lexer.TokenType
	Right Expr
}

func (b *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinary(b)
}

type Logical struct {
	BaseNode
	Left  Expr
	Token lexer.TokenType
	Right Expr
}

func (b *Logical) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLogical(b)
}

type Unary struct {
	BaseNode
	Token lexer.TokenType
	Expr  Expr
}

func (b *Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnary(b)
}

type Comma struct {
	BaseNode
	Left  Expr
	Right Expr
}

func (b *Comma) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitComma(b)
}

type Ternary struct {
	BaseNode
	Left   Expr
	Middle Expr
	Right  Expr
}

func (b *Ternary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTernary(b)
}

type Grouping struct {
	BaseNode
	Expr Expr
}

func (b *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGrouping(b)
}

type Literal struct {
	BaseNode
	Token lexer.Token
}

func (b *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteral(b)
}

type GetExpr struct {
	BaseNode
	From     Expr
	Property lexer.Token
}

func (g *GetExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitGet(g)
}

type SetExpr struct {
	BaseNode
	Object Expr
	Name   lexer.Token
	Value  Expr
}

func (s *SetExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitSet(s)
}

type This struct {
	BaseNode
}

func (t *This) Accept(v Visitor) (interface{}, error) {
	return v.VisitThis(t)
}

type Variable struct {
	BaseNode
	Name string
}

func (b *Variable) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitVariable(b)
}

type Assign struct {
	BaseNode
	Name  string
	Value Expr
}

func (b *Assign) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAssign(b)
}

// Helper functions for working with expression positions

// GetExprStartPos returns the start position of any expression
func GetExprStartPos(expr Expr) lexer.Pos {
	switch e := expr.(type) {
	case *Binary:
		return e.GetStartPos()
	case *Unary:
		return e.GetStartPos()
	case *Comma:
		return e.GetStartPos()
	case *Ternary:
		return e.GetStartPos()
	case *Grouping:
		return e.GetStartPos()
	case *Literal:
		return e.GetStartPos()
	case *Variable:
		return e.GetStartPos()
	case *Assign:
		return e.GetStartPos()
	default:
		return lexer.Pos{}
	}
}

// GetExprEndPos returns the end position of any expression
func GetExprEndPos(expr Expr) lexer.Pos {
	switch e := expr.(type) {
	case *Binary:
		return e.GetEndPos()
	case *Unary:
		return e.GetEndPos()
	case *Comma:
		return e.GetEndPos()
	case *Ternary:
		return e.GetEndPos()
	case *Grouping:
		return e.GetEndPos()
	case *Literal:
		return e.GetEndPos()
	case *Variable:
		return e.GetEndPos()
	case *Assign:
		return e.GetEndPos()
	default:
		return lexer.Pos{}
	}
}

// NewBaseNodeSpanning creates a BaseNode that spans from the start of one position to the end of another
func NewBaseNodeSpanning(start, end lexer.Pos) BaseNode {
	return BaseNode{
		StartPos: start,
		EndPos:   end,
	}
}

// NewBaseNodeFromExprs creates a BaseNode that spans from the start of the first expression to the end of the last expression
func NewBaseNodeFromExprs(first, last Expr) BaseNode {
	return BaseNode{
		StartPos: GetExprStartPos(first),
		EndPos:   GetExprEndPos(last),
	}
}
