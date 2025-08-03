package evaluator

import "mars/ast"

type StackFrame struct {
	Function string
	Location ast.Position
	Context  string // "function call", "if statement", etc.
}
