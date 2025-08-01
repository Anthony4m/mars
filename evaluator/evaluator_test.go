package evaluator

import (
	"mars/ast"
	"strings"
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

		errObj, ok := result.(*RuntimeError)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", result, result)
			continue
		}

		if errObj.Detail.Message != tc.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tc.expectedMessage, errObj.Detail.Message)
		}
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		input         ast.Node
		expectedError string
		expectedCode  string
		hasStackTrace bool
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
			expectedError: "type mismatch: cannot add INTEGER and BOOLEAN",
			expectedCode:  "E001",
			hasStackTrace: true,
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
			expectedError: "division by zero",
			expectedCode:  "E001", // using 001 for now because I havn't implemented the other error codes
			hasStackTrace: true,
		},
	}

	for _, tc := range tests {
		eval := New()
		result := eval.Eval(tc.input)

		errObj, ok := result.(*RuntimeError)
		if !ok {
			t.Fatalf("no error object returned. got=%T(%+v)", result, result)
		}

		if errObj.Detail.Message != tc.expectedError {
			t.Errorf("wrong error message. got=%q, want=%q",
				errObj.Detail.Message, tc.expectedError)
		}

		if errObj.Detail.ErrorCode != tc.expectedCode {
			t.Errorf("wrong error code. got=%q, want=%q",
				errObj.Detail.ErrorCode, tc.expectedCode)
		}

		if tc.hasStackTrace && len(errObj.StackTrace) == 0 {
			t.Error("expected stack trace but got none")
		}
	}
}

func TestConditionalStatements(t *testing.T) {
	tests := []struct {
		name            string
		input           ast.Node
		expectedMessage string
	}{
		{
			name: "if true returns int",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "10",
									Value:    int64(10),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "if true with else returns int",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "10",
									Value:    int64(10),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
						Alternative: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "20",
									Value:    int64(20),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "nested if false returns int",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.IfStatement{
								Condition: &ast.Literal{
									Token:    BOOLEAN_TYPE,
									Value:    false,
									Position: ast.Position{Line: 1, Column: 1},
								},
								Consequence: &ast.BlockStatement{Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Literal{
											Token:    "30",
											Value:    int64(30),
											Position: ast.Position{Line: 1, Column: 1},
										},
									},
								}},
							},
						}},
					},
				},
			},
			expectedMessage: NULL_TYPE,
		},
		{
			name: "if with string condition errors",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    "string",
							Value:    "invalid",
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "50",
									Value:    int64(50),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "IfWithFalseConditionElseBranch",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    false,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
						Alternative: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "110",
									Value:    int64(110),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "IfWithInvalidConditionType",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    "string",
							Value:    "",
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "50",
									Value:    int64(50),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: NULL_TYPE,
		},
		{
			name: "IfWithMultipleStatementsInConsequence",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "60",
									Value:    int64(60),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "70",
									Value:    int64(70),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "IfWithEmptyConsequence",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{}},
					},
				},
			},
			expectedMessage: NULL_TYPE,
		},
		{
			name: "IfWithNilCondition",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: nil,
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "120",
									Value:    int64(120),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: NULL_TYPE,
		},
		{
			name: "IfWithMixedTypesInConsequence",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "80",
									Value:    int64(80),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "string",
									Value:    "",
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: STRING_TYPE,
		},
		{
			name: "IfFalseNoElse",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    false,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "90",
									Value:    int64(90),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						}},
					},
				},
			},
			expectedMessage: NULL_TYPE,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eval := New()
			result := eval.Eval(tc.input)
			if tc.expectedMessage != "" {
				if result == nil {
					t.Fatalf("Expected %s but got nil result", tc.expectedMessage)
				}
				if result.Type() != tc.expectedMessage {
					t.Errorf("Expected %s but got %s", tc.expectedMessage, result.Type())
				}
			} else if result != nil && result.Type() != "" {
				t.Errorf("Expected empty result but got %s", result.Type())
			}
		})
	}
}

func TestBlockStatement(t *testing.T) {
	tests := []struct {
		name            string
		input           ast.Node
		expectedMessage string
	}{
		{
			name: "SimpleBlockWithMultipleStatements",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "200",
									Value:    int64(200),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE, // Should return last statement result
		},
		{
			name: "BlockWithIfStatementTrueBranch",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.IfStatement{
								Condition: &ast.Literal{
									Token:    BOOLEAN_TYPE,
									Value:    true,
									Position: ast.Position{Line: 1, Column: 1},
								},
								Consequence: &ast.BlockStatement{Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Literal{
											Token:    "10",
											Value:    int64(10),
											Position: ast.Position{Line: 1, Column: 1},
										},
									},
								}},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "BlockWithIfStatementFalseBranch",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.IfStatement{
								Condition: &ast.Literal{
									Token:    BOOLEAN_TYPE,
									Value:    false,
									Position: ast.Position{Line: 1, Column: 1},
								},
								Consequence: &ast.BlockStatement{Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Literal{
											Token:    "10",
											Value:    int64(10),
											Position: ast.Position{Line: 1, Column: 1},
										},
									},
								}},
								Alternative: &ast.BlockStatement{Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Literal{
											Token:    "20",
											Value:    int64(20),
											Position: ast.Position{Line: 1, Column: 1},
										},
									},
								}},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "BlockWithIfStatementFalseNoAlternative",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.IfStatement{
								Condition: &ast.Literal{
									Token:    BOOLEAN_TYPE,
									Value:    false,
									Position: ast.Position{Line: 1, Column: 1},
								},
								Consequence: &ast.BlockStatement{Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Literal{
											Token:    "10",
											Value:    int64(10),
											Position: ast.Position{Line: 1, Column: 1},
										},
									},
								}},
							},
						},
					},
				},
			},
			expectedMessage: NULL_TYPE, // Should return empty since last statement returns null
		},
		{
			name: "EmptyBlockStatement",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{},
					},
				},
			},
			expectedMessage: NULL_TYPE,
		},
		{
			name: "BlockWithReturnStatement",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ReturnStatement{
								Value: &ast.Literal{
									Token:    "50",
									Value:    int64(50),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "200",
									Value:    int64(200),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE, // Should return 50 from return statement
		},
		{
			name: "BlockWithNestedBlockStatements",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.BlockStatement{
								Statements: []ast.Statement{
									&ast.ExpressionStatement{
										Expression: &ast.Literal{
											Token:    "30",
											Value:    int64(30),
											Position: ast.Position{Line: 1, Column: 1},
										},
									},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "40",
									Value:    int64(40),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "BlockWithVariableDeclaration",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.VarDecl{
								Name: &ast.Identifier{
									Name:     "x",
									Position: ast.Position{Line: 1, Column: 1},
								},
								Value: &ast.Literal{
									Token:    "123",
									Value:    int64(123),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Identifier{
									Name:     "x",
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "BlockWithMixedTypesLastInteger",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "string",
									Value:    "hello",
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "42",
									Value:    int64(42),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "BlockWithErrorInMiddle",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "100",
									Value:    int64(100),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.BinaryExpression{
									Left: &ast.Literal{
										Token:    "10",
										Value:    int64(10),
										Position: ast.Position{Line: 1, Column: 1},
									},
									Operator: "/",
									Right: &ast.Literal{
										Token:    "0",
										Value:    int64(0),
										Position: ast.Position{Line: 1, Column: 1},
									},
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Literal{
									Token:    "200",
									Value:    int64(200),
									Position: ast.Position{Line: 1, Column: 1},
								},
							},
						},
					},
				},
			},
			expectedMessage: ERROR_TYPE, // Assuming you have error type handling
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eval := New()
			result := eval.Eval(tc.input)

			if tc.expectedMessage != "" {
				if result == nil {
					t.Fatalf("Expected %s but got nil result", tc.expectedMessage)
				}
				if tc.expectedMessage != result.Type() {
					t.Errorf("expected %s got %s", tc.expectedMessage, result.Type())
				}
			} else if result != nil && result.Type() != "" {
				t.Errorf("Expected empty result but got %s", result.Type())
			}
		})
	}
}

func TestVariableDecl(t *testing.T) {
	tests := []struct {
		name            string
		input           ast.Node
		expectedMessage string
	}{
		{
			name: "StringVarWithValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  false,
						Name:     &ast.Identifier{Name: "x"},
						Type:     &ast.Type{BaseType: STRING_TYPE},
						Value:    &ast.Literal{Token: "42", Value: "42"},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: STRING_TYPE,
		},
		{
			name: "IntVarWithValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  true,
						Name:     &ast.Identifier{Name: "y"},
						Type:     &ast.Type{BaseType: INTEGER_TYPE},
						Value:    &ast.Literal{Token: "100", Value: int64(100)},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "BoolVarWithValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  false,
						Name:     &ast.Identifier{Name: "z"},
						Type:     &ast.Type{BaseType: BOOLEAN_TYPE},
						Value:    &ast.Literal{Token: BOOLEAN_TYPE, Value: true},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: BOOLEAN_TYPE,
		},
		{
			name: "VarWithoutTypeButWithValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  true,
						Name:     &ast.Identifier{Name: "a"},
						Type:     nil,
						Value:    &ast.Literal{Token: "200", Value: int64(200)},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "VarWithTypeButNoValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  false,
						Name:     &ast.Identifier{Name: "b"},
						Type:     &ast.Type{BaseType: INTEGER_TYPE},
						Value:    nil,
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: INTEGER_TYPE, // Should initialize to zero value
		},
		{
			name: "StringVarWithTypeButNoValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  true,
						Name:     &ast.Identifier{Name: "c"},
						Type:     &ast.Type{BaseType: STRING_TYPE},
						Value:    nil,
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: STRING_TYPE, // Should initialize to empty string
		},
		{
			name: "BoolVarWithTypeButNoValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  false,
						Name:     &ast.Identifier{Name: "d"},
						Type:     &ast.Type{BaseType: BOOLEAN_TYPE},
						Value:    nil,
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: BOOLEAN_TYPE, // Should initialize to false
		},
		{
			name: "VarWithIncompatibleTypes",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  true,
						Name:     &ast.Identifier{Name: "e"},
						Type:     &ast.Type{BaseType: INTEGER_TYPE},
						Value:    &ast.Literal{Token: "hello", Value: "hello"},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "type mismatch: cannot assign STRING to INTEGER",
		},
		{
			name: "VarWithoutName",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  false,
						Name:     nil,
						Type:     &ast.Type{BaseType: INTEGER_TYPE},
						Value:    &ast.Literal{Token: "300", Value: int64(300)},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "variable declaration missing name",
		},
		{
			name: "VarWithoutTypeAndValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  true,
						Name:     &ast.Identifier{Name: "f"},
						Type:     nil,
						Value:    nil,
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "variable 'f' needs type or initial value",
		},
		{
			name: "NestedVarDeclarationInBlock",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.VarDecl{
								Mutable:  false,
								Name:     &ast.Identifier{Name: "g"},
								Type:     &ast.Type{BaseType: INTEGER_TYPE},
								Value:    &ast.Literal{Token: "400", Value: int64(400)},
								Position: ast.Position{Line: 1, Column: 1},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Identifier{Name: "g"},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "MultipleVarDeclarations",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.VarDecl{
								Mutable:  true,
								Name:     &ast.Identifier{Name: "h1"},
								Type:     &ast.Type{BaseType: INTEGER_TYPE},
								Value:    &ast.Literal{Token: "500", Value: int64(500)},
								Position: ast.Position{Line: 1, Column: 1},
							},
							&ast.VarDecl{
								Mutable:  false,
								Name:     &ast.Identifier{Name: "h2"},
								Type:     &ast.Type{BaseType: STRING_TYPE},
								Value:    &ast.Literal{Token: "test", Value: "test"},
								Position: ast.Position{Line: 1, Column: 1},
							},
							&ast.ExpressionStatement{
								Expression: &ast.Identifier{Name: "h2"},
							},
						},
					},
				},
			},
			expectedMessage: STRING_TYPE,
		},
		{
			name: "VarWithFloatValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable:  true,
						Name:     &ast.Identifier{Name: "i"},
						Type:     &ast.Type{BaseType: FLOAT_TYPE},
						Value:    &ast.Literal{Token: "3.14", Value: 3.14},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: FLOAT_TYPE,
		},
		{
			name: "VarDeclarationInIfStatement",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.IfStatement{
						Condition: &ast.Literal{
							Token:    BOOLEAN_TYPE,
							Value:    true,
							Position: ast.Position{Line: 1, Column: 1},
						},
						Consequence: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.VarDecl{
									Mutable:  false,
									Name:     &ast.Identifier{Name: "j"},
									Type:     &ast.Type{BaseType: INTEGER_TYPE},
									Value:    &ast.Literal{Token: "600", Value: int64(600)},
									Position: ast.Position{Line: 1, Column: 1},
								},
								&ast.ExpressionStatement{
									Expression: &ast.Identifier{Name: "j"},
								},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "VarWithExpressionValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Mutable: true,
						Name:    &ast.Identifier{Name: "k"},
						Type:    &ast.Type{BaseType: INTEGER_TYPE},
						Value: &ast.BinaryExpression{
							Left:     &ast.Literal{Token: "10", Value: int64(10)},
							Operator: "+",
							Right:    &ast.Literal{Token: "20", Value: int64(20)},
							Position: ast.Position{Line: 1, Column: 1},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eval := New()
			result := eval.Eval(tc.input)

			if strings.Contains(tc.expectedMessage, "error") ||
				strings.Contains(tc.expectedMessage, "mismatch") ||
				strings.Contains(tc.expectedMessage, "needs") ||
				strings.Contains(tc.expectedMessage, "missing") {
				// Expecting an error message
				if result == nil {
					t.Fatalf("Expected error but got nil result")
				}
				errobj, _ := result.(*RuntimeError)
				if !strings.Contains(errobj.Detail.Message, tc.expectedMessage) {
					t.Errorf("expected error containing '%s' got '%s'", tc.expectedMessage, errobj.Detail.Message)
				}
			} else if tc.expectedMessage != "" {
				// Expecting a specific type
				if result == nil {
					t.Fatalf("Expected %s but got nil result", tc.expectedMessage)
				}
				if tc.expectedMessage != result.Type() {
					t.Errorf("expected %s got %s", tc.expectedMessage, result.Type())
				}
			} else if result != nil && result.Type() != "" {
				t.Errorf("Expected empty result but got %s", result.Type())
			}
		})
	}
}

func TestEvaluator_EvalAssignment(t *testing.T) {
	tests := []struct {
		name            string
		input           ast.Node
		expectedMessage string
	}{
		{
			name: "SimpleAssignment",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "x"},
						Value:   &ast.Literal{Token: "5", Value: int64(5)},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "x"},
						Value: &ast.Literal{Token: "10", Value: int64(10)},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "AssignmentToStringVariable",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "name"},
						Value:   &ast.Literal{Token: "initial", Value: "initial"},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "name"},
						Value: &ast.Literal{Token: "updated", Value: "updated"},
					},
				},
			},
			expectedMessage: STRING_TYPE,
		},
		{
			name: "AssignmentToBooleanVariable",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "flag"},
						Value:   &ast.Literal{Token: BOOLEAN_TYPE, Value: false},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "flag"},
						Value: &ast.Literal{Token: BOOLEAN_TYPE, Value: true},
					},
				},
			},
			expectedMessage: BOOLEAN_TYPE,
		},
		{
			name: "AssignmentWithExpressionValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "a"},
						Value:   &ast.Literal{Token: "15", Value: int64(15)},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name: &ast.Identifier{Name: "a"},
						Value: &ast.BinaryExpression{
							Left:     &ast.Literal{Token: "10", Value: int64(10)},
							Operator: "+",
							Right:    &ast.Literal{Token: "20", Value: int64(20)},
							Position: ast.Position{Line: 1, Column: 1},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "AssignmentToImmutableVariable",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "y"},
						Value:   &ast.Literal{Token: "5", Value: int64(5)},
						Mutable: false, // Immutable variable
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "y"},
						Value: &ast.Literal{Token: "10", Value: int64(10)},
					},
				},
			},
			expectedMessage: "cannot assign to immutable variable 'y'",
		},
		{
			name: "AssignmentToUndeclaredVariable",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "undeclared"},
						Value: &ast.Literal{Token: "42", Value: int64(42)},
					},
				},
			},
			expectedMessage: "undefined variable 'undeclared'",
		},
		{
			name: "AssignmentWithIncompatibleTypes",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "num"},
						Value:   &ast.Literal{Token: "100", Value: int64(100)},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "num"},
						Value: &ast.Literal{Token: "hello", Value: "hello"}, // String to int
					},
				},
			},
			expectedMessage: "type mismatch: cannot assign STRING to INTEGER",
		},
		{
			name: "MultipleAssignments",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "x"},
						Value:   &ast.Literal{Token: "1", Value: int64(1)},
						Mutable: true,
					},
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "y"},
						Value:   &ast.Literal{Token: "2", Value: int64(2)},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "x"},
						Value: &ast.Literal{Token: "10", Value: int64(10)},
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "y"},
						Value: &ast.Literal{Token: "20", Value: int64(20)},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "AssignmentInBlockStatement",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "z"},
						Value:   &ast.Literal{Token: "0", Value: int64(0)},
						Mutable: true,
					},
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.AssignmentStatement{
								Name:  &ast.Identifier{Name: "z"},
								Value: &ast.Literal{Token: "99", Value: int64(99)},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "AssignmentInIfStatement",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "conditionVar"},
						Value:   &ast.Literal{Token: BOOLEAN_TYPE, Value: false},
						Mutable: true,
					},
					&ast.IfStatement{
						Condition: &ast.Literal{Token: BOOLEAN_TYPE, Value: true},
						Consequence: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.AssignmentStatement{
									Name:  &ast.Identifier{Name: "conditionVar"},
									Value: &ast.Literal{Token: BOOLEAN_TYPE, Value: true},
								},
							},
						},
					},
				},
			},
			expectedMessage: BOOLEAN_TYPE,
		},
		{
			name: "AssignmentToZeroValueVariable",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "zeroInt"},
						Type:    &ast.Type{BaseType: INTEGER_TYPE},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "zeroInt"},
						Value: &ast.Literal{Token: "42", Value: int64(42)},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "AssignmentWithNilIdentifier",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.AssignmentStatement{
						Name:  nil, // Nil identifier
						Value: &ast.Literal{Token: "42", Value: int64(42)},
					},
				},
			},
			expectedMessage: "assignment missing variable name",
		},
		{
			name: "AssignmentWithNilValue",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "nilTest"},
						Value:   &ast.Literal{Token: "1", Value: int64(1)},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "nilTest"},
						Value: nil, // Nil value
					},
				},
			},
			expectedMessage: "assignment missing value",
		},
		{
			name: "ChainedAssignments",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "first"},
						Value:   &ast.Literal{Token: "100", Value: int64(100)},
						Mutable: true,
					},
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "second"},
						Value:   &ast.Literal{Token: "200", Value: int64(200)},
						Mutable: true,
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "first"},
						Value: &ast.Literal{Token: "500", Value: int64(500)},
					},
					&ast.AssignmentStatement{
						Name:  &ast.Identifier{Name: "second"},
						Value: &ast.Identifier{Name: "first"}, // Assign first's value to second
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "AssignmentInNestedScope",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.VarDecl{
						Name:    &ast.Identifier{Name: "outer"},
						Value:   &ast.Literal{Token: "1", Value: int64(1)},
						Mutable: true,
					},
					&ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.AssignmentStatement{
								Name:  &ast.Identifier{Name: "outer"},
								Value: &ast.Literal{Token: "1000", Value: int64(1000)},
							},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eval := New()
			result := eval.Eval(tc.input)

			// Check if we're expecting an error message
			if strings.Contains(tc.expectedMessage, "error") ||
				strings.Contains(tc.expectedMessage, "cannot") ||
				strings.Contains(tc.expectedMessage, "undefined") ||
				strings.Contains(tc.expectedMessage, "missing") ||
				strings.Contains(tc.expectedMessage, "mismatch") {
				// Expecting an error
				if result == nil {
					t.Fatalf("[%s] Expected error but got nil result", tc.name)
				}
				errobj, _ := result.(*RuntimeError)
				if !strings.Contains(errobj.Detail.Message, tc.expectedMessage) {
					t.Errorf("expected error containing '%s' got '%s'", tc.expectedMessage, errobj.Detail.Message)
				}
				if !strings.Contains(errobj.Detail.Message, tc.expectedMessage) &&
					result.Type() != tc.expectedMessage {
					t.Errorf("[%s] Expected error containing '%s' but got '%s'",
						tc.name, tc.expectedMessage, errobj.Detail.Message)
				}
			} else if tc.expectedMessage != "" {
				// Expecting a successful result with specific type
				if result == nil {
					t.Fatalf("[%s] Expected %s but got nil result", tc.name, tc.expectedMessage)
				}
				if result.Type() != tc.expectedMessage {
					t.Errorf("[%s] Expected %s but got %s", tc.name, tc.expectedMessage, result.Type())
				}
			} else if result != nil && result.Type() != "" {
				t.Errorf("[%s] Expected empty result but got %s", tc.name, result.Type())
			}
		})
	}
}

func TestFuncDecl(t *testing.T) {
	tests := []struct {
		name            string
		input           ast.Node
		expectedMessage string
	}{
		{
			name: "SimpleFunctionDeclaration",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "add"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.Literal{
										Token:    "42",
										Value:    int64(42),
										Position: ast.Position{Line: 2, Column: 5},
									},
									Position: ast.Position{Line: 2, Column: 1},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
		{
			name: "FunctionWithParameters",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "multiply"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "a"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
								{
									Name: &ast.Identifier{Name: "b"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "a"},
										Operator: "*",
										Right:    &ast.Identifier{Name: "b"},
										Position: ast.Position{Line: 2, Column: 9},
									},
									Position: ast.Position{Line: 2, Column: 5},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
		{
			name: "FunctionWithoutReturnType",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "printHello"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: nil,
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Expression: &ast.Literal{
										Token:    "hello",
										Value:    "hello",
										Position: ast.Position{Line: 2, Column: 5},
									},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
		{
			name: "FunctionWithoutName",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: nil,
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.Literal{
										Token:    "10",
										Value:    int64(10),
										Position: ast.Position{Line: 2, Column: 5},
									},
									Position: ast.Position{Line: 2, Column: 1},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "function declaration missing name",
		},
		{
			name: "FunctionWithVoidReturn",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "doNothing"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: nil,
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
		{
			name: "FunctionWithComplexBody",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "complexFunc"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "x"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.VarDecl{
									Name: &ast.Identifier{Name: "result"},
									Value: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "x"},
										Operator: "*",
										Right:    &ast.Literal{Token: "2", Value: int64(2)},
										Position: ast.Position{Line: 2, Column: 15},
									},
									Mutable:  false,
									Position: ast.Position{Line: 2, Column: 5},
								},
								&ast.IfStatement{
									Condition: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "result"},
										Operator: ">",
										Right:    &ast.Literal{Token: "10", Value: int64(10)},
										Position: ast.Position{Line: 3, Column: 15},
									},
									Consequence: &ast.BlockStatement{
										Statements: []ast.Statement{
											&ast.ReturnStatement{
												Value:    &ast.Identifier{Name: "result"},
												Position: ast.Position{Line: 4, Column: 9},
											},
										},
									},
									Position: ast.Position{Line: 3, Column: 5},
								},
								&ast.ReturnStatement{
									Value:    &ast.Literal{Token: "0", Value: int64(0)},
									Position: ast.Position{Line: 6, Column: 5},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
		{
			name: "FunctionWithStringParameter",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "greet"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "name"},
									Type: &ast.Type{BaseType: STRING_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: STRING_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Literal{Token: "Hello, ", Value: "Hello, "},
										Operator: "+",
										Right:    &ast.Identifier{Name: "name"},
										Position: ast.Position{Line: 2, Column: 16},
									},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
		{
			name: "MultipleFunctionDeclarations",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "first"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value:    &ast.Literal{Token: "1", Value: int64(1)},
									Position: ast.Position{Line: 2, Column: 5},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "second"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 4, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value:    &ast.Literal{Token: "2", Value: int64(2)},
									Position: ast.Position{Line: 5, Column: 5},
								},
							},
						},
						Position: ast.Position{Line: 4, Column: 1},
					},
				},
			},
			expectedMessage: "FUNCTION",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eval := New()
			result := eval.Eval(tc.input)

			// Check if we're expecting an error message
			if strings.Contains(tc.expectedMessage, "error") ||
				strings.Contains(tc.expectedMessage, "missing") {
				// Expecting an error
				if result == nil {
					t.Fatalf("[%s] Expected error but got nil result", tc.name)
				}
				if !strings.Contains(result.String(), tc.expectedMessage) &&
					result.Type() != tc.expectedMessage {
					t.Errorf("[%s] Expected error containing '%s' but got '%s'",
						tc.name, tc.expectedMessage, result.String())
				}
			} else if tc.expectedMessage != "" {
				// Expecting a successful result with specific type
				if result == nil {
					t.Fatalf("[%s] Expected %s but got nil result", tc.name, tc.expectedMessage)
				}
				if result.Type() != tc.expectedMessage {
					t.Errorf("[%s] Expected %s but got %s", tc.name, tc.expectedMessage, result.Type())
				}
			} else if result != nil && result.Type() != "" {
				t.Errorf("[%s] Expected empty result but got %s", tc.name, result.Type())
			}
		})
	}
}

func TestFuncCall(t *testing.T) {
	tests := []struct {
		name            string
		input           ast.Node
		expectedMessage string
	}{
		{
			name: "SimpleFunctionCall",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					// First declare the function
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "simpleFunc"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.Literal{
										Token:    "42",
										Value:    int64(42),
										Position: ast.Position{Line: 2, Column: 9},
									},
									Position: ast.Position{Line: 2, Column: 5},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					// Then call it
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function:  &ast.Identifier{Name: "simpleFunc"},
							Arguments: []ast.Expression{},
							Position:  ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "FunctionCallWithArguments",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					// Declare function with parameters
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "add"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "a"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
								{
									Name: &ast.Identifier{Name: "b"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "a"},
										Operator: "+",
										Right:    &ast.Identifier{Name: "b"},
										Position: ast.Position{Line: 2, Column: 16},
									},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					// Call the function
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function: &ast.Identifier{Name: "add"},
							Arguments: []ast.Expression{
								&ast.Literal{Token: "5", Value: int64(5)},
								&ast.Literal{Token: "3", Value: int64(3)},
							},
							Position: ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "FunctionCallWithStringReturn",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "greet"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "name"},
									Type: &ast.Type{BaseType: STRING_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: STRING_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Literal{Token: "Hello, ", Value: "Hello, "},
										Operator: "+",
										Right:    &ast.Identifier{Name: "name"},
										Position: ast.Position{Line: 2, Column: 16},
									},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function: &ast.Identifier{Name: "greet"},
							Arguments: []ast.Expression{
								&ast.Literal{Token: "World", Value: "World"},
							},
							Position: ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: STRING_TYPE,
		},
		{
			name: "FunctionCallUndefinedFunction",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function:  &ast.Identifier{Name: "nonExistentFunc"},
							Arguments: []ast.Expression{},
							Position:  ast.Position{Line: 1, Column: 1},
						},
					},
				},
			},
			expectedMessage: "undefined variable 'nonExistentFunc",
		},
		{
			name: "FunctionCallWithWrongArgumentCount",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "singleParam"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "x"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value:    &ast.Identifier{Name: "x"},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function: &ast.Identifier{Name: "singleParam"},
							Arguments: []ast.Expression{
								&ast.Literal{Token: "1", Value: int64(1)},
								&ast.Literal{Token: "2", Value: int64(2)}, // Extra argument
							},
							Position: ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: "function 'singleParam' expects 1 arguments, got 2",
		},
		{
			name: "FunctionCallWithTypeMismatch",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "expectInt"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "num"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value:    &ast.Identifier{Name: "num"},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function: &ast.Identifier{Name: "expectInt"},
							Arguments: []ast.Expression{
								&ast.Literal{Token: "hello", Value: "hello"}, // String instead of int
							},
							Position: ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: "type mismatch: cannot assign string to int",
		},
		{
			name: "VoidFunctionCall",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "voidFunc"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{},
							ReturnType: nil, // Void return
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Expression: &ast.Literal{
										Token:    "executed",
										Value:    "executed",
										Position: ast.Position{Line: 2, Column: 9},
									},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function:  &ast.Identifier{Name: "voidFunc"},
							Arguments: []ast.Expression{},
							Position:  ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: NULL_TYPE, // Void functions return nothing/nil
		},
		{
			name: "NestedFunctionCall",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "doubleValue"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "x"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "x"},
										Operator: "*",
										Right:    &ast.Literal{Token: "2", Value: int64(2)},
										Position: ast.Position{Line: 2, Column: 16},
									},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "addTen"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "x"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: INTEGER_TYPE},
							Position:   ast.Position{Line: 5, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "x"},
										Operator: "+",
										Right:    &ast.Literal{Token: "10", Value: int64(10)},
										Position: ast.Position{Line: 6, Column: 16},
									},
									Position: ast.Position{Line: 6, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 5, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function: &ast.Identifier{Name: "addTen"},
							Arguments: []ast.Expression{
								&ast.FunctionCall{
									Function: &ast.Identifier{Name: "doubleValue"},
									Arguments: []ast.Expression{
										&ast.Literal{Token: "5", Value: int64(5)},
									},
									Position: ast.Position{Line: 9, Column: 13},
								},
							},
							Position: ast.Position{Line: 9, Column: 1},
						},
					},
				},
			},
			expectedMessage: INTEGER_TYPE,
		},
		{
			name: "FunctionCallWithBooleanReturn",
			input: &ast.Program{
				Declarations: []ast.Declaration{
					&ast.FuncDecl{
						Name: &ast.Identifier{Name: "isPositive"},
						Signature: &ast.FunctionSignature{
							Parameters: []*ast.Parameter{
								{
									Name: &ast.Identifier{Name: "num"},
									Type: &ast.Type{BaseType: INTEGER_TYPE},
								},
							},
							ReturnType: &ast.Type{BaseType: BOOLEAN_TYPE},
							Position:   ast.Position{Line: 1, Column: 1},
						},
						Body: &ast.BlockStatement{
							Statements: []ast.Statement{
								&ast.ReturnStatement{
									Value: &ast.BinaryExpression{
										Left:     &ast.Identifier{Name: "num"},
										Operator: ">",
										Right:    &ast.Literal{Token: "0", Value: int64(0)},
										Position: ast.Position{Line: 2, Column: 16},
									},
									Position: ast.Position{Line: 2, Column: 9},
								},
							},
						},
						Position: ast.Position{Line: 1, Column: 1},
					},
					&ast.ExpressionStatement{
						Expression: &ast.FunctionCall{
							Function: &ast.Identifier{Name: "isPositive"},
							Arguments: []ast.Expression{
								&ast.Literal{Token: "15", Value: int64(15)},
							},
							Position: ast.Position{Line: 4, Column: 1},
						},
					},
				},
			},
			expectedMessage: BOOLEAN_TYPE,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eval := New()
			result := eval.Eval(tc.input)

			// Check if we're expecting an error message
			if strings.Contains(tc.expectedMessage, "error") ||
				strings.Contains(tc.expectedMessage, "undefined") ||
				strings.Contains(tc.expectedMessage, "wrong") ||
				strings.Contains(tc.expectedMessage, "expects") ||
				strings.Contains(tc.expectedMessage, "mismatch") {
				// Expecting an error
				if result == nil {
					t.Fatalf("[%s] Expected error but got nil result", tc.name)
				}
				if !strings.Contains(result.String(), tc.expectedMessage) &&
					result.Type() != tc.expectedMessage {
					t.Errorf("[%s] Expected error containing '%s' but got '%s'",
						tc.name, tc.expectedMessage, result.String())
				}
			} else if tc.expectedMessage != "" {
				// Expecting a successful result with specific type
				if result == nil {
					t.Fatalf("[%s] Expected %s but got nil result", tc.name, tc.expectedMessage)
				}
				if result.Type() != tc.expectedMessage {
					t.Errorf("[%s] Expected %s but got %s", tc.name, tc.expectedMessage, result.Type())
				}
			} else if result != nil && result.Type() != "" {
				// Expecting empty result (void function)
				t.Errorf("[%s] Expected empty result but got %s", tc.name, result.Type())
			}
		})
	}
}
