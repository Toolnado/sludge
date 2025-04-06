package interpreter

import (
	"github.com/Toolnado/sludge/ast"
)

type Peaky struct{}

type Interpreter struct{}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) any {
	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) any {
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) any {
	return nil

}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) any {
	return nil
}
