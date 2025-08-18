package interpreter

import (
	"errors"
	"fmt"
)

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, val interface{}) {
	e.values[name] = val
}

func (e *Environment) Assign(name string, val interface{}) error {
	if _, ok := e.values[name]; ok {
		e.values[name] = val
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, val)
	}

	return errors.New("Undefined variable '" + name + "'.")
}

func (e *Environment) Get(name string) (interface{}, error) {
	if val, ok := e.values[name]; ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("undefined variable '%s'", name)
}

func (e *Environment) GetAt(depth int, name string) (interface{}, error) {
	env := e
	for d := 0; d < depth; d++ {
		env = env.enclosing
		if env == nil {
			return nil, fmt.Errorf("no env at depth %d", depth)
		}
	}

	return env.Get(name)
}
