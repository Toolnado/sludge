package ast

type IASTVisitor interface {
	IExprVisitor
	IStmtVisitor
}
type IExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitVariableExpr(expr *VariableExpr) (any, error)
	VisitAssignExpr(expr *AssignExpr) (any, error)
}

type IStmtVisitor interface {
	VisitPrintStmt(expr *PrintStmt) (any, error)
	VisitExprStmt(expr *ExprStmt) (any, error)
	VisitVarStmt(expr *VarStmt) (any, error)
}
