// errors/errors.go
package errors

import (
	"fmt"
	"strings"
)

// Error represents a compilation error with position information
type Error struct {
	Message  string
	Line     int
	Column   int
	Severity ErrorSeverity
	Code     string
	Help     string
}

// ErrorSeverity represents the severity level of an error
type ErrorSeverity int

const (
	ErrorSeverityError ErrorSeverity = iota
	ErrorSeverityWarning
	ErrorSeverityInfo
)

func (s ErrorSeverity) String() string {
	switch s {
	case ErrorSeverityError:
		return "error"
	case ErrorSeverityWarning:
		return "warning"
	case ErrorSeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.String()
}

// String returns a formatted error message with position information
func (e *Error) String() string {
	var sb strings.Builder

	// Error header with severity and code
	sb.WriteString(fmt.Sprintf("%s[%s]: %s\n", e.Severity, e.Code, e.Message))

	// Position information
	sb.WriteString(fmt.Sprintf("  --> line %d, column %d\n", e.Line, e.Column))

	// Help message if provided
	if e.Help != "" {
		sb.WriteString(fmt.Sprintf("  help: %s\n", e.Help))
	}

	return sb.String()
}

// NewError creates a new error with the given information
func NewError(message string, line, column int) *Error {
	return &Error{
		Message:  message,
		Line:     line,
		Column:   column,
		Severity: ErrorSeverityError,
		Code:     "E0001",
	}
}

// NewWarning creates a new warning with the given information
func NewWarning(message string, line, column int) *Error {
	return &Error{
		Message:  message,
		Line:     line,
		Column:   column,
		Severity: ErrorSeverityWarning,
		Code:     "W0001",
	}
}

// WithCode sets the error code
func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

// WithHelp sets the help message
func (e *Error) WithHelp(help string) *Error {
	e.Help = help
	return e
}

// WithSeverity sets the error severity
func (e *Error) WithSeverity(severity ErrorSeverity) *Error {
	e.Severity = severity
	return e
}

// ErrorList represents a collection of errors
type ErrorList struct {
	errors []*Error
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		errors: make([]*Error, 0),
	}
}

// Add adds an error to the list
func (el *ErrorList) Add(err *Error) {
	el.errors = append(el.errors, err)
}

// AddError adds a simple error message
func (el *ErrorList) AddError(message string, line, column int) {
	el.Add(NewError(message, line, column))
}

// AddWarning adds a warning message
func (el *ErrorList) AddWarning(message string, line, column int) {
	el.Add(NewWarning(message, line, column))
}

// Errors returns all errors in the list
func (el *ErrorList) Errors() []*Error {
	return el.errors
}

// HasErrors returns true if there are any errors
func (el *ErrorList) HasErrors() bool {
	return len(el.errors) > 0
}

// HasWarnings returns true if there are any warnings
func (el *ErrorList) HasWarnings() bool {
	for _, err := range el.errors {
		if err.Severity == ErrorSeverityWarning {
			return true
		}
	}
	return false
}

// String returns a formatted string representation of all errors
func (el *ErrorList) String() string {
	if len(el.errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range el.errors {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.String())
	}
	return sb.String()
}

// Error implements the error interface
func (el *ErrorList) Error() string {
	return el.String()
}

// Common error codes
const (
	ErrCodeSyntaxError    = "E0001"
	ErrCodeTypeError      = "E0002"
	ErrCodeUndefinedVar   = "E0003"
	ErrCodeDuplicateDecl  = "E0004"
	ErrCodeInvalidType    = "E0005"
	ErrCodeUnsafeError    = "E0006"
	ErrCodeImmutableError = "E0007"
	ErrCodeUndefinedField = "E0008"
	ErrCodeUndefinedType  = "E0009"
	ErrCodeImmutable      = "E0010"

	WarnCodeUnusedVar    = "W0001"
	WarnCodeUnusedImport = "W0002"
	WarnCodeDeprecated   = "W0003"
)

// Common error constructors
func NewSyntaxError(message string, line, column int) *Error {
	return NewError(message, line, column).WithCode(ErrCodeSyntaxError)
}

func NewTypeError(message string, line, column int) *Error {
	return NewError(message, line, column).WithCode(ErrCodeTypeError)
}

func NewUndefinedVarError(varName string, line, column int) *Error {
	return NewError(fmt.Sprintf("undefined variable '%s'", varName), line, column).
		WithCode(ErrCodeUndefinedVar).
		WithHelp("declare the variable before using it")
}

func NewDuplicateDeclError(name string, line, column int) *Error {
	return NewError(fmt.Sprintf("duplicate declaration of '%s'", name), line, column).
		WithCode(ErrCodeDuplicateDecl).
		WithHelp("each variable/function can only be declared once in the same scope")
}

func NewImmutableError(varName string, line, column int) *Error {
	return NewError(fmt.Sprintf("cannot assign to immutable variable '%s'", varName), line, column).
		WithCode(ErrCodeImmutableError).
		WithHelp("use 'mut' keyword to declare a mutable variable")
}
