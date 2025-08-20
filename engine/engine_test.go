package engine

import (
	"os"
	"testing"
)

func TestEngine(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{name: "function", filename: "testdata/function.lox"},
		{name: "classes", filename: "testdata/classes.lox"},
		{name: "inheritance (1)", filename: "testdata/inheritance.lox"},
		{name: "inheritance (2)", filename: "testdata/inheritance_2.lox"},
		{name: "inheritance (3)", filename: "testdata/inheritance_3.lox"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.OpenFile(tt.filename, os.O_RDONLY, 0666)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			defer file.Close()

			err = Run(file, nil)
			if err != nil {
				t.Fatalf("Failed to run test file: %v", err)
			}
		})
	}
}
