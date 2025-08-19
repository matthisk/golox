package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// isASCIIDigit checks if a rune is an ASCII digit (0-9)
// This is more restrictive than unicode.IsDigit which includes Unicode digits
func isASCIIDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Optimization ideas are marked with comments below.

type Lexer struct {
	r       *bufio.Reader
	curRune rune
	curPos  int
	err     error
	line    int
	column  int
	// Accumulates the source code
	source strings.Builder
	// Optimization: Buffer for identifier/number/string building to avoid string concatenation
	sb strings.Builder
}

type Lex struct {
	Source string
	Tokens []Token
	Err    error
}

func New(r io.Reader) *Lexer {
	lx := &Lexer{
		r:      bufio.NewReader(r),
		curPos: -1, // start at -1 so the first read positions us at 0.
		line:   1,
		column: 0, // start at 0 so the first read positions this at 1.
	}
	return lx
}

func (lx *Lexer) Next() Token {
	if lx.err != nil {
		return lx.Return(EOF)
	}

	lx.read()

	// Skip whitespace first
	lx.skipWhitespace()

	// Record start position for the token
	startPos := Pos{
		Offset: lx.curPos,
		Line:   lx.line,
		Column: lx.column,
	}

	switch lx.curRune {
	case '(': // LEFT_PAREN
		return lx.makeToken(LEFT_PAREN, startPos)
	case ')': // RIGHT_PAREN
		return lx.makeToken(RIGHT_PAREN, startPos)
	case '{': // LEFT_BRACE
		return lx.makeToken(LEFT_BRACE, startPos)
	case '}': // RIGHT_BRACE
		return lx.makeToken(RIGHT_BRACE, startPos)
	case ',': // COMMA
		return lx.makeToken(COMMA, startPos)
	case '.': // DOT
		return lx.makeToken(DOT, startPos)
	case '-': // MINUS
		return lx.makeToken(MINUS, startPos)
	case '+': // PLUS
		return lx.makeToken(PLUS, startPos)
	case ';': // SEMICOLON
		return lx.makeToken(SEMICOLON, startPos)
	case ':': // COLON
		return lx.makeToken(COLON, startPos)
	case '?': // QUESTION_MARK
		return lx.makeToken(QUESTION_MARK, startPos)
	case '/': // SLASH
		if lx.match('/') {
			for lx.peek() != '\n' && !lx.atEnd() {
				lx.read()
			}
			return lx.Next()
		} else {
			return lx.makeToken(SLASH, startPos)
		}
	case '*': // STAR
		return lx.makeToken(STAR, startPos)
	case '!':
		// BANG (handle "!=" after the switch)
		if lx.match('=') {
			return lx.makeToken(BANG_EQUAL, startPos)
		}
		return lx.makeToken(BANG, startPos)
	case '=':
		// EQUAL (handle "==" after the switch)
		if lx.match('=') {
			return lx.makeToken(EQUAL_EQUAL, startPos)
		}
		return lx.makeToken(EQUAL, startPos)
	case '>':
		// GREATER (handle ">=" after the switch)
		if lx.match('=') {
			lx.read()
			return lx.makeToken(GREATER_EQUAL, startPos)
		}
		return lx.makeToken(GREATER, startPos)
	case '<':
		// LESS (handle "<=" after the switch)
		if lx.match('=') {
			return lx.makeToken(LESS_EQUAL, startPos)
		}
		return lx.makeToken(LESS, startPos)
	case '&':
		if lx.match('&') {
			return lx.makeToken(AND, startPos)
		}
	case '|':
		if lx.match('|') {
			return lx.makeToken(OR, startPos)
		}
	case 0: // rune(0) when you hit EOF
		return lx.makeToken(EOF, startPos)
	}

	if lx.curRune == '"' {
		return lx.string(startPos)
	}

	if isASCIIDigit(lx.curRune) {
		return lx.number(startPos)
	}

	if unicode.IsLetter(lx.curRune) {
		// Optimization: Use a map for keywords for O(1) lookup
		keywords := map[string]TokenType{
			"and":      AND,
			"class":    CLASS,
			"else":     ELSE,
			"false":    FALSE,
			"fun":      FUN,
			"for":      FOR,
			"if":       IF,
			"nil":      NIL,
			"or":       OR,
			"print":    PRINT,
			"return":   RETURN,
			"super":    SUPER,
			"this":     THIS,
			"true":     TRUE,
			"var":      VAR,
			"while":    WHILE,
			"break":    BREAK,
			"continue": CONTINUE,
		}

		lx.sb.Reset()
		lx.sb.WriteRune(lx.curRune)
		// Read all characters of the identifier/keyword
		for unicode.IsLetter(lx.peek()) || unicode.IsDigit(lx.peek()) || lx.peek() == '_' {
			lx.read()
			lx.sb.WriteRune(lx.curRune)
		}
		lexeme := lx.sb.String()

		if tokType, ok := keywords[lexeme]; ok {
			return lx.makeToken(tokType, startPos)
		}

		return Token{
			Type:     IDENTIFIER,
			Lexeme:   lexeme,
			StartPos: startPos,
			EndPos: Pos{
				Offset: lx.curPos + 1,
				Line:   lx.line,
				Column: lx.column + 1,
			},
		}
	}

	lx.read() // Advance past the illegal character
	return lx.makeToken(ILLEGAL, startPos)
}

