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

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(Statements []Stmt) *BlockStmt {
	return &BlockStmt{
		Statements: Statements,
	}
}

func (b *BlockStmt) Accept(v IASTVisitor) (any, error) { return v.VisitBlockStmt(b)}

type IfStmt struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(Condition Expr, ThenBranch Stmt, ElseBranch Stmt) *IfStmt {
	return &IfStmt{
		Condition: Condition,
		ThenBranch: ThenBranch,
		ElseBranch: ElseBranch,
	}
}

func (i *IfStmt) Accept(v IASTVisitor) (any, error) { return v.VisitIfStmt(i)}

type WhileStmt struct {
	Condition Expr
	Body Stmt
}

func NewWhileStmt(Condition Expr, Body Stmt) *WhileStmt {
	return &WhileStmt{
		Condition: Condition,
		Body: Body,
	}
}

func (w *WhileStmt) Accept(v IASTVisitor) (any, error) { return v.VisitWhileStmt(w)}

type FunctionStmt struct {
	Name token.Token
	Params []token.Token
	Body []Stmt
}

func NewFunctionStmt(Name token.Token, Params []token.Token, Body []Stmt) *FunctionStmt {
	return &FunctionStmt{
		Name: Name,
		Params: Params,
		Body: Body,
	}
}

func (f *FunctionStmt) Accept(v IASTVisitor) (any, error) { return v.VisitFunctionStmt(f)}

