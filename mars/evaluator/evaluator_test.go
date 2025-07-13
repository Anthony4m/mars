package evaluator

import (
	"mars/ast"
	"testing"
)

func TestIntegerValue_Type(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected any
		typ      string
	}{
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.Literal{
							Token:    "5",
							Value:    int64(5),
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: 5,
			typ:      "int",
		}, {
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.Literal{
							Token:    "true",
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
			typ:      "bool",
		},
	}

	for _, tc := range tests {
		eval := New()
		result := eval.Eval(tc.input)
		// Check if evaluation returned a result
		if result == nil {
			t.Fatalf("Evaluator returned nil for input: 5")
		}

		switch tc.typ {
		case "int":
			// Type assert to IntegerValue
			integerValue, ok := result.(*IntegerValue)
			if !ok {
				t.Fatalf("Expected IntegerValue, got %T", result)
			}

			if integerValue.Value != int64(tc.expected.(int)) {
				t.Errorf("Integer object has wrong value. got=%d, want=%d", integerValue.Value, tc.expected)
			}
		case "bool":
			// Type assert to BooleanValue
			booleanValue, ok := result.(*BooleanValue)
			if !ok {
				t.Fatalf("Expected BooleanValue, got %T", result)
			}
			if booleanValue.Value != tc.expected {
				t.Errorf("Boolean object has wrong value. got=%t, want=%t", booleanValue.Value, tc.expected)
			}
		}

	}

}
