package lexer

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	data, err := os.ReadFile("testdata/main.lox")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected, err := os.ReadFile("testdata/main.lexed")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	lx := New(bytes.NewReader(data))
	lex := Consume(lx)

	if lex.Err != nil {
		t.Fatalf("Unexpected error %v", lex.Err)
		return
	}

	exp := strings.Split(string(expected), "\n")

	// Verify token types
	for i, token := range lex.Tokens {
		if i >= len(exp) {
			t.Fatalf("Not enough Tokens in lexer result")
		}

		e := strings.TrimSpace(strings.Split(exp[i], "/")[0])
		if e != token.Type.String() {
			t.Errorf("%d %s %s", i, token.Type, exp[i])
		}
	}

	// Verify token positions
	source := string(data)
	verifyTokenPositions(t, lex.Tokens, source)
}

// verifyTokenPositions checks that each token's StartPos and EndPos are correct
func verifyTokenPositions(t *testing.T, tokens []Token, source string) {
	for i, token := range tokens {
		if token.Type == EOF {
			// EOF token should be positioned where lexing actually ended
			// (which may be before the literal end of file if there are comments)
			// We just need to verify it's within valid bounds
			if token.StartPos.Offset < 0 || token.StartPos.Offset > len(source) {
				t.Errorf("Token %d (EOF): offset %d is out of bounds [0, %d]", i, token.StartPos.Offset, len(source))
			}
			continue
		}

		// Verify that StartPos offset is within source bounds
		if token.StartPos.Offset < 0 || token.StartPos.Offset >= len(source) {
			t.Errorf("Token %d (%s): start offset %d is out of bounds [0, %d)", i, token.Type, token.StartPos.Offset, len(source))
			continue
		}

		// Verify that EndPos offset is within source bounds
		if token.EndPos.Offset < 0 || token.EndPos.Offset > len(source) {
			t.Errorf("Token %d (%s): end offset %d is out of bounds [0, %d]", i, token.Type, token.EndPos.Offset, len(source))
			continue
		}

		// Verify that StartPos comes before EndPos
		if token.StartPos.Offset >= token.EndPos.Offset {
			t.Errorf("Token %d (%s): start offset %d should be less than end offset %d", i, token.Type, token.StartPos.Offset, token.EndPos.Offset)
			continue
		}

		// Extract the actual text from source using positions
		actualText := source[token.StartPos.Offset:token.EndPos.Offset]

		// Verify position correctness based on token type
		switch token.Type {
		case IDENTIFIER:
			expectedText := token.Lexeme.(string)
			if actualText != expectedText {
				t.Errorf("Token %d (IDENTIFIER): position mismatch - expected '%s', got '%s' at [%d:%d]",
					i, expectedText, actualText, token.StartPos.Offset, token.EndPos.Offset)
			}
		case STRING:
			// For strings, the positions should include the quotes
			expectedText := fmt.Sprintf(`"%s"`, token.Lexeme.(string))
			if actualText != expectedText {
				t.Errorf("Token %d (STRING): position mismatch - expected '%s', got '%s' at [%d:%d]",
					i, expectedText, actualText, token.StartPos.Offset, token.EndPos.Offset)
			}
		case NUMBER:
			// For numbers, convert back to string representation
			expectedText := fmt.Sprintf("%g", token.Lexeme.(float64))
			if actualText != expectedText {
				t.Errorf("Token %d (NUMBER): position mismatch - expected '%s', got '%s' at [%d:%d]",
					i, expectedText, actualText, token.StartPos.Offset, token.EndPos.Offset)
			}
		default:
			// For operators and keywords, verify against expected text
			var expectedText string
			switch token.Type {
			case LEFT_PAREN:
				expectedText = "("
			case RIGHT_PAREN:
				expectedText = ")"
			case LEFT_BRACE:
				expectedText = "{"
			case RIGHT_BRACE:
				expectedText = "}"
			case COMMA:
				expectedText = ","
			case DOT:
				expectedText = "."
			case MINUS:
				expectedText = "-"
			case PLUS:
				expectedText = "+"
			case SEMICOLON:
				expectedText = ";"
			case SLASH:
				expectedText = "/"
			case STAR:
				expectedText = "*"
			case COLON:
				expectedText = ":"
			case QUESTION_MARK:
				expectedText = "?"
			case BANG:
				expectedText = "!"
			case BANG_EQUAL:
				expectedText = "!="
			case EQUAL:
				expectedText = "="
			case EQUAL_EQUAL:
				expectedText = "=="
			case GREATER:
				expectedText = ">"
			case GREATER_EQUAL:
				expectedText = ">="
			case LESS:
				expectedText = "<"
			case LESS_EQUAL:
				expectedText = "<="
			case AND:
				expectedText = "and"
			case CLASS:
				expectedText = "class"
			case ELSE:
				expectedText = "else"
			case FALSE:
				expectedText = "false"
			case FUN:
				expectedText = "fun"
			case FOR:
				expectedText = "for"
			case IF:
				expectedText = "if"
			case NIL:
				expectedText = "nil"
			case OR:
				expectedText = "or"
			case PRINT:
				expectedText = "print"
			case RETURN:
				expectedText = "return"
			case SUPER:
				expectedText = "super"
			case THIS:
				expectedText = "this"
			case TRUE:
				expectedText = "true"
			case VAR:
				expectedText = "var"
			case WHILE:
				expectedText = "while"
			}

			if expectedText != "" && actualText != expectedText {
				t.Errorf("Token %d (%s): position mismatch - expected '%s', got '%s' at [%d:%d]",
					i, token.Type, expectedText, actualText, token.StartPos.Offset, token.EndPos.Offset)
			}
		}

		// Verify line and column positions make sense
		if token.StartPos.Line < 1 {
			t.Errorf("Token %d (%s): invalid start line %d (should be >= 1)", i, token.Type, token.StartPos.Line)
		}
		if token.StartPos.Column < 1 {
			t.Errorf("Token %d (%s): invalid start column %d (should be >= 1)", i, token.Type, token.StartPos.Column)
		}
		if token.EndPos.Line < 1 {
			t.Errorf("Token %d (%s): invalid end line %d (should be >= 1)", i, token.Type, token.EndPos.Line)
		}
		if token.EndPos.Column < 1 {
			t.Errorf("Token %d (%s): invalid end column %d (should be >= 1)", i, token.Type, token.EndPos.Column)
		}

		// Verify line progression makes sense
		if token.EndPos.Line < token.StartPos.Line {
			t.Errorf("Token %d (%s): end line %d should not be less than start line %d",
				i, token.Type, token.EndPos.Line, token.StartPos.Line)
		}
	}
}

func TestScanNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "whole number",
			input:    "42",
			expected: 42.0,
		},
		{
			name:     "floating point number",
			input:    "3.14159",
			expected: 3.14159,
		},
		{
			name:     "zero",
			input:    "0",
			expected: 0.0,
		},
		{
			name:     "decimal starting with zero",
			input:    "0.5",
			expected: 0.5,
		},
		{
			name:     "large number",
			input:    "123456789",
			expected: 123456789.0,
		},
		{
			name:     "very small decimal",
			input:    "0.001",
			expected: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lx := New(strings.NewReader(tt.input))

			token := lx.Next()

			if token.Type != NUMBER {
				t.Errorf("Expected to scan NUMBER found %s", token.Type)
			}
			if token.Lexeme != tt.expected {
				t.Errorf("Expected lexeme %f found %f", tt.expected, token.Lexeme)
			}
		})
	}
}

func TestScanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    `"hello"`,
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: "",
		},
		{
			name:     "string with spaces",
			input:    `"hello world"`,
			expected: "hello world",
		},
		{
			name:     "string with special characters",
			input:    `"hello@#$%^&*()world"`,
			expected: "hello@#$%^&*()world",
		},
		{
			name:     "string with numbers",
			input:    `"test123"`,
			expected: "test123",
		},
		{
			name:     "string with newline",
			input:    "\"hello\nworld\"",
			expected: "hello\nworld",
		},
		{
			name:     "string with multiple newlines",
			input:    "\"line1\nline2\nline3\"",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "string with tab character",
			input:    "\"hello\tworld\"",
			expected: "hello\tworld",
		},
		{
			name:     "string with mixed whitespace",
			input:    "\"hello\n\tworld  \"",
			expected: "hello\n\tworld  ",
		},
		{
			name:     "string with punctuation",
			input:    `"Hello, World! How are you?"`,
			expected: "Hello, World! How are you?",
		},
		{
			name:     "string with comment inside",
			input:    `"Hello, World! // How are you?"`,
			expected: "Hello, World! // How are you?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lx := New(strings.NewReader(tt.input))

			token := lx.Next()

			if token.Type != STRING {
				t.Errorf("Expected to scan STRING found %s", token.Type)
			}
			if token.Lexeme != tt.expected {
				t.Errorf("Expected lexeme %s found %s", tt.expected, token.Lexeme)
			}
		})
	}
}

