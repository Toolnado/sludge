package interpreter

import (
	"errors"
	"fmt"

	"github.com/Toolnado/sludge/token"
)

func (i *Interpreter) logicalNot(value any) bool {
	return !i.isTruthy(value)
}

func (i *Interpreter) isTruthy(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func (i *Interpreter) negate(value any) (any, error) {
	switch v := value.(type) {
	case int64:
		return -v, nil
	case float64:
		return -v, nil
	default:
		return nil, errors.New("unary '-' expects number")
	}
}

func (i *Interpreter) performNumericOp(op token.Token, left, right any) (any, error) {
	// Пробуем int → int
	if li, lok := left.(int64); lok {
		if ri, rok := right.(int64); rok {
			switch op.Type {
			case token.PLUS:
				return li + ri, nil
			case token.MINUS:
				return li - ri, nil
			case token.STAR:
				return li * ri, nil
			case token.SLASH:
				if ri == 0 {
					return nil, errors.New("division by zero")
				}
				return float64(li) / float64(ri), nil // Деление всё равно float
			case token.PERCENT:
				if ri == 0 {
					return nil, errors.New("modulo by zero")
				}
				return li % ri, nil
			}
		}
	}

	// Fallback to float64
	lf, err := i.toFloat64(left)
	if err != nil {
		return nil, errors.New("left operand is not a number")
	}

	rf, err := i.toFloat64(right)
	if err != nil {
		return nil, errors.New("right operand is not a number")
	}

	switch op.Type {
	case token.PLUS:
		return lf + rf, nil
	case token.MINUS:
		return lf - rf, nil
	case token.STAR:
		return lf * rf, nil
	case token.SLASH:
		if rf == 0 {
			return nil, errors.New("division by zero")
		}
		return lf / rf, nil
	default:
		return nil, errors.New("unsupported numeric op on float")
	}
}

func (i *Interpreter) compareValues(op token.Token, left, right any) (any, error) {
	switch op.Type {
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	case token.BANG_EQUAL:
		return !i.isEqual(left, right), nil
	}

	lf, err := i.toFloat64(left)
	if err != nil {
		return nil, errors.New("left not number for comparison")
	}
	rf, err := i.toFloat64(right)
	if err != nil {
		return nil, errors.New("right not number for comparison")
	}

	switch op.Type {
	case token.LESS:
		return lf < rf, nil
	case token.LESS_EQUAL:
		return lf <= rf, nil
	case token.GREATER:
		return lf > rf, nil
	case token.GREATER_EQUAL:
		return lf >= rf, nil
	default:
		return nil, errors.New("unknown comparison op")
	}
}

func (i *Interpreter) isEqual(a, b any) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		}
	case string:
		if y, ok := b.(string); ok {
			return x == y
		}
	case bool:
		if y, ok := b.(bool); ok {
			return x == y
		}
	}
	return false
}

func (i *Interpreter) toFloat64(v any) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, fmt.Errorf("value %v (type %T) is not numeric", v, v)
	}
}

func (i *Interpreter) add(left, right any, op token.Token) (any, error) {
	switch l := left.(type) {
	case float64, int64:
		return i.performNumericOp(op, left, right)
	case string:
		r, ok := right.(string)
		if !ok {
			return nil, errors.New("cannot concatenate string")
		}
		return l + r, nil
	default:
		return nil, errors.New("unsupported operand types for '+'")
	}
}
