package parser

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
}
