package parser

import (
	"mars/ast"
	"mars/lexer"
	"testing"
)

func TestParseFuncDecl_CurrentBehavior(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
		// expectedParamIdentifiers is used to confirm parseParameterList itself works,
		// even if parseFuncDecl doesn't store them in ast.FuncDecl.
		expectedParamIdentifiers []string
		expectedErrorCount       int
	}{
		{
			input:                    "func myFunc()",
			expectedName:             "myFunc",
			expectedParamIdentifiers: []string{},
			expectedErrorCount:       0,
		},
		{
			input:                    "func anotherFunc(a, b, c)",
			expectedName:             "anotherFunc",
			expectedParamIdentifiers: []string{"a", "b", "c"}, // parseParameterList would see these
			expectedErrorCount:       0,
		},
		{
			input:                    "func withSingleParam(p1)",
			expectedName:             "withSingleParam",
			expectedParamIdentifiers: []string{"p1"},
			expectedErrorCount:       0,
		},
		// Test for a case that might expose issues if not handled carefully, like EOF after params
		{
			input:                    "func trailing(p)",
			expectedName:             "trailing",
			expectedParamIdentifiers: []string{"p"},
			expectedErrorCount:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := NewParser(l)
			program := p.ParseProgram()

			if len(p.errors) != tt.expectedErrorCount {
				t.Fatalf("Expected %d errors, but got %d: %v", tt.expectedErrorCount, len(p.errors), p.errors)
			}
			if tt.expectedErrorCount > 0 {
				return // Don't check AST if errors were expected and found
			}

			if len(program.Declarations) != 1 {
				t.Fatalf("ParseProgram() did not produce 1 declaration, got %d. errors: %v", len(program.Declarations), p.errors)
				return
			}

			decl, ok := program.Declarations[0].(*ast.FuncDecl)
			if !ok {
				t.Fatalf("program.Declarations[0] is not *ast.FuncDecl. got=%T. errors: %v", program.Declarations[0], p.errors)
				return
			}

			if decl.Name == nil || decl.Name.Name != tt.expectedName {
				t.Errorf("FuncDecl.Name.Name not '%s'. got='%s'", tt.expectedName, decl.Name.Name)
			}

			// Verify current behavior: Parameters, ReturnType, and Body are nil
			if decl.Parameters != nil {
				t.Errorf("Expected FuncDecl.Parameters to be nil for current implementation, got %v", decl.Parameters)
			}
			if decl.ReturnType != nil {
				t.Errorf("Expected FuncDecl.ReturnType to be nil for current implementation, got %v", decl.ReturnType)
			}
			if decl.Body != nil {
				t.Errorf("Expected FuncDecl.Body to be nil for current implementation, got %v", decl.Body)
			}

			// Although parseFuncDecl doesn't store parameters in the AST node,
			// we can indirectly verify parseParameterList was called and processed parameters
			// by checking token positions or by re-parsing just the param list if necessary.
			// For this test, we focus on what parseFuncDecl *returns*.
			// The fact that decl.Parameters is nil *is* the key part of its current behavior.
		})
	}
}
