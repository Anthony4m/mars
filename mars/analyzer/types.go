package analyzer

import (
	"fmt"
	"mars/ast"
)

// TypeChecker performs type checking on the AST
type TypeChecker struct {
	errors []Error
}

// NewTypeChecker creates a new type checker instance
func NewTypeChecker() *TypeChecker {
	return &TypeChecker{}
}

// CheckType verifies that an expression's type matches the expected type
func (tc *TypeChecker) CheckType(expr ast.Expression, expected ast.Type) error {
	actual := tc.InferType(expr)
	if !tc.typesCompatible(actual, expected) {
		return fmt.Errorf("type mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

// InferType determines the type of an expression
func (tc *TypeChecker) InferType(expr ast.Expression) ast.Type {
	switch e := expr.(type) {
	case *ast.Literal:
		return tc.inferLiteralType(e)
	case *ast.Identifier:
		return tc.inferIdentifierType(e)
	case *ast.BinaryExpr:
		return tc.inferBinaryExprType(e)
	case *ast.CallExpr:
		return tc.inferCallExprType(e)
	default:
		return &ast.UnknownType{}
	}
}

// typesCompatible checks if two types are compatible
func (tc *TypeChecker) typesCompatible(t1, t2 ast.Type) bool {
	// TODO: Implement type compatibility rules
	// - Basic types must match exactly
	// - Arrays must have compatible element types
	// - Structs must have compatible fields
	// - Functions must have compatible signatures
	return false
}

// inferLiteralType determines the type of a literal expression
func (tc *TypeChecker) inferLiteralType(lit *ast.Literal) ast.Type {
	switch lit.Value.(type) {
	case int:
		return &ast.BasicType{Name: "int"}
	case float64:
		return &ast.BasicType{Name: "float"}
	case string:
		return &ast.BasicType{Name: "string"}
	case bool:
		return &ast.BasicType{Name: "bool"}
	default:
		return &ast.UnknownType{}
	}
}

// inferIdentifierType determines the type of an identifier
func (tc *TypeChecker) inferIdentifierType(id *ast.Identifier) ast.Type {
	// TODO: Look up identifier in symbol table
	return &ast.UnknownType{}
}

// inferBinaryExprType determines the type of a binary expression
func (tc *TypeChecker) inferBinaryExprType(expr *ast.BinaryExpr) ast.Type {
	leftType := tc.InferType(expr.Left)
	rightType := tc.InferType(expr.Right)

	// TODO: Implement binary operation type rules
	// - Arithmetic ops: numeric types
	// - Comparison ops: boolean result
	// - Logical ops: boolean operands and result
	return &ast.UnknownType{}
}

// inferCallExprType determines the type of a function call
func (tc *TypeChecker) inferCallExprType(expr *ast.CallExpr) ast.Type {
	// TODO: Look up function in symbol table and return its return type
	return &ast.UnknownType{}
}
