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
	VisitLogicalExpr(expr *LogicalExpr) (any, error)
}

type IStmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt) (any, error)
	VisitExprStmt(stmt *ExprStmt) (any, error)
	VisitVarStmt(stmt *VarStmt) (any, error)
	VisitBlockStmt(stmt *BlockStmt) (any, error)
	VisitIfStmt(stmt *IfStmt) (any, error)
}
