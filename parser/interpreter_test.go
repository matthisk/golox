package parser

import (
	"bytes"
	"testing"

	"github.com/matthisk/lox/lexer"
)

func TestInterpreter_SimpleArithmetic(t *testing.T) {
	// Represents the expression: 1 + 2
	left := &Literal{
		token: lexer.Token{Type: lexer.NUMBER, Lexeme: float64(1)},
	}
	right := &Literal{
		token: lexer.Token{Type: lexer.NUMBER, Lexeme: float64(2)},
	}
	expr := &Binary{
		left:  left,
		token: lexer.PLUS,
		right: right,
	}

	interpreter := Interpreter{}
	result, err := interpreter.evaluate(expr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := float64(3)
	if res, ok := result.(float64); !ok || res != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestEndToEnd_LexerParserInterpreter(t *testing.T) {
	source := "1 + 2"
	lx := lexer.New(bytes.NewBufferString(source))
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		t.Fatalf("Lexer error: %v", lexResult.Err)
	}

	parser := New(lexResult.Tokens)
	expr, err := parser.expression()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	interpreter := Interpreter{}
	result, err := interpreter.evaluate(expr)
	if err != nil {
		t.Fatalf("Interpreter error: %v", err)
	}

	expected := float64(3)
	if res, ok := result.(float64); !ok || res != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
