package parser

import (
	"testing"

	"github.com/matthisk/lox/lexer"
)

func TestPrint(t *testing.T) {
	result, err := Print(&Binary{
		Left: &Unary{
			Token: lexer.BANG,
			Expr: &Literal{
				Token: lexer.Token{
					Type:   lexer.NUMBER,
					Lexeme: "5",
				},
			},
		},
		Token: lexer.PLUS,
		Right: &Literal{
			BaseNode: BaseNode{},
			Token: lexer.Token{
				Type:   lexer.NUMBER,
				Lexeme: "10.05",
			},
		},
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "(PLUS (BANG 5) 10.05)"
	if result != expected {
		t.Errorf("Expected %s got %s", expected, result)
	}
}
