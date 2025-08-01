package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matthisk/lox/lexer"
)

func TestInterpreter_WithStmts_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		logs    []string
		wantErr bool
	}{
		// Basic arithmetic expressions
		{"print statement", "print \"hello world\";", []string{"hello world"}, false},
		{"print statement", "print 5 + 5;", []string{"10"}, false},
		{"print statement", "print 5 + 5 * 10;", []string{"55"}, false},
		{"print statement", "print 5 + 5 * 10; print \"hello\";", []string{"55", "hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := runLox(tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !SliceEqual(tt.logs, logs) {
				t.Errorf("Interpreter printed unexpected results %v expected %v", logs, tt.logs)
			}
		})
	}
}

func TestInterpreter_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected interface{}
		wantErr  bool
	}{
		// Basic arithmetic expressions
		{"simple addition", "1 + 2", float64(3), false},
		{"simple multiplication", "3 * 4", float64(12), false},
		{"string concatenation", `"hello" + "world"`, "helloworld", false},

		// Arithmetic operations
		{"subtraction", "5 - 3", float64(2), false},
		{"division", "8 / 2", float64(4), false},

		// Unary operators
		{"unary minus", "-5", float64(-5), false},
		{"unary not true", "!true", false, false},
		{"unary not false", "!false", true, false},
		{"unary not truthy", "!123", false, false},
		{"unary not nil", "!nil", true, false},

		// Comparison operators (these need number operands)
		{"greater than true", "5 > 3", true, false},
		{"greater than false", "3 > 5", false, false},
		{"greater equal true", "5 >= 5", true, false},
		{"greater equal false", "3 >= 5", false, false},
		{"less than true", "3 < 5", true, false},
		{"less than false", "5 < 3", false, false},
		{"less equal true", "3 <= 3", true, false},
		{"less equal false", "5 <= 3", false, false},

		// Equality operators
		{"equal numbers", "5 == 5", true, false},
		{"not equal numbers", "5 != 3", true, false},
		{"equal strings", `"hello" == "hello"`, true, false},
		{"not equal strings", `"hello" != "world"`, true, false},
		{"equal booleans", "true == true", true, false},
		{"not equal booleans", "true != false", true, false},
		{"nil equality", "nil == nil", true, false},
		{"nil inequality", "nil != nil", false, false},
		{"different types not equal", "5 != true", true, false},

		// Grouping expressions
		{"grouped addition", "(1 + 2)", float64(3), false},
		{"grouped multiplication", "(3 * 4)", float64(12), false},
		{"nested grouping", "((5))", float64(5), false},

		// Mixed operator precedence
		{"addition and multiplication", "2 + 3 * 4", float64(14), false},
		{"multiplication and addition with grouping", "(2 + 3) * 4", float64(20), false},
		{"unary and binary", "-2 + 3", float64(1), false},
		{"comparison and equality", "3 > 2 == true", true, false},

		// Ternary operator
		{"ternary true condition", "true ? 1 : 2", float64(1), false},
		{"ternary false condition", "false ? 1 : 2", float64(2), false},
		{"ternary with nil", "nil ? 1 : 2", float64(2), false},
		{"ternary with truthy", "5 ? 1 : 2", float64(1), false},
		{"nested ternary", "true ? (false ? 1 : 2) : 3", float64(2), false},

		// Comma operator
		{"simple comma", "1, 2", float64(2), false},
		{"comma with expressions", "1 + 2, 3 * 4", float64(12), false},
		{"multiple comma", "1, 2, 3", float64(3), false},

		// Complex nested expressions
		{"complex arithmetic", "2 + 3 * 4 - 1", float64(13), false},
		{"complex with grouping", "(2 + 3) * (4 - 1)", float64(15), false},
		{"complex with ternary", "2 > 1 ? 3 + 4 : 5 * 6", float64(7), false},
		{"complex with comma", "1 + 2, 3 * 4, 5 > 3", true, false},
		{"deeply nested", "((2 + 3) * 4) > 15 ? true : false", true, false},

		// Edge cases with different types
		{"string and number not equal", `"5" != 5`, true, false},
		{"boolean and number not equal", "true != 1", true, false},
		{"nil and false not equal", "nil != false", true, false},

		// More complex truthiness tests
		{"zero is truthy", "!0", false, false},
		{"empty string is truthy", `!""`, false, false},
		{"non-zero number is truthy", "!42", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runLoxExpression(tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expression %q: expected %v (%T), got %v (%T)",
					tt.expr, tt.expected, tt.expected, result, result)
			}
		})
	}
}

func runLoxExpression(source string) (interface{}, error) {
	lx := lexer.New(bytes.NewBufferString(source))
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		return nil, lexResult.Err
	}

	parser := New(lexResult.Tokens)
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	interpreter := Interpreter{}
	return interpreter.evaluate(expr)
}

func runLox(source string) ([]string, error) {
	lx := lexer.New(bytes.NewBufferString(source))
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		return nil, lexResult.Err
	}

	parser := New(lexResult.Tokens)
	stmts, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	printer := &SpyPrinter{}
	interpreter := NewInterpreterWithPrinter(printer)
	err = interpreter.Run(stmts)
	if err != nil {
		return nil, err
	}

	return printer.log, nil
}

type SpyPrinter struct {
	log []string
}

func (s *SpyPrinter) Print(value interface{}) {
	s.log = append(s.log, fmt.Sprintf("%v", value))
}

// SliceEqual compares two slices for equality.
// Returns true if both slices have the same length and all elements are equal.
func SliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
