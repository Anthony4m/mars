package main

import (
	"fmt"
	"io"
	"io/fs"
	"mars/evaluator"
	"mars/lexer"
	"mars/parser"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TestResult struct {
	Name     string
	Passed   bool
	Expected string
	Actual   string
	Error    string
	Duration time.Duration
}

type TestFile struct {
	Path     string
	Content  string
	Expected string
}

func runTests() {
	fmt.Println("Running Mars tests...")
	fmt.Println()

	// Look for tests directory
	testDir := "tests"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		fmt.Printf("Error: Tests directory '%s' not found\n", testDir)
		fmt.Println("Create a 'tests' directory with .mars files to run tests")
		os.Exit(1)
	}

	// Find all .mars files in tests directory
	testFiles, err := findTestFiles(testDir)
	if err != nil {
		fmt.Printf("Error finding test files: %v\n", err)
		os.Exit(1)
	}

	if len(testFiles) == 0 {
		fmt.Printf("No .mars test files found in '%s' directory\n", testDir)
		fmt.Println("Create .mars files in the tests directory to run tests")
		os.Exit(1)
	}

	// Run tests
	var results []TestResult
	passed := 0
	failed := 0

	for _, testFile := range testFiles {
		result := runTest(testFile)
		results = append(results, result)

		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	// Print results
	printTestResults(results, passed, failed)

	// Exit with appropriate code
	if failed > 0 {
		os.Exit(1)
	}
}

func findTestFiles(testDir string) ([]TestFile, error) {
	var testFiles []TestFile

	err := filepath.WalkDir(testDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".mars") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			expected := extractExpectedOutput(string(content))
			testFiles = append(testFiles, TestFile{
				Path:     path,
				Content:  string(content),
				Expected: expected,
			})
		}

		return nil
	})

	return testFiles, err
}

func extractExpectedOutput(content string) string {
	// Look for EXPECTED: or EXPECT: comments
	lines := strings.Split(content, "\n")
	var expectedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "// EXPECTED:") || strings.HasPrefix(trimmed, "// EXPECT:") {
			expected := strings.TrimSpace(strings.TrimPrefix(trimmed, "// EXPECTED:"))
			expected = strings.TrimSpace(strings.TrimPrefix(expected, "// EXPECT:"))
			expectedLines = append(expectedLines, expected)
		}
	}

	return strings.Join(expectedLines, "\n")
}

func runTest(testFile TestFile) TestResult {
	start := time.Now()

	// Capture stdout
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return TestResult{
			Name:   testFile.Path,
			Passed: false,
			Error:  fmt.Sprintf("Failed to capture stdout: %v", err),
		}
	}
	os.Stdout = w

	// Run the test
	var actualOutput string
	var runtimeError string

	func() {
		defer func() {
			if r := recover(); r != nil {
				runtimeError = fmt.Sprintf("Panic: %v", r)
			}
		}()

		// Lexical analysis
		l := lexer.New(testFile.Content)

		// Parsing
		p := parser.NewParser(l)
		program := p.ParseProgram()

		// Check for parser errors
		errors := p.GetErrors()
		if errors != nil && errors.HasErrors() {
			runtimeError = "Parse errors:\n"
			for _, err := range errors.Errors() {
				runtimeError += fmt.Sprintf("  %s\n", err)
			}
			return
		}

		// Create evaluator
		eval := evaluator.New()

		// Evaluate the program
		result := eval.Eval(program)

		// Check for evaluation errors
		if result != nil && result.Type() == "ERROR" {
			runtimeError = fmt.Sprintf("Runtime error: %s", result.String())
			return
		}

		// Get final result if it's not null
		if result != nil && result.Type() != "NULL" {
			fmt.Printf("Result: %s\n", result.String())
		}
	}()

	// Restore stdout and get captured output
	w.Close()
	os.Stdout = originalStdout

	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	if err != nil {
		return TestResult{
			Name:   testFile.Path,
			Passed: false,
			Error:  fmt.Sprintf("Failed to read captured output: %v", err),
		}
	}
	actualOutput = strings.TrimSpace(buf.String())

	duration := time.Since(start)

	// Check for runtime errors
	if runtimeError != "" {
		return TestResult{
			Name:     testFile.Path,
			Passed:   false,
			Expected: testFile.Expected,
			Actual:   actualOutput,
			Error:    runtimeError,
			Duration: duration,
		}
	}

	// Compare output
	passed := actualOutput == testFile.Expected

	return TestResult{
		Name:     testFile.Path,
		Passed:   passed,
		Expected: testFile.Expected,
		Actual:   actualOutput,
		Duration: duration,
	}
}

func printTestResults(results []TestResult, passed, failed int) {
	fmt.Println("Test Results:")
	fmt.Println("=============")

	for _, result := range results {
		if result.Passed {
			fmt.Printf("✅ %s (%v)\n", result.Name, result.Duration)
		} else {
			fmt.Printf("❌ %s (%v)\n", result.Name, result.Duration)
			if result.Error != "" {
				fmt.Printf("   Error: %s\n", result.Error)
			}
			if result.Expected != "" {
				fmt.Printf("   Expected: %s\n", result.Expected)
			}
			if result.Actual != "" {
				fmt.Printf("   Actual:   %s\n", result.Actual)
			}
		}
	}

	fmt.Println()
	fmt.Printf("Summary: %d passed, %d failed\n", passed, failed)

	if failed > 0 {
		fmt.Println("\nSome tests failed. Check the output above for details.")
	} else {
		fmt.Println("\nAll tests passed! 🎉")
	}
}
