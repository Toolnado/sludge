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
