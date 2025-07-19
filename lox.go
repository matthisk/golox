package main

import (
	"fmt"
	"os"

	"github.com/matthisk/lox/lexer"
	"github.com/matthisk/lox/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	fmt.Fprintf(os.Stdout, "Running on %s\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	lx := lexer.New(file)
	res := lexer.Consume(lx)
	if res.Err != nil {
		panic(res.Err)
	}
	ps := parser.New(res.Tokens)
	expr, err := ps.Parse()
	if err != nil {
		panic(err)
	}
	interpreter := parser.Interpreter{}
	val, err := interpreter.Run(expr)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "%v\n", val)
}
