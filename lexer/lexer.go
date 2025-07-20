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

// Optimization ideas are marked with comments below.

type Lexer struct {
	r       *bufio.Reader
	curRune rune
	curPos  int
	err     error
	line    int
	// Optimization: Buffer for identifier/number/string building to avoid string concatenation
	sb strings.Builder
}

type Lex struct {
	Tokens []Token
	Err    error
}

func New(r io.Reader) *Lexer {
	return &Lexer{
		r:      bufio.NewReader(r),
		curPos: 0,
		line:   0,
	}
}

func (lx *Lexer) Next() Token {
	if lx.err != nil {
		return lx.Return(EOF)
	}

	// Skip whitespace first
	lx.skipWhitespace()

	switch lx.curRune {
	case '(': // LEFT_PAREN
		lx.read() // Advance past the token
		return lx.Return(LEFT_PAREN)
	case ')': // RIGHT_PAREN
		lx.read() // Advance past the token
		return lx.Return(RIGHT_PAREN)
	case '{': // LEFT_BRACE
		lx.read() // Advance past the token
		return lx.Return(LEFT_BRACE)
	case '}': // RIGHT_BRACE
		lx.read() // Advance past the token
		return lx.Return(RIGHT_BRACE)
	case ',': // COMMA
		lx.read() // Advance past the token
		return lx.Return(COMMA)
	case '.': // DOT
		lx.read() // Advance past the token
		return lx.Return(DOT)
	case '-': // MINUS
		lx.read() // Advance past the token
		return lx.Return(MINUS)
	case '+': // PLUS
		lx.read() // Advance past the token
		return lx.Return(PLUS)
	case ';': // SEMICOLON
		lx.read() // Advance past the token
		return lx.Return(SEMICOLON)
	case ':': // COLON
		lx.read() // Advance past the token
		return lx.Return(COLON)
	case '?': // QUESTION_MARK
		lx.read() // Advance past the token
		return lx.Return(QUESTION_MARK)
	case '/': // SLASH
		if lx.match('/') {
			for lx.curRune != '\n' {
				lx.read()
				if lx.curRune == '\n' {
					lx.line++
				}
				if lx.atEnd() {
					return lx.Return(EOF)
				}
			}
			return lx.Next()

		} else {
			lx.read() // Advance past the token
			return lx.Return(SLASH)
		}
	case '*': // STAR
		lx.read() // Advance past the token
		return lx.Return(STAR)
	case '!':
		// BANG (handle "!=" after the switch)
		if lx.match('=') {
			return lx.Return(BANG_EQUAL)
		}
		lx.read() // Advance past the token
		return lx.Return(BANG)
	case '=':
		// EQUAL (handle "==" after the switch)
		if lx.match('=') {
			return lx.Return(EQUAL_EQUAL)
		}
		lx.read() // Advance past the token
		return lx.Return(EQUAL)
	case '>':
		// GREATER (handle ">=" after the switch)
		if lx.match('=') {
			return lx.Return(GREATER_EQUAL)
		}
		lx.read() // Advance past the token
		return lx.Return(GREATER)
	case '<':
		// LESS (handle "<=" after the switch)
		if lx.match('=') {
			return lx.Return(LESS_EQUAL)
		}
		lx.read() // Advance past the token
		return lx.Return(LESS)
	case 0: // rune(0) when you hit EOF
		return lx.Return(EOF)
	}

	if lx.curRune == '"' {
		lx.read() // Advance past the opening quote
		return lx.string()
	}

	if unicode.IsDigit(lx.curRune) {
		return lx.number()
	}

	if unicode.IsLetter(lx.curRune) {
		// Optimization: Use a map for keywords for O(1) lookup
		keywords := map[string]TokenType{
			"and":    AND,
			"class":  CLASS,
			"else":   ELSE,
			"false":  FALSE,
			"fun":    FUN,
			"for":    FOR,
			"if":     IF,
			"nil":    NIL,
			"or":     OR,
			"print":  PRINT,
			"return": RETURN,
			"super":  SUPER,
			"this":   THIS,
			"true":   TRUE,
			"var":    VAR,
			"while":  WHILE,
		}

		lx.sb.Reset()
		// Read all characters of the identifier/keyword
		for ; unicode.IsLetter(lx.curRune) || unicode.IsDigit(lx.curRune) || lx.curRune == '_'; lx.read() {
			lx.sb.WriteRune(lx.curRune)
		}
		lexeme := lx.sb.String()

		if tokType, ok := keywords[lexeme]; ok {
			return lx.Return(tokType)
		}

		return Token{
			Type:     IDENTIFIER,
			Lexeme:   lexeme,
			Position: lx.curPos,
		}
	}

	lx.read() // Advance past the illegal character
	return lx.Return(ILLEGAL)
}

func (lx *Lexer) number() Token {
	lx.sb.Reset()
	for ; unicode.IsDigit(lx.curRune); lx.read() {
		lx.sb.WriteRune(lx.curRune)
	}

	if lx.curRune == '.' && unicode.IsDigit(lx.peek()) {
		lx.sb.WriteRune('.')
		lx.read()
		for ; unicode.IsDigit(lx.curRune); lx.read() {
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
		Position: lx.curPos,
	}
}

func (lx *Lexer) string() Token {
	lx.sb.Reset()
	for lx.curRune != '"' && !lx.atEnd() {
		if lx.curRune == '\n' {
			lx.line++
		}
		lx.sb.WriteRune(lx.curRune)
		lx.read()
	}

	if lx.atEnd() {
		lx.err = lx.Err("Unterminated string.")
		return lx.Return(EOF)
	}

	// Advance past the closing quote
	lx.read()

	return Token{
		Type:     STRING,
		Lexeme:   lx.sb.String(),
		Position: lx.curPos,
	}
}

func (lx *Lexer) Err(m string) error {
	return errors.New(m)
}

func (lx *Lexer) skipWhitespace() {
	// If curRune is 0, we need to read the first character
	if lx.curRune == 0 {
		lx.read()
	}

	for {
		if !unicode.IsSpace(lx.curRune) {
			break
		}
		lx.read()
		if lx.curRune == '\n' {
			lx.line++
		}
	}
}

func (lx *Lexer) Return(token TokenType) Token {
	return Token{
		Type:     token,
		Lexeme:   "",
		Position: lx.curPos,
	}
}

func (lx *Lexer) read() {
	readRune, _, err := lx.r.ReadRune()

	lx.curPos++
	lx.curRune = readRune

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

func Consume(lx *Lexer) Lex {
	// Optimization: Preallocate slice with a reasonable capacity
	tokens := make([]Token, 0, 128)

	token := lx.Next()
	for token.Type != EOF {
		tokens = append(tokens, token)
		token = lx.Next()
	}

	return Lex{
		Tokens: tokens,
		Err:    lx.err,
	}
}
