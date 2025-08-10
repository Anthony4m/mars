package analyzer

import (
	"fmt"
	"mars/ast"
	"mars/lexer"
	"mars/parser"
	"strings"
	"testing"
)

func TestBasicVariableCollection(t *testing.T) {
	input := `
    x := 42
    mut name : string = "Mars"
    `

	// Parse to AST (you already have this)
	lexer := lexer.New(input)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()

	// Create analyzer
	analyzer := New(input, "test.mars")

	// Test declaration collection
	err := analyzer.collectDeclarations(program)
	if err != nil {
		t.Fatalf("Collection failed: %v", err)
	}

	// Verify symbols were collected
	xSymbol, err := analyzer.symbols.Resolve("x")
	if err != nil {
		t.Fatalf("Symbol 'x' not found: %v", err)
	}

	if xSymbol.Type.BaseType != "int" {
		t.Errorf("Expected x to be int, got %s", xSymbol.Type.BaseType)
	}

	nameSymbol, err := analyzer.symbols.Resolve("name")
	if err != nil {
		t.Fatalf("Symbol 'name' not found: %v", err)
	}

	if !nameSymbol.IsMutable {
		t.Errorf("Expected 'name' to be mutable")
	}
}

func TestFunctionDeclarationCollection(t *testing.T) {
	input := `
	func add(a : int, b : int) -> int {
		return a + b
	}
	`
	// Parse to AST (you already have this)
	lexer := lexer.New(input)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()

	// Create analyzer
	analyzer := New(input, "test.mars")

	// Test declaration collection
	err := analyzer.collectDeclarations(program)
	if err != nil {
		t.Fatalf("Collection failed: %v", err)
	}

	// Verify function was collected
	_, err = analyzer.symbols.Resolve("add")
	if err != nil {
		t.Fatalf("Symbol 'add' not found: %v", err)
	}
	if stmt, ok := program.Declarations[0].(*ast.FuncDecl); ok {
		if !ok {
			t.Fatalf("Expected first declaration to be a function declaration")
		}
		if stmt.Signature.Parameters[0].Type.BaseType != "int" {
			t.Errorf("Expected first parameter to be int, got %s", stmt.Signature.Parameters[0].Type.BaseType)
		}
		if stmt.Signature.Parameters[1].Type.BaseType != "int" {
			t.Errorf("Expected second parameter to be int, got %s", stmt.Signature.Parameters[1].Type.BaseType)
		}
		if stmt.Signature.ReturnType.BaseType != "int" {
			t.Errorf("Expected return type to be int, got %s", stmt.Signature.ReturnType.BaseType)
		}
		if len(stmt.Signature.Parameters) != 2 {
			t.Errorf("Expected 2 parameters, got %d", len(stmt.Signature.Parameters))
		}
	}

}

func TestStructDeclarationCollection(t *testing.T) {
	input := `
	struct Point {
		x : int
		y : int
	}
	`
	// Parse to AST (you already have this)
	lexer := lexer.New(input)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgram()

	// Create analyzer
	analyzer := New(input, "test.mars")

	// Test declaration collection
	err := analyzer.collectDeclarations(program)
	if err != nil {
		t.Fatalf("Collection failed: %v", err)
	}

	// Verify struct was collected
	structSymbol, err := analyzer.symbols.Resolve("Point")
	if err != nil {
		t.Fatalf("Symbol 'Point' not found: %v", err)
	}

	if structSymbol.Type.StructName != "Point" {
		t.Errorf("Expected Point to be a struct, got %s", structSymbol.Type.BaseType)
	}
	if len(structSymbol.Type.StructFields) != 2 {
		t.Errorf("Expected 2 fields in Point, got %d", len(structSymbol.Type.StructFields))
	}
}

func TestStructDeclaration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Test Valid struct",
			input:    "struct Point {\n\tx : int;}\n",
			expected: "Point",
			wantErr:  false,
		},
		{
			name:     "Test Invalid struct",
			input:    "struct Point {\n\t: int;}\n",
			expected: "",
			wantErr:  true,
		},
	}

	for _, item := range tests {
		l := lexer.New(item.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		analyzer := New(item.input, "test.mars")
		_ = analyzer.collectDeclarations(program)

		parserHasErr := len(p.GetErrors().Errors()) > 0
		analyzerHasErr := analyzer.errors.HasErrors()
		gotErr := parserHasErr || analyzerHasErr

		if item.wantErr {
			if !gotErr {
				t.Errorf("Expected error but got none")
			}
		} else {
			if gotErr {
				t.Errorf("Unexpected error: parser errors: %v, analyzer errors: %v", p.GetErrors().Errors(), analyzer.errors.String())
			}
		}
	}
}

// testAnalyze is a helper function that takes source code, runs the analyzer,
// and returns a string of all reported errors.
func testAnalyze(code string) string {
	l := lexer.New(code)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	if len(p.GetErrors().Errors()) > 0 {
		// For analyzer tests, we assume the parser is correct.
		// This helps isolate issues to the analyzer.
		// We can return a specific error to indicate a parser failure during testing.
		return fmt.Sprintf("parser error: %s", p.GetErrors().Error())
	}

	analyzer := New(code, "test.mars")
	err := analyzer.Analyze(program)
	if err != nil {
		return err.Error()
	}
	if analyzer.errors.HasErrors() {
		return analyzer.errors.String()
	}

	return ""
}

// assertErrorContains is a test helper to check for specific error messages.
func assertErrorContains(t *testing.T, errStr, substr string) {
	t.Helper()
	if errStr == "" {
		t.Errorf("expected error containing %q, but got no error", substr)
		return
	}
	if !strings.Contains(errStr, substr) {
		t.Errorf("expected error message containing %q, got %q", substr, errStr)
	}
}

