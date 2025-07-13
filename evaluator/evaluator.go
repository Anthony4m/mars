package evaluator

import (
	"fmt"
	"mars/ast"
)

type Evaluator struct {
	env *Environment
}

func New() *Evaluator {
	return &Evaluator{env: NewEnvironment()}
}

func (e *Evaluator) Eval(node ast.Node) Value {
	fmt.Printf("Evaluating: %T\n", node)
	switch n := node.(type) {
	case *ast.Program:
		// Evaluate the last declaration/statement and return its value
		if len(n.Declarations) == 0 {
			return nil
		}
		return e.Eval(n.Declarations[len(n.Declarations)-1])
	case *ast.ExpressionStatement:
		return e.Eval(n.Expression)
	case *ast.Literal:
		return e.evalLiteral(n)
	case ast.Expression:
		return e.Eval(n.(ast.Node))
	default:
		return nil
	}
}

func (e *Evaluator) evalLiteral(lit *ast.Literal) Value {
	switch v := lit.Value.(type) {
	case int64:
		return &IntegerValue{Value: v}
	case int:
		return &IntegerValue{Value: int64(v)}
	case string:
		return &StringValue{Value: v}
	case bool:
		// Use singleton values to save allocations
		if v {
			return TRUE
		}
		return FALSE
	case float64:
		return &FloatValue{Value: v}
	default:
		return newError("unknown literal type: %T", lit.Value)
	}
}

func isError(obj Value) bool {
	if obj != nil {
		return obj.Type() == ERROR_TYPE
	}
	return false
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}
