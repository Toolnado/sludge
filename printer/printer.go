package printer

import (
	"fmt"
	"strings"

	"github.com/Toolnado/sludge/ast"
)

type AstPrinter struct{}

func New() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr ast.Expr) {
	fmt.Println(expr.Accept(a).(string))
}

func (a *AstPrinter) VisitBinaryExpr(expr *ast.Binary) any {
	return a.parenthesize(expr.Operator.Literal, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) any {
	return a.parenthesize("group", expr.Expession)
}

func (a *AstPrinter) VisitLiteralExpr(expr *ast.Literal) any {
	if expr == nil {
		return "nil"
	}
	return expr.Value
}

func (a *AstPrinter) VisitUnaryExpr(expr *ast.Unary) any {
	return a.parenthesize(expr.Operator.Literal, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	b := &strings.Builder{}
	b.WriteString("(")
	b.WriteString(name)
	for _, e := range exprs {
		b.WriteString(" ")
		b.WriteString(e.Accept(a).(string))
	}
	b.WriteString(")")
	return b.String()
}
