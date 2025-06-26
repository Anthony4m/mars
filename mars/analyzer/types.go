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
	actual := tc.inferType(expr)
	if !tc.typesCompatible(actual, &expected) {
		return fmt.Errorf("type mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

// InferType determines the type of an expression
func (tc *TypeChecker) inferType(expr interface{}) *ast.Type {
	switch e := expr.(type) {
	case *ast.Literal:
		switch e.Value.(type) {
		case int, int64:
			return &ast.Type{BaseType: "int"}
		case float64:
			return &ast.Type{BaseType: "float"}
		case string:
			return &ast.Type{BaseType: "string"}
		case bool:
			return &ast.Type{BaseType: "bool"}
		}
	}
	// For now, default to unknown
	return &ast.Type{BaseType: "unknown"}
}

// typesCompatible checks if two types are compatible
func (tc *TypeChecker) typesCompatible(actual, expected *ast.Type) bool {
	if expected == nil || actual == nil {
		return false
	}

	// For now, simple base type comparison
	if expected.BaseType != actual.BaseType {
		return false
	}

	// Special handling for function types
	if expected.IsFunctionType() && actual.IsFunctionType() {
		return tc.functionSignaturesCompatible(
			expected.GetFunctionSignature(),
			actual.GetFunctionSignature(),
		)
	}

	return true
}

// functionSignaturesCompatible checks if two function signatures are compatible
func (tc *TypeChecker) functionSignaturesCompatible(expected, actual *ast.FunctionSignature) bool {
	if expected == nil || actual == nil {
		return false
	}

	// Check parameter count
	if len(expected.Parameters) != len(actual.Parameters) {
		return false
	}

	// Check parameter types
	for i, expectedParam := range expected.Parameters {
		actualParam := actual.Parameters[i]
		if !tc.typesCompatible(expectedParam.Type, actualParam.Type) {
			return false
		}
	}
	// Check return type
	if expected.ReturnType == nil && actual.ReturnType == nil {
		return true
	}
	if expected.ReturnType == nil || actual.ReturnType == nil {
		return false
	}

	return tc.typesCompatible(expected.ReturnType, actual.ReturnType)
}

// inferIdentifierType determines the type of an identifier
func (tc *TypeChecker) inferIdentifierType(id *ast.Identifier) *ast.Type {
	// TODO: Look up identifier in symbol table
	return &ast.Type{BaseType: "unknown"}
}

// inferBinaryExprType determines the type of a binary expression
func (tc *TypeChecker) inferBinaryExprType(expr *ast.BinaryExpression) *ast.Type {
	leftType := tc.inferType(expr.Left)
	rightType := tc.inferType(expr.Right)

	// TODO: Implement binary operation type rules
	// - Arithmetic ops: numeric types
	// - Comparison ops: boolean result
	// - Logical ops: boolean operands and result

	// For now, return the left type as a placeholder
	_ = rightType // Suppress unused variable warning
	return leftType
}

// inferCallExprType determines the type of a function call
func (tc *TypeChecker) inferCallExprType(expr *ast.FunctionCall) *ast.Type {
	// TODO: Look up function in symbol table and return its return type
	return &ast.Type{BaseType: "unknown"}
}
