package evaluator

import (
	"fmt"
	"math"
	"strings"
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
	"toInt": {
		Name:       "toInt",
		Parameters: []string{"value"},
		Function:   builtinInt,
	},
	"toFloat": {
		Name:       "toFloat",
		Parameters: []string{"value"},
		Function:   builtinFloat,
	},
	"toString": {
		Name:       "toString",
		Parameters: []string{"value"},
		Function:   builtinString,
	},
	"getType": {
		Name:       "getType",
		Parameters: []string{"value"},
		Function:   builtinType,
	},
	"abs": {
		Name:       "abs",
		Parameters: []string{"value"},
		Function:   builtinAbs,
	},
	"min": {
		Name:       "min",
		Parameters: []string{"a", "b"},
		Function:   builtinMin,
	},
	"max": {
		Name:       "max",
		Parameters: []string{"a", "b"},
		Function:   builtinMax,
	},
	"isInt": {
		Name:       "isInt",
		Parameters: []string{"value"},
		Function:   builtinIsInt,
	},
	"isFloat": {
		Name:       "isFloat",
		Parameters: []string{"value"},
		Function:   builtinIsFloat,
	},
	"isString": {
		Name:       "isString",
		Parameters: []string{"value"},
		Function:   builtinIsString,
	},
	"isArray": {
		Name:       "isArray",
		Parameters: []string{"value"},
		Function:   builtinIsArray,
	},
	"isBool": {
		Name:       "isBool",
		Parameters: []string{"value"},
		Function:   builtinIsBool,
	},
	"pow": {
		Name:       "pow",
		Parameters: []string{"base", "exponent"},
		Function:   builtinPow,
	},
	"floor": {
		Name:       "floor",
		Parameters: []string{"value"},
		Function:   builtinFloor,
	},
	"ceil": {
		Name:       "ceil",
		Parameters: []string{"value"},
		Function:   builtinCeil,
	},
	"push": {
		Name:       "push",
		Parameters: []string{"array", "value"},
		Function:   builtinPush,
	},
	"pop": {
		Name:       "pop",
		Parameters: []string{"array"},
		Function:   builtinPop,
	},
	"reverse": {
		Name:       "reverse",
		Parameters: []string{"array"},
		Function:   builtinReverse,
	},
	"join": {
		Name:       "join",
		Parameters: []string{"array", "separator"},
		Function:   builtinJoin,
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

// builtinInt converts a value to integer
func builtinInt(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("int() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	switch arg.Type() {
	case INTEGER_TYPE:
		return arg
	case FLOAT_TYPE:
		return &IntegerValue{Value: int64(arg.(*FloatValue).Value)}
	case STRING_TYPE:
		// Try to parse string as integer
		var value int64
		_, err := fmt.Sscanf(arg.(*StringValue).Value, "%d", &value)
		if err != nil {
			return &Error{Message: fmt.Sprintf("cannot convert string '%s' to int", arg.(*StringValue).Value)}
		}
		return &IntegerValue{Value: value}
	default:
		return &Error{Message: fmt.Sprintf("cannot convert %s to int", arg.Type())}
	}
}

// builtinFloat converts a value to float
func builtinFloat(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("float() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	switch arg.Type() {
	case FLOAT_TYPE:
		return arg
	case INTEGER_TYPE:
		return &FloatValue{Value: float64(arg.(*IntegerValue).Value)}
	case STRING_TYPE:
		// Try to parse string as float
		var value float64
		_, err := fmt.Sscanf(arg.(*StringValue).Value, "%f", &value)
		if err != nil {
			return &Error{Message: fmt.Sprintf("cannot convert string '%s' to float", arg.(*StringValue).Value)}
		}
		return &FloatValue{Value: value}
	default:
		return &Error{Message: fmt.Sprintf("cannot convert %s to float", arg.Type())}
	}
}

// builtinString converts a value to string
func builtinString(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("string() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	return &StringValue{Value: arg.String()}
}

// builtinType returns the type of a value as a string
func builtinType(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("type() expects 1 argument, got %d", len(args))}
	}

	return &StringValue{Value: args[0].Type()}
}

// builtinAbs returns the absolute value
func builtinAbs(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("abs() expects 1 argument, got %d", len(args))}
	}

	arg := args[0]
	switch arg.Type() {
	case INTEGER_TYPE:
		value := arg.(*IntegerValue).Value
		if value < 0 {
			return &IntegerValue{Value: -value}
		}
		return arg
	case FLOAT_TYPE:
		return &FloatValue{Value: math.Abs(arg.(*FloatValue).Value)}
	default:
		return &Error{Message: fmt.Sprintf("abs() not supported for type %s", arg.Type())}
	}
}

