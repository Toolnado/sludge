package printer

import (
	"fmt"

	"github.com/Toolnado/sludge/ast"
)

type AstPrinter struct{}

func New() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr ast.Expr) {
	expr.Accept(a)
}

func (a *AstPrinter) Visit(expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.Binary:
		a.VisitBinaryExpr(x)
	case *ast.Grouping:
		a.VisitGroupingExpr(x)
	case *ast.Literal:
		a.VisitLiteralExpr(x)
	case *ast.Unary:
		a.VisitUnaryExpr(x)
	default:
		fmt.Println("unknown expression")
	}
}

func (a *AstPrinter) VisitBinaryExpr(expr *ast.Binary) {
	a.parenthesize(expr.Operator.Literal, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) {
	a.parenthesize("group", expr.Expession)
}

func (a *AstPrinter) VisitLiteralExpr(expr *ast.Literal) {
	if expr == nil {
		fmt.Print("nil")
		return
	}
	fmt.Print(expr.Value)
}

func (a *AstPrinter) VisitUnaryExpr(expr *ast.Unary) {
	a.parenthesize(expr.Operator.Literal, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) {
	fmt.Print("(" + name)
	for _, e := range exprs {
		fmt.Print(" ")
		e.Accept(a)
	}
	fmt.Print(")")
}
