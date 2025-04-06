package parser

import "github.com/Toolnado/sludge/token"

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
	return ""
}
