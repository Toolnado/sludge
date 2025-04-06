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
	val, _ := expr.Accept(a)
	fmt.Println(val.(string))
}

func (a *AstPrinter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	return a.parenthesize(expr.Operator.Literal, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return a.parenthesize("group", expr.Expession)
}

func (a *AstPrinter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	if expr == nil {
		return "nil", nil
	}
	return expr.Value, nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	return a.parenthesize(expr.Operator.Literal, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) (string, error) {
	b := &strings.Builder{}
	b.WriteString("(")
	b.WriteString(name)
	for _, e := range exprs {
		b.WriteString(" ")
		val, err := e.Accept(a)
		if err != nil {
			return "", err
		}
		b.WriteString(val.(string))
	}
	b.WriteString(")")
	return b.String(), nil
}
