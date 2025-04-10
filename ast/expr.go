package ast

import "github.com/Toolnado/sludge/token"

type GroupingExpr struct {
	Expession Expr
}

func NewGroupingExpr(Expession Expr) *GroupingExpr {
	return &GroupingExpr{
		Expession: Expession,
	}
}

func (g *GroupingExpr) Accept(v IASTVisitor) (any, error) { return v.VisitGroupingExpr(g)}

type VariableExpr struct {
	Name token.Token
}

func NewVariableExpr(Name token.Token) *VariableExpr {
	return &VariableExpr{
		Name: Name,
	}
}

func (va *VariableExpr) Accept(v IASTVisitor) (any, error) { return v.VisitVariableExpr(va)}

type AssignExpr struct {
	Name token.Token
	Value Expr
}

func NewAssignExpr(Name token.Token, Value Expr) *AssignExpr {
	return &AssignExpr{
		Name: Name,
		Value: Value,
	}
}

func (a *AssignExpr) Accept(v IASTVisitor) (any, error) { return v.VisitAssignExpr(a)}

type LogicalExpr struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewLogicalExpr(Left Expr, Operator token.Token, Right Expr) *LogicalExpr {
	return &LogicalExpr{
		Left: Left,
		Operator: Operator,
		Right: Right,
	}
}

func (l *LogicalExpr) Accept(v IASTVisitor) (any, error) { return v.VisitLogicalExpr(l)}

type BinaryExpr struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinaryExpr(Left Expr, Operator token.Token, Right Expr) *BinaryExpr {
	return &BinaryExpr{
		Left: Left,
		Operator: Operator,
		Right: Right,
	}
}

func (b *BinaryExpr) Accept(v IASTVisitor) (any, error) { return v.VisitBinaryExpr(b)}

type UnaryExpr struct {
	Operator token.Token
	Right Expr
}

func NewUnaryExpr(Operator token.Token, Right Expr) *UnaryExpr {
	return &UnaryExpr{
		Operator: Operator,
		Right: Right,
	}
}

func (u *UnaryExpr) Accept(v IASTVisitor) (any, error) { return v.VisitUnaryExpr(u)}

type LiteralExpr struct {
	Value any
}

func NewLiteralExpr(Value any) *LiteralExpr {
	return &LiteralExpr{
		Value: Value,
	}
}

func (l *LiteralExpr) Accept(v IASTVisitor) (any, error) { return v.VisitLiteralExpr(l)}

