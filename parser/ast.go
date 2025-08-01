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
}

type Stmt interface {
	Accept(v Visitor) (interface{}, error)
}

type VarDecl struct {
	BaseNode
	name        string
	initializer Expr
}

func (vd *VarDecl) Accept(v Visitor) (interface{}, error) {
	return v.VisitVarDecl(vd)
}

type ExprStmt struct {
	BaseNode
	expr Expr
}

func (e *ExprStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitExprStmt(e)
}

type PrintStmt struct {
	BaseNode
	expr Expr
}

func (p *PrintStmt) Accept(v Visitor) (interface{}, error) {
	return v.VisitPrintStmt(p)
}

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type BaseNode struct {
	startPos lexer.Pos
	endPos   lexer.Pos
}

// GetStartPos returns the start position of the BaseNode
func (b *BaseNode) GetStartPos() lexer.Pos {
	return b.startPos
}

// GetEndPos returns the end position of the BaseNode
func (b *BaseNode) GetEndPos() lexer.Pos {
	return b.endPos
}

// SetPos sets both start and end positions of the BaseNode
func (b *BaseNode) SetPos(start, end lexer.Pos) {
	b.startPos = start
	b.endPos = end
}

type Binary struct {
	BaseNode
	left  Expr
	token lexer.TokenType
	right Expr
}

func (b *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinary(b)
}

type Unary struct {
	BaseNode
	token lexer.TokenType
	expr  Expr
}

func (b *Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnary(b)
}

type Comma struct {
	BaseNode
	left  Expr
	right Expr
}

func (b *Comma) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitComma(b)
}

type Ternary struct {
	BaseNode
	left   Expr
	middle Expr
	right  Expr
}

func (b *Ternary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTernary(b)
}

type Grouping struct {
	BaseNode
	expr Expr
}

func (b *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGrouping(b)
}

type Literal struct {
	BaseNode
	token lexer.Token
}

func (b *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteral(b)
}

type Variable struct {
	BaseNode
	name string
}

func (b *Variable) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitVariable(b)
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
	default:
		return lexer.Pos{}
	}
}

// NewBaseNodeSpanning creates a BaseNode that spans from the start of one position to the end of another
func NewBaseNodeSpanning(start, end lexer.Pos) BaseNode {
	return BaseNode{
		startPos: start,
		endPos:   end,
	}
}

// NewBaseNodeFromExprs creates a BaseNode that spans from the start of the first expression to the end of the last expression
func NewBaseNodeFromExprs(first, last Expr) BaseNode {
	return BaseNode{
		startPos: GetExprStartPos(first),
		endPos:   GetExprEndPos(last),
	}
}