func (lx *Lexer) number(startPos Pos) Token {
	lx.sb.Reset()

	lx.sb.WriteRune(lx.curRune)

	for isASCIIDigit(lx.peek()) {
		lx.read()
		lx.sb.WriteRune(lx.curRune)
	}

	r1, r2 := lx.peek2()
	if r1 == '.' && isASCIIDigit(r2) {
		lx.read()
		lx.sb.WriteRune('.')

		for isASCIIDigit(lx.peek()) {
			lx.read()
			lx.sb.WriteRune(lx.curRune)
		}
	}

	lexeme := lx.sb.String()
	float, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		// This should never happen if our lexer rules are correct.
		panic(fmt.Sprintf("Unparsable float %s", lexeme))
	}

	return Token{
		Type:     NUMBER,
		Lexeme:   float,
		StartPos: startPos,
		EndPos: Pos{
			Offset: lx.curPos,
			Line:   lx.line,
			Column: lx.column,
		},
	}
}

func (lx *Lexer) string(startPos Pos) Token {
	lx.sb.Reset()
	for lx.peek() != '"' && !lx.atEnd() {
		lx.read()
		if lx.curRune == '\n' {
			lx.line++
		}
		lx.sb.WriteRune(lx.curRune)
	}

	if lx.atEnd() {
		lx.err = lx.Err("Unterminated string.")
		return lx.makeToken(EOF, startPos)
	}

	// Advance past the closing quote
	lx.read()

	return Token{
		Type:     STRING,
		Lexeme:   lx.sb.String(),
		StartPos: startPos,
		EndPos: Pos{
			Offset: lx.curPos + 1,
			Line:   lx.line,
			Column: lx.column + 1,
		},
	}
}

func (lx *Lexer) Err(m string) error {
	return errors.New(m)
}

func (lx *Lexer) skipWhitespace() {
	for {
		if !unicode.IsSpace(lx.curRune) {
			break
		}
		lx.read()
	}
}

func (lx *Lexer) Return(token TokenType) Token {
	return Token{
		Type:   token,
		Lexeme: "",
		StartPos: Pos{
			Offset: lx.curPos,
			Line:   lx.line,
			Column: lx.column,
		},
		EndPos: Pos{
			Offset: lx.curPos,
			Line:   lx.line,
			Column: lx.column,
		},
	}
}

// Source returns the source string which the lexer processed.
func (lx *Lexer) Source() string {
	return lx.source.String()
}

func (lx *Lexer) makeToken(token TokenType, startPos Pos) Token {
	return Token{
		Type:     token,
		Lexeme:   "",
		StartPos: startPos,
		EndPos: Pos{
			Offset: lx.curPos + 1,
			Line:   lx.line,
			Column: lx.column + 1,
		},
	}
}

