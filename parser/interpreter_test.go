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
		// Variable declaration tests
		{"basic var declaration and use", "var a = 1; var b = 2; print a + b;", []string{"3"}, false},
		{"var with string", "var name = \"John\"; print name;", []string{"John"}, false},
		{"var with boolean", "var flag = true; print flag;", []string{"true"}, false},
		{"var with nil", "var empty = nil; print empty;", []string{"<nil>"}, false},
		{"var with expression", "var result = 2 + 3 * 4; print result;", []string{"14"}, false},
		{"var with comparison", "var isGreater = 5 > 3; print isGreater;", []string{"true"}, false},
		{"var with string concat", "var greeting = \"Hello\" + \" World\"; print greeting;", []string{"Hello World"}, false},
		{"var with ternary", "var value = true ? 42 : 0; print value;", []string{"42"}, false},
		{"var with unary", "var negative = -10; print negative;", []string{"-10"}, false},
		{"var with logical not", "var opposite = !false; print opposite;", []string{"true"}, false},

		// Multiple variable operations
		{"multiple vars same line", "var x = 5; var y = 10; var z = x * y; print z;", []string{"50"}, false},
		{"var reuse in expression", "var a = 3; var b = a + a; print b;", []string{"6"}, false},
		{"vars in complex expression", "var x = 2; var y = 3; print (x + y) * (x - y);", []string{"-5"}, false},
		{"vars with grouping", "var a = 2; var b = 3; var c = 4; print a + (b * c);", []string{"14"}, false},
		{"vars in ternary", "var a = 5; var b = 3; print a > b ? a : b;", []string{"5"}, false},
		{"vars in comma expression", "var a = 1; var b = 2; var c = 3; print a, b, c;", []string{"3"}, false},

		// Variable reassignment through redeclaration
		{"var redeclaration", "var x = 1; print x; var x = 2; print x;", []string{"1", "2"}, false},
		{"var redeclaration with different type", "var v = 42; print v; var v = \"hello\"; print v;", []string{"42", "hello"}, false},
		{"assignment error", "x = 2;", []string{}, true},
		{"print assignemnt statement", "var x = 1; print x = 2;", []string{"2"}, false},
		{"var redeclaration with assignment", "var x = 1; print x; x = 2; print x;", []string{"1", "2"}, false},
		{"var redeclaration with assignment and different type", "var v = 42; print v; v = \"hello\"; print v;", []string{"42", "hello"}, false},

		// Variables with all operators
		{"vars with arithmetic", "var a = 10; var b = 3; print a + b; print a - b; print a * b; print a / b;", []string{"13", "7", "30", "3.3333333333333335"}, false},
		{"vars with comparison", "var x = 5; var y = 3; print x > y; print x >= y; print x < y; print x <= y;", []string{"true", "true", "false", "false"}, false},
		{"vars with equality", "var a = 5; var b = 5; var c = 3; print a == b; print a != c;", []string{"true", "true"}, false},

		// Edge cases and complex scenarios
		{"var with decimal", "var pi = 3.14159; print pi * 2;", []string{"6.28318"}, false},
		{"var with large number", "var big = 999999; print big + 1;", []string{"1e+06"}, false},
		{"var with zero", "var zero = 0; print zero == 0;", []string{"true"}, false},
		{"var with empty string", "var empty = \"\"; print empty == \"\";", []string{"true"}, false},

		// Variables in nested expressions
		{"deeply nested vars", "var a = 1; var b = 2; var c = 3; print ((a + b) * c) > (a * (b + c));", []string{"true"}, false},
		{"vars with mixed types", "var num = 42; var str = \"answer\"; var bool = true; print bool ? num : str;", []string{"42"}, false},
		// Print statement tests
		{"print string literal", "print \"hello world\";", []string{"hello world"}, false},
		{"print arithmetic", "print 5 + 5;", []string{"10"}, false},
		{"print complex arithmetic", "print 5 + 5 * 10;", []string{"55"}, false},
		{"multiple print statements", "print 5 + 5 * 10; print \"hello\";", []string{"55", "hello"}, false},
		{"print boolean true", "print true;", []string{"true"}, false},
		{"print boolean false", "print false;", []string{"false"}, false},
		{"print nil", "print nil;", []string{"<nil>"}, false},
		{"print string concatenation", "print \"hello\" + \" \" + \"world\";", []string{"hello world"}, false},
		{"print comparison result", "print 5 > 3;", []string{"true"}, false},
		{"print equality result", "print 5 == 5;", []string{"true"}, false},
		{"print unary negation", "print -42;", []string{"-42"}, false},
		{"print logical not", "print !true;", []string{"false"}, false},
		{"print ternary result", "print true ? \"yes\" : \"no\";", []string{"yes"}, false},
		{"print comma expression", "print 1, 2, 3;", []string{"3"}, false},
		{"print grouped expression", "print (2 + 3) * 4;", []string{"20"}, false},

		// Block statement
		{"simple block scoping", "var x = 0; var y = 1; var z = 2; { var x = 1; var y = 2; { var x = 2; print x; print y; print z; } } print x;", []string{"2", "2", "2", "0"}, false},

		// Expression statement tests (no output expected)
		{"expression statement arithmetic", "5 + 5;", []string{}, false},
		{"expression statement string", "\"hello\";", []string{}, false},
		{"expression statement boolean", "true;", []string{}, false},
		{"expression statement comparison", "5 > 3;", []string{}, false},
		{"expression statement function call", "!false;", []string{}, false},

		// Mixed statement tests
		{"print and expression mixed", "print \"start\"; 5 + 5; print \"end\";", []string{"start", "end"}, false},
		{"multiple expression statements", "1 + 1; 2 + 2; 3 + 3;", []string{}, false},
		{"complex mixed statements", "print \"result:\"; print 2 * 3; 4 + 4; print \"done\";", []string{"result:", "6", "done"}, false},

		// Edge cases
		{"empty program", "", []string{}, false},
		{"print empty string", "print \"\";", []string{""}, false},
		{"print zero", "print 0;", []string{"0"}, false},
		{"print negative zero", "print -0;", []string{"-0"}, false},
		{"print large number", "print 999999;", []string{"999999"}, false},
		{"print decimal", "print 3.14159;", []string{"3.14159"}, false},
		{"print simple string", "print \"simple string\";", []string{"simple string"}, false},

		// Complex expressions in print statements
		{"print nested ternary", "print true ? (false ? 1 : 2) : 3;", []string{"2"}, false},
		{"print chained comparisons", "print 1 < 2 == 2 > 1;", []string{"true"}, false},
		{"print mixed operators", "print 2 + 3 * 4 == 14;", []string{"true"}, false},
		{"print complex grouping", "print ((1 + 2) * (3 + 4)) / 7;", []string{"3"}, false},
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
