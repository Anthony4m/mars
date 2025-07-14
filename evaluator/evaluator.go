package evaluator

import (
	"fmt"
	"mars/ast"
)

type Evaluator struct {
	env *Environment
}

type binaryOpFn func(left, right Value) Value

var binaryOps = map[string]binaryOpFn{
	"+":  add,
	"-":  subtract,
	"*":  multiply,
	"/":  divide,
	"==": equal,
	"!=": notEqual,
	"<":  lessThan,
	">":  greaterThan,
	"<=": lessThanOrEqual,
	">=": greaterThanOrEqual,
}

// Handler for the '+' operator
func add(left, right Value) Value {
	// Handle integer addition
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*IntegerValue).Value
		rv := right.(*IntegerValue).Value
		return &IntegerValue{Value: lv + rv}
	}
	// Handle float addition
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		lv := left.(*FloatValue).Value
		rv := right.(*FloatValue).Value
		return &FloatValue{Value: lv + rv}
	}
	// Handle float addition
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*FloatValue).Value
		rv := right.(*IntegerValue).Value
		return &FloatValue{Value: lv + float64(rv)}
	}
	// Handle float addition
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		lv := left.(*IntegerValue).Value
		rv := right.(*FloatValue).Value
		return &FloatValue{Value: float64(lv) + rv}
	}
	// Handle string concatenation
	if left.Type() == STRING_TYPE && right.Type() == STRING_TYPE {
		lv := left.(*StringValue).Value
		rv := right.(*StringValue).Value
		return &StringValue{Value: lv + rv}
	}
	// If no rule matches, return an error
	return newError("type mismatch: cannot add %s and %s", left.Type(), right.Type())
}

// Handler for the '-' operator
func subtract(left, right Value) Value {
	// Handle integer subtraction
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*IntegerValue).Value
		rv := right.(*IntegerValue).Value
		return &IntegerValue{Value: lv - rv}
	}
	// Handle float subtraction
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		lv := left.(*FloatValue).Value
		rv := right.(*FloatValue).Value
		return &FloatValue{Value: lv - rv}
	}
	// Handle mixed int/float subtraction
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*FloatValue).Value
		rv := float64(right.(*IntegerValue).Value)
		return &FloatValue{Value: lv - rv}
	}
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		lv := float64(left.(*IntegerValue).Value)
		rv := right.(*FloatValue).Value
		return &FloatValue{Value: lv - rv}
	}
	return newError("type mismatch: cannot subtract %s from %s", right.Type(), left.Type())
}

// Handler for the '*' operator
func multiply(left, right Value) Value {
	// Handle integer multiplication
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*IntegerValue).Value
		rv := right.(*IntegerValue).Value
		return &IntegerValue{Value: lv * rv}
	}
	// Handle float multiplication
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		lv := left.(*FloatValue).Value
		rv := right.(*FloatValue).Value
		return &FloatValue{Value: lv * rv}
	}
	// Handle mixed int/float multiplication
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*FloatValue).Value
		rv := float64(right.(*IntegerValue).Value)
		return &FloatValue{Value: lv * rv}
	}
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		lv := float64(left.(*IntegerValue).Value)
		rv := right.(*FloatValue).Value
		return &FloatValue{Value: lv * rv}
	}
	return newError("type mismatch: cannot multiply %s and %s", left.Type(), right.Type())
}

// Handler for the '/' operator
func divide(left, right Value) Value {
	// Handle integer division
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*IntegerValue).Value
		rv := right.(*IntegerValue).Value
		if rv == 0 {
			return newError("division by zero")
		}
		return &IntegerValue{Value: lv / rv}
	}
	// Handle float division
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		lv := left.(*FloatValue).Value
		rv := right.(*FloatValue).Value
		if rv == 0.0 {
			return newError("division by zero")
		}
		return &FloatValue{Value: lv / rv}
	}
	// Handle mixed int/float division
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*FloatValue).Value
		rv := float64(right.(*IntegerValue).Value)
		if rv == 0.0 {
			return newError("division by zero")
		}
		return &FloatValue{Value: lv / rv}
	}
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		lv := float64(left.(*IntegerValue).Value)
		rv := right.(*FloatValue).Value
		if rv == 0.0 {
			return newError("division by zero")
		}
		return &FloatValue{Value: lv / rv}
	}
	return newError("type mismatch: cannot divide %s by %s", left.Type(), right.Type())
}

