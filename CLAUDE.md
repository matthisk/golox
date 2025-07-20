# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build
```bash
go build -o lox .
```

### Test
```bash
go test ./...
```

### Run
```bash
./lox <filename.lox>
```

## Architecture

This is a Lox language interpreter implementation in Go, following the structure from "Crafting Interpreters". The codebase is organized into three main packages:

### Core Components

**lexer/** - Tokenizes Lox source code
- `lexer.go`: Main lexer with rune-based scanning
- `token.go`: Token type definitions
- Uses buffered reader with string builder optimization
- Handles keywords, operators, strings, numbers, and identifiers

**parser/** - Recursive descent parser generating AST
- `ast.go`: AST node definitions using visitor pattern
- `parser.go`: Recursive descent parser following grammar precedence
- `interpreter.go`: Tree-walking interpreter implementing visitor interface
- `printer.go`: AST printer for debugging
- Grammar supports expressions with ternary, comma, binary operators, unary operators, and groupings

**Main Entry** - `lox.go` coordinates lexing, parsing, and interpretation

### Data Flow
1. Lexer consumes source file and produces tokens
2. Parser consumes tokens and builds expression AST
3. Interpreter walks AST using visitor pattern and evaluates expressions

### Testing Strategy
- Unit tests for each component in corresponding `_test.go` files
- Test data in `lexer/testdata/` with `.lox` source and `.lexed` expected output
- End-to-end tests combining lexer, parser, and interpreter

Note: Current test suite has parser boundary check issues that need fixing.