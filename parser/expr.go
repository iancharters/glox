package parser

import (
	"github.com/iancharters/glox/token"
)

type Expr interface {
	Accept(visitor Visitor) (interface{}, error)
}

type BinaryExpr struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func NewBinaryExpr(left Expr, operator *token.Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left, operator, right}
}

func (b *BinaryExpr) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression}
}

func (g *GroupingExpr) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value interface{}
}

func NewLiteralExpr(value interface{}) Expr {
	return &LiteralExpr{value}
}

func (l *LiteralExpr) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator *token.Token
	Right    Expr
}

func NewUnaryExpr(operator *token.Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator, right}
}

func (u *UnaryExpr) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(u)
}
