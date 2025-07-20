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
	column  int
	// Optimization: Buffer for identifier/number/string building to avoid string concatenation
	sb strings.Builder
}

type Lex struct {
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
	// Read the first rune from the buffer to start the lexing process.
	lx.read()
	return lx
}

func (lx *Lexer) Next() Token {
	if lx.err != nil {
		return lx.Return(EOF)
	}

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
		lx.read() // Advance past the token
		return lx.makeToken(LEFT_PAREN, startPos)
	case ')': // RIGHT_PAREN
		lx.read() // Advance past the token
		return lx.makeToken(RIGHT_PAREN, startPos)
	case '{': // LEFT_BRACE
		lx.read() // Advance past the token
		return lx.makeToken(LEFT_BRACE, startPos)
	case '}': // RIGHT_BRACE
		lx.read() // Advance past the token
		return lx.makeToken(RIGHT_BRACE, startPos)
	case ',': // COMMA
		lx.read() // Advance past the token
		return lx.makeToken(COMMA, startPos)
	case '.': // DOT
		lx.read() // Advance past the token
		return lx.makeToken(DOT, startPos)
	case '-': // MINUS
		lx.read() // Advance past the token
		return lx.makeToken(MINUS, startPos)
	case '+': // PLUS
		lx.read() // Advance past the token
		return lx.makeToken(PLUS, startPos)
	case ';': // SEMICOLON
		lx.read() // Advance past the token
		return lx.makeToken(SEMICOLON, startPos)
	case ':': // COLON
		lx.read() // Advance past the token
		return lx.makeToken(COLON, startPos)
	case '?': // QUESTION_MARK
		lx.read() // Advance past the token
		return lx.makeToken(QUESTION_MARK, startPos)
	case '/': // SLASH
		if lx.match('/') {
			for lx.curRune != '\n' {
				lx.read()
				if lx.atEnd() {
					return lx.makeToken(EOF, startPos)
				}
			}
			return lx.Next()

		} else {
			lx.read() // Advance past the token
			return lx.makeToken(SLASH, startPos)
		}
	case '*': // STAR
		lx.read() // Advance past the token
		return lx.makeToken(STAR, startPos)
	case '!':
		// BANG (handle "!=" after the switch)
		if lx.match('=') {
			return lx.makeToken(BANG_EQUAL, startPos)
		}
		lx.read() // Advance past the token
		return lx.makeToken(BANG, startPos)
	case '=':
		// EQUAL (handle "==" after the switch)
		if lx.match('=') {
			return lx.makeToken(EQUAL_EQUAL, startPos)
		}
		lx.read() // Advance past the token
		return lx.makeToken(EQUAL, startPos)
	case '>':
		// GREATER (handle ">=" after the switch)
		if lx.match('=') {
			return lx.makeToken(GREATER_EQUAL, startPos)
		}
		lx.read() // Advance past the token
		return lx.makeToken(GREATER, startPos)
	case '<':
		// LESS (handle "<=" after the switch)
		if lx.match('=') {
			return lx.makeToken(LESS_EQUAL, startPos)
		}
		lx.read() // Advance past the token
		return lx.makeToken(LESS, startPos)
	case 0: // rune(0) when you hit EOF
		return lx.makeToken(EOF, startPos)
	}

	if lx.curRune == '"' {
		lx.read() // Advance past the opening quote
		return lx.string(startPos)
	}

	if unicode.IsDigit(lx.curRune) {
		return lx.number(startPos)
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
			return lx.makeToken(tokType, startPos)
		}

		return Token{
			Type:     IDENTIFIER,
			Lexeme:   lexeme,
			StartPos: startPos,
			EndPos: Pos{
				Offset: lx.curPos,
				Line:   lx.line,
				Column: lx.column,
			},
		}
	}

	lx.read() // Advance past the illegal character
	return lx.makeToken(ILLEGAL, startPos)
}

func (lx *Lexer) number(startPos Pos) Token {
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
	for lx.curRune != '"' && !lx.atEnd() {
		if lx.curRune == '\n' {
			lx.line++
		}
		lx.sb.WriteRune(lx.curRune)
		lx.read()
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
			Offset: lx.curPos,
			Line:   lx.line,
			Column: lx.column,
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

func (lx *Lexer) makeToken(token TokenType, startPos Pos) Token {
	return Token{
		Type:     token,
		Lexeme:   "",
		StartPos: startPos,
		EndPos: Pos{
			Offset: lx.curPos,
			Line:   lx.line,
			Column: lx.column,
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

	// Add the EOF token to the end
	tokens = append(tokens, token)

	return Lex{
		Tokens: tokens,
		Err:    lx.err,
	}
}
