package interpreter

import (
	"fmt"

	"github.com/Toolnado/sludge/token"
)

type InterpreterError struct {
	pos     token.Position
	message string
}

func NewError(message string, pos token.Position) InterpreterError {
	return InterpreterError{
		pos:     pos,
		message: message,
	}
}

func (t InterpreterError) Error() string {
	filename := "<input>"
	if t.pos.Filename != "" {
		filename = t.pos.Filename
	}
	return fmt.Sprintf("%s\n%s:%d:%d", t.message, filename, t.pos.Line, t.pos.Column)
}
