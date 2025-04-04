package ast

type IASTVisitor interface {
	Visit(Expr)
}
