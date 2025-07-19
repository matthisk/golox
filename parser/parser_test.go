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
				{Type: lexer.NUMBER, Lexeme: 42, Position: 0},
				{Type: lexer.EOF},
			},
			expected: "42",
			wantErr:  false,
		},
		{
			name: "string literal",
			tokens: []lexer.Token{
				{Type: lexer.STRING, Lexeme: "hello", Position: 0},
				{Type: lexer.EOF},
			},
			expected: "hello",
			wantErr:  false,
		},
		{
			name: "true literal",
			tokens: []lexer.Token{
				{Type: lexer.TRUE, Lexeme: "true", Position: 0},
				{Type: lexer.EOF},
			},
			expected: "true",
			wantErr:  false,
		},
		{
			name: "false literal",
			tokens: []lexer.Token{
				{Type: lexer.FALSE, Lexeme: "false", Position: 0},
				{Type: lexer.EOF},
			},
			expected: "false",
			wantErr:  false,
		},
		{
			name: "nil literal",
			tokens: []lexer.Token{
				{Type: lexer.NIL, Lexeme: "nil", Position: 0},
				{Type: lexer.EOF},
			},
			expected: "nil",
			wantErr:  false,
		},
		{
			name: "decimal number",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "3.14", Position: 0},
				{Type: lexer.EOF},
			},
			expected: "3.14",
			wantErr:  false,
		},

		// Grouping tests
		{
			name: "simple grouping",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, Position: 0},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 1},
				{Type: lexer.RIGHT_PAREN, Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(group 5)",
			wantErr:  false,
		},
		{
			name: "nested grouping",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, Position: 0},
				{Type: lexer.LEFT_PAREN, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "10", Position: 2},
				{Type: lexer.RIGHT_PAREN, Position: 3},
				{Type: lexer.RIGHT_PAREN, Position: 4},
				{Type: lexer.EOF},
			},
			expected: "(group (group 10))",
			wantErr:  false,
		},

		// Unary tests
		{
			name: "unary minus",
			tokens: []lexer.Token{
				{Type: lexer.MINUS, Position: 0},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 1},
				{Type: lexer.EOF},
			},
			expected: "(MINUS 5)",
			wantErr:  false,
		},
		{
			name: "unary bang",
			tokens: []lexer.Token{
				{Type: lexer.BANG, Position: 0},
				{Type: lexer.TRUE, Lexeme: "true", Position: 1},
				{Type: lexer.EOF},
			},
			expected: "(BANG true)",
			wantErr:  false,
		},
		{
			name: "double unary",
			tokens: []lexer.Token{
				{Type: lexer.MINUS, Position: 0},
				{Type: lexer.MINUS, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(MINUS (MINUS 3))",
			wantErr:  false,
		},
		{
			name: "unary with grouping",
			tokens: []lexer.Token{
				{Type: lexer.BANG, Position: 0},
				{Type: lexer.LEFT_PAREN, Position: 1},
				{Type: lexer.FALSE, Lexeme: "false", Position: 2},
				{Type: lexer.RIGHT_PAREN, Position: 3},
				{Type: lexer.EOF},
			},
			expected: "(BANG (group false))",
			wantErr:  false,
		},

		// Binary arithmetic tests
		{
			name: "simple addition",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", Position: 0},
				{Type: lexer.PLUS, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(PLUS 5 3)",
			wantErr:  false,
		},
		{
			name: "simple subtraction",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "10", Position: 0},
				{Type: lexer.MINUS, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "4", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(MINUS 10 4)",
			wantErr:  false,
		},
		{
			name: "simple multiplication",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "6", Position: 0},
				{Type: lexer.STAR, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "7", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(STAR 6 7)",
			wantErr:  false,
		},
		{
			name: "simple division",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "8", Position: 0},
				{Type: lexer.SLASH, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "2", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(SLASH 8 2)",
			wantErr:  false,
		},

		// Binary comparison tests
		{
			name: "equality",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", Position: 0},
				{Type: lexer.EQUAL_EQUAL, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(EQUAL_EQUAL 5 5)",
			wantErr:  false,
		},
		{
			name: "inequality",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", Position: 0},
				{Type: lexer.BANG_EQUAL, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(BANG_EQUAL 5 3)",
			wantErr:  false,
		},
		{
			name: "less than",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "3", Position: 0},
				{Type: lexer.LESS, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(LESS 3 5)",
			wantErr:  false,
		},
		{
			name: "less than or equal",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "3", Position: 0},
				{Type: lexer.LESS_EQUAL, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(LESS_EQUAL 3 5)",
			wantErr:  false,
		},
		{
			name: "greater than",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "7", Position: 0},
				{Type: lexer.GREATER, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(GREATER 7 3)",
			wantErr:  false,
		},
		{
			name: "greater than or equal",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "7", Position: 0},
				{Type: lexer.GREATER_EQUAL, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(GREATER_EQUAL 7 3)",
			wantErr:  false,
		},

		// Complex expressions (operator precedence)
		{
			name: "addition and multiplication",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "2", Position: 0},
				{Type: lexer.PLUS, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.STAR, Position: 3},
				{Type: lexer.NUMBER, Lexeme: "4", Position: 4},
				{Type: lexer.EOF},
			},
			expected: "(PLUS 2 (STAR 3 4))",
			wantErr:  false,
		},
		{
			name: "multiplication and addition",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "2", Position: 0},
				{Type: lexer.STAR, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.PLUS, Position: 3},
				{Type: lexer.NUMBER, Lexeme: "4", Position: 4},
				{Type: lexer.EOF},
			},
			expected: "(PLUS (STAR 2 3) 4)",
			wantErr:  false,
		},
		{
			name: "comparison and arithmetic",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", Position: 0},
				{Type: lexer.PLUS, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 2},
				{Type: lexer.GREATER, Position: 3},
				{Type: lexer.NUMBER, Lexeme: "2", Position: 4},
				{Type: lexer.STAR, Position: 5},
				{Type: lexer.NUMBER, Lexeme: "4", Position: 6},
				{Type: lexer.EOF},
			},
			expected: "(GREATER (PLUS 5 3) (STAR 2 4))",
			wantErr:  false,
		},

		// Complex expressions with grouping
		{
			name: "grouping affects precedence",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, Position: 0},
				{Type: lexer.NUMBER, Lexeme: "2", Position: 1},
				{Type: lexer.PLUS, Position: 2},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 3},
				{Type: lexer.RIGHT_PAREN, Position: 4},
				{Type: lexer.STAR, Position: 5},
				{Type: lexer.NUMBER, Lexeme: "4", Position: 6},
				{Type: lexer.EOF},
			},
			expected: "(STAR (group (PLUS 2 3)) 4)",
			wantErr:  false,
		},
		{
			name: "complex nested expression",
			tokens: []lexer.Token{ // - ( 5 + 3 ) * 2
				{Type: lexer.MINUS, Position: 0},
				{Type: lexer.LEFT_PAREN, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 2},
				{Type: lexer.PLUS, Position: 3},
				{Type: lexer.NUMBER, Lexeme: "3", Position: 4},
				{Type: lexer.RIGHT_PAREN, Position: 5},
				{Type: lexer.STAR, Position: 6},
				{Type: lexer.NUMBER, Lexeme: "2", Position: 7},
				{Type: lexer.EOF},
			},
			expected: "(STAR (MINUS (group (PLUS 5 3))) 2)",
			wantErr:  false,
		},

		// Mixed types
		{
			name: "string concatenation",
			tokens: []lexer.Token{
				{Type: lexer.STRING, Lexeme: "hello", Position: 0},
				{Type: lexer.PLUS, Position: 1},
				{Type: lexer.STRING, Lexeme: "world", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(PLUS hello world)",
			wantErr:  false,
		},
		{
			name: "boolean comparison",
			tokens: []lexer.Token{
				{Type: lexer.TRUE, Lexeme: "true", Position: 0},
				{Type: lexer.EQUAL_EQUAL, Position: 1},
				{Type: lexer.FALSE, Lexeme: "false", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(EQUAL_EQUAL true false)",
			wantErr:  false,
		},
		{
			name: "nil comparison",
			tokens: []lexer.Token{
				{Type: lexer.NIL, Lexeme: "nil", Position: 0},
				{Type: lexer.BANG_EQUAL, Position: 1},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 2},
				{Type: lexer.EOF},
			},
			expected: "(BANG_EQUAL nil 5)",
			wantErr:  false,
		},

		// Error cases
		{
			name: "missing closing parenthesis",
			tokens: []lexer.Token{
				{Type: lexer.LEFT_PAREN, Position: 0},
				{Type: lexer.NUMBER, Lexeme: "5", Position: 1},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "missing operand",
			tokens: []lexer.Token{
				{Type: lexer.NUMBER, Lexeme: "5", Position: 0},
				{Type: lexer.PLUS, Position: 1},
				{Type: lexer.EOF},
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "invalid token",
			tokens: []lexer.Token{
				{Type: lexer.SEMICOLON, Position: 0},
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
