# Lox Interpreter

A toy implementation of the Lox programming language from Robert Nystrom's excellent book ["Crafting Interpreters"](https://craftinginterpreters.com/), written in Go.

## About

This project implements a tree-walking interpreter for the Lox language, a dynamically-typed scripting language with features like:

- Variables and basic data types (numbers, strings, booleans, nil)
- Arithmetic and logical expressions
- Control flow (if/else, while, for loops)
- Functions with closures
- Classes with inheritance
- Object-oriented programming with methods and fields

## Building

Build the interpreter binary:

```bash
go build -o lox .
```

## Usage

Run a Lox program from a file:

```bash
./lox <filename.lox>
```

Example:

```bash
./lox engine/testdata/fibonacci.lox
```

## Testing

Run the complete test suite:

```bash
go test ./...
```

The test suite includes:
- Unit tests for the lexer, parser, and interpreter components
- End-to-end integration tests with sample Lox programs
- Test programs in `engine/testdata/` covering various language features

## Architecture

The interpreter is organized into three main components:

- **Lexer** (`lexer/`) - Tokenizes Lox source code into a stream of tokens
- **Parser** (`parser/`) - Builds an Abstract Syntax Tree (AST) from tokens using recursive descent parsing
- **Interpreter** (`interpreter/`) - Evaluates the AST using the visitor pattern

## Language Features

This implementation supports the full Lox language specification including:

- Expression evaluation with proper operator precedence
- Variable declaration and assignment
- Functions with parameters and return values
- Closures and lexical scoping
- Classes with methods and constructors
- Inheritance with `super` keyword support
- Built-in functions like `print` and `clock`

## Example Programs

See the `engine/testdata/` directory for example Lox programs demonstrating various language features.