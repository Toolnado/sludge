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

func (a *AstPrinter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	return a.parenthesize(expr.Operator.Literal.(string), expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return a.parenthesize("group", expr.Expession)
}

func (a *AstPrinter) VisitExprStmt(stmt *ast.ExprStmt) (any, error) {
	a.parenthesize("exprStmt", stmt.Expession)
	return nil, nil
}

func (a *AstPrinter) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	a.parenthesize("printStmt", stmt.Expession)
	return nil, nil
}

func (a *AstPrinter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	if expr == nil {
		return "nil", nil
	}
	return expr.Value, nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	return a.parenthesize(expr.Operator.Literal.(string), expr.Right)
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

func (a *AstPrinter) VisitVariableExpr(expr *ast.VariableExpr) (any, error) { return nil, nil }
func (a *AstPrinter) VisitVarStmt(expr *ast.VarStmt) (any, error)           { return nil, nil }
func (a *AstPrinter) VisitAssignExpr(expr *ast.AssignExpr) (any, error)     { return nil, nil }
func (a *AstPrinter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error)   { return nil, nil }
func (a *AstPrinter) VisitBlockStmt(stmt *ast.BlockStmt) (any, error)       { return nil, nil }
func (a *AstPrinter) VisitIfStmt(stmt *ast.IfStmt) (any, error)             { return nil, nil }
func (a *AstPrinter) VisitWhileStmt(stmt *ast.WhileStmt) (any, error)       { return nil, nil }
func (i *AstPrinter) VisitCallExpr(expr *ast.CallExpr) (any, error)         { return nil, nil }