// Handler for the '==' operator
func equal(left, right Value) Value {
	if left.Type() != right.Type() {
		return FALSE
	}
	switch left.Type() {
	case INTEGER_TYPE:
		return boolToValue(left.(*IntegerValue).Value == right.(*IntegerValue).Value)
	case FLOAT_TYPE:
		return boolToValue(left.(*FloatValue).Value == right.(*FloatValue).Value)
	case STRING_TYPE:
		return boolToValue(left.(*StringValue).Value == right.(*StringValue).Value)
	case BOOLEAN_TYPE:
		return boolToValue(left.(*BooleanValue).Value == right.(*BooleanValue).Value)
	case NULL_TYPE:
		return TRUE // null == null
	}
	return FALSE
}

// Handler for the '!=' operator
func notEqual(left, right Value) Value {
	result := equal(left, right)
	if isError(result) {
		return result
	}
	return boolToValue(!result.IsTruthy())
}

// Handler for the '<' operator
func lessThan(left, right Value) Value {
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		return boolToValue(left.(*IntegerValue).Value < right.(*IntegerValue).Value)
	}
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		return boolToValue(left.(*FloatValue).Value < right.(*FloatValue).Value)
	}
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		return boolToValue(left.(*FloatValue).Value < float64(right.(*IntegerValue).Value))
	}
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		return boolToValue(float64(left.(*IntegerValue).Value) < right.(*FloatValue).Value)
	}
	return newError("type mismatch: cannot compare %s < %s", left.Type(), right.Type())
}

// Handler for the '>' operator
func greaterThan(left, right Value) Value {
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		return boolToValue(left.(*IntegerValue).Value > right.(*IntegerValue).Value)
	}
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		return boolToValue(left.(*FloatValue).Value > right.(*FloatValue).Value)
	}
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		return boolToValue(left.(*FloatValue).Value > float64(right.(*IntegerValue).Value))
	}
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		return boolToValue(float64(left.(*IntegerValue).Value) > right.(*FloatValue).Value)
	}
	return newError("type mismatch: cannot compare %s > %s", left.Type(), right.Type())
}

// Handler for the '<=' operator
func lessThanOrEqual(left, right Value) Value {
	lt := lessThan(left, right)
	if isError(lt) {
		return lt
	}
	eq := equal(left, right)
	if isError(eq) {
		return eq
	}
	return boolToValue(lt.IsTruthy() || eq.IsTruthy())
}

// Handler for the '>=' operator
func greaterThanOrEqual(left, right Value) Value {
	gt := greaterThan(left, right)
	if isError(gt) {
		return gt
	}
	eq := equal(left, right)
	if isError(eq) {
		return eq
	}
	return boolToValue(gt.IsTruthy() || eq.IsTruthy())
}

// Helper to convert bool to Value
func boolToValue(b bool) Value {
	if b {
		return TRUE
	}
	return FALSE
}

func New() *Evaluator {
	return &Evaluator{env: NewEnvironment()}
}

func (e *Evaluator) Eval(node ast.Node) Value {
	fmt.Printf("Evaluating: %T\n", node)
	switch n := node.(type) {
	case *ast.Program:
		var result Value
		for _, decl := range n.Declarations {
			result = e.Eval(decl)
			if isError(result) {
				return result // Stop on first error
			}
		}
		return result // Return last non-error value
	case *ast.ExpressionStatement:
		return e.Eval(n.Expression)
	case *ast.Literal:
		return e.evalLiteral(n)
	case *ast.BinaryExpression:
		left := e.Eval(n.Left)
		right := e.Eval(n.Right)

		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}
		if handler, ok := binaryOps[n.Operator]; ok {
			// 3. If found, call the handler function.
			return handler(left, right)
		}
		// 4. If the operator doesn't exist in the map, it's an unknown operator.
		return newError("unknown operator: %s %s %s", left.Type(), n.Operator, right.Type())
	case *ast.UnaryExpression:
		right := e.Eval(n.Right)
		if isError(right) {
			return right
		}
		return e.evalUnary(n.Operator, right)
	case ast.Expression:
		return e.Eval(n.(ast.Node))
	default:
		return nil
	}
}

func (e *Evaluator) evalUnary(operator string, right Value) Value {
	switch operator {
	case "!":
		return boolToValue(!right.IsTruthy())
	case "-":
		if right.Type() == INTEGER_TYPE {
			return &IntegerValue{Value: -right.(*IntegerValue).Value}
		}
		if right.Type() == FLOAT_TYPE {
			return &FloatValue{Value: -right.(*FloatValue).Value}
		}
		return newError("unknown operator: %s%s", operator, right.Type())
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
