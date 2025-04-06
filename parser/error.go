package parser

import (
	"fmt"

	"github.com/Toolnado/sludge/token"
)

type TokenError struct {
	token   token.Token
	message string
}

func NewError(t token.Token, message string) TokenError {
	return TokenError{
		token:   t,
		message: message,
	}
}

func (t TokenError) Error() string {
	filename := "<input>"
	if t.token.Position.Filename != "" {
		filename = t.token.Position.Filename
	}
	return fmt.Sprintf("%s\n%s:%d:%d", t.message, filename, t.token.Position.Line, t.token.Position.Column)
}
