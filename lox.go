package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/matthisk/lox/engine"
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

	err = engine.Run(file, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runREPL() {
	fmt.Println("Welcome to the Lox REPL!")
	fmt.Println("Type Lox expressions or statements. Press Ctrl+C to exit.")
	fmt.Println()

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
		evaluateREPLInput(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func evaluateREPLInput(source string) {
	input := bytes.NewBufferString(source)

	// Try parsing as statements first
	err := engine.Run(input, nil)
	if err == nil {
		return
	}

	// If that failed, try parsing as an expression
	input = bytes.NewBufferString(source)
	result, err := engine.EvaluateExpr(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
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
