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
	return &Evaluator{env: NewEnvironment()}
}

/*
			TODO: 1. **Variable Declaration**
		  - When you see a `let` or `var` statement, you must bind the new name into the *current* environment.
		  - 2. **Identifier Lookup**
		  - Whenever you evaluate an `Identifier` node, you must call `env.Get(name)` to fetch its current value (or signal an undefined-name error).
		  - 3. **Assignment / Mutation**
		  - On an assignment expression (`x = expr`), first evaluate the right-hand side, then use `env.Update(name, newValue)` (or `Set` if you allow shadowing) to change the binding in the correct scope.
		  - 4. **Function Literal (Closure) Creation**
		  - When you evaluate a function literal, you capture the *current* environment in the closure object so that when the function later runs, it still “sees” those outer variables.
		  - 5. **Function Call**
		  - Just before executing a function body, you create a **new enclosed environment** whose `outer` points at the function’s defining (closure) environment.
		  - You then bind the call’s parameters in that new environment and evaluate the body there.
		  - 6. **Block / Nested Scope**
		  - If your language has explicit block scopes (e.g. `{ … }` inside a function), you would push a fresh environment at block entry and pop it at exit so that temporaries don’t leak out.
		  - 7. **Built-in or Standard Library Calls**
		  - Even built-ins that modify state (e.g. `push(array, x)`) might need to read or write into the environment, depending on how they’re wired up.
		  - 8. **Runtime Error Reporting**
		  - If a variable lookup fails (Get returns not-found), or an update fails, you typically turn that into a runtime error—so your error‐handling paths also inspect the environment.
		  - Those are the “where” — every time you introduce, lookup, or mutate a name in your AST, you reach for the `Environment`.
	      - 9. **Enhance Error Handling with Stack Traces**
		  - Modify the error reporting mechanism to capture and display call stacks for runtime errors, similar to how Rust handles panics or unrecoverable errors.
		  - 10. **Consider NaN-boxing for Value Representation**
		  - Investigate and potentially implement NaN-boxing for more efficient value representation in `evaluator/value.go` if performance becomes a critical concern. This is a complex optimization and should be approached carefully.
*/
func (e *Evaluator) Eval(node ast.Node) Value {
	fmt.Printf("Evaluating: %T\n", node)
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
	case *ast.IfStatement:
		return e.EvalConditional(n)
	case *ast.ForStatement:
		return e.EvalForStatement(n)
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
	if valueType != varType {
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
	if expectedType != actualType {
		return false
	}
	return true
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
		paramType := param.Type.BaseType
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
	default:
		return v.Type()
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
