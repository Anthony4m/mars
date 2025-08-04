package main

import (
	"fmt"
	"mars/evaluator"
	"mars/lexer"
	"mars/parser"
)

func main() {
	fmt.Println("=== Mars Built-in Functions Test ===")

	// Test cases for built-in functions
	testCases := []string{
		`len("hello world")`,
		`len([1, 2, 3, 4, 5])`,
		`append([1, 2, 3], 4)`,
		`print("Hello from print()")`,
		`println("Hello from println()")`,
		`printf("Value: %s", "test")`,
		`sin(0)`,
		`cos(0)`,
		`sqrt(16)`,
		`now()`,
	}

	for _, testCase := range testCases {
		fmt.Printf("\n--- Testing: %s ---\n", testCase)

		// Lexical analysis
		l := lexer.New(testCase)

		// Parsing
		p := parser.NewParser(l)
		ast := p.ParseProgram()

		if ast == nil {
			fmt.Println("❌ Parser failed: no AST")
			continue
		}

		// Evaluation
		eval := evaluator.New()
		result := eval.Eval(ast)

		if result == nil {
			fmt.Println("❌ Evaluator failed: no result")
			continue
		}

		fmt.Printf("✅ Result: %s\n", result.String())
	}

	fmt.Println("\n=== Built-in Functions Test Complete ===")
}
