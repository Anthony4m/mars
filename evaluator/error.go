package evaluator

import (
	"fmt"
	"mars/ast"
	"strings"
)

type ErrorDetail struct {
	Message   string
	Location  ast.Position
	Hint      string
	ErrorCode string
}

type RuntimeError struct {
	Detail     ErrorDetail
	StackTrace []StackFrame
}

func (e *RuntimeError) Type() string   { return ERROR_TYPE }
func (e *RuntimeError) String() string { return e.formatError() }
func (e *RuntimeError) IsTruthy() bool { return false }

func (e *RuntimeError) formatError() string {
	var sb strings.Builder

	//Error header
	sb.WriteString(fmt.Sprintf("\033[31merror[%s]\033[0m: %s\n", e.Detail.ErrorCode, e.Detail.Message))

	//Location if available
	if e.Detail.Location.Line > 0 {
		sb.WriteString(fmt.Sprintf("  \\033[34m-->\\033[0m %d:%d\\n", e.Detail.Location.Line, e.Detail.Location.Column))
	}

	//Hint if provided
	if e.Detail.Hint != "" {
		sb.WriteString(fmt.Sprintf("  \033[32mhint:\033[0m %s\n", e.Detail.Hint))
	}

	// stack trace
	if len(e.StackTrace) > 0 {
		sb.WriteString("\n\033[34mstack trace:\033[0m\n")
		for _, frame := range e.StackTrace {
			sb.WriteString(fmt.Sprintf("  at %s (%d:%d)\n",
				frame.Function, frame.Location.Line, frame.Location.Column))
		}
	}

	return sb.String()
}
