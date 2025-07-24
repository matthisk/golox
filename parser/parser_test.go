package parser

import (
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
			expected: "hello",
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
			expected: "(PLUS hello world)",
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
		// Literal tests
		{
			name: "print statement",
			tokens: []lexer.Token{
				{Type: lexer.PRINT, StartPos: lexer.Pos{Offset: 0, Line: 1, Column: 1}},
				{Type: lexer.NUMBER, Lexeme: 42},
				{Type: lexer.SEMICOLON},
				{Type: lexer.EOF},
			},
			expected: "print 42;",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(tt.tokens)
			stmt, err := parser.statement()

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
