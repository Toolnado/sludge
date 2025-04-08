package environment

import (
	"fmt"

	"github.com/Toolnado/sludge/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func New(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
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
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}

func (e *Environment) Assign(name token.Token, value any) (any, error) {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}