func TestScanIdentifier(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedType   TokenType
		expectedLexeme string
	}{
		{
			name:           "simple identifier",
			input:          "variable",
			expectedType:   IDENTIFIER,
			expectedLexeme: "variable",
		},
		{
			name:           "identifier with underscore",
			input:          "my_variable",
			expectedType:   IDENTIFIER,
			expectedLexeme: "my_variable",
		},
		{
			name:           "identifier with numbers",
			input:          "var123",
			expectedType:   IDENTIFIER,
			expectedLexeme: "var123",
		},
		{
			name:           "single letter identifier",
			input:          "x",
			expectedType:   IDENTIFIER,
			expectedLexeme: "x",
		},
		{
			name:           "identifier starting with uppercase",
			input:          "MyClass",
			expectedType:   IDENTIFIER,
			expectedLexeme: "MyClass",
		},
		{
			name:           "keyword and",
			input:          "and",
			expectedType:   AND,
			expectedLexeme: "",
		},
		{
			name:           "keyword class",
			input:          "class",
			expectedType:   CLASS,
			expectedLexeme: "",
		},
		{
			name:           "keyword else",
			input:          "else",
			expectedType:   ELSE,
			expectedLexeme: "",
		},
		{
			name:           "keyword false",
			input:          "false",
			expectedType:   FALSE,
			expectedLexeme: "",
		},
		{
			name:           "keyword fun",
			input:          "fun",
			expectedType:   FUN,
			expectedLexeme: "",
		},
		{
			name:           "keyword for",
			input:          "for",
			expectedType:   FOR,
			expectedLexeme: "",
		},
		{
			name:           "keyword if",
			input:          "if",
			expectedType:   IF,
			expectedLexeme: "",
		},
		{
			name:           "keyword nil",
			input:          "nil",
			expectedType:   NIL,
			expectedLexeme: "",
		},
		{
			name:           "keyword or",
			input:          "or",
			expectedType:   OR,
			expectedLexeme: "",
		},
		{
			name:           "keyword print",
			input:          "print",
			expectedType:   PRINT,
			expectedLexeme: "",
		},
		{
			name:           "keyword return",
			input:          "return",
			expectedType:   RETURN,
			expectedLexeme: "",
		},
		{
			name:           "keyword super",
			input:          "super",
			expectedType:   SUPER,
			expectedLexeme: "",
		},
		{
			name:           "keyword this",
			input:          "this",
			expectedType:   THIS,
			expectedLexeme: "",
		},
		{
			name:           "keyword true",
			input:          "true",
			expectedType:   TRUE,
			expectedLexeme: "",
		},
		{
			name:           "keyword var",
			input:          "var",
			expectedType:   VAR,
			expectedLexeme: "",
		},
		{
			name:           "keyword while",
			input:          "while",
			expectedType:   WHILE,
			expectedLexeme: "",
		},
		{
			name:           "identifier that looks like keyword",
			input:          "variable_and",
			expectedType:   IDENTIFIER,
			expectedLexeme: "variable_and",
		},
		{
			name:           "identifier with keyword prefix",
			input:          "classes",
			expectedType:   IDENTIFIER,
			expectedLexeme: "classes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lx := New(strings.NewReader(tt.input))

			token := lx.Next()

			if token.Type != tt.expectedType {
				t.Errorf("Expected to scan %s found %s", tt.expectedType, token.Type)
			}
			if token.Lexeme != tt.expectedLexeme {
				t.Errorf("Expected lexeme %s found %s", tt.expectedLexeme, token.Lexeme)
			}
		})
	}
}
