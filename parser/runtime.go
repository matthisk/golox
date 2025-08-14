package parser

import (
	"fmt"
	"time"
)

type LoxCallable interface {
	Arity() int
	Call(i *Interpreter, args []interface{}) (interface{}, error)
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (l *LoxInstance) ToString() string {
	return fmt.Sprintf("%s instance", l.class.name)
}

func (l *LoxInstance) Get(property string) (interface{}, error) {
	if v, ok := l.fields[property]; ok {
		return v, nil
	}

	if m := l.class.FindMethod(property); m != nil {
		return m, nil
	}

	return nil, fmt.Errorf("undefined property %s", property)
}

func (l *LoxInstance) Set(property string, val interface{}) {
	l.fields[property] = val
}

type LoxClass struct {
	name    string
	methods map[string]LoxCallable
}

func (l *LoxClass) Arity() int {
	return 0
}

func (l *LoxClass) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	instance := NewLoxInstance(l)
	return instance, nil
}

func (l *LoxClass) ToString() string {
	return l.name
}

func (l *LoxClass) FindMethod(property string) LoxCallable {
	return l.methods[property]
}

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.params)
}

func (l *LoxFunction) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	env := NewEnvironment(l.closure)

	for j, param := range l.declaration.params {
		env.Define(param.Lexeme.(string), args[j])
	}

	return i.executeBlock(l.declaration.body, env)
}

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()), nil
}
