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

func TestArithmetic(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected interface{}
		typ      string
	}{
		// Integer arithmetic
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
			expected: int64(5),
			typ:      "int",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: int64(10),
			typ:      "int",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.BinaryExpression{
								Left: &ast.Literal{
									Token:    "2",
									Value:    int64(2),
									Position: ast.Position{Line: 1, Column: 1},
								},
								Operator: "*",
								Right: &ast.Literal{
									Token:    "2",
									Value:    int64(2),
									Position: ast.Position{Line: 1, Column: 5},
								},
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "*",
							Right: &ast.Literal{
								Token:    "2",
								Value:    int64(2),
								Position: ast.Position{Line: 1, Column: 9},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: int64(8),
			typ:      "int",
		},
		// Float arithmetic
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.Literal{
							Token:    "5.5",
							Value:    float64(5.5),
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: float64(5.5),
			typ:      "float",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "2.5",
								Value:    float64(2.5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "2.5",
								Value:    float64(2.5),
								Position: ast.Position{Line: 1, Column: 6},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: float64(5.0),
			typ:      "float",
		},
		// Mixed arithmetic
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "2.5",
								Value:    float64(2.5),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: float64(7.5),
			typ:      "float",
		},
	}

	for _, tc := range tests {
		eval := New()
		result := eval.Eval(tc.input)

		if result == nil {
			t.Fatalf("Evaluator returned nil for input")
		}

		switch tc.typ {
		case "int":
			integerValue, ok := result.(*IntegerValue)
			if !ok {
				t.Fatalf("Expected IntegerValue, got %T", result)
			}
			if integerValue.Value != tc.expected.(int64) {
				t.Errorf("Integer object has wrong value. got=%d, want=%d", integerValue.Value, tc.expected)
			}
		case "float":
			floatValue, ok := result.(*FloatValue)
			if !ok {
				t.Fatalf("Expected FloatValue, got %T", result)
			}
			if floatValue.Value != tc.expected.(float64) {
				t.Errorf("Float object has wrong value. got=%f, want=%f", floatValue.Value, tc.expected)
			}
		}
	}
}

