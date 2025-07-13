// errors/reporter.go
package errors

import (
	"fmt"
	"mars/ast"
	"strings"
)

// DiagnosticError represents a rich error with source context
type DiagnosticError struct {
	*Error     // Embed your existing Error
	SourceCode string
	EndPos     ast.Position // For multi-token errors
}

// MarsReporter provides beautiful error formatting
type MarsReporter struct {
	sourceCode string
	filename   string
	errors     []*DiagnosticError
}

func NewMarsReporter(sourceCode, filename string) *MarsReporter {
	return &MarsReporter{
		sourceCode: sourceCode,
		filename:   filename,
		errors:     []*DiagnosticError{},
	}
}

// AddError creates a diagnostic error from position
func (cr *MarsReporter) AddError(pos ast.Position, code, message string) {
	cr.errors = append(cr.errors, &DiagnosticError{
		Error: &Error{
			Message:  message,
			Line:     pos.Line,
			Column:   pos.Column,
			Severity: ErrorSeverityError,
			Code:     code,
		},
		SourceCode: cr.sourceCode,
	})
}

// AddErrorWithHelp adds an error with help text
func (cr *MarsReporter) AddErrorWithHelp(pos ast.Position, code, message, help string) {
	cr.errors = append(cr.errors, &DiagnosticError{
		Error: &Error{
			Message:  message,
			Line:     pos.Line,
			Column:   pos.Column,
			Severity: ErrorSeverityError,
			Code:     code,
			Help:     help,
		},
		SourceCode: cr.sourceCode,
	})
}

// AddErrorWithSpan adds an error spanning multiple tokens
func (cr *MarsReporter) AddErrorWithSpan(startPos, endPos ast.Position, code, message, help string) {
	cr.errors = append(cr.errors, &DiagnosticError{
		Error: &Error{
			Message:  message,
			Line:     startPos.Line,
			Column:   startPos.Column,
			Severity: ErrorSeverityError,
			Code:     code,
			Help:     help,
		},
		SourceCode: cr.sourceCode,
		EndPos:     endPos,
	})
}

// HasErrors returns true if there are any errors
func (cr *MarsReporter) HasErrors() bool {
	for _, e := range cr.errors {
		if e.Severity == ErrorSeverityError {
			return true
		}
	}
	return false
}

// String formats all errors in detailed style
func (cr *MarsReporter) String() string {
	if len(cr.errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range cr.errors {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(err.MarsFormat(cr.filename))
	}

	// Summary
	errorCount := 0
	warningCount := 0
	for _, e := range cr.errors {
		switch e.Severity {
		case ErrorSeverityError:
			errorCount++
		case ErrorSeverityWarning:
			warningCount++
		}
	}

	if errorCount > 0 || warningCount > 0 {
		sb.WriteString("\n\n")
		if errorCount > 0 {
			sb.WriteString(fmt.Sprintf("\033[31merror\033[0m: aborting due to %d previous error(s)", errorCount))
		}
		if warningCount > 0 {
			if errorCount > 0 {
				sb.WriteString("; ")
			}
			sb.WriteString(fmt.Sprintf("%d warning(s) emitted", warningCount))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// MarsFormat formats a single error in detailed style
func (d *DiagnosticError) MarsFormat(filename string) string {
	var sb strings.Builder

	// Color codes
	colors := map[ErrorSeverity]string{
		ErrorSeverityError:   "\033[31m", // red
		ErrorSeverityWarning: "\033[33m", // yellow
		ErrorSeverityInfo:    "\033[36m", // cyan
	}

	// Error header
	color := colors[d.Severity]
	sb.WriteString(fmt.Sprintf("%s%s[%s]\033[0m: %s\n",
		color, d.Severity, d.Code, d.Message))
	sb.WriteString(fmt.Sprintf(" \033[34m-->\033[0m %s:%d:%d\n",
		filename, d.Line, d.Column))

	// Source context
	lines := strings.Split(d.SourceCode, "\n")
	if d.Line > 0 && d.Line <= len(lines) {
		// Calculate line number width for alignment
		maxLine := d.Line + 1
		if maxLine > len(lines) {
			maxLine = len(lines)
		}
		lineNumWidth := len(fmt.Sprintf("%d", maxLine))

		sb.WriteString(fmt.Sprintf("%s \033[34m|\033[0m\n", strings.Repeat(" ", lineNumWidth)))

		// Show previous line for context
		if d.Line > 1 {
			sb.WriteString(fmt.Sprintf(" \033[34m%*d |\033[0m %s\n",
				lineNumWidth, d.Line-1, lines[d.Line-2]))
		}

		// Error line
		sb.WriteString(fmt.Sprintf(" \033[34m%*d |\033[0m %s\n",
			lineNumWidth, d.Line, lines[d.Line-1]))

		// Error pointer
		sb.WriteString(fmt.Sprintf("%s \033[34m|\033[0m ", strings.Repeat(" ", lineNumWidth)))

		// Calculate underline
		underlineStart := d.Column - 1
		underlineLen := 1
		if d.EndPos.Line == d.Line && d.EndPos.Column > d.Column {
			underlineLen = d.EndPos.Column - d.Column
		}

		sb.WriteString(strings.Repeat(" ", underlineStart))
		sb.WriteString(fmt.Sprintf("%s%s\033[0m", color, strings.Repeat("^", underlineLen)))

		// Inline help if short
		if d.Help != "" && len(d.Help) < 40 {
			sb.WriteString(fmt.Sprintf(" %s\033[0m", d.Help))
		}
		sb.WriteString("\n")

		// Show next line for context
		if d.Line < len(lines) {
			sb.WriteString(fmt.Sprintf(" \033[34m%*d |\033[0m %s\n",
				lineNumWidth, d.Line+1, lines[d.Line]))
		}

		sb.WriteString(fmt.Sprintf("%s \033[34m|\033[0m\n", strings.Repeat(" ", lineNumWidth)))
	}

	// Help text (if not shown inline)
	if d.Help != "" && len(d.Help) >= 40 {
		sb.WriteString(fmt.Sprintf(" \033[34m=\033[0m \033[1mhelp\033[0m: %s\n", d.Help))
	}

	return sb.String()
}
