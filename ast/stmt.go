package ast

import "github.com/Toolnado/sludge/token"

type PrintStmt struct {
	Expession Expr
}

func NewPrintStmt(Expession Expr) *PrintStmt {
	return &PrintStmt{
		Expession: Expession,
	}
}

func (p *PrintStmt) Accept(v IASTVisitor) (any, error) { return v.VisitPrintStmt(p)}

type ExprStmt struct {
	Expession Expr
}

func NewExprStmt(Expession Expr) *ExprStmt {
	return &ExprStmt{
		Expession: Expession,
	}
}

func (e *ExprStmt) Accept(v IASTVisitor) (any, error) { return v.VisitExprStmt(e)}

type VarStmt struct {
	Name token.Token
	Initializer Expr
}

func NewVarStmt(Name token.Token, Initializer Expr) *VarStmt {
	return &VarStmt{
		Name: Name,
		Initializer: Initializer,
	}
}

func (va *VarStmt) Accept(v IASTVisitor) (any, error) { return v.VisitVarStmt(va)}

