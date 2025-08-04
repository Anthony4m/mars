package evaluator

import (
	"fmt"
	"math"
	"time"
)

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name       string
	Parameters []string
	Function   func(args []Value) Value
}

// BuiltinFunctions holds all registered built-in functions
var BuiltinFunctions = map[string]*BuiltinFunction{
	"len": {
		Name:       "len",
		Parameters: []string{"value"},
		Function:   builtinLen,
	},
	"append": {
		Name:       "append",
		Parameters: []string{"slice", "value"},
		Function:   builtinAppend,
	},
	"print": {
		Name:       "print",
		Parameters: []string{"value"},
		Function:   builtinPrint,
	},
	"println": {
		Name:       "println",
		Parameters: []string{"value"},
		Function:   builtinPrintln,
	},
	"printf": {
		Name:       "printf",
		Parameters: []string{"format", "values..."},
		Function:   builtinPrintf,
	},
	"sin": {
		Name:       "sin",
		Parameters: []string{"angle"},
		Function:   builtinSin,
	},
	"cos": {
		Name:       "cos",
		Parameters: []string{"angle"},
		Function:   builtinCos,
	},
	"sqrt": {
		Name:       "sqrt",
		Parameters: []string{"value"},
		Function:   builtinSqrt,
	},
	"now": {
		Name:       "now",
		Parameters: []string{},
		Function:   builtinNow,
	},
}

// builtinLen returns the length of a string or array
func builtinLen(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("len() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	switch arg.Type() {
	case STRING_TYPE:
		return &IntegerValue{Value: int64(len(arg.(*StringValue).Value))}
	case ARRAY_TYPE:
		return &IntegerValue{Value: int64(len(arg.(*ArrayValue).Elements))}
	default:
		return &Error{Message: fmt.Sprintf("len() not supported for type %s", arg.Type())}
	}
}

// builtinAppend appends a value to an array
func builtinAppend(args []Value) Value {
	if len(args) != 2 {
		return &Error{Message: fmt.Sprintf("append() expects 2 arguments, got %d", len(args))}
	}

	slice := args[0]
	value := args[1]

	if slice.Type() != ARRAY_TYPE {
		return &Error{Message: fmt.Sprintf("append() first argument must be array, got %s", slice.Type())}
	}

	array := slice.(*ArrayValue)
	newElements := make([]Value, len(array.Elements)+1)
	copy(newElements, array.Elements)
	newElements[len(array.Elements)] = value

	return &ArrayValue{Elements: newElements}
}

// builtinPrint prints a value without newline
func builtinPrint(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("print() expects 1 argument, got %d", len(args))}
	}

	fmt.Print(formatValueForOutput(args[0]))
	return NULL
}

// builtinPrintln prints a value with newline
func builtinPrintln(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("println() expects 1 argument, got %d", len(args))}
	}

	fmt.Println(formatValueForOutput(args[0]))
	return NULL
}

// builtinPrintf prints formatted output
func builtinPrintf(args []Value) Value {
	if len(args) < 1 {
		return &Error{Message: fmt.Sprintf("printf() expects at least 1 argument, got %d", len(args))}
	}

	format := args[0]
	if format.Type() != STRING_TYPE {
		return &Error{Message: "printf() first argument must be string"}
	}

	formatStr := format.(*StringValue).Value
	formatArgs := make([]interface{}, len(args)-1)

	for i, arg := range args[1:] {
		formatArgs[i] = formatValueForOutput(arg)
	}

	fmt.Printf(formatStr, formatArgs...)
	return NULL
}

// builtinSin returns the sine of an angle in radians
func builtinSin(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("sin() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	var angle float64

	switch arg.Type() {
	case INTEGER_TYPE:
		angle = float64(arg.(*IntegerValue).Value)
	case FLOAT_TYPE:
		angle = arg.(*FloatValue).Value
	default:
		return &Error{Message: fmt.Sprintf("sin() expects number, got %s", arg.Type())}
	}

	return &FloatValue{Value: math.Sin(angle)}
}

// builtinCos returns the cosine of an angle in radians
func builtinCos(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("cos() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	var angle float64

	switch arg.Type() {
	case INTEGER_TYPE:
		angle = float64(arg.(*IntegerValue).Value)
	case FLOAT_TYPE:
		angle = arg.(*FloatValue).Value
	default:
		return &Error{Message: fmt.Sprintf("cos() expects number, got %s", arg.Type())}
	}

	return &FloatValue{Value: math.Cos(angle)}
}

// builtinSqrt returns the square root of a number
func builtinSqrt(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("sqrt() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	var value float64

	switch arg.Type() {
	case INTEGER_TYPE:
		value = float64(arg.(*IntegerValue).Value)
	case FLOAT_TYPE:
		value = arg.(*FloatValue).Value
	default:
		return &Error{Message: fmt.Sprintf("sqrt() expects number, got %s", arg.Type())}
	}

	if value < 0 {
		return &Error{Message: "sqrt() of negative number"}
	}

	return &FloatValue{Value: math.Sqrt(value)}
}

// builtinNow returns the current time as a string
func builtinNow(args []Value) Value {
	if len(args) != 0 {
		return &Error{Message: fmt.Sprintf("now() expects 0 arguments, got %d", len(args))}
	}

	return &StringValue{Value: time.Now().Format(time.RFC3339)}
}
