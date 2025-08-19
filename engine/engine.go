package engine

import (
	"fmt"
	"io"

	"github.com/matthisk/lox/interpreter"
	"github.com/matthisk/lox/lexer"
	"github.com/matthisk/lox/parser"
)

func Run(in io.Reader, printer interpreter.Printer) error {
	lexResult, err := lex(in)
	if err != nil {
		return err
	}

	ps := parser.New(lexResult.Tokens, lexResult.Source)
	stmts, err := ps.Parse()
	if err != nil {
		fmt.Println(ps.ReportErrors())
		return err
	}

	var it *interpreter.Interpreter
	if printer != nil {
		it = interpreter.NewInterpreterWithPrinter(printer)
	} else {
		it = interpreter.NewInterpreter()
	}

	resolver := interpreter.NewResolver(it)
	err = resolver.Resolve(stmts)
	if err != nil {
		return err
	}

	return it.Run(stmts)
}

func EvaluateExpr(in io.Reader) (interface{}, error) {
	lexResult, err := lex(in)
	if err != nil {
		return nil, err
	}

	ps := parser.New(lexResult.Tokens, lexResult.Source)
	expr, err := ps.Expression()
	if err != nil {
		fmt.Println(ps.ReportErrors())
		return nil, err
	}

	return interpreter.NewInterpreter().EvaluateExpression(expr)
}

func lex(in io.Reader) (*lexer.Lex, error) {
	lx := lexer.New(in)
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		return nil, lexResult.Err
	}

	for i := range lexResult.Tokens {
		if lexResult.Tokens[i].Type == lexer.ILLEGAL {
			return nil, fmt.Errorf("encountered illegal token at index %d", i)
		}
	}
	return &lexResult, nil
}
