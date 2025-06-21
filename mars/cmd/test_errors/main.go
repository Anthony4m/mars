// cmd/test_errors/main.go
package main

import (
	"fmt"
	"mars/errors"
	"mars/lexer"
	"mars/parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test_errors <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create lexer and parser
	l := lexer.New(string(content))
	p := parser.NewParser(l)

	// Parse the program
	program := p.ParseProgram()

	// Get errors
	errorList := p.GetErrors()

	// Print the program structure (for debugging)
	fmt.Printf("Parsed program with %d declarations\n", len(program.Declarations))

	// Print errors if any
	if errorList.HasErrors() {
		fmt.Printf("\n=== Compilation Errors ===\n")
		fmt.Println(errorList.String())
		fmt.Printf("\nTotal errors: %d\n", len(errorList.Errors()))
	} else {
		fmt.Println("No errors found!")
	}

	// Print warnings if any
	if errorList.HasWarnings() {
		fmt.Printf("\n=== Warnings ===\n")
		for _, err := range errorList.Errors() {
			if err.Severity == errors.ErrorSeverityWarning {
				fmt.Println(err.String())
			}
		}
	}

	// Exit with error code if there are errors
	if errorList.HasErrors() {
		os.Exit(1)
	}
}