// assertNoError is a test helper to check for the absence of errors.
func assertNoError(t *testing.T, errStr string) {
	t.Helper()
	if errStr != "" {
		t.Errorf("unexpected error: %s", errStr)
	}
}

func TestVariableDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		errorMsg string // Substring of the expected error
	}{
		{"valid declaration with type and value", "x: int = 10;", ""},
		{"valid declaration with type only", "x: int;", ""},
		{"valid type inference", "x := 10;", ""},
		{"valid mutable declaration", "mut x := 10;", ""},
		{"mismatched types", "var x: int = \"hello\";", "mismatched types: expected int, found string"},
		{"redeclaration in same scope", "var x: int; var x: string;", "variable 'x' is already defined"},
		{"declaration without type or initializer", "mut x;", "expected ':' or ':=' in variable declaration"},
		{"inference failure from undeclared var", "x := y;", "undefined"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := testAnalyze(tt.code)
			if tt.errorMsg != "" {
				assertErrorContains(t, errStr, tt.errorMsg)
			} else {
				assertNoError(t, errStr)
			}
		})
	}
}

func TestAssignments(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		errorMsg string
	}{
		{"valid assignment to mutable var", "mut x := 10; x = 20;", ""},
		{"assignment to immutable var", "x := 10; x = 20;", "cannot assign to immutable variable"},
		{"assignment to undeclared var", "x = 20;", "undefined symbol 'x'"},
		{"type mismatch in assignment", "mut x : int = 10; x = \"hello\";", "cannot assign string to int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := testAnalyze(tt.code)
			if tt.errorMsg != "" {
				assertErrorContains(t, errStr, tt.errorMsg)
			} else {
				assertNoError(t, errStr)
			}
		})
	}
}

func TestForInitMutability(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		errorMsg string
	}{
		{
			name:     "for-init without mut should error on post increment",
			code:     "func main(){ for i := 0; i < 2; i = i + 1 { } }",
			errorMsg: "cannot assign to immutable variable",
		},
		{
			name:     "for-init with mut should be ok",
			code:     "func main(){ for mut i := 0; i < 2; i = i + 1 { } }",
			errorMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := testAnalyze(tt.code)
			if tt.errorMsg != "" {
				assertErrorContains(t, errStr, tt.errorMsg)
			} else {
				assertNoError(t, errStr)
			}
		})
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		errorMsg string
	}{
		{"valid function call", "func add(a: int, b: int) -> int { return a + b; } add(1, 2);", ""},
		{"wrong number of arguments (too few)", "func foo(a: int) {} foo();", "wrong number of arguments"},
		{"wrong number of arguments (too many)", "func foo() {} foo(1);", "wrong number of arguments"},
		{"wrong argument type", "func foo(a: int) {} foo(\"hello\");", "cannot use 'string' as type 'int'"},
		{"calling a non-function", "x := 10; x();", "'x' is not a function"},
		{"calling an undefined function", "foo();", "undefined symbol 'foo'"},
		{"valid mutual recursion", "func a() { b(); } func b() { a(); } a();", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := testAnalyze(tt.code)
			if tt.errorMsg != "" {
				assertErrorContains(t, errStr, tt.errorMsg)
			} else {
				assertNoError(t, errStr)
			}
		})
	}
}

func TestStructLiterals(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		errorMsg string
	}{
		{"valid struct literal", "struct Point { x: int; y: int; } p: Point = Point{x: 1, y: 2};", ""},
		{"unknown struct type", "p: Point = Point{x: 1};", "unknown type \"Point\""},
		{"not a struct type", "x: int; y: x = x{};", "\"x\" is not a struct type"},
		{"non-existent field", "struct Point { x: int; } p: Point = Point{y: 1};", "field \"y\" does not exist on Point"},
		{"duplicate field in literal", "struct Point { x: int; } p: Point = Point{x: 1, x: 2};", "duplicate field \"x\" in literal"},
		{"mismatched field type", "struct Point { x: int; } p: Point = Point{x: \"hello\"};", "cannot use string to initialize field \"x\""},
		{"missing a field is valid", "struct Point { x: int; y: int; } p: Point = Point{x: 1};", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := testAnalyze(tt.code)
			if tt.errorMsg != "" {
				assertErrorContains(t, errStr, tt.errorMsg)
			} else {
				assertNoError(t, errStr)
			}
		})
	}
}

func TestBreakContinueValidation(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "break inside for loop - valid",
			code:        "mut i := 0; for ; i < 10; ; { i = i + 1; }",
			shouldError: false,
		},
		{
			name:        "continue inside for loop - valid",
			code:        "mut i := 0; for ; i < 10; ; { if i % 2 == 0 { continue; } i = i + 1; }",
			shouldError: false,
		},
		{
			name:        "break outside loop - invalid",
			code:        "break;",
			shouldError: true,
			errorMsg:    "break statement outside loop",
		},
		{
			name:        "continue outside loop - invalid",
			code:        "continue;",
			shouldError: true,
			errorMsg:    "continue statement outside loop",
		},
		{
			name:        "break inside if but outside loop - invalid",
			code:        "if true { break; }",
			shouldError: true,
			errorMsg:    "break statement outside loop",
		},
		{
			name:        "continue inside if but outside loop - invalid",
			code:        "if true { continue; }",
			shouldError: true,
			errorMsg:    "continue statement outside loop",
		},
		{
			name:        "simple block test",
			code:        "{ i := 0; }",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := testAnalyze(tt.code)

			if tt.shouldError {
				if errStr == "" {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(errStr, tt.errorMsg) {
					t.Errorf("expected error message containing '%s', got '%s'", tt.errorMsg, errStr)
				}
			} else {
				if errStr != "" {
					t.Errorf("unexpected error: %s", errStr)
				}
			}
		})
	}
}
