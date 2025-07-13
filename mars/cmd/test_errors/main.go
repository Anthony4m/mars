// cmd/test_errors/main.go
package main

import (
	"fmt"
	"mars/lexer"
	"mars/parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// Test with a simple number literal
		input := "5;"
		l := lexer.New(input)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		fmt.Printf("Input: %s\n", input)
		fmt.Printf("Program: %+v\n", program)

		if len(program.Declarations) > 0 {
			fmt.Printf("First declaration: %T - %+v\n", program.Declarations[0], program.Declarations[0])
		}

		errors := p.GetErrors()
		if errors != nil && errors.HasErrors() {
			fmt.Printf("Errors: %v\n", errors)
		}
		return
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.NewParser(l)
	program := p.ParseProgram()

	fmt.Printf("Program: %+v\n", program)

	errors := p.GetErrors()
	if errors != nil && errors.HasErrors() {
		fmt.Printf("Errors: %v\n", errors)
	}
}
