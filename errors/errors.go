// errors/errors.go
package errors

import (
	"fmt"
	"strings"
)

// Error represents a compilation error with position information
type Error struct {
	Message    string
	Line       int
	Column     int
	Severity   ErrorSeverity
	Code       string
	Help       string
	SourceLine string // The actual line of source code
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

	// Source line if provided
	if e.SourceLine != "" {
		sb.WriteString(fmt.Sprintf("  %s\n", e.SourceLine))
		// Add a caret pointing to the error position
		if e.Column > 0 && e.Column <= len(e.SourceLine) {
			caret := strings.Repeat(" ", e.Column-1) + "^"
			sb.WriteString(fmt.Sprintf("  %s\n", caret))
		}
	}

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

// WithSourceLine sets the source line for better error reporting
func (e *Error) WithSourceLine(sourceLine string) *Error {
	e.SourceLine = sourceLine
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
	ErrCodeSyntaxError       = "E0001"
	ErrCodeTypeError         = "E0002"
	ErrCodeUndefinedVar      = "E0003"
	ErrCodeDuplicateDecl     = "E0004"
	ErrCodeInvalidType       = "E0005"
	ErrCodeUnsafeError       = "E0006"
	ErrCodeImmutableError    = "E0007"
	ErrCodeUndefinedField    = "E0008"
	ErrCodeUndefinedType     = "E0009"
	ErrCodeImmutable         = "E0010"
	ErrCodeParserState       = "E0011"
	ErrCodeUnexpectedToken   = "E0012"
	ErrCodeMissingToken      = "E0013"
	ErrCodeInvalidExpression = "E0014"
	ErrCodeArrayIndexError   = "E0015"
	ErrCodeFunctionCallError = "E0016"
	ErrCodeControlFlowError  = "E0017"

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

func NewParserStateError(message string, line, column int) *Error {
	return NewError(fmt.Sprintf("parser state error: %s", message), line, column).
		WithCode(ErrCodeParserState).
		WithHelp("this may be a parser bug. Try simplifying the expression or check for missing semicolons/braces")
}

func NewUnexpectedTokenError(expected, got string, line, column int) *Error {
	// Create more descriptive error messages based on context
	var message, help string

	// Convert token names to user-friendly symbols
	expectedSymbol := tokenToSymbol(expected)
	gotSymbol := tokenToSymbol(got)

	switch expected {
	case "RBRACE":
		message = fmt.Sprintf("missing closing brace '}'")
		help = "add a closing brace '}' to complete this block. Check for matching opening braces."
	case "LBRACE":
		message = fmt.Sprintf("missing opening brace '{'")
		help = "add an opening brace '{' to start this block."
	case "SEMICOLON":
		message = fmt.Sprintf("missing semicolon ';'")
		help = "add a semicolon ';' to end this statement."
	case "RPAREN":
		message = fmt.Sprintf("missing closing parenthesis ')'")
		help = "add a closing parenthesis ')' to complete this expression."
	case "LPAREN":
		message = fmt.Sprintf("missing opening parenthesis '('")
		help = "add an opening parenthesis '(' to start this expression."
	case "RBRACKET":
		message = fmt.Sprintf("missing closing bracket ']'")
		help = "add a closing bracket ']' to complete this array access."
	case "LBRACKET":
		message = fmt.Sprintf("missing opening bracket '['")
		help = "add an opening bracket '[' to start this array access."
	case "COLON":
		message = fmt.Sprintf("missing colon ':'")
		help = "add a colon ':' for type declaration or struct field."
	case "COLONEQ":
		message = fmt.Sprintf("missing ':=' for variable declaration")
		help = "use ':=' to declare and initialize a variable."
	case "EQ":
		message = fmt.Sprintf("missing assignment operator '='")
		help = "use '=' to assign a value to a variable."
	case "IDENT":
		message = fmt.Sprintf("expected identifier (variable or function name)")
		help = "provide a valid identifier name (letters, numbers, underscore)."
	case "NUMBER":
		message = fmt.Sprintf("expected number")
		help = "provide a numeric value (integer or float)."
	case "STRING":
		message = fmt.Sprintf("expected string")
		help = "provide a string value in quotes."
	case "EOF":
		message = fmt.Sprintf("unexpected end of file")
		help = "check for missing closing braces, parentheses, or semicolons."
	default:
		message = fmt.Sprintf("expected %s, got %s", expectedSymbol, gotSymbol)
		help = "check syntax around this position. Common issues: missing braces, semicolons, or incorrect operator usage"
	}

	return NewError(message, line, column).
		WithCode(ErrCodeUnexpectedToken).
		WithHelp(help)
}

// tokenToSymbol converts token names to user-friendly symbols
func tokenToSymbol(token string) string {
	switch token {
	case "RBRACE":
		return "'}'"
	case "LBRACE":
		return "'{'"
	case "SEMICOLON":
		return "';'"
	case "RPAREN":
		return "')'"
	case "LPAREN":
		return "'('"
	case "RBRACKET":
		return "']'"
	case "LBRACKET":
		return "'['"
	case "COLON":
		return "':'"
	case "COLONEQ":
		return "':='"
	case "EQ":
		return "'='"
	case "PLUS":
		return "'+'"
	case "MINUS":
		return "'-'"
	case "ASTERISK":
		return "'*'"
	case "SLASH":
		return "'/'"
	case "PERCENT":
		return "'%'"
	case "BANG":
		return "'!'"
	case "LT":
		return "'<'"
	case "GT":
		return "'>'"
	case "LTEQ":
		return "'<='"
	case "GTEQ":
		return "'>='"
	case "EQEQ":
		return "'=='"
	case "BANGEQ":
		return "'!='"
	case "AND":
		return "'&&'"
	case "OR":
		return "'||'"
	case "COMMA":
		return "','"
	case "DOT":
		return "'.'"
	case "EOF":
		return "end of file"
	case "FUNC":
		return "function keyword"
	case "RETURN":
		return "return keyword"
	case "IF":
		return "if keyword"
	case "ELSE":
		return "else keyword"
	case "FOR":
		return "for keyword"
	case "WHILE":
		return "while keyword"
	case "MUT":
		return "mut keyword"
	case "STRUCT":
		return "struct keyword"
	case "INT":
		return "int keyword"
	case "FLOAT":
		return "float keyword"
	case "STRING_KW":
		return "string keyword"
	case "BOOL":
		return "bool keyword"
	case "TRUE":
		return "true"
	case "FALSE":
		return "false"
	case "NIL":
		return "nil"
	case "IDENT":
		return "identifier"
	case "NUMBER":
		return "number"
	case "STRING":
		return "string literal"
	default:
		return token
	}
}

func NewMissingTokenError(token string, line, column int) *Error {
	return NewError(fmt.Sprintf("missing %s", token), line, column).
		WithCode(ErrCodeMissingToken).
		WithHelp(fmt.Sprintf("add the missing %s token", token))
}

func NewArrayIndexError(message string, line, column int) *Error {
	return NewError(fmt.Sprintf("array indexing error: %s", message), line, column).
		WithCode(ErrCodeArrayIndexError).
		WithHelp("ensure the array exists and the index is valid. Check for proper array declaration and bounds")
}

func NewFunctionCallError(message string, line, column int) *Error {
	return NewError(fmt.Sprintf("function call error: %s", message), line, column).
		WithCode(ErrCodeFunctionCallError).
		WithHelp("check function name, parameter count, and parameter types")
}

func NewControlFlowError(message string, line, column int) *Error {
	return NewError(fmt.Sprintf("control flow error: %s", message), line, column).
		WithCode(ErrCodeControlFlowError).
		WithHelp("check loop syntax, conditional statements, and control flow keywords")
}

// NewMissingBraceError creates a specific error for missing braces
func NewMissingBraceError(braceType string, line, column int) *Error {
	var message, help string
	if braceType == "closing" {
		message = "missing closing brace '}'"
		help = "add a closing brace '}' to complete this block. Check for matching opening braces."
	} else {
		message = "missing opening brace '{'"
		help = "add an opening brace '{' to start this block."
	}

	return NewError(message, line, column).
		WithCode(ErrCodeUnexpectedToken).
		WithHelp(help)
}

// NewMissingSemicolonError creates a specific error for missing semicolons
func NewMissingSemicolonError(line, column int) *Error {
	return NewError("missing semicolon ';'", line, column).
		WithCode(ErrCodeUnexpectedToken).
		WithHelp("add a semicolon ';' to end this statement.")
}

// NewUnexpectedEndOfFileError creates a specific error for unexpected EOF
func NewUnexpectedEndOfFileError(line, column int) *Error {
	return NewError("unexpected end of file", line, column).
		WithCode(ErrCodeUnexpectedToken).
		WithHelp("check for missing closing braces, parentheses, or semicolons.")
}

// NewInvalidSyntaxError creates a specific error for invalid syntax
func NewInvalidSyntaxError(context string, line, column int) *Error {
	return NewError(fmt.Sprintf("invalid syntax in %s", context), line, column).
		WithCode(ErrCodeSyntaxError).
		WithHelp("check the syntax rules for this construct.")
}
