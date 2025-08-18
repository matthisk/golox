package parser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/matthisk/lox/lexer"
)

func TestParserTD(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected string
		wantErr  bool
	}{
		// Literal tests
		{
			name: "number literal",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: 42, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "42",
			wantErr:  false,
		},
		{
			name: "string literal",
			tokens: []lexer.Token{
				{Type: lexer.STRING, Lexeme: "hello", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "\"hello\"",
			wantErr:  false,
		},
		{
			name: "true literal",
			tokens: []lexer.Token{
				{Type: lexer.TRUE, Lexeme: "true", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "true",
			wantErr:  false,
		},
		{
			name: "false literal",
			tokens: []lexer.Token{
				{Type: lexer.FALSE, Lexeme: "false", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "false",
			wantErr:  false,
		},
		{
			name: "nil literal",
			tokens: []lexer.Token{
				{Type: lexer.NIL, Lexeme: "nil", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "nil",
			wantErr:  false,
		},
		{
			name: "decimal number",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "3.14", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "3.14",
			wantErr:  false,
		},

		// Grouping tests
		{
			name: "simple grouping",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.RIGHT_PAREN, StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(group 5)",
			wantErr:  false,
		},
		{
			name: "nested grouping",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "10", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.RIGHT_PAREN, StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.RIGHT_PAREN, StartPos: lexer.Pos{Offset: 4, Line: 1, Column: 5}},
				{Type: lexer.EOF},
			},
			expected: "(group (group 10))",
			wantErr:  false,
		},

		// Unary tests
		{
			name: "unary minus",
			tokens: []lexer.Token{
				{Type: lexer.MINUS, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.EOF},
			},
			expected: "(MINUS 5)",
			wantErr:  false,
		},
		{
			name: "unary bang",
			tokens: []lexer.Token{
				{Type: lexer.BANG, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.TRUE, Lexeme: "true", StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.EOF},
			},
			expected: "(BANG true)",
			wantErr:  false,
		},
		{
			name: "double unary",
			tokens: []lexer.Token{
				{Type: lexer.MINUS, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.MINUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(MINUS (MINUS 3))",
			wantErr:  false,
		},
		{
			name: "unary with grouping",
			tokens: []lexer.Token{
				{Type: lexer.BANG, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.FALSE, Lexeme: "false", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.RIGHT_PAREN, StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.EOF},
			},
			expected: "(BANG (group false))",
			wantErr:  false,
		},

		// Binary arithmetic tests
		{
			name: "simple addition",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(PLUS 5 3)",
			wantErr:  false,
		},
		{
			name: "simple subtraction",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "10", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.MINUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "4", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(MINUS 10 4)",
			wantErr:  false,
		},
		{
			name: "simple multiplication",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "6", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.STAR, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "7", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(STAR 6 7)",
			wantErr:  false,
		},
		{
			name: "simple division",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "8", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.SLASH, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "2", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(SLASH 8 2)",
			wantErr:  false,
		},

		// Binary comparison tests
		{
			name: "equality",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EQUAL_EQUAL, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(EQUAL_EQUAL 5 5)",
			wantErr:  false,
		},
		{
			name: "inequality",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.BANG_EQUAL, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(BANG_EQUAL 5 3)",
			wantErr:  false,
		},
		{
			name: "less than",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.LESS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(LESS 3 5)",
			wantErr:  false,
		},
		{
			name: "less than or equal",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.LESS_EQUAL, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(LESS_EQUAL 3 5)",
			wantErr:  false,
		},
		{
			name: "greater than",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "7", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.GREATER, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(GREATER 7 3)",
			wantErr:  false,
		},
		{
			name: "greater than or equal",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "7", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.GREATER_EQUAL, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(GREATER_EQUAL 7 3)",
			wantErr:  false,
		},

		// Complex expressions (operator precedence)
		{
			name: "addition and multiplication",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "2", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.STAR, StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.NUMBER, Lexeme: "4", StartPos: lexer.Pos{Offset: 4, Line: 1, Column: 5}},
				{Type: lexer.EOF},
			},
			expected: "(PLUS 2 (STAR 3 4))",
			wantErr:  false,
		},
		{
			name: "multiplication and addition",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "2", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.STAR, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.NUMBER, Lexeme: "4", StartPos: lexer.Pos{Offset: 4, Line: 1, Column: 5}},
				{Type: lexer.EOF},
			},
			expected: "(PLUS (STAR 2 3) 4)",
			wantErr:  false,
		},
		{
			name: "comparison and arithmetic",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.GREATER, StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.NUMBER, Lexeme: "2", StartPos: lexer.Pos{Offset: 4, Line: 1, Column: 5}},
				{Type: lexer.STAR, StartPos: lexer.Pos{Offset: 5, Line: 1, Column: 6}},
				{Type: lexer.NUMBER, Lexeme: "4", StartPos: lexer.Pos{Offset: 6, Line: 1, Column: 7}},
				{Type: lexer.EOF},
			},
			expected: "(GREATER (PLUS 5 3) (STAR 2 4))",
			wantErr:  false,
		},

		// Complex expressions with grouping
		{
			name: "grouping affects precedence",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "2", StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.RIGHT_PAREN, StartPos: lexer.Pos{Offset: 4, Line: 1, Column: 5}},
				{Type: lexer.STAR, StartPos: lexer.Pos{Offset: 5, Line: 1, Column: 6}},
				{Type: lexer.NUMBER, Lexeme: "4", StartPos: lexer.Pos{Offset: 6, Line: 1, Column: 7}},
				{Type: lexer.EOF},
			},
			expected: "(STAR (group (PLUS 2 3)) 4)",
			wantErr:  false,
		},
		{
			name: "complex nested expression",
			tokens: []lexer.Token{ // - ( 5 + 3 ) * 2
				{Type: lexer.MINUS, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 3, Line: 1, Column: 4}},
				{Type: lexer.NUMBER, Lexeme: "3", StartPos: lexer.Pos{Offset: 4, Line: 1, Column: 5}},
				{Type: lexer.RIGHT_PAREN, StartPos: lexer.Pos{Offset: 5, Line: 1, Column: 6}},
				{Type: lexer.STAR, StartPos: lexer.Pos{Offset: 6, Line: 1, Column: 7}},
				{Type: lexer.NUMBER, Lexeme: "2", StartPos: lexer.Pos{Offset: 7, Line: 1, Column: 8}},
				{Type: lexer.EOF},
			},
			expected: "(STAR (MINUS (group (PLUS 5 3))) 2)",
			wantErr:  false,
		},

		// Mixed types
		{
			name: "string concatenation",
			tokens: []lexer.Token{
				{Type: lexer.STRING, Lexeme: "hello", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.STRING, Lexeme: "world", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(PLUS \"hello\" \"world\")",
			wantErr:  false,
		},
		{
			name: "boolean comparison",
			tokens: []lexer.Token{
				{Type: lexer.TRUE, Lexeme: "true", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EQUAL_EQUAL, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.FALSE, Lexeme: "false", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(EQUAL_EQUAL true false)",
			wantErr:  false,
		},
		{
			name: "nil comparison",
			tokens: []lexer.Token{
				{Type: lexer.NIL, Lexeme: "nil", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.BANG_EQUAL, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 2, Line: 1, Column: 3}},
				{Type: lexer.EOF},
			},
			expected: "(BANG_EQUAL nil 5)",
			wantErr:  false,
		},

		// Error cases
		{
			name: "missing closing parenthesis",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "missing operand",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.PLUS, StartPos: lexer.Pos{Offset: 1, Line: 1, Column: 2}},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "invalid token",
			tokens: []lexer.Token{
				{Type: lexer.SEMICOLON, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "empty expression",
			tokens: []lexer.Token{
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(tt.tokens)
			expr, err := parser.expression()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			result, err2 := Print(expr)
			if err2 != nil {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParserStatementsTD(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected string
		wantErr  bool
	}{
		// Print statement tests
		{
			name: "print number",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: 42},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print 42;",
			wantErr:  false,
		},
		{
			name: "print string",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.STRING, Lexeme: "Hello World"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print \"Hello World\";",
			wantErr:  false,
		},
		{
			name: "print boolean true",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print true;",
			wantErr:  false,
		},
		{
			name: "print boolean false",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.FALSE, Lexeme: "false"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print false;",
			wantErr:  false,
		},
		{
			name: "print nil",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NIL, Lexeme: "nil"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print nil;",
			wantErr:  false,
		},
		{
			name: "print arithmetic expression",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "5"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "3"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print (PLUS 5 3);",
			wantErr:  false,
		},
		{
			name: "print comparison expression",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "10"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "5"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print (GREATER 10 5);",
			wantErr:  false,
		},
		{
			name: "print grouped expression",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "3"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.STAR},
				{Type: lexer.NUMBER, Lexeme: "4"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print (STAR (group (PLUS 2 3)) 4);",
			wantErr:  false,
		},

		// Expression statement tests
		{
			name: "expression statement with number",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "42"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "42;",
			wantErr:  false,
		},
		{
			name: "expression statement with string",
			tokens: []lexer.Token{
				{Type: lexer.STRING, Lexeme: "test"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "\"test\";",
			wantErr:  false,
		},
		{
			name: "expression statement with arithmetic",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "7"},
				{Type: lexer.STAR},
				{Type: lexer.NUMBER, Lexeme: "6"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "(STAR 7 6);",
			wantErr:  false,
		},
		{
			name: "expression statement with unary",
			tokens: []lexer.Token{
				{Type: lexer.MINUS},
				{Type: lexer.NUMBER, Lexeme: "15"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "(MINUS 15);",
			wantErr:  false,
		},
		{
			name: "expression statement with grouping",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "(group true);",
			wantErr:  false,
		},
		{
			name: "block statement with grouping",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.PRINT},
				{Type: lexer.FALSE, Lexeme: "false"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "{\nprint true;\nprint false;\n}",
			wantErr:  false,
		},

		// If/else statement tests
		{
			name: "simple if statement",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if (true) print \"hello\";",
			wantErr:  false,
		},
		{
			name: "if with expression condition",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.NUMBER, Lexeme: "5"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "3"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "greater"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if ((GREATER 5 3)) print \"greater\";",
			wantErr:  false,
		},
		{
			name: "if with block statement",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "in block"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "if (true) {\nprint \"in block\";\n}",
			wantErr:  false,
		},
		{
			name: "if-else statement",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.FALSE, Lexeme: "false"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "if true"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.ELSE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "else block"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if (false) print \"if true\"; else print \"else block\";",
			wantErr:  false,
		},
		{
			name: "if-else with blocks",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.EQUAL_EQUAL},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "equal"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.ELSE},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "not equal"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "if ((EQUAL_EQUAL 1 2)) {\nprint \"equal\";\n} else {\nprint \"not equal\";\n}",
			wantErr:  false,
		},
		{
			name: "nested if statements",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.FALSE, Lexeme: "false"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "nested"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if (true) if (false) print \"nested\";",
			wantErr:  false,
		},
		{
			name: "if-else-if chain",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "one"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.ELSE},
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "two"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.ELSE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "other"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if (1) print \"one\"; else if (2) print \"two\"; else print \"other\";",
			wantErr:  false,
		},
		{
			name: "if with variable condition",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "condition"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "variable condition"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if (condition) print \"variable condition\";",
			wantErr:  false,
		},
		{
			name: "if with complex condition",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.BANG_EQUAL},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "complex"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "if ((BANG_EQUAL (GREATER x 0) true)) print \"complex\";",
			wantErr:  false,
		},

		// While statement tests
		{
			name: "simple while statement",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "looping"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while (true) print \"looping\";",
			wantErr:  false,
		},
		{
			name: "while with expression condition",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "counter"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "10"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "counter"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while ((LESS counter 10)) print counter;",
			wantErr:  false,
		},
		{
			name: "while with block statement",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "running"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "in loop"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "while (running) {\nprint \"in loop\";\n}",
			wantErr:  false,
		},
		{
			name: "while with complex condition",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.AND},
				{Type: lexer.IDENTIFIER, Lexeme: "y"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "100"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "looping"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while ((AND (GREATER x 0) (LESS y 100))) print \"looping\";",
			wantErr:  false,
		},
		{
			name: "while with variable condition",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "keepGoing"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "still going"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while (keepGoing) print \"still going\";",
			wantErr:  false,
		},
		{
			name: "nested while statement",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "outer"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "inner"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "nested"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while (outer) while (inner) print \"nested\";",
			wantErr:  false,
		},

		// Error cases for while statements
		{
			name: "while without condition parentheses",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "while with missing opening paren",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "while with missing closing paren",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "while with missing condition",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "while with missing statement",
			tokens: []lexer.Token{
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},

		// Error cases for if statements
		{
			name: "if without condition parentheses",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "if with missing opening paren",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "if with missing closing paren",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "if with missing condition",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "if with missing statement",
			tokens: []lexer.Token{
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},

		// For statement tests
		{
			name: "basic for loop with all parts",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "10"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "{\nvar i = 0;\nwhile ((LESS i 10)) {\nprint i;\n(ASSIGN i (PLUS i 1));\n}\n}",
			wantErr:  false,
		},
		{
			name: "for loop with expression initializer",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "count"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "5"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "count"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "count"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "count"},
				{Type: lexer.MINUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "count"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "{\n(ASSIGN count 5);\nwhile ((GREATER count 0)) {\nprint count;\n(ASSIGN count (MINUS count 1));\n}\n}",
			wantErr:  false,
		},
		{
			name: "for loop with no initializer",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "running"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "counter"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "counter"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "looping"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while (running) {\nprint \"looping\";\n(ASSIGN counter (PLUS counter 1));\n}",
			wantErr:  false,
		},
		{
			name: "for loop with no condition",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "{\nvar i = 0;\nwhile (true) {\nprint i;\n(ASSIGN i (PLUS i 1));\n}\n}",
			wantErr:  false,
		},
		{
			name: "for loop with no increment",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.LESS_EQUAL},
				{Type: lexer.NUMBER, Lexeme: "5"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "{\nvar x = 1;\nwhile ((LESS_EQUAL x 5)) print x;\n}",
			wantErr:  false,
		},
		{
			name: "empty for loop (infinite loop syntax)",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "infinite"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "while (true) print \"infinite\";",
			wantErr:  false,
		},
		{
			name: "for loop with block statement",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "3"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "iteration"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "{\nvar i = 0;\nwhile ((LESS i 3)) {\n{\nprint \"iteration\";\nprint i;\n}\n(ASSIGN i (PLUS i 1));\n}\n}",
			wantErr:  false,
		},
		{
			name: "for loop with complex expressions",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.STAR},
				{Type: lexer.NUMBER, Lexeme: "3"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "10"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "5"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.STAR},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "{\nvar start = (STAR 2 3);\nwhile ((LESS start (PLUS 10 5))) {\nprint start;\n(ASSIGN start (STAR start 2));\n}\n}",
			wantErr:  false,
		},
		{
			name: "nested for loops",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "j"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "j"},
				{Type: lexer.LESS},
				{Type: lexer.NUMBER, Lexeme: "2"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "j"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "j"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "nested"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "{\nvar i = 0;\nwhile ((LESS i 2)) {\n{\nvar j = 0;\nwhile ((LESS j 2)) {\nprint \"nested\";\n(ASSIGN j (PLUS j 1));\n}\n}\n(ASSIGN i (PLUS i 1));\n}\n}",
			wantErr:  false,
		},

		// Error cases for for statements
		{
			name: "for without opening paren",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "for without closing paren",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.SEMICOLON},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "hello"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "for without first semicolon",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "for without second semicolon",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.TRUE, Lexeme: "true"},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "i"},
				{Type: lexer.PLUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "for without statement Body",
			tokens: []lexer.Token{
				{Type: lexer.FOR},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},

		// Function statement tests
		{
			name: "function with no parameters",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "sayHello"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "Hello!"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun sayHello() {\nprint \"Hello!\";\n}",
			wantErr:  false,
		},
		{
			name: "function with single parameter",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "greet"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "Name"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "Name"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun greet(Name) {\nprint Name;\n}",
			wantErr:  false,
		},
		{
			name: "function with multiple parameters",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "add"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "a"},
				{Type: lexer.COMMA},
				{Type: lexer.IDENTIFIER, Lexeme: "b"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "a"},
				{Type: lexer.PLUS},
				{Type: lexer.IDENTIFIER, Lexeme: "b"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun add(a, b) {\nprint (PLUS a b);\n}",
			wantErr:  false,
		},
		{
			name: "function with three parameters",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "calculate"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.COMMA},
				{Type: lexer.IDENTIFIER, Lexeme: "y"},
				{Type: lexer.COMMA},
				{Type: lexer.IDENTIFIER, Lexeme: "z"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun calculate(x, y, z) {\nprint x;\n}",
			wantErr:  false,
		},
		{
			name: "function with empty Body",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "empty"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun empty() {\n}",
			wantErr:  false,
		},
		{
			name: "function with complex Body",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "complex"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "n"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.IF},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "n"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "positive"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.ELSE},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.STRING, Lexeme: "non-positive"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun complex(n) {\nif ((GREATER n 0)) {\nprint \"positive\";\n} else {\nprint \"non-positive\";\n}\n}",
			wantErr:  false,
		},
		{
			name: "function with variable declarations",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "withVars"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.VAR},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.EQUAL},
				{Type: lexer.NUMBER, Lexeme: "42"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun withVars() {\nvar x = 42;\nprint x;\n}",
			wantErr:  false,
		},
		{
			name: "call function on member select",
			tokens: []lexer.Token{
				{Type: lexer.IDENTIFIER, Lexeme: "x"},
				{Type: lexer.DOT},
				{Type: lexer.IDENTIFIER, Lexeme: "y"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "x.y();",
			wantErr:  false,
		},
		{
			name: "function with loops",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "countdown"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.WHILE},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.GREATER},
				{Type: lexer.NUMBER, Lexeme: "0"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.PRINT},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.EQUAL},
				{Type: lexer.IDENTIFIER, Lexeme: "start"},
				{Type: lexer.MINUS},
				{Type: lexer.NUMBER, Lexeme: "1"},
				{Type: lexer.SEMICOLON},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "fun countdown(start) {\nwhile ((GREATER start 0)) {\nprint start;\n(ASSIGN start (MINUS start 1));\n}\n}",
			wantErr:  false,
		},

		// Error cases for function statements
		{
			name: "function without Name",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "function without opening paren",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "test"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "function without closing paren",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "test"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "param"},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "function without Body",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "test"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "function with invalid parameter",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "test"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.NUMBER, Lexeme: "42"},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "function with trailing comma in parameters",
			tokens: []lexer.Token{
				{Type: lexer.FUN},
				{Type: lexer.IDENTIFIER, Lexeme: "test"},
				{Type: lexer.LEFT_PAREN},
				{Type: lexer.IDENTIFIER, Lexeme: "a"},
				{Type: lexer.COMMA},
				{Type: lexer.RIGHT_PAREN},
				{Type: lexer.LEFT_BRACE},
				{Type: lexer.RIGHT_BRACE},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},

		// Error cases for statements
		{
			name: "print without semicolon",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: "42"},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "expression without semicolon",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "42"},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "print with missing expression",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(tt.tokens)
			stmt, err := parser.declStatement()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			result, err2 := Print(stmt)
			if err2 != nil {
				t.Fatalf("Unexpected error: %v", err2)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	data, err := os.ReadFile("testdata/function_errors.lox")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	l := lexer.New(bytes.NewReader(data))
	toks := lexer.Consume(l)
	if toks.Err != nil {
		t.Fatalf("Failed to lex test file: %v", toks.Err)
	}

	parser := New(toks.Tokens)
	_, _ = parser.Parse()

	fmt.Print(parser.ReportErrors())
}
