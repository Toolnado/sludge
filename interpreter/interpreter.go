package interpreter

import (
	"fmt"

	"github.com/Toolnado/sludge/ast"
	"github.com/Toolnado/sludge/environment"
	"github.com/Toolnado/sludge/syslib/time"
	"github.com/Toolnado/sludge/token"
)

type Interpreter struct {
	globals     *environment.Environment
	environment *environment.Environment
}

func New() *Interpreter {
	base := environment.New(nil)
	i := &Interpreter{
		globals:     base,
		environment: base,
	}

	i.globals.Define("clock", time.New())
	return i
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) (any, error) {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) (any, error) {
	return i.evaluate(expr.Expession)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) (any, error) {
	return i.evaluate(stmt.Expession)
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) (any, error) {
	for {
		value, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !i.isTruthy(value) {
			break
		}
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expession)
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) (any, error) {
	return i.excecuteBlock(stmt.Statements, environment.New(i.environment))
}

func (i *Interpreter) excecuteBlock(stmts []ast.Stmt, env *environment.Environment) (any, error) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = env
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.AssignExpr) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, NewError(err.Error(), expr.Name.Position)
	}
	_, err = i.environment.Assign(expr.Name, value)
	if err != nil {
		return nil, NewError(err.Error(), expr.Name.Position)
	}
	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) (any, error) {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) (any, error) {
	var value any
	if stmt.Initializer != nil {
		value, _ = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) (any, error) {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		_, err = i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		_, err = i.execute(stmt.ElseBranch)
	}
	return nil, err
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, NewError(err.Error(), expr.Operator.Position)
	}

	switch expr.Operator.Type {
	case token.MINUS:
		value, err := i.negate(right)
		if err != nil {
			return nil, NewError(err.Error(), expr.Operator.Position)
		}
		return value, nil
	case token.BANG:
		return i.logicalNot(right), nil
	default:
		return nil, NewError("unsupported unary operator", expr.Operator.Position)
	}
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, NewError(err.Error(), expr.Operator.Position)
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, NewError(err.Error(), expr.Operator.Position)
	}

	switch expr.Operator.Type {
	case token.PLUS:
		value, err := i.add(left, right, expr.Operator)
		if err != nil {
			return nil, NewError(err.Error(), expr.Operator.Position)
		}
		return value, nil
	case token.MINUS, token.STAR, token.SLASH, token.PERCENT:
		value, err := i.performNumericOp(expr.Operator, left, right)
		if err != nil {
			return nil, NewError(err.Error(), expr.Operator.Position)
		}
		return value, nil

	case token.EQUAL_EQUAL, token.BANG_EQUAL,
		token.GREATER, token.GREATER_EQUAL,
		token.LESS, token.LESS_EQUAL:
		value, err := i.compareValues(expr.Operator, left, right)
		if err != nil {
			return nil, NewError(err.Error(), expr.Operator.Position)
		}
		return value, nil

	default:
		return nil, NewError("unsupported binary operator", expr.Operator.Position)
	}
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, NewError(err.Error(), expr.Paren.Position)
	}

	args := make([]any, len(expr.Arguments))
	for indx, item := range expr.Arguments {
		result, err := i.evaluate(item)
		if err != nil {
			return nil, NewError(err.Error(), expr.Paren.Position)
		}
		args[indx] = result
	}

	function, ok := callee.(ast.Callable)
	if !ok {
		return nil, NewError("can only call functions and classes", expr.Paren.Position)
	}

	if len(args) != function.Arity() {
		return nil, NewError(
			fmt.Sprintf("expected %d arguments, but got %d", function.Arity(), len(args)),
			expr.Paren.Position,
		)
	}

	return function.Call(i, args), nil
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}
