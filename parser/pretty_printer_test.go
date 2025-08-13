package parser

import (
	"strings"
	"testing"

	"github.com/matthisk/lox/lexer"
)

func TestPrettyPrinter(t *testing.T) {
	unformattedLox := `var a=1+2*3;print a;if(a>5){print "big";}else{print "small";}fun factorial(n){if(n<=1){return 1;}else{return n*factorial(n-1);}}var result=factorial(5);print result;`

	expectedFormatted := `var a = 1 + 2 * 3;
print a;
if (a > 5) {
  print "big";
} else {
  print "small";
}
fun factorial(n) {
  if (n <= 1) {
    return 1;
  } else {
    return n * factorial(n - 1);
  }
}
var result = factorial(5);
print result;`

	// Lex the unformatted code
	l := lexer.New(strings.NewReader(unformattedLox))
	lexResult := lexer.Consume(l)
	if lexResult.Err != nil {
		t.Fatalf("Failed to lex input: %v", lexResult.Err)
	}

	// Parse the tokens into AST
	p := New(lexResult.Tokens)
	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse input: %v", err)
	}

	// Pretty print the AST
	formatted, err := PrettyPrint(stmts)
	if err != nil {
		t.Fatalf("Failed to pretty print: %v", err)
	}

	// Compare the result
	if strings.TrimSpace(formatted) != strings.TrimSpace(expectedFormatted) {
		t.Errorf("Pretty printer output doesn't match expected.\nGot:\n%s\n\nExpected:\n%s", formatted, expectedFormatted)
	}
}

func TestPrettyPrinterComplexNesting(t *testing.T) {
	complexLox := `fun makeCounter(){var count=0;fun increment(){count=count+1;return count;}return increment;}var counter=makeCounter();for(var i=0;i<3;i=i+1){print counter();}`

	expectedFormatted := `fun makeCounter() {
  var count = 0;
  fun increment() {
    count = count + 1;
    return count;
  }
  return increment;
}
var counter = makeCounter();
{
  var i = 0;
  while (i < 3) {
    {
      print counter();
    }
    i = i + 1;
  }
}`

	// Lex the unformatted code
	l := lexer.New(strings.NewReader(complexLox))
	lexResult := lexer.Consume(l)
	if lexResult.Err != nil {
		t.Fatalf("Failed to lex input: %v", lexResult.Err)
	}

	// Parse the tokens into AST
	p := New(lexResult.Tokens)
	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse input: %v", err)
	}

	// Pretty print the AST
	formatted, err := PrettyPrint(stmts)
	if err != nil {
		t.Fatalf("Failed to pretty print: %v", err)
	}

	// Compare the result
	if strings.TrimSpace(formatted) != strings.TrimSpace(expectedFormatted) {
		t.Errorf("Pretty printer output doesn't match expected.\nGot:\n%s\n\nExpected:\n%s", formatted, expectedFormatted)
	}
}