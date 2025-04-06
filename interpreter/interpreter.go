package interpreter

import (
	"github.com/Toolnado/sludge/ast"
	"github.com/Toolnado/sludge/token"
)

type Interpreter struct {
	// hadRuntimeError bool
}

func New() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr ast.Expr) (any, error) {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
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

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
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

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}
