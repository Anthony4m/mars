package evaluator

import (
	"bytes"
	"fmt"
	"mars/ast"
	"strings"
)

// Type constants
const (
	INTEGER_TYPE  = "INTEGER"
	BOOLEAN_TYPE  = "BOOLEAN"
	STRING_TYPE   = "STRING"
	NULL_TYPE     = "NULL"
	ERROR_TYPE    = "ERROR"
	FUNCTION_TYPE = "FUNCTION"
	RETURN_TYPE   = "RETURN"
	FLOAT_TYPE    = "FLOAT"
	BREAK_TYPE    = "BREAK"
	CONTINUE_TYPE = "CONTINUE"
	ARRAY_TYPE    = "ARRAY"
)

// Value interface - all runtime values implement this
type Value interface {
	Type() string
	String() string
	IsTruthy() bool
}

// IntegerValue represents integer values
type IntegerValue struct {
	Value int64
}

func (i *IntegerValue) Type() string   { return INTEGER_TYPE }
func (i *IntegerValue) String() string { return fmt.Sprintf("%d", i.Value) }
func (i *IntegerValue) IsTruthy() bool { return i.Value != 0 }

// BooleanValue represents boolean values
type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) Type() string   { return BOOLEAN_TYPE }
func (b *BooleanValue) String() string { return fmt.Sprintf("%t", b.Value) }
func (b *BooleanValue) IsTruthy() bool { return b.Value }

// StringValue represents string values
type StringValue struct {
	Value string
}

func (s *StringValue) Type() string   { return STRING_TYPE }
func (s *StringValue) String() string { return s.Value }
func (s *StringValue) IsTruthy() bool { return len(s.Value) > 0 }

// FloatValue represents floating point values
type FloatValue struct {
	Value float64
}

func (f *FloatValue) Type() string   { return FLOAT_TYPE }
func (f *FloatValue) String() string { return fmt.Sprintf("%g", f.Value) }
func (f *FloatValue) IsTruthy() bool { return f.Value != 0.0 }

// NullValue represents null/nil values
type NullValue struct{}

func (n *NullValue) Type() string   { return NULL_TYPE }
func (n *NullValue) String() string { return "null" }
func (n *NullValue) IsTruthy() bool { return false }

// Error represents runtime errors
type Error struct {
	Message string
}

func (e *Error) Type() string   { return ERROR_TYPE }
func (e *Error) String() string { return "ERROR: " + e.Message }
func (e *Error) IsTruthy() bool { return false }

// ReturnValue wraps a return value
type ReturnValue struct {
	Value Value
}

func (r *ReturnValue) Type() string   { return RETURN_TYPE }
func (r *ReturnValue) String() string { return r.Value.String() }
func (r *ReturnValue) IsTruthy() bool { return r.Value.IsTruthy() }

type ContinueValue struct {
	Value    Value
	Position ast.Position
}

func (i *ContinueValue) Type() string   { return CONTINUE_TYPE }
func (i *ContinueValue) String() string { return fmt.Sprintf("%d", i.Value) }
func (i *ContinueValue) IsTruthy() bool { return false }

type BreakValue struct {
	Value    Value
	Position ast.Position
}
type FunctionValue struct {
	Name       string
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	ReturnType *ast.Type
	Env        *Environment // For closure support
	Position   ast.Position
	IsBuiltin  bool                     // True if this is a built-in function
	BuiltinFn  func(args []Value) Value // Built-in function implementation
}

func (fv *FunctionValue) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fv.Parameters {
		params = append(params, p.Name.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fv.Body.String())
	out.WriteString("\n}")
	return out.String()
}

func (fv *FunctionValue) Type() string {
	return "FUNCTION"
}

func (fv *FunctionValue) IsTruthy() bool {
	return true
}

// ArrayValue represents array values
type ArrayValue struct {
	Elements []Value
}

func (a *ArrayValue) Type() string { return ARRAY_TYPE }
func (a *ArrayValue) String() string {
	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, elem.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}
func (a *ArrayValue) IsTruthy() bool { return len(a.Elements) > 0 }

func (i *BreakValue) Type() string   { return BREAK_TYPE }
func (i *BreakValue) String() string { return fmt.Sprintf("%d", i.Value) }
func (i *BreakValue) IsTruthy() bool { return false }

// Singleton values to save allocations
var (
	TRUE  = &BooleanValue{Value: true}
	FALSE = &BooleanValue{Value: false}
	NULL  = &NullValue{}
)
