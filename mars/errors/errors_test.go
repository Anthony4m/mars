// errors/errors_test.go
package errors

import (
	"testing"
)

func TestErrorFormatting(t *testing.T) {
	err := NewSyntaxError("unexpected token '}'", 5, 12)

	expected := "error[E0001]: unexpected token '}'\n  --> line 5, column 12\n"
	if err.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, err.String())
	}
}

func TestErrorWithHelp(t *testing.T) {
	err := NewUndefinedVarError("x", 3, 8)

	expected := "error[E0003]: undefined variable 'x'\n  --> line 3, column 8\n  help: declare the variable before using it\n"
	if err.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, err.String())
	}
}

func TestWarningFormatting(t *testing.T) {
	warn := NewWarning("unused variable 'y'", 7, 15)

	expected := "warning[W0001]: unused variable 'y'\n  --> line 7, column 15\n"
	if warn.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, warn.String())
	}
}

func TestErrorList(t *testing.T) {
	el := NewErrorList()

	el.AddError("syntax error", 1, 5)
	el.AddWarning("unused variable", 2, 10)
	el.Add(NewImmutableError("x", 3, 15))

	if !el.HasErrors() {
		t.Error("Expected error list to have errors")
	}

	if !el.HasWarnings() {
		t.Error("Expected error list to have warnings")
	}

	if len(el.Errors()) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(el.Errors()))
	}

	// Test that the error list string contains all errors
	result := el.String()
	if len(result) == 0 {
		t.Error("Expected non-empty error list string")
	}

	// Should contain all three error messages
	expectedStrings := []string{
		"syntax error",
		"unused variable",
		"cannot assign to immutable variable 'x'",
	}

	for _, expected := range expectedStrings {
		if !contains(result, expected) {
			t.Errorf("Expected error list to contain '%s'", expected)
		}
	}
}

func TestErrorSeverity(t *testing.T) {
	err := NewError("test error", 1, 1)
	if err.Severity != ErrorSeverityError {
		t.Errorf("Expected ErrorSeverityError, got %v", err.Severity)
	}

	warn := NewWarning("test warning", 1, 1)
	if warn.Severity != ErrorSeverityWarning {
		t.Errorf("Expected ErrorSeverityWarning, got %v", warn.Severity)
	}
}

func TestErrorChaining(t *testing.T) {
	err := NewError("base error", 1, 1).
		WithCode("CUSTOM001").
		WithHelp("custom help message").
		WithSeverity(ErrorSeverityWarning)

	if err.Code != "CUSTOM001" {
		t.Errorf("Expected code CUSTOM001, got %s", err.Code)
	}

	if err.Help != "custom help message" {
		t.Errorf("Expected help message, got %s", err.Help)
	}

	if err.Severity != ErrorSeverityWarning {
		t.Errorf("Expected warning severity, got %v", err.Severity)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
