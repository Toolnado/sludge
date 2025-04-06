package ast

import "github.com/Toolnado/sludge/token"

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewBinary(Left Expr, Operator token.Token, Right Expr) *Binary {
	return &Binary{
		Left:     Left,
		Operator: Operator,
		Right:    Right,
	}
}

func (b *Binary) Accept(v IASTVisitor) (any, error) { return v.VisitBinaryExpr(b) }

type Unary struct {
	Operator token.Token
	Right    Expr
}

func NewUnary(Operator token.Token, Right Expr) *Unary {
	return &Unary{
		Operator: Operator,
		Right:    Right,
	}
}

func (u *Unary) Accept(v IASTVisitor) (any, error) { return v.VisitUnaryExpr(u) }

type Literal struct {
	Value any
}

func NewLiteral(Value any) *Literal {
	return &Literal{
		Value: Value,
	}
}

func (l *Literal) Accept(v IASTVisitor) (any, error) { return v.VisitLiteralExpr(l) }

type Grouping struct {
	Expession Expr
}

func NewGrouping(Expession Expr) *Grouping {
	return &Grouping{
		Expession: Expession,
	}
}

func (g *Grouping) Accept(v IASTVisitor) (any, error) { return v.VisitGroupingExpr(g) }
