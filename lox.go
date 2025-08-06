package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/matthisk/lox/lexer"
	"github.com/matthisk/lox/parser"
)

func main() {
	if len(os.Args) == 1 {
		// No arguments - start REPL
		runREPL()
	} else if len(os.Args) == 2 {
		// One argument - run file
		runFile(os.Args[1])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [filename]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  No arguments: Start interactive REPL\n")
		fmt.Fprintf(os.Stderr, "  filename: Execute Lox file\n")
		os.Exit(1)
	}
}

func runFile(filename string) {
	fmt.Fprintf(os.Stdout, "Running %s\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	lx := lexer.New(file)
	res := lexer.Consume(lx)
	if res.Err != nil {
		fmt.Fprintf(os.Stderr, "Lexer error: %v\n", res.Err)
		os.Exit(1)
	}

	ps := parser.New(res.Tokens)
	stmts, err := ps.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parser error: %v\n", err)
		os.Exit(1)
	}

	interpreter := parser.NewInterpreter()
	err = interpreter.Run(stmts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}

func runREPL() {
	fmt.Println("Welcome to the Lox REPL!")
	fmt.Println("Type Lox expressions or statements. Press Ctrl+C to exit.")
	fmt.Println()

	interpreter := parser.NewInterpreter()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("lox> ")

		if !scanner.Scan() {
			// EOF or Ctrl+C
			fmt.Println("\nGoodbye!")
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Special commands
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Try to evaluate the input
		evaluateREPLInput(line, interpreter)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func evaluateREPLInput(source string, interpreter *parser.Interpreter) {
	// Lexical analysis
	lx := lexer.New(bytes.NewBufferString(source))
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		fmt.Printf("Lexer error: %v\n", lexResult.Err)
		return
	}

	// Try parsing as statements first
	p := parser.New(lexResult.Tokens)
	stmts, err := p.Parse()
	if err == nil && len(stmts) > 0 {
		// Successfully parsed as statements - execute them
		err = interpreter.Run(stmts)
		if err != nil {
			fmt.Printf("Runtime error: %v\n", err)
		}
		return
	}

	tokens := lexResult.Tokens
	p = parser.New(tokens)
	expr, err := p.Expression()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Evaluate the expression and print the result
	result, err := interpreter.EvaluateExpression(expr)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		return
	}

	// Print the result (unless it's nil from a statement)
	if result != nil {
		fmt.Printf("%v\n", formatValue(result))
	}
}

func formatValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		// Check if it's a whole number for cleaner display
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
