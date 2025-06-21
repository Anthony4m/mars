package analyzer

import (
	"mars/ast"
	"mars/lexer"
	"mars/parser"
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

	if xSymbol.Type.BaseType != "float" {
		t.Errorf("Expected x to be float, got %s", xSymbol.Type.BaseType)
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
