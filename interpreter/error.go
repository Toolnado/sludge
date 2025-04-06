package interpreter

import "github.com/Toolnado/sludge/ast"

type InterpreterError struct {
	expr    []ast.Expr
	message string
}

func NewError(message string, exprs ...ast.Expr) InterpreterError {
	return InterpreterError{
		expr:    exprs,
		message: message,
	}
}

func (t InterpreterError) Error() string {
	return ""
}
