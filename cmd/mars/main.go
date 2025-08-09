package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "repl":
		runREPL()
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: 'run' command requires a file path")
			fmt.Println("Usage: mars run <file.mars>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "fmt":
		if len(os.Args) < 3 {
			fmt.Println("Error: 'fmt' command requires a file path")
			fmt.Println("Usage: mars fmt <file.mars>")
			os.Exit(1)
		}
		formatFile(os.Args[2])
	case "test":
		runTests()
	case "version", "-v", "--version":
		fmt.Printf("Mars Programming Language v%s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Mars Programming Language")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  mars repl                    Start interactive REPL")
	fmt.Println("  mars run <file.mars>         Parse and evaluate a file")
	fmt.Println("  mars fmt <file.mars>         Format a Mars file")
	fmt.Println("  mars test                    Run tests in tests/ directory")
	fmt.Println("  mars version                 Show version information")
	fmt.Println("  mars help                    Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mars repl")
	fmt.Println("  mars run hello.mars")
	fmt.Println("  mars fmt program.mars")
	fmt.Println("  mars test")
}
