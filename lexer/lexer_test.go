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

	for i, token := range lex.Tokens {
		if i >= len(exp) {
			t.Fatalf("Not enough Tokens in lexer result")
		}

		e := strings.TrimSpace(strings.Split(exp[i], "/")[0])
		if e != token.Type.String() {
			t.Errorf("%s %s", token.Type, exp[i])
		} else {
			fmt.Printf("%s (%s) == %s\n", token.Type, token.Lexeme, exp[i])
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
				t.Errorf("Expected lexeme %s found %s", tt.expected, token.Lexeme)
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
