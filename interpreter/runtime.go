package interpreter

import (
	"fmt"
	"time"

	"github.com/matthisk/lox/parser"
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
		// Only LoxFunction's can be set as class methods.
		// Panic in case we find an illegal LoxCallable as a Class' method.
		return m.(*LoxFunction).Bind(l), nil
	}

	if l.class.super != nil {
		if m := l.class.super.FindMethod(property); m != nil {
			return m.(*LoxFunction).Bind(l), nil
		}
	}

	return nil, fmt.Errorf("undefined property %s", property)
}

func (l *LoxInstance) Set(property string, val interface{}) {
	l.fields[property] = val
}

type LoxClass struct {
	name    string
	super   *LoxClass
	methods map[string]LoxCallable
}

func (l *LoxClass) Arity() int {
	initializer := l.FindMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (l *LoxClass) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	instance := NewLoxInstance(l)
	initializer := l.FindMethod("init")
	if initializer != nil {
		_, err := initializer.(*LoxFunction).Bind(instance).Call(i, args)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (l *LoxClass) ToString() string {
	return l.name
}

func (l *LoxClass) FindMethod(property string) LoxCallable {
	return l.methods[property]
}

type LoxFunction struct {
	declaration   *parser.Function
	closure       *Environment
	isInitializer bool
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	env := NewEnvironment(l.closure)

	for j, param := range l.declaration.Params {
		env.Define(param.Lexeme.(string), args[j])
	}

	result, err := i.executeBlock(l.declaration.Body, env)
	if err != nil {
		return nil, err
	}

	if l.isInitializer {
		return l.closure.GetAt(0, "this")
	}

	return result, nil
}

func (l *LoxFunction) Bind(i *LoxInstance) *LoxFunction {
	env := NewEnvironment(l.closure)
	env.Define("this", i)
	return &LoxFunction{l.declaration, env, l.isInitializer}
}

type Clock struct{}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(i *Interpreter, args []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()), nil
}

type Printer interface {
	Print(value interface{})
}

type DefaultPrinter struct{}

func (p DefaultPrinter) Print(value interface{}) {
	fmt.Println(value)
}
