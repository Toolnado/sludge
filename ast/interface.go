package ast

//go:generate go run ../tools/generator.go

type Expr interface {
	Accept(IASTVisitor) (any, error)
}

type Stmt interface {
	Expr
}

type Callable interface {
	Call(interpreter IASTVisitor, arguments []any) any
	Arity() int
}