// builtinMin returns the minimum of two values
func builtinMin(args []Value) Value {
	if len(args) != 2 {
		return &Error{Message: fmt.Sprintf("min() expects 2 arguments, got %d", len(args))}
	}

	a, b := args[0], args[1]

	// Handle integer comparison
	if a.Type() == INTEGER_TYPE && b.Type() == INTEGER_TYPE {
		valA := a.(*IntegerValue).Value
		valB := b.(*IntegerValue).Value
		if valA < valB {
			return a
		}
		return b
	}

	// Handle float comparison
	if a.Type() == FLOAT_TYPE && b.Type() == FLOAT_TYPE {
		valA := a.(*FloatValue).Value
		valB := b.(*FloatValue).Value
		if valA < valB {
			return a
		}
		return b
	}

	// Handle mixed int/float
	if a.Type() == INTEGER_TYPE && b.Type() == FLOAT_TYPE {
		valA := float64(a.(*IntegerValue).Value)
		valB := b.(*FloatValue).Value
		if valA < valB {
			return &FloatValue{Value: valA}
		}
		return b
	}

	if a.Type() == FLOAT_TYPE && b.Type() == INTEGER_TYPE {
		valA := a.(*FloatValue).Value
		valB := float64(b.(*IntegerValue).Value)
		if valA < valB {
			return a
		}
		return &FloatValue{Value: valB}
	}

	return &Error{Message: fmt.Sprintf("min() not supported for types %s and %s", a.Type(), b.Type())}
}

// builtinMax returns the maximum of two values
func builtinMax(args []Value) Value {
	if len(args) != 2 {
		return &Error{Message: fmt.Sprintf("max() expects 2 arguments, got %d", len(args))}
	}

	a, b := args[0], args[1]

	// Handle integer comparison
	if a.Type() == INTEGER_TYPE && b.Type() == INTEGER_TYPE {
		valA := a.(*IntegerValue).Value
		valB := b.(*IntegerValue).Value
		if valA > valB {
			return a
		}
		return b
	}

	// Handle float comparison
	if a.Type() == FLOAT_TYPE && b.Type() == FLOAT_TYPE {
		valA := a.(*FloatValue).Value
		valB := b.(*FloatValue).Value
		if valA > valB {
			return a
		}
		return b
	}

	// Handle mixed int/float
	if a.Type() == INTEGER_TYPE && b.Type() == FLOAT_TYPE {
		valA := float64(a.(*IntegerValue).Value)
		valB := b.(*FloatValue).Value
		if valA > valB {
			return &FloatValue{Value: valA}
		}
		return b
	}

	if a.Type() == FLOAT_TYPE && b.Type() == INTEGER_TYPE {
		valA := a.(*FloatValue).Value
		valB := float64(b.(*IntegerValue).Value)
		if valA > valB {
			return a
		}
		return &FloatValue{Value: valB}
	}

	return &Error{Message: fmt.Sprintf("max() not supported for types %s and %s", a.Type(), b.Type())}
}

// builtinIsInt checks if a value is an integer
func builtinIsInt(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("isInt() expects 1 argument, got %d", len(args))}
	}
	return &BooleanValue{Value: args[0].Type() == INTEGER_TYPE}
}

// builtinIsFloat checks if a value is a float
func builtinIsFloat(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("isFloat() expects 1 argument, got %d", len(args))}
	}
	return &BooleanValue{Value: args[0].Type() == FLOAT_TYPE}
}

// builtinIsString checks if a value is a string
func builtinIsString(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("isString() expects 1 argument, got %d", len(args))}
	}
	return &BooleanValue{Value: args[0].Type() == STRING_TYPE}
}

// builtinIsArray checks if a value is an array
func builtinIsArray(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("isArray() expects 1 argument, got %d", len(args))}
	}
	return &BooleanValue{Value: args[0].Type() == ARRAY_TYPE}
}

// builtinIsBool checks if a value is a boolean
func builtinIsBool(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("isBool() expects 1 argument, got %d", len(args))}
	}
	return &BooleanValue{Value: args[0].Type() == BOOLEAN_TYPE}
}

// builtinPow calculates base raised to the power of exponent
func builtinPow(args []Value) Value {
	if len(args) != 2 {
		return &Error{Message: fmt.Sprintf("pow() expects 2 arguments, got %d", len(args))}
	}

	base, exponent := args[0], args[1]

	// Handle integer power
	if base.Type() == INTEGER_TYPE && exponent.Type() == INTEGER_TYPE {
		baseVal := base.(*IntegerValue).Value
		expVal := exponent.(*IntegerValue).Value
		result := math.Pow(float64(baseVal), float64(expVal))
		// If result is a whole number, return as integer
		if result == math.Floor(result) {
			return &IntegerValue{Value: int64(result)}
		}
		return &FloatValue{Value: result}
	}

	// Handle float power
	if base.Type() == FLOAT_TYPE && exponent.Type() == FLOAT_TYPE {
		baseVal := base.(*FloatValue).Value
		expVal := exponent.(*FloatValue).Value
		return &FloatValue{Value: math.Pow(baseVal, expVal)}
	}

	// Handle mixed int/float
	if base.Type() == INTEGER_TYPE && exponent.Type() == FLOAT_TYPE {
		baseVal := float64(base.(*IntegerValue).Value)
		expVal := exponent.(*FloatValue).Value
		return &FloatValue{Value: math.Pow(baseVal, expVal)}
	}

	if base.Type() == FLOAT_TYPE && exponent.Type() == INTEGER_TYPE {
		baseVal := base.(*FloatValue).Value
		expVal := float64(exponent.(*IntegerValue).Value)
		return &FloatValue{Value: math.Pow(baseVal, expVal)}
	}

	return &Error{Message: fmt.Sprintf("pow() not supported for types %s and %s", base.Type(), exponent.Type())}
}

