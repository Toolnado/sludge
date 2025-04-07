package environment

import (
	"fmt"

	"github.com/Toolnado/sludge/token"
)

type Environment struct {
	values map[string]any
}

func New() Environment {
	return Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name token.Token) (any, error) {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}

func (e *Environment) Assign(name token.Token, value any) (any, error) {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil, nil
	}
	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}
