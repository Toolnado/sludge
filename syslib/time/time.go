package time

import (
	"time"

	"github.com/Toolnado/sludge/ast"
)

type Time struct {
	arity int
}

func New() Time {
	return Time{}
}

var start = time.Now()

func (t Time) Call(interpreter ast.IASTVisitor, arguments []any) any {
	return time.Since(start)
}
func (t Time) Arity() int { return t.arity }