// builtinFloor returns the floor of a number
func builtinFloor(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("floor() expects 1 argument, got %d", len(args))}
	}

	value := args[0]

	if value.Type() == INTEGER_TYPE {
		return value // Integer is already "floored"
	}

	if value.Type() == FLOAT_TYPE {
		val := value.(*FloatValue).Value
		result := math.Floor(val)
		// If result is a whole number, return as integer
		if result == math.Floor(result) {
			return &IntegerValue{Value: int64(result)}
		}
		return &FloatValue{Value: result}
	}

	return &Error{Message: fmt.Sprintf("floor() not supported for type %s", value.Type())}
}

// builtinCeil returns the ceiling of a number
func builtinCeil(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("ceil() expects 1 argument, got %d", len(args))}
	}

	value := args[0]

	if value.Type() == INTEGER_TYPE {
		return value // Integer is already "ceiled"
	}

	if value.Type() == FLOAT_TYPE {
		val := value.(*FloatValue).Value
		result := math.Ceil(val)
		// If result is a whole number, return as integer
		if result == math.Floor(result) {
			return &IntegerValue{Value: int64(result)}
		}
		return &FloatValue{Value: result}
	}

	return &Error{Message: fmt.Sprintf("ceil() not supported for type %s", value.Type())}
}

// builtinPush adds a value to the end of an array
func builtinPush(args []Value) Value {
	if len(args) != 2 {
		return &Error{Message: fmt.Sprintf("push() expects 2 arguments, got %d", len(args))}
	}

	array, value := args[0], args[1]

	if array.Type() != ARRAY_TYPE {
		return &Error{Message: fmt.Sprintf("push() expects array as first argument, got %s", array.Type())}
	}

	arr := array.(*ArrayValue)
	arr.Elements = append(arr.Elements, value)
	return array // Return the modified array
}

// builtinPop removes and returns the last element of an array
func builtinPop(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("pop() expects 1 argument, got %d", len(args))}
	}

	array := args[0]

	if array.Type() != ARRAY_TYPE {
		return &Error{Message: fmt.Sprintf("pop() expects array as argument, got %s", array.Type())}
	}

	arr := array.(*ArrayValue)
	if len(arr.Elements) == 0 {
		return &Error{Message: "pop() called on empty array"}
	}

	// Get the last element
	lastElement := arr.Elements[len(arr.Elements)-1]

	// Remove the last element
	arr.Elements = arr.Elements[:len(arr.Elements)-1]

	return lastElement
}

// builtinReverse reverses an array in place
func builtinReverse(args []Value) Value {
	if len(args) != 1 {
		return &Error{Message: fmt.Sprintf("reverse() expects 1 argument, got %d", len(args))}
	}

	array := args[0]

	if array.Type() != ARRAY_TYPE {
		return &Error{Message: fmt.Sprintf("reverse() expects array as argument, got %s", array.Type())}
	}

	arr := array.(*ArrayValue)

	// Reverse the array in place
	for i, j := 0, len(arr.Elements)-1; i < j; i, j = i+1, j-1 {
		arr.Elements[i], arr.Elements[j] = arr.Elements[j], arr.Elements[i]
	}

	return array // Return the modified array
}

// builtinJoin joins array elements with a separator
func builtinJoin(args []Value) Value {
	if len(args) != 2 {
		return &Error{Message: fmt.Sprintf("join() expects 2 arguments, got %d", len(args))}
	}

	array, separator := args[0], args[1]

	if array.Type() != ARRAY_TYPE {
		return &Error{Message: fmt.Sprintf("join() expects array as first argument, got %s", array.Type())}
	}

	if separator.Type() != STRING_TYPE {
		return &Error{Message: fmt.Sprintf("join() expects string as second argument, got %s", separator.Type())}
	}

	arr := array.(*ArrayValue)
	sep := separator.(*StringValue).Value

	if len(arr.Elements) == 0 {
		return &StringValue{Value: ""}
	}

	// Convert all elements to strings and join them
	var result strings.Builder
	for i, element := range arr.Elements {
		if i > 0 {
			result.WriteString(sep)
		}
		result.WriteString(element.String())
	}

	return &StringValue{Value: result.String()}
}
