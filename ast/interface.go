package ast

//go:generate go run ../tools/generator.go

type Expr interface {
	Accept(IASTVisitor) (any, error)
}

type Stmt interface {
	Expr
}
