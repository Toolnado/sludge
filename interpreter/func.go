package interpreter

import (
	"github.com/Toolnado/sludge/ast"
	"github.com/Toolnado/sludge/environment"
)

type Callable interface {
	Call(interpreter *Interpreter, arguments []any) any
	Arity() int
}

type Function struct {
	declaration ast.FunctionStmt
}

func NewFunction(declaration ast.FunctionStmt) Function {
	return Function{
		declaration: declaration,
	}
}

func (f Function) Call(interpreter *Interpreter, arguments []any) any {
	environment := environment.New(interpreter.globals)
	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}
	interpreter.excecuteBlock(f.declaration.Body, environment)
	return nil
}

func (f Function) Arity() int {
	return len(f.declaration.Params)
}
