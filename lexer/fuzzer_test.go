package lexer

import (
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// FuzzInput represents a single fuzzing test case
type FuzzInput struct {
	input       string
	description string
}

// TestLexerFuzzer runs fuzzy tests against the lexer
func TestLexerFuzzer(t *testing.T) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	fuzzer := NewLexerFuzzer()

	// Run different types of fuzzing tests
	t.Run("RandomASCII", func(t *testing.T) {
		fuzzer.testRandomASCII(t, 100)
	})

	t.Run("RandomUnicode", func(t *testing.T) {
		fuzzer.testRandomUnicode(t, 50)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		fuzzer.testEdgeCases(t)
	})

	t.Run("LongInputs", func(t *testing.T) {
		fuzzer.testLongInputs(t, 20)
	})

	t.Run("NestedStructures", func(t *testing.T) {
		fuzzer.testNestedStructures(t, 30)
	})

	t.Run("MalformedTokens", func(t *testing.T) {
		fuzzer.testMalformedTokens(t, 50)
	})
}

// LexerFuzzer contains methods for fuzzing the lexer
type LexerFuzzer struct {
	// Common Lox tokens and keywords for realistic fuzzing
	keywords   []string
	operators  []string
	delimiters []string
	literals   []string
}

// NewLexerFuzzer creates a new fuzzer instance
func NewLexerFuzzer() *LexerFuzzer {
	return &LexerFuzzer{
		keywords: []string{
			"and", "class", "else", "false", "fun", "for", "if", "nil", "or",
			"print", "return", "super", "this", "true", "var", "while", "break", "continue",
		},
		operators: []string{
			"+", "-", "*", "/", "=", "==", "!=", "<", "<=", ">", ">=", "!", "?", ":",
		},
		delimiters: []string{
			"(", ")", "{", "}", "[", "]", ",", ";", ".",
		},
		literals: []string{
			`"hello"`, `"world"`, `""`, `"test string"`,
			"123", "45.67", "0", "999.999",
			"true", "false", "nil",
		},
	}
}

// testRandomASCII generates random ASCII strings
func (f *LexerFuzzer) testRandomASCII(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		input := f.generateRandomASCII(rand.Intn(200) + 1)
		f.runFuzzTest(t, input, "RandomASCII")
	}
}

// testRandomUnicode generates random Unicode strings
func (f *LexerFuzzer) testRandomUnicode(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		input := f.generateRandomUnicode(rand.Intn(100) + 1)
		f.runFuzzTest(t, input, "RandomUnicode")
	}
}

// testEdgeCases tests specific edge cases
func (f *LexerFuzzer) testEdgeCases(t *testing.T) {
	edgeCases := []FuzzInput{
		{"", "empty string"},
		{" ", "single space"},
		{"\n", "single newline"},
		{"\t", "single tab"},
		{"\r\n", "CRLF"},
		{"//", "comment start only"},
		{"/* */", "empty block comment"},
		{"/*", "unterminated block comment"},
		{"\"", "single quote"},
		{"\"\"\"", "triple quotes"},
		{"\\", "single backslash"},
		{"\\n\\t\\r", "escape sequences"},
		{"123.456.789", "multiple dots in number"},
		{"+++++", "many operators"},
		{"((((((", "many open parens"},
		{"))))))", "many close parens"},
		{"........", "many dots"},
		{"_", "single underscore"},
		{"____", "many underscores"},
		{"0x", "incomplete hex"},
		{"1e", "incomplete exponential"},
		{string([]byte{0, 1, 2, 3, 255}), "control characters"},
	}

	for _, tc := range edgeCases {
		f.runFuzzTest(t, tc.input, tc.description)
	}
}

// testLongInputs generates very long inputs
func (f *LexerFuzzer) testLongInputs(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		length := rand.Intn(5000) + 1000 // 1000-6000 characters
		input := f.generateRandomLoxCode(length)
		f.runFuzzTest(t, input, "LongInput")
	}
}

// testNestedStructures generates deeply nested structures
func (f *LexerFuzzer) testNestedStructures(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		depth := rand.Intn(50) + 10
		input := f.generateNestedStructure(depth)
		f.runFuzzTest(t, input, "NestedStructure")
	}
}

// testMalformedTokens generates intentionally malformed tokens
func (f *LexerFuzzer) testMalformedTokens(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		input := f.generateMalformedTokens()
		f.runFuzzTest(t, input, "MalformedTokens")
	}
}

// generateRandomASCII creates a random ASCII string
func (f *LexerFuzzer) generateRandomASCII(length int) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		// Generate printable ASCII characters (32-126) plus some whitespace
		char := byte(32 + rand.Intn(95))
		if rand.Float32() < 0.1 { // 10% chance of whitespace
			char = byte([]byte{' ', '\t', '\n', '\r'}[rand.Intn(4)])
		}
		builder.WriteByte(char)
	}
	return builder.String()
}

