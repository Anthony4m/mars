package main

import (
	"fmt"
	"mars/evaluator"
	"mars/lexer"
	"mars/parser"
	"os"
	"path/filepath"
)

func runFile(filename string) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", filename)
		os.Exit(1)
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename)
		os.Exit(1)
	}

	// Check file extension
	if filepath.Ext(filename) != ".mars" {
		fmt.Printf("Warning: File '%s' doesn't have .mars extension\n", filename)
	}

	// Lexical analysis
	l := lexer.New(string(content))

	// Parsing
	p := parser.NewParser(l)
	program := p.ParseProgram()

	// Check for parser errors
	errors := p.GetErrors()
	if errors != nil && errors.HasErrors() {
		fmt.Printf("Parse errors in '%s':\n", filename)
		for _, err := range errors.Errors() {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	// Create evaluator
	eval := evaluator.New()

	// Evaluate the program
	result := eval.Eval(program)

	// Check for evaluation errors
	if result != nil && result.Type() == "ERROR" {
		fmt.Printf("Runtime error in '%s': %s\n", filename, result.String())
		os.Exit(1)
	}

	// Print final result if it's not null
	if result != nil && result.Type() != "NULL" {
		fmt.Printf("Result: %s\n", result.String())
	}
}
