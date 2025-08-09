package evaluator

import (
	"fmt"
	"mars/ast"
	"strings"
)

// Error codes
const (
	ErrTypeMismatch   = "E001"
	ErrUndefinedVar   = "E002"
	ErrDivisionByZero = "E003"
	ErrNotAFunction   = "E004"
	ErrWrongArgCount  = "E005"
	ErrSyntaxError    = "E006"
	ErrImmutable      = "E007"
	ErrUndefined      = "E008"
	ErrRuntimeError   = "E009"
)

type Evaluator struct {
	env        *Environment
	callStack  []StackFrame
	sourceCode string // For showing code snippets
}

type binaryOpFn func(left, right Value) Value

var binaryOps = map[string]binaryOpFn{
	"+":  add,
	"-":  subtract,
	"*":  multiply,
	"/":  divide,
	"%":  modulo,
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

// Handler for the '%' operator
func modulo(left, right Value) Value {
	// Handle integer modulo
	if left.Type() == INTEGER_TYPE && right.Type() == INTEGER_TYPE {
		lv := left.(*IntegerValue).Value
		rv := right.(*IntegerValue).Value
		if rv == 0 {
			return newError("modulo by zero")
		}
		return &IntegerValue{Value: lv % rv}
	}
	// Handle float modulo (convert to int64 for now)
	if left.Type() == FLOAT_TYPE && right.Type() == FLOAT_TYPE {
		lv := int64(left.(*FloatValue).Value)
		rv := int64(right.(*FloatValue).Value)
		if rv == 0 {
			return newError("modulo by zero")
		}
		return &IntegerValue{Value: lv % rv}
	}
	// Handle mixed int/float modulo
	if left.Type() == FLOAT_TYPE && right.Type() == INTEGER_TYPE {
		lv := int64(left.(*FloatValue).Value)
		rv := right.(*IntegerValue).Value
		if rv == 0 {
			return newError("modulo by zero")
		}
		return &IntegerValue{Value: lv % rv}
	}
	if left.Type() == INTEGER_TYPE && right.Type() == FLOAT_TYPE {
		lv := left.(*IntegerValue).Value
		rv := int64(right.(*FloatValue).Value)
		if rv == 0 {
			return newError("modulo by zero")
		}
		return &IntegerValue{Value: lv % rv}
	}
	return newError("type mismatch: cannot compute modulo of %s and %s", left.Type(), right.Type())
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

func (e *Evaluator) pushFrame(name string, pos ast.Position, context string) {
	frame := StackFrame{
		Function: name,
		Location: pos,
		Context:  context,
	}
	e.callStack = append(e.callStack, frame)
}

func (e *Evaluator) popFrame() {
	if len(e.callStack) > 0 {
		e.callStack = e.callStack[:len(e.callStack)-1]
	}
}

func (e *Evaluator) captureStackTrace() []StackFrame {
	// Return a copy to avoid mutations
	trace := make([]StackFrame, len(e.callStack))
	copy(trace, e.callStack)
	return trace
}

// Update your newError function
func (e *Evaluator) newError(pos ast.Position, code, format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Detail: ErrorDetail{
			Message:   fmt.Sprintf(format, args...),
			Location:  pos,
			ErrorCode: code,
		},
		StackTrace: e.captureStackTrace(),
	}
}

// Helper for type mismatch with hint
func (e *Evaluator) typeMismatchError(pos ast.Position, op string, left, right Value) *RuntimeError {
	err := e.newError(pos, ErrTypeMismatch,
		"type mismatch: cannot %s %s and %s", op, left.Type(), right.Type())

	// Add helpful hints
	if op == "+" && left.Type() == STRING_TYPE && right.Type() == INTEGER_TYPE {
		err.Detail.Hint = "convert the integer to string first"
	}

	return err
}

func New() *Evaluator {
	evaluator := &Evaluator{env: NewEnvironment()}

	// Register builtin functions
	for name, builtin := range BuiltinFunctions {
		// Create a FunctionValue for the builtin function
		function := &FunctionValue{
			Name:       name,
			Parameters: []*ast.Parameter{}, // Builtins handle their own parameter validation
			Body:       nil,                // Builtins don't have AST bodies
			ReturnType: nil,                // Builtins can return different types
			Env:        evaluator.env,
			Position:   ast.Position{Line: 0, Column: 0},
			IsBuiltin:  true,
			BuiltinFn:  builtin.Function,
		}

		// Store the builtin function in the environment
		evaluator.env.Set(name, function, false)
	}

	return evaluator
}

func (e *Evaluator) Eval(node ast.Node) Value {
	//fmt.Printf("Evaluating: %T\n", node)
	switch n := node.(type) {
	case *ast.Program:
		e.pushFrame("main", n.Position, "program")
		defer e.popFrame()

		var result Value
		for _, decl := range n.Declarations {
			result = e.Eval(decl)
			if isError(result) {
				return result
			}
			if ret, ok := result.(*ReturnValue); ok {
				return ret.Value
			}
		}
		return result
	case *ast.ExpressionStatement:
		return e.Eval(n.Expression)
	case *ast.Literal:
		return e.evalLiteral(n)
	case *ast.BinaryExpression:
		// TODO: Refactor this into its own function, possibly extract the sort circuiting part
		if n.Operator == "&&" {
			left := e.Eval(n.Left)
			if isError(left) {
				return left
			}
			if !left.IsTruthy() {
				return FALSE
			}
			right := e.Eval(n.Right)
			if isError(right) {
				return right
			}
			return boolToValue(right.IsTruthy())
		}
		if n.Operator == "||" {
			left := e.Eval(n.Left)
			if isError(left) {
				return left
			}
			if left.IsTruthy() {
				return TRUE
			}
			right := e.Eval(n.Right)
			if isError(right) {
				return right
			}
			return boolToValue(right.IsTruthy())
		}
		left := e.Eval(n.Left)
		right := e.Eval(n.Right)

		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}
		if handler, ok := binaryOps[n.Operator]; ok {
			result := handler(left, right)
			// Convert old errors to new format
			if err, ok := result.(*Error); ok {
				return e.newError(n.Position, ErrTypeMismatch, err.Message)
			}
			return result
		}

		return e.newError(n.Position, ErrTypeMismatch,
			"unknown operator: %s %s %s", left.Type(), n.Operator, right.Type())
	case *ast.UnaryExpression:
		right := e.Eval(n.Right)
		if isError(right) {
			return right
		}
		return e.evalUnary(n.Operator, n.Position, right)
	case *ast.VarDecl:
		return e.EvalVariableDecl(n)
	case *ast.AssignmentStatement:
		return e.EvalAssignment(n)
	case *ast.IndexAssignmentStatement:
		return e.EvalIndexAssignment(n)
	case *ast.IfStatement:
		return e.EvalConditional(n)
	case *ast.ForStatement:
		return e.EvalForStatement(n)
	case *ast.WhileStatement:
		return e.EvalWhileStatement(n)
	case *ast.BlockStatement:
		e.pushFrame("main", n.Position, "block")
		defer e.popFrame()
		return e.evalBlock(n)
	case *ast.Identifier:
		e.pushFrame("main", n.Position, "program")
		defer e.popFrame()
		return e.EvalIdentifier(n)
	case *ast.ReturnStatement:
		return e.evalReturnStatement(n)
	case *ast.ContinueStatement:
		return e.evalContinueStatement(n)
	case *ast.BreakStatement:
		return e.evalBreak(n)
	case *ast.FuncDecl:
		return e.evalFunctionDecl(n)
	case *ast.FunctionCall:
		return e.evalFunctionCall(n)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(n)
	case *ast.StructLiteral:
		return e.evalStructLiteral(n)
	case *ast.MemberExpression:
		return e.evalMemberExpression(n)
	case *ast.IndexExpression:
		return e.evalIndexExpression(n)
	case *ast.SliceExpression:
		return e.evalSliceExpression(n)
	case *ast.PrintStatement:
		if n.Expression == nil {
			fmt.Println("null")
			return NULL
		}

		// Evaluate the expression to print
		value := e.Eval(n.Expression)
		if isError(value) {
			return value
		}

		// Print the value
		fmt.Println(formatValueForOutput(value))
		return NULL
	case ast.Expression:
		return e.Eval(n.(ast.Node))
	default:
		return nil
	}
}