// generateRandomUnicode creates a random Unicode string
func (f *LexerFuzzer) generateRandomUnicode(length int) string {
	var builder strings.Builder
	for i := 0; i < length; i++ {
		// Generate random Unicode code points
		var r rune
		for {
			r = rune(rand.Intn(0x10000)) // Basic Multilingual Plane
			if utf8.ValidRune(r) && r != 0 {
				break
			}
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

// generateRandomLoxCode creates semi-realistic Lox code
func (f *LexerFuzzer) generateRandomLoxCode(targetLength int) string {
	var builder strings.Builder
	
	for builder.Len() < targetLength {
		switch rand.Intn(6) {
		case 0: // Add a keyword
			builder.WriteString(f.keywords[rand.Intn(len(f.keywords))])
			builder.WriteString(" ")
		case 1: // Add an operator
			builder.WriteString(f.operators[rand.Intn(len(f.operators))])
		case 2: // Add a delimiter
			builder.WriteString(f.delimiters[rand.Intn(len(f.delimiters))])
		case 3: // Add a literal
			builder.WriteString(f.literals[rand.Intn(len(f.literals))])
		case 4: // Add an identifier
			builder.WriteString(f.generateRandomIdentifier())
		case 5: // Add whitespace or comment
			if rand.Float32() < 0.7 {
				builder.WriteString(" ")
			} else {
				builder.WriteString("// random comment\n")
			}
		}
	}
	
	return builder.String()
}

// generateRandomIdentifier creates a random valid identifier
func (f *LexerFuzzer) generateRandomIdentifier() string {
	length := rand.Intn(10) + 1
	var builder strings.Builder
	
	// First character must be letter or underscore
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	builder.WriteByte(chars[rand.Intn(len(chars))])
	
	// Subsequent characters can be letters, digits, or underscores
	chars += "0123456789"
	for i := 1; i < length; i++ {
		builder.WriteByte(chars[rand.Intn(len(chars))])
	}
	
	return builder.String()
}

// generateNestedStructure creates deeply nested parentheses/braces
func (f *LexerFuzzer) generateNestedStructure(depth int) string {
	var builder strings.Builder
	
	// Opening brackets
	for i := 0; i < depth; i++ {
		switch rand.Intn(3) {
		case 0:
			builder.WriteString("(")
		case 1:
			builder.WriteString("{")
		case 2:
			builder.WriteString("\"")
		}
		if rand.Float32() < 0.3 {
			builder.WriteString(" ")
		}
	}
	
	// Some content in the middle
	builder.WriteString("content")
	
	// Closing brackets (in reverse order for some)
	for i := 0; i < depth; i++ {
		if rand.Float32() < 0.3 {
			builder.WriteString(" ")
		}
		switch rand.Intn(3) {
		case 0:
			builder.WriteString(")")
		case 1:
			builder.WriteString("}")
		case 2:
			builder.WriteString("\"")
		}
	}
	
	return builder.String()
}

// generateMalformedTokens creates intentionally broken tokens
func (f *LexerFuzzer) generateMalformedTokens() string {
	malformed := []string{
		"123abc",           // number followed by letters
		"\"unclosed string", // unclosed string
		"/* unclosed comment", // unclosed comment
		"===",              // too many equals
		"!!!=",             // mixed operators
		"...",              // multiple dots
		"123.456.789",      // multiple decimal points
		"0x",               // incomplete hex
		"_123_abc_",        // mixed identifier
		"@#$%^&",           // invalid operators
		"func ion",         // split keyword
		"tru e",            // split literal
		"nil123",           // literal with number
	}
	
	var builder strings.Builder
	count := rand.Intn(5) + 1
	
	for i := 0; i < count; i++ {
		builder.WriteString(malformed[rand.Intn(len(malformed))])
		builder.WriteString(" ")
	}
	
	return builder.String()
}

// runFuzzTest executes a single fuzz test
func (f *LexerFuzzer) runFuzzTest(t *testing.T, input, testType string) {
	// The fuzzer should not panic - that's the main requirement
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PANIC in %s fuzzer with input %q: %v", testType, truncateString(input, 50), r)
		}
	}()
	
	// Create lexer and try to tokenize
	lexer := New(strings.NewReader(input))
	result := Consume(lexer)
	
	// We don't care about errors - fuzzing should find crashes, not logical errors
	// But we do want to make sure the result structure is valid
	if result.Tokens == nil {
		t.Errorf("Lexer returned nil tokens for input: %q", truncateString(input, 50))
	}
	
	// Basic sanity check - if no error, should have at least EOF token
	if result.Err == nil && len(result.Tokens) == 0 {
		t.Errorf("Lexer returned no tokens and no error for input: %q", truncateString(input, 50))
	}
	
	// Verify all tokens have valid positions
	for i, token := range result.Tokens {
		if token.StartPos.Line < 1 || token.StartPos.Column < 1 {
			t.Errorf("Invalid token position at index %d: %+v for input: %q", 
				i, token.StartPos, truncateString(input, 50))
		}
	}
}

// truncateString truncates a string for display purposes
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// BenchmarkLexerFuzzer benchmarks the fuzzer performance
func BenchmarkLexerFuzzer(b *testing.B) {
	fuzzer := NewLexerFuzzer()
	
	b.Run("ShortInputs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			input := fuzzer.generateRandomLoxCode(100)
			lexer := New(strings.NewReader(input))
			Consume(lexer)
		}
	})
	
	b.Run("LongInputs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			input := fuzzer.generateRandomLoxCode(1000)
			lexer := New(strings.NewReader(input))
			Consume(lexer)
		}
	})
}

// TestSpecificFuzzCases tests cases that have caused issues before
func TestSpecificFuzzCases(t *testing.T) {
	// Add specific test cases here that have been found to cause problems
	testCases := []struct {
		name  string
		input string
	}{
		{"Empty input", ""},
		{"Only whitespace", "   \n\t\r  "},
		{"Only comments", "// comment\n/* block */"},
		{"Unclosed string", "\"hello world"},
		{"Unclosed comment", "/* this never ends"},
		{"Invalid UTF-8", string([]byte{0xff, 0xfe, 0xfd})},
		{"Very long identifier", strings.Repeat("a", 1000)},
		{"Very long string", "\"" + strings.Repeat("x", 1000) + "\""},
		{"Many operators", strings.Repeat("!=", 100)},
		{"Deep nesting", strings.Repeat("(", 100) + strings.Repeat(")", 100)},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fuzzer := NewLexerFuzzer()
			fuzzer.runFuzzTest(t, tc.input, tc.name)
		})
	}
}