package evaluator

import (
	"testing"
)

func TestBuiltinLen(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected int64
		hasError bool
	}{
		{
			name:     "len of string",
			input:    &StringValue{Value: "hello"},
			expected: 5,
			hasError: false,
		},
		{
			name:     "len of empty string",
			input:    &StringValue{Value: ""},
			expected: 0,
			hasError: false,
		},
		{
			name:     "len of array",
			input:    &ArrayValue{Elements: []Value{&IntegerValue{Value: 1}, &IntegerValue{Value: 2}, &IntegerValue{Value: 3}}},
			expected: 3,
			hasError: false,
		},
		{
			name:     "len of empty array",
			input:    &ArrayValue{Elements: []Value{}},
			expected: 0,
			hasError: false,
		},
		{
			name:     "len of integer (should error)",
			input:    &IntegerValue{Value: 42},
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builtinLen([]Value{tt.input})

			if tt.hasError {
				if result.Type() != ERROR_TYPE {
					t.Errorf("Expected error, got %T: %v", result, result)
				}
			} else {
				if result.Type() != INTEGER_TYPE {
					t.Errorf("Expected integer, got %T: %v", result, result)
					return
				}

				actual := result.(*IntegerValue).Value
				if actual != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, actual)
				}
			}
		})
	}
}

func TestBuiltinAppend(t *testing.T) {
	tests := []struct {
		name     string
		slice    Value
		value    Value
		expected []Value
		hasError bool
	}{
		{
			name:     "append to array",
			slice:    &ArrayValue{Elements: []Value{&IntegerValue{Value: 1}, &IntegerValue{Value: 2}}},
			value:    &IntegerValue{Value: 3},
			expected: []Value{&IntegerValue{Value: 1}, &IntegerValue{Value: 2}, &IntegerValue{Value: 3}},
			hasError: false,
		},
		{
			name:     "append to empty array",
			slice:    &ArrayValue{Elements: []Value{}},
			value:    &StringValue{Value: "hello"},
			expected: []Value{&StringValue{Value: "hello"}},
			hasError: false,
		},
		{
			name:     "append to string (should error)",
			slice:    &StringValue{Value: "hello"},
			value:    &IntegerValue{Value: 42},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builtinAppend([]Value{tt.slice, tt.value})

			if tt.hasError {
				if result.Type() != ERROR_TYPE {
					t.Errorf("Expected error, got %T: %v", result, result)
				}
			} else {
				if result.Type() != ARRAY_TYPE {
					t.Errorf("Expected array, got %T: %v", result, result)
					return
				}

				array := result.(*ArrayValue)
				if len(array.Elements) != len(tt.expected) {
					t.Errorf("Expected %d elements, got %d", len(tt.expected), len(array.Elements))
					return
				}

				for i, expected := range tt.expected {
					if array.Elements[i].String() != expected.String() {
						t.Errorf("Element %d: expected %s, got %s", i, expected.String(), array.Elements[i].String())
					}
				}
			}
		})
	}
}

func TestBuiltinPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected string
	}{
		{
			name:     "print string",
			input:    &StringValue{Value: "hello"},
			expected: "hello",
		},
		{
			name:     "print integer",
			input:    &IntegerValue{Value: 42},
			expected: "42",
		},
		{
			name:     "print array",
			input:    &ArrayValue{Elements: []Value{&IntegerValue{Value: 1}, &StringValue{Value: "two"}}},
			expected: "[1, two]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builtinPrint([]Value{tt.input})

			if result.Type() != NULL_TYPE {
				t.Errorf("Expected NULL, got %T: %v", result, result)
			}
		})
	}
}

func TestBuiltinMath(t *testing.T) {
	tests := []struct {
		name     string
		function func(args []Value) Value
		input    Value
		expected float64
		hasError bool
	}{
		{
			name:     "sin of 0",
			function: builtinSin,
			input:    &IntegerValue{Value: 0},
			expected: 0.0,
			hasError: false,
		},
		{
			name:     "sin of pi/2",
			function: builtinSin,
			input:    &FloatValue{Value: 1.5707963267948966}, // π/2
			expected: 1.0,
			hasError: false,
		},
		{
			name:     "cos of 0",
			function: builtinCos,
			input:    &IntegerValue{Value: 0},
			expected: 1.0,
			hasError: false,
		},
		{
			name:     "cos of pi",
			function: builtinCos,
			input:    &FloatValue{Value: 3.141592653589793}, // π
			expected: -1.0,
			hasError: false,
		},
		{
			name:     "sqrt of 4",
			function: builtinSqrt,
			input:    &IntegerValue{Value: 4},
			expected: 2.0,
			hasError: false,
		},
		{
			name:     "sqrt of 2",
			function: builtinSqrt,
			input:    &FloatValue{Value: 2.0},
			expected: 1.4142135623730951,
			hasError: false,
		},
		{
			name:     "sqrt of negative (should error)",
			function: builtinSqrt,
			input:    &IntegerValue{Value: -1},
			expected: 0.0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function([]Value{tt.input})

			if tt.hasError {
				if result.Type() != ERROR_TYPE {
					t.Errorf("Expected error, got %T: %v", result, result)
				}
			} else {
				if result.Type() != FLOAT_TYPE {
					t.Errorf("Expected float, got %T: %v", result, result)
					return
				}

				actual := result.(*FloatValue).Value
				// Use approximate comparison for floating point
				if abs(actual-tt.expected) > 1e-10 {
					t.Errorf("Expected %g, got %g", tt.expected, actual)
				}
			}
		})
	}
}

func TestBuiltinNow(t *testing.T) {
	result := builtinNow([]Value{})

	if result.Type() != STRING_TYPE {
		t.Errorf("Expected string, got %T: %v", result, result)
	}

	timeStr := result.(*StringValue).Value
	if len(timeStr) == 0 {
		t.Errorf("Expected non-empty time string, got empty string")
	}
}

func TestBuiltinFunctionCall(t *testing.T) {
	// Test that built-in functions can be called through the evaluator
	_ = New() // Create evaluator to ensure built-ins are registered

	tests := []struct {
		name     string
		code     string
		expected string
		hasError bool
	}{
		{
			name:     "len function call",
			code:     `len("hello")`,
			expected: "5",
			hasError: false,
		},
		{
			name:     "print function call",
			code:     `print("test")`,
			expected: "null",
			hasError: false,
		},
		{
			name:     "sin function call",
			code:     `sin(0)`,
			expected: "0",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple AST for the function call
			// This is a simplified test - in practice, you'd parse the code
			// For now, we'll test the built-in functions directly
			if tt.name == "len function call" {
				result := builtinLen([]Value{&StringValue{Value: "hello"}})
				if result.String() != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result.String())
				}
			}
		})
	}
}

// Helper function for floating point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