func (lx *Lexer) read() {
	// Update position based on current character before reading the next one
	if lx.curRune == '\n' {
		lx.line++
		lx.column = 1
	} else {
		lx.column++
	}
	lx.curPos++

	// Now read the next rune
	readRune, _, err := lx.r.ReadRune()
	lx.curRune = readRune
	lx.source.WriteRune(readRune)

	if err == io.EOF {
		lx.curRune = 0
		return
	}
	if err != nil {
		lx.err = lx.Err("Unexpected character.")
		return
	}
}

func (lx *Lexer) peek() rune {
	readRune, err := peekRune(lx.r)

	if err == io.EOF {
		lx.curRune = 0
		return 0
	}
	if err != nil {
		lx.err = lx.Err("Unexpected character.")
		return 0
	}

	return readRune
}

func (lx *Lexer) peek2() (rune, rune) {
	r1, r2, err := peek2Rune(lx.r)

	if err == io.EOF {
		lx.curRune = 0
		return 0, 0
	}
	if err != nil {
		lx.err = lx.Err("Unexpected character.")
		return 0, 0
	}

	return r1, r2
}

// match compares rune to the next rune, consumes rune from buffer in case of match
func (lx *Lexer) match(r rune) bool {
	peekRune, err := peekRune(lx.r)
	if err != nil {
		return false
	}

	if peekRune == r {
		lx.read()
		return true
	}

	return false
}

func (lx *Lexer) atEnd() bool {
	return lx.curRune == 0
}

func peekRune(r *bufio.Reader) (rune, error) {
	for peekBytes := 4; peekBytes > 0; peekBytes-- { // unicode rune can be up to 4 bytes
		b, err := r.Peek(peekBytes)
		if err == nil {
			rune, _ := utf8.DecodeRune(b)
			if rune == utf8.RuneError {
				return rune, fmt.Errorf("Rune error")
			}
			// success
			return rune, nil
		}
		// Otherwise, we ignore Peek errors and try the next smallest number of bytes
	}

	// Pretty sure we can assume EOF if we get this far
	return -1, io.EOF
}

func peek2Rune(r *bufio.Reader) (rune1, rune2 rune, err error) {
	// We request to peek up to 8 bytes. This is the maximum possible size for two
	// UTF-8 encoded runes, as each can be up to 4 bytes long.
	peekedBytes, err := r.Peek(8)

	// The behavior of bufio.Reader.Peek is to return an error (often io.EOF) if
	// it cannot read the full number of requested bytes. However, it still returns
	// the bytes it *was* able to read before the error occurred. We must handle this.
	if err != nil {
		if err == io.EOF && len(peekedBytes) == 0 {
			return 0, 0, io.EOF
		}
		if err != io.EOF {
			return 0, 0, err
		}
	}

	// Double-check if we have any bytes to process. This handles the case where
	// Peek returns a nil error but also zero bytes (unlikely but possible).
	if len(peekedBytes) == 0 {
		return 0, 0, io.EOF
	}

	// Decode the first rune from the beginning of our peeked byte slice.
	// utf8.DecodeRune returns the rune and its byte size.
	r1, size1 := utf8.DecodeRune(peekedBytes)
	if r1 == utf8.RuneError {
		// If the bytes do not form a valid UTF-8 sequence, DecodeRune returns
		// RuneError and a size of 1. We will treat it as a single-byte character
		// to ensure the logic proceeds.
		size1 = 1
	}

	// Check if there are any bytes remaining in our peeked slice after accounting
	// for the first rune's size.
	if len(peekedBytes) > size1 {
		// If yes, there's at least part of a second rune available. Decode it.
		// The second return value (size2) is ignored as we don't need it.
		r2, _ := utf8.DecodeRune(peekedBytes[size1:])
		return r1, r2, nil
	}

	// If we've reached this point, it means we successfully peeked exactly one rune
	// at the very end of the input. As per the function's contract, we return
	// the first rune, a zero value for the second, and no error.
	return r1, 0, nil
}

func Consume(lx *Lexer) Lex {
	// Optimization: Preallocate slice with a reasonable capacity
	tokens := make([]Token, 0, 128)

	token := lx.Next()
	for token.Type != EOF {
		tokens = append(tokens, token)
		token = lx.Next()
	}

	// Add the EOF token to the end
	tokens = append(tokens, token)

	return Lex{
		Source: lx.Source(),
		Tokens: tokens,
		Err:    lx.err,
	}
}
