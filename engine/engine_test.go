package engine

import (
	"os"
	"testing"
)

func TestEngine(t *testing.T) {
	file, err := os.OpenFile("testdata/classes.lox", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	err = Run(file, nil)
	if err != nil {
		t.Fatalf("Failed to run test file: %v", err)
	}
}
