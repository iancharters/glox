package parser

import (
	"fmt"
	"strings"
)

type Printer struct{}

func (p *Printer) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *Printer) VisitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *Printer) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}

	return fmt.Sprintf("%v", expr.Value), nil
}

func (p *Printer) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p *Printer) parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString("(")
	sb.WriteString(name)

	for _, ex := range exprs {
		sb.WriteString(" ")
		val, _ := ex.Accept(p)
		sb.WriteString(fmt.Sprint(val))
	}

	sb.WriteString(")")

	return sb.String()
}

func (p *Printer) Print(expr Expr) string {
	output, _ := expr.Accept(p)
	return fmt.Sprintf("%v", output)
}
