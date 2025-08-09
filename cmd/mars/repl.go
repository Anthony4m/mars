package main

import (
	"bufio"
	"fmt"
	"mars/evaluator"
	"mars/lexer"
	"mars/parser"
	"os"
	"strings"
)

type REPL struct {
	evaluator *evaluator.Evaluator
	history   []string
	lineNum   int
}

func runREPL() {
	fmt.Println("Mars Programming Language REPL")
	fmt.Printf("Version %s\n", version)
	fmt.Println("Type 'exit' or 'quit' to exit, 'help' for help")
	fmt.Println()

	repl := &REPL{
		evaluator: evaluator.New(),
		history:   make([]string, 0),
		lineNum:   1,
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("mars> ")

	var currentInput strings.Builder
	parenCount := 0
	braceCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Handle special commands
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			break
		}
		if line == "help" {
			printREPLHelp()
			fmt.Print("mars> ")
			continue
		}
		if line == "clear" {
			fmt.Print("\033[H\033[2J") // Clear screen
			fmt.Print("mars> ")
			continue
		}
		if line == "history" {
			repl.printHistory()
			fmt.Print("mars> ")
			continue
		}
		if strings.HasPrefix(line, "!") {
			repl.executeHistoryCommand(line)
			fmt.Print("mars> ")
			continue
		}

		// Count brackets to determine if we need more input
		parenCount += strings.Count(line, "(") - strings.Count(line, ")")
		braceCount += strings.Count(line, "{") - strings.Count(line, "}")

		if currentInput.Len() > 0 {
			currentInput.WriteString(" ")
		}
		currentInput.WriteString(line)

		// If we have balanced brackets and the line ends with semicolon or is complete
		if parenCount == 0 && braceCount == 0 && (strings.HasSuffix(line, ";") || isCompleteStatement(line)) {
			input := currentInput.String()
			if strings.TrimSpace(input) != "" {
				repl.evaluateInput(input)
			}
			currentInput.Reset()
			parenCount = 0
			braceCount = 0
		} else {
			// Continue reading for multi-line input
			fmt.Print("  ")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

func (r *REPL) evaluateInput(input string) {
	// Add to history
	r.history = append(r.history, input)
	if len(r.history) > 100 {
		r.history = r.history[1:] // Keep only last 100 commands
	}

	// Lexical analysis
	l := lexer.New(input)

	// Parsing
	p := parser.NewParser(l)
	program := p.ParseProgram()

	// Check for parser errors
	errors := p.GetErrors()
	if errors != nil && errors.HasErrors() {
		fmt.Println("Parse errors:")
		for _, err := range errors.Errors() {
			fmt.Printf("  %s\n", err)
		}
		return
	}

	// Evaluation
	result := r.evaluator.Eval(program)

	// Print result (but not for void functions or statements)
	if result != nil && result.Type() != "NULL" {
		fmt.Printf("= %s\n", result.String())
	}
}

func (r *REPL) printHistory() {
	fmt.Println("Command History:")
	for i, cmd := range r.history {
		fmt.Printf("%3d: %s\n", i+1, cmd)
	}
}

func (r *REPL) executeHistoryCommand(cmd string) {
	if len(cmd) < 2 {
		fmt.Println("Usage: !<number> to execute command from history")
		return
	}

	// Parse the number
	var num int
	_, err := fmt.Sscanf(cmd[1:], "%d", &num)
	if err != nil {
		fmt.Printf("Invalid history number: %s\n", cmd[1:])
		return
	}

	if num < 1 || num > len(r.history) {
		fmt.Printf("History number %d out of range (1-%d)\n", num, len(r.history))
		return
	}

	// Execute the command from history
	historyCmd := r.history[num-1]
	fmt.Printf("Executing: %s\n", historyCmd)
	r.evaluateInput(historyCmd)
}

func isCompleteStatement(line string) bool {
	// Check if the line looks like a complete statement
	trimmed := strings.TrimSpace(line)

	// Empty line
	if trimmed == "" {
		return true
	}

	// Ends with semicolon
	if strings.HasSuffix(trimmed, ";") {
		return true
	}

	// Function declaration
	if strings.HasPrefix(trimmed, "func ") && strings.Contains(trimmed, "{") {
		return true
	}

	// Struct declaration
	if strings.HasPrefix(trimmed, "struct ") && strings.Contains(trimmed, "{") {
		return true
	}

	// Variable declaration with semicolon
	if strings.Contains(trimmed, ":=") && strings.HasSuffix(trimmed, ";") {
		return true
	}

	// Expression that doesn't need semicolon (like function calls)
	if strings.Contains(trimmed, "(") && strings.HasSuffix(trimmed, ")") {
		return true
	}

	return false
}

func printREPLHelp() {
	fmt.Println("Mars REPL Commands:")
	fmt.Println("  exit, quit          Exit the REPL")
	fmt.Println("  help                Show this help")
	fmt.Println("  clear               Clear the screen")
	fmt.Println("  history             Show command history")
	fmt.Println("  !<number>           Execute command from history")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  x := 42;")
	fmt.Println("  func add(a: int, b: int) -> int { return a + b; }")
	fmt.Println("  log(add(5, 3));")
	fmt.Println("  !1                  Execute the first command in history")
}