func (e *Evaluator) evalReturnStatement(n *ast.ReturnStatement) Value {
	if n.Value == nil {
		return NULL
	}
	var value Value
	value = e.Eval(n.Value)
	if isError(value) {
		return value
	}
	return &ReturnValue{Value: value}
}

func (e *Evaluator) evalUnary(operator string, position ast.Position, right Value) Value {
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
		return e.newError(position, ErrTypeMismatch, "unknown operator: %s%s", operator, right.Type())
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

func (e *Evaluator) evalFunctionDecl(n *ast.FuncDecl) Value {
	// Check for identifier
	if n.Name == nil {
		return e.newError(n.Position, ErrSyntaxError, "function declaration missing name")
	}
	e.pushFrame(n.Name.Name, n.Position, n.Signature.String())
	defer e.popFrame()

	// Create a function value that encapsulates the function definition
	function := &FunctionValue{
		Name:       n.Name.Name,
		Parameters: n.Signature.Parameters,
		Body:       n.Body,
		ReturnType: n.Signature.ReturnType,
		Env:        e.env, // Capture current environment for closures
		Position:   n.Position,
	}

	const functionIsMutable = false
	// Store the function in the environment
	e.env.Set(n.Name.Name, function, functionIsMutable)

	// Function declarations typically return nil/void
	// or the function itself for REPL convenience
	return function
}

func (e *Evaluator) EvalVariableDecl(n *ast.VarDecl) Value {
	// Check for identifier if not return error
	if n.Name == nil {
		return e.newError(n.Position, ErrSyntaxError, "variable declaration missing name")
	}

	var value Value

	// Case 1: Has initializing value
	if n.Value != nil {
		value = e.Eval(n.Value)
		if isError(value) {
			return value
		}

		// If type is specified, check compatibility
		if n.Type != nil {
			expectedType := n.Type.BaseType
			actualType := getValueType(value)
			if !e.TypesCompatible(expectedType, actualType) {
				return e.newError(n.Position, ErrTypeMismatch, "type mismatch: cannot assign %s to %s",
					actualType, expectedType)
			}
		}
	} else if n.Type != nil {
		// Case 2: Only type, no value (x: int)
		// Initialize with zero value
		value = e.initializeToZero(n.Type)
	} else {
		// Case 3: Neither type nor value - error!
		return e.newError(n.Position, ErrSyntaxError,
			"variable '%s' needs type or initial value", n.Name.Name)
	}

	// Store in environment
	e.env.Set(n.Name.Name, value, n.Mutable)

	// Variable declarations typically return nil/void
	// or the value for REPL convenience
	return value
}

// For Assignment (x = 50)
func (e *Evaluator) EvalAssignment(n *ast.AssignmentStatement) Value {
	// Validate AST structure
	if n.Name == nil {
		return e.newError(n.Position, ErrSyntaxError, "assignment missing variable name")
	}

	if n.Value == nil {
		return e.newError(n.Position, ErrSyntaxError, "assignment missing value")
	}

	// Evaluate the value to be assigned
	value := e.Eval(n.Value)
	if isError(value) {
		return value
	}

	// Check if variable exists (optional but good for clarity)
	bind, exist := e.env.Get(n.Name.Name)
	if !exist {
		return e.newError(n.Position, ErrUndefined, "undefined variable '%s'", n.Name.Name)
	}

	// Check if variable is mutable (if your environment supports this)
	if !bind.IsMutable {
		return e.newError(n.Position, ErrImmutable,
			"cannot assign to immutable variable '%s'", n.Name.Name)
	}

	if !bind.IsMutable {
		return e.newError(n.Position, ErrImmutable,
			"cannot assign to immutable variable '%s'", n.Name.Name)
	}

	valueType := getValueType(value)
	varType := getValueType(bind.Value)
	if !e.TypesCompatible(varType, valueType) {
		return e.newError(n.Position, ErrTypeMismatch,
			"type mismatch: cannot assign %s to %s", valueType, varType)
	}

	// Perform the assignment
	err := e.env.Update(n.Name.Name, value)
	if err != nil {
		return e.newError(n.Position, ErrRuntimeError, "assignment failed: %s", err.Error())
	}

	return value
}

// For Index Assignment (arr[0] = 42)
func (e *Evaluator) EvalIndexAssignment(n *ast.IndexAssignmentStatement) Value {
	// Validate AST structure
	if n.Object == nil {
		return e.newError(n.Position, ErrSyntaxError, "index assignment missing object")
	}

	if n.Index == nil {
		return e.newError(n.Position, ErrSyntaxError, "index assignment missing index")
	}

	if n.Value == nil {
		return e.newError(n.Position, ErrSyntaxError, "index assignment missing value")
	}

	// Evaluate the object being indexed
	object := e.Eval(n.Object)
	if isError(object) {
		return object
	}

	// Evaluate the index
	index := e.Eval(n.Index)
	if isError(index) {
		return index
	}

	// Evaluate the value to be assigned
	value := e.Eval(n.Value)
	if isError(value) {
		return value
	}

	// Check if index is an integer
	if index.Type() != INTEGER_TYPE {
		return e.newError(n.Position, ErrTypeMismatch, "array index must be integer, got %s", index.Type())
	}

	indexValue := index.(*IntegerValue).Value

	// Handle array assignment
	if object.Type() == ARRAY_TYPE {
		array := object.(*ArrayValue)
		if indexValue < 0 || indexValue >= int64(len(array.Elements)) {
			return e.newError(n.Position, ErrRuntimeError, "index out of bounds: %d", indexValue)
		}

		// Check type compatibility
		elementType := getValueType(array.Elements[indexValue])
		valueType := getValueType(value)
		if elementType != valueType {
			return e.newError(n.Position, ErrTypeMismatch,
				"type mismatch: cannot assign %s to array element of type %s", valueType, elementType)
		}

		// Perform the assignment
		array.Elements[indexValue] = value
		return value
	}

	// Handle string assignment (if we want to support it)
	if object.Type() == STRING_TYPE {
		return e.newError(n.Position, ErrRuntimeError, "cannot assign to string elements (strings are immutable)")
	}

	return e.newError(n.Position, ErrTypeMismatch, "cannot assign to index of type %s", object.Type())
}

func (e *Evaluator) initializeToZero(t *ast.Type) Value {
	v := t.BaseType
	switch v {
	case "STRING":
		return &StringValue{Value: ""}
	case "INTEGER":
		return &IntegerValue{Value: 0}
	case "FLOAT":
		return &FloatValue{Value: 0.00}
	case "BOOLEAN":
		return &BooleanValue{Value: false}
	default:
		return NULL
	}
}

func (e *Evaluator) TypesCompatible(expectedType string, actualType string) bool {
	// Handle case-insensitive type matching
	expected := strings.ToLower(expectedType)
	actual := strings.ToLower(actualType)

	// Direct match
	if expected == actual {
		return true
	}

	// Handle unknown type - allow assignment to unknown
	if expected == "unknown" || actual == "unknown" {
		return true
	}

	// Handle array types
	if strings.HasPrefix(expected, "[]") && strings.HasPrefix(actual, "[]") {
		// Extract element types and compare them
		expectedElement := strings.TrimPrefix(expected, "[]")
		actualElement := strings.TrimPrefix(actual, "[]")
		return e.TypesCompatible(expectedElement, actualElement)
	}

	// Handle fixed array types
	if strings.HasPrefix(expected, "[") && strings.HasPrefix(actual, "[") {
		// For now, allow any fixed array to match any fixed array
		// In the future, we could check sizes and element types
		return true
	}

	// Handle common type aliases
	if (expected == "int" && actual == "integer") ||
		(expected == "integer" && actual == "int") {
		return true
	}

	if (expected == "float" && actual == "float64") ||
		(expected == "float64" && actual == "float") {
		return true
	}

	if (expected == "bool" && actual == "boolean") ||
		(expected == "boolean" && actual == "bool") {
		return true
	}

	return false
}

func (e *Evaluator) EvalConditional(n *ast.IfStatement) Value {
	if n.Condition == nil {
		return NULL
	}
	condition := e.Eval(n.Condition)
	if isError(condition) {
		return condition
	}

	if condition.IsTruthy() {
		return e.Eval(n.Consequence)
	} else if n.Alternative != nil {
		return e.Eval(n.Alternative)
	}

	return NULL
}

func (e *Evaluator) EvalForStatement(n *ast.ForStatement) Value {
	e.pushFrame("for-loop", n.Position, "for statement")
	defer e.popFrame()

	// New scope for init/post variables
	e.env = NewEnclosedEnvironment(e.env)
	defer func() { e.env = e.env.outer }()

	// Init (if any)
	if n.Init != nil {
		if v := e.Eval(n.Init); isError(v) {
			return v
		}
	}

	// Main loop
	for {
		// Condition (if any)
		if n.Condition != nil {
			cond := e.Eval(n.Condition)
			if isError(cond) {
				return cond
			}
			if !cond.IsTruthy() {
				break
			}
		}

		// Body
		bodyVal := e.Eval(n.Body)
		if isError(bodyVal) {
			return bodyVal
		}

		switch bodyVal.Type() {
		case RETURN_TYPE:
			return bodyVal // propagate return
		case BREAK_TYPE:
			return NULL // exit loop
		case CONTINUE_TYPE:
			// Execute post, then continue
			if n.Post != nil {
				if v := e.Eval(n.Post); isError(v) {
					return v
				}
			}
			continue
		}

		// Post-statement (normal flow)
		if n.Post != nil {
			if v := e.Eval(n.Post); isError(v) {
				return v
			}
		}
	}

	return NULL
}

func (e *Evaluator) EvalWhileStatement(n *ast.WhileStatement) Value {
	e.pushFrame("while-loop", n.Position, "while statement")
	defer e.popFrame()

	for {
		// Evaluate condition
		cond := e.Eval(n.Condition)
		if isError(cond) {
			return cond
		}
		if !cond.IsTruthy() {
			break
		}

		// Execute body
		bodyVal := e.Eval(n.Body)
		if isError(bodyVal) {
			return bodyVal
		}

		// Handle control flow
		switch bodyVal.Type() {
		case RETURN_TYPE:
			return bodyVal // propagate return
		case BREAK_TYPE:
			return NULL // exit loop
		case CONTINUE_TYPE:
			continue // continue to next iteration
		}
	}

	return NULL
}

func (e *Evaluator) evalBlock(n *ast.BlockStatement) Value {
	//create new scopes for this block
	e.env = NewEnclosedEnvironment(e.env)
	defer func() { e.env = e.env.outer }()

	var result Value = NULL

	for _, stmt := range n.Statements {
		result = e.Eval(stmt)
		if isError(result) {
			return result
		}
		// After evaluating a statement, check if it was a "flow-breaking" one.
		if result != nil {
			rt := result.Type()
			// If we get a ReturnValue or BreakValue, stop executing this block
			// and pass the signal up the call stack immediately.
			if rt == RETURN_TYPE || rt == BREAK_TYPE || rt == CONTINUE_TYPE || rt == ERROR_TYPE {
				return result
			}
		}
	}
	return result
}

func (e *Evaluator) EvalIdentifier(n *ast.Identifier) Value {
	// Look up the variable in the environment
	binding, exists := e.env.Get(n.Name)
	if !exists {
		return e.newError(n.Position, ErrUndefinedVar,
			"undefined variable '%s'", n.Name)
	}

	return binding.Value
}

func (e *Evaluator) evalContinueStatement(n *ast.ContinueStatement) Value {
	return &ContinueValue{
		Position: n.Position,
	}
}

func (e *Evaluator) evalBreak(n *ast.BreakStatement) Value {
	return &BreakValue{
		Position: n.Position,
	}
}

func (e *Evaluator) evalFunctionCall(n *ast.FunctionCall) Value {
	function := e.Eval(n.Function)
	if isError(function) {
		return function
	}

	var results []Value
	for _, args := range n.Arguments {
		evaluated := e.Eval(args)

		if isError(evaluated) {
			return evaluated
		}
		results = append(results, evaluated)
	}

	if len(results) == 1 && isError(results[0]) {
		return results[0]
	}

	isFunction, ok := function.(*FunctionValue)
	if !ok {
		return e.newError(n.Position, ErrNotAFunction,
			"'%s' is not a function", function.Type())
	}

	// Handle built-in functions
	if isFunction.IsBuiltin {
		e.pushFrame(isFunction.Name, n.Position, "builtin")
		defer e.popFrame()
		return isFunction.BuiltinFn(results)
	}

	// Handle user-defined functions
	e.pushFrame(isFunction.Name, n.Position, "call")
	defer e.popFrame()

	oldEnv := e.env
	e.env = NewEnclosedEnvironment(isFunction.Env)
	defer func() { e.env = oldEnv }()

	if len(results) != len(isFunction.Parameters) {
		return e.newError(n.Position, ErrWrongArgCount,
			"function '%s' expects %d arguments, got %d",
			isFunction.Name, len(isFunction.Parameters), len(results))
	}

	for paramIdx, param := range isFunction.Parameters {
		argValue := results[paramIdx]
		paramType := getTypeString(param.Type)
		argType := getValueType(argValue)

		if !e.TypesCompatible(paramType, argType) {
			return e.newError(n.Position, ErrTypeMismatch,
				"type mismatch: cannot assign %s to %s",
				strings.ToLower(argType), strings.ToLower(paramType))
		}

		e.env.Set(param.Name.Name, argValue, false)
	}
	execution := e.Eval(isFunction.Body)
	if isError(execution) {
		return execution
	}
	if returnValue, ok := execution.(*ReturnValue); ok {
		return returnValue.Value
	}

	// If no explicit return, return NULL
	return NULL
}

func (e *Evaluator) evalArrayLiteral(n *ast.ArrayLiteral) Value {
	elements := make([]Value, 0, len(n.Elements))

	for _, element := range n.Elements {
		evaluated := e.Eval(element)
		if isError(evaluated) {
			return evaluated
		}
		elements = append(elements, evaluated)
	}

	return &ArrayValue{Elements: elements}
}

func (e *Evaluator) evalStructLiteral(n *ast.StructLiteral) Value {
	if n.Type == nil {
		return e.newError(n.Position, ErrRuntimeError, "struct literal missing type")
	}

	fields := make(map[string]Value)
	for _, f := range n.Fields {
		if f == nil || f.Name == nil || f.Value == nil {
			return e.newError(n.Position, ErrRuntimeError, "invalid struct field initialization")
		}
		v := e.Eval(f.Value)
		if isError(v) {
			return v
		}
		fields[f.Name.Name] = v
	}

	return &StructValue{TypeName: n.Type.Name, Fields: fields}
}

func (e *Evaluator) evalMemberExpression(n *ast.MemberExpression) Value {
	obj := e.Eval(n.Object)
	if isError(obj) {
		return obj
	}
	if obj.Type() == STRUCT_TYPE {
		sv := obj.(*StructValue)
		if val, ok := sv.Fields[n.Property.Name]; ok {
			return val
		}
		return e.newError(n.Position, ErrRuntimeError, "field '%s' not found on %s", n.Property.Name, sv.TypeName)
	}
	return e.newError(n.Position, ErrRuntimeError, "cannot access member on type %s", obj.Type())
}

func getValueType(v Value) string {
	switch v.Type() {
	case INTEGER_TYPE:
		return "INTEGER"
	case FLOAT_TYPE:
		return "FLOAT"
	case STRING_TYPE:
		return "STRING"
	case BOOLEAN_TYPE:
		return "BOOLEAN"
	case ARRAY_TYPE:
		// For arrays, we need to determine the element type
		array := v.(*ArrayValue)
		if len(array.Elements) > 0 {
			elementType := getValueType(array.Elements[0])
			// Convert to lowercase to match our type system
			switch elementType {
			case "INTEGER":
				return "[]int"
			case "FLOAT":
				return "[]float"
			case "STRING":
				return "[]string"
			case "BOOLEAN":
				return "[]bool"
			default:
				return "[]unknown"
			}
		}
		return "[]unknown"
	default:
		return v.Type()
	}
}

// getTypeString converts an AST Type to a string representation
func getTypeString(t *ast.Type) string {
	if t == nil {
		return "unknown"
	}

	if t.ArrayType != nil {
		if t.ArraySize != nil {
			return fmt.Sprintf("[%d]%s", *t.ArraySize, getTypeString(t.ArrayType))
		} else {
			return fmt.Sprintf("[]%s", getTypeString(t.ArrayType))
		}
	}
	if t.PointerType != nil {
		return fmt.Sprintf("*%s", getTypeString(t.PointerType))
	}
	if t.StructName != "" {
		return t.StructName
	}
	return t.BaseType
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

func (e *Evaluator) evalIndexExpression(n *ast.IndexExpression) Value {
	// Evaluate the object being indexed
	object := e.Eval(n.Object)
	if isError(object) {
		return object
	}

	// Evaluate the index
	index := e.Eval(n.Index)
	if isError(index) {
		return index
	}

	// Check if index is an integer
	if index.Type() != INTEGER_TYPE {
		return e.newError(n.Position, ErrTypeMismatch, "array index must be integer, got %s", index.Type())
	}

	indexValue := index.(*IntegerValue).Value

	// Handle array indexing
	if object.Type() == ARRAY_TYPE {
		array := object.(*ArrayValue)
		if indexValue < 0 || indexValue >= int64(len(array.Elements)) {
			return e.newError(n.Position, ErrRuntimeError, "index out of bounds: %d", indexValue)
		}
		return array.Elements[indexValue]
	}

	// Handle string indexing
	if object.Type() == STRING_TYPE {
		str := object.(*StringValue).Value
		if indexValue < 0 || indexValue >= int64(len(str)) {
			return e.newError(n.Position, ErrRuntimeError, "index out of bounds: %d", indexValue)
		}
		// Return a single character as a string
		return &StringValue{Value: string(str[indexValue])}
	}

	return e.newError(n.Position, ErrTypeMismatch, "cannot index type %s", object.Type())
}

func (e *Evaluator) evalSliceExpression(n *ast.SliceExpression) Value {
	// Evaluate the object being sliced
	object := e.Eval(n.Object)
	if isError(object) {
		return object
	}

	// Evaluate start index (can be nil for [:end])
	var startIndex int64 = 0
	if n.Start != nil {
		start := e.Eval(n.Start)
		if isError(start) {
			return start
		}
		if start.Type() != INTEGER_TYPE {
			return e.newError(n.Position, ErrTypeMismatch, "slice start index must be integer, got %s", start.Type())
		}
		startIndex = start.(*IntegerValue).Value
	}

	// Evaluate end index (can be nil for [start:])
	var endIndex int64
	if n.End != nil {
		end := e.Eval(n.End)
		if isError(end) {
			return end
		}
		if end.Type() != INTEGER_TYPE {
			return e.newError(n.Position, ErrTypeMismatch, "slice end index must be integer, got %s", end.Type())
		}
		endIndex = end.(*IntegerValue).Value
	}

	// Handle string slicing
	if object.Type() == STRING_TYPE {
		str := object.(*StringValue).Value
		strLen := int64(len(str))

		// Handle negative indices (Python-style)
		if startIndex < 0 {
			startIndex = strLen + startIndex
		}
		if n.End != nil && endIndex < 0 {
			endIndex = strLen + endIndex
		}

		// Bounds checking
		if startIndex < 0 {
			startIndex = 0
		}
		if startIndex > strLen {
			startIndex = strLen
		}
		if n.End != nil {
			if endIndex < 0 {
				endIndex = 0
			}
			if endIndex > strLen {
				endIndex = strLen
			}
			if startIndex > endIndex {
				startIndex = endIndex
			}
		} else {
			endIndex = strLen
		}

		return &StringValue{Value: str[startIndex:endIndex]}
	}

	// Handle array slicing
	if object.Type() == ARRAY_TYPE {
		array := object.(*ArrayValue)
		arrayLen := int64(len(array.Elements))

		// Handle negative indices (Python-style)
		if startIndex < 0 {
			startIndex = arrayLen + startIndex
		}
		if n.End != nil && endIndex < 0 {
			endIndex = arrayLen + endIndex
		}

		// Bounds checking
		if startIndex < 0 {
			startIndex = 0
		}
		if startIndex > arrayLen {
			startIndex = arrayLen
		}
		if n.End != nil {
			if endIndex < 0 {
				endIndex = 0
			}
			if endIndex > arrayLen {
				endIndex = arrayLen
			}
			if startIndex > endIndex {
				startIndex = endIndex
			}
		} else {
			endIndex = arrayLen
		}

		// Create new array with sliced elements
		slicedElements := make([]Value, 0, endIndex-startIndex)
		for i := startIndex; i < endIndex; i++ {
			slicedElements = append(slicedElements, array.Elements[i])
		}
		return &ArrayValue{Elements: slicedElements}
	}

	return e.newError(n.Position, ErrTypeMismatch, "cannot slice type %s", object.Type())
}

func formatValueForOutput(value Value) string {
	switch v := value.(type) {
	case *StringValue:
		return v.Value // Don't add quotes for log output
	case *IntegerValue:
		return fmt.Sprintf("%d", v.Value)
	case *FloatValue:
		return fmt.Sprintf("%g", v.Value)
	case *BooleanValue:
		return fmt.Sprintf("%t", v.Value)
	case *NullValue:
		return "null"
	case *FunctionValue:
		return fmt.Sprintf("%s", v.Name)
	case *ArrayValue:
		return v.String()
	case *RuntimeError:
		return fmt.Sprintf("ERROR: %s", v.Detail.Message)
	default:
		return value.String()
	}
}

// GetEnvironment returns the current environment
func (e *Evaluator) GetEnvironment() *Environment {
	return e.env
}