func TestComparison(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected bool
	}{
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "1",
								Value:    int64(1),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "<",
							Right: &ast.Literal{
								Token:    "2",
								Value:    int64(2),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "1",
								Value:    int64(1),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: ">",
							Right: &ast.Literal{
								Token:    "2",
								Value:    int64(2),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: false,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "1",
								Value:    int64(1),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "==",
							Right: &ast.Literal{
								Token:    "1",
								Value:    int64(1),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "1",
								Value:    int64(1),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "!=",
							Right: &ast.Literal{
								Token:    "2",
								Value:    int64(2),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		// Mixed comparisons
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: ">",
							Right: &ast.Literal{
								Token:    "3.5",
								Value:    float64(3.5),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "2.5",
								Value:    float64(2.5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "<",
							Right: &ast.Literal{
								Token:    "3",
								Value:    int64(3),
								Position: ast.Position{Line: 1, Column: 6},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "3.0",
								Value:    float64(3.0),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "==",
							Right: &ast.Literal{
								Token:    "3",
								Value:    int64(3),
								Position: ast.Position{Line: 1, Column: 6},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: false, // 3.0 == 3 should be false since they're different types
		},
	}

	for _, tc := range tests {
		eval := New()
		result := eval.Eval(tc.input)

		if result == nil {
			t.Fatalf("Evaluator returned nil for input")
		}

		booleanValue, ok := result.(*BooleanValue)
		if !ok {
			t.Fatalf("Expected BooleanValue, got %T", result)
		}
		if booleanValue.Value != tc.expected {
			t.Errorf("Boolean object has wrong value. got=%t, want=%t", booleanValue.Value, tc.expected)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    ast.Node
		expected bool
	}{
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.Literal{
								Token:    "true",
								Value:    true,
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: false,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.Literal{
								Token:    "false",
								Value:    false,
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: false,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.UnaryExpression{
								Operator: "!",
								Right: &ast.Literal{
									Token:    "true",
									Value:    true,
									Position: ast.Position{Line: 1, Column: 3},
								},
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.UnaryExpression{
								Operator: "!",
								Right: &ast.Literal{
									Token:    "false",
									Value:    false,
									Position: ast.Position{Line: 1, Column: 3},
								},
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: false,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.UnaryExpression{
								Operator: "!",
								Right: &ast.Literal{
									Token:    "5",
									Value:    int64(5),
									Position: ast.Position{Line: 1, Column: 3},
								},
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.Literal{
								Token:    "0",
								Value:    int64(0),
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.Literal{
								Token:    "\"\"",
								Value:    "",
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: true,
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "!",
							Right: &ast.Literal{
								Token:    "\"hello\"",
								Value:    "hello",
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		eval := New()
		result := eval.Eval(tc.input)

		if result == nil {
			t.Fatalf("Evaluator returned nil for input")
		}

		booleanValue, ok := result.(*BooleanValue)
		if !ok {
			t.Fatalf("Expected BooleanValue, got %T", result)
		}
		if booleanValue.Value != tc.expected {
			t.Errorf("Boolean object has wrong value. got=%t, want=%t", booleanValue.Value, tc.expected)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           ast.Node
		expectedMessage string
	}{
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "true",
								Value:    true,
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "type mismatch: cannot add INTEGER and BOOLEAN",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "true",
								Value:    true,
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.Literal{
							Token:    "5",
							Value:    int64(5),
							Position: ast.Position{Line: 1, Column: 12},
						},
						Position: ast.Position{Line: 1, Column: 12},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "type mismatch: cannot add INTEGER and BOOLEAN",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.UnaryExpression{
							Operator: "-",
							Right: &ast.Literal{
								Token:    "true",
								Value:    true,
								Position: ast.Position{Line: 1, Column: 2},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "unknown operator: -BOOLEAN",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "true",
								Value:    true,
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "false",
								Value:    false,
								Position: ast.Position{Line: 1, Column: 7},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "type mismatch: cannot add BOOLEAN and BOOLEAN",
		},
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
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "true",
								Value:    true,
								Position: ast.Position{Line: 1, Column: 4},
							},
							Operator: "+",
							Right: &ast.Literal{
								Token:    "false",
								Value:    false,
								Position: ast.Position{Line: 1, Column: 10},
							},
							Position: ast.Position{Line: 1, Column: 4},
						},
						Position: ast.Position{Line: 1, Column: 4},
					},
					&ast.ExpressionStatement{
						Expression: &ast.Literal{
							Token:    "5",
							Value:    int64(5),
							Position: ast.Position{Line: 1, Column: 17},
						},
						Position: ast.Position{Line: 1, Column: 17},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "type mismatch: cannot add BOOLEAN and BOOLEAN",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "\"Hello\"",
								Value:    "Hello",
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "-",
							Right: &ast.Literal{
								Token:    "\"World\"",
								Value:    "World",
								Position: ast.Position{Line: 1, Column: 9},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "type mismatch: cannot subtract STRING from STRING",
		},
		{
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.BinaryExpression{
							Left: &ast.Literal{
								Token:    "5",
								Value:    int64(5),
								Position: ast.Position{Line: 1, Column: 1},
							},
							Operator: "/",
							Right: &ast.Literal{
								Token:    "0",
								Value:    int64(0),
								Position: ast.Position{Line: 1, Column: 5},
							},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
				Position: ast.Position{Line: 1, Column: 1},
			},
			expectedMessage: "division by zero",
		},
	}

	for _, tc := range tests {
		eval := New()
		result := eval.Eval(tc.input)

		errObj, ok := result.(*Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", result, result)
			continue
		}

		if errObj.Message != tc.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tc.expectedMessage, errObj.Message)
		}
	}
}
