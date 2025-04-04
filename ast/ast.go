package ast

import "github.com/Toolnado/sludge/token"

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinary(Left Expr, Operator token.Token, Right Expr) *Binary {
	return &Binary{
		Left: Left,
		Operator: Operator,
		Right: Right,
	}
}

func (b *Binary) Accept(v IASTVisitor) {v.Visit(b)}

type Unary struct {
	Operator token.Token
	Right Expr
}

func NewUnary(Operator token.Token, Right Expr) *Unary {
	return &Unary{
		Operator: Operator,
		Right: Right,
	}
}

func (u *Unary) Accept(v IASTVisitor) {v.Visit(u)}

type Literal struct {
	Value string
}

func NewLiteral(Value string) *Literal {
	return &Literal{
		Value: Value,
	}
}

func (l *Literal) Accept(v IASTVisitor) {v.Visit(l)}

type Grouping struct {
	Expession Expr
}

func NewGrouping(Expession Expr) *Grouping {
	return &Grouping{
		Expession: Expession,
	}
}

func (g *Grouping) Accept(v IASTVisitor) {v.Visit(g)}

