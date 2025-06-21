package analyzer

import (
	"fmt"
	"mars/ast"
	"mars/errors"
)

// Analyzer performs semantic analysis on the AST
type Analyzer struct {
	errors     *errors.MarsReporter
	symbols    *SymbolTable
	types      *TypeChecker
	immutable  *ImmutabilityChecker
	sourceCode string
	filename   string
}

// New creates a new analyzer instance
func New(sourceCode, filename string) *Analyzer {
	return &Analyzer{
		errors:     errors.NewMarsReporter(sourceCode, filename),
		symbols:    NewSymbolTable(),
		types:      NewTypeChecker(),
		immutable:  NewImmutabilityChecker(),
		sourceCode: sourceCode,
		filename:   filename,
	}
}

// Analyze performs semantic analysis on the given AST
func (a *Analyzer) Analyze(node ast.Node) error {
	// First pass: collect all declarations
	if err := a.collectDeclarations(node); err != nil {
		return err
	}

	// Check if we have errors from first pass
	if a.errors.HasErrors() {
		return fmt.Errorf("%s", a.errors.String())
	}

	// Second pass: type checking and immutability analysis
	if err := a.checkTypes(node); err != nil {
		return err
	}

	if err := a.checkImmutability(node); err != nil {
		return err
	}

	return nil
}

// Error represents a semantic analysis error
type Error struct {
	Line   int
	Column int
	Msg    string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Column, e.Msg)
}

// collectDeclarations performs the first pass of semantic analysis,
// collecting all declarations and building the symbol table
func (a *Analyzer) collectDeclarations(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		//start here - traverse the program and collect declarations
		for _, decl := range n.Declarations {
			if err := a.collectDeclarations(decl); err != nil {
				return err
			}
		}
	case *ast.VarDecl:
		//collect variable declarations
		return a.collectVariableDeclaration(n)
	case *ast.FuncDecl:
		//collect function declarations
		return a.collectFunctionDeclaration(n)
	case *ast.StructDecl:
		//collect struct declarations
		return a.collectStructDeclaration(n)
	case *ast.UnsafeBlock:
		//collect unsafe blocks
		return a.collectUnsafeBlock(n)
	default:
		return nil // ✅ Ignore non-declaration nodes
	}

	return nil
}

func (a *Analyzer) collectVariableDeclaration(decl *ast.VarDecl) error {
	var varType *ast.Type

	if decl.Type != nil && decl.Value != nil {
		// Both explicit type and value - check compatibility
		varType = decl.Type
		inferredType := a.inferType(decl.Value)

		if !a.typesCompatible(varType, inferredType) {
			help := fmt.Sprintf("change type to '%s' or cast the value", inferredType.BaseType)
			a.errors.AddErrorWithHelp(
				decl.Name.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("mismatched types: expected '%s', found '%s'",
					varType.BaseType, inferredType.BaseType),
				help,
			)
		}
	} else if decl.Type != nil {
		varType = decl.Type
	} else if decl.Value != nil {
		varType = a.inferType(decl.Value)
		if varType == nil || varType.BaseType == "unknown" {
			a.errors.AddError(
				decl.Name.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("cannot infer type for variable '%s'", decl.Name.Name),
			)
		}
	} else {
		a.errors.AddError(
			decl.Name.Position,
			errors.ErrCodeSyntaxError,
			fmt.Sprintf("variable '%s' needs either a type annotation or an initial value",
				decl.Name.Name),
		)
	}

	// Check for duplicate declarations
	if varType != nil {
		if err := a.symbols.Define(decl.Name.Name, *varType, decl.Mutable, false, decl); err != nil {
			a.errors.AddErrorWithHelp(
				decl.Name.Position,
				errors.ErrCodeDuplicateDecl,
				fmt.Sprintf("variable '%s' is already defined in this scope", decl.Name.Name),
				"give this variable a different name",
			)

			// Add note about original declaration
			if original, _ := a.symbols.Resolve(decl.Name.Name); original != nil && original.DeclaredAt != nil {
				//origPos := (*original.DeclaredAt).Pos()
				// You could enhance this to show the original location
			}
		}
	}

	return nil
}

func (a *Analyzer) collectFunctionDeclaration(decl *ast.FuncDecl) error {
	// 1. Create the function's type directly from the signature in the AST.
	//    It's guaranteed to be there and be complete.
	funcType := ast.NewFunctionType(decl.Signature)

	// 2. Add the function with its full signature to the symbol table.
	if err := a.symbols.Define(decl.Name.Name, *funcType, false, true, decl); err != nil {
		a.errors.AddErrorWithHelp(
			decl.Name.Position,
			errors.ErrCodeDuplicateDecl,
			fmt.Sprintf("function '%s' is already defined in this scope", decl.Name.Name),
			"give this function a different name",
		)
	}

	// 3. Defer checking the body until the second pass.
	return nil
}

func (a *Analyzer) checkFunctionBody(decl *ast.FuncDecl) error {
	// Enter function scope
	a.symbols.EnterScope()
	defer a.symbols.ExitScope()

	// Add parameters to local scope
	for _, param := range decl.Signature.Parameters {
		err := a.symbols.Define(param.Name.Name, *param.Type, false, false, param.Name)
		if err != nil {
			return err
		}
	}

	// Now check the body - recursion is safe because
	// the function name is already in the symbol table
	return a.checkTypes(decl.Body)
}

func (a *Analyzer) collectStructDeclaration(decl *ast.StructDecl) error {
	structType := &ast.Type{BaseType: "struct"}
	if decl.Name == nil {
		return fmt.Errorf("struct declaration must have a name")
	}
	// 1. Add struct type to global scope
	err := a.symbols.Define(decl.Name.Name, *structType, false, false, decl)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyzer) collectUnsafeBlock(block *ast.UnsafeBlock) error {
	if block.Body == nil {
		return fmt.Errorf("unsafe block must have a body")
	}

	// ✅ Traverse INTO the unsafe block to find declarations
	return a.collectDeclarations(block.Body)

	// Note: We don't create any symbol for the unsafe block itself
	// We just look for declarations inside it
}

// checkTypes performs type checking on the AST
func (a *Analyzer) checkTypes(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, decl := range n.Declarations {
			if err := a.checkTypes(decl); err != nil {
				return err
			}
		}

	case *ast.FuncDecl:
		return a.checkFunctionBody(n)
	default:
		return nil
	}
	return nil
}

// checkImmutability verifies immutability rules
func (a *Analyzer) checkImmutability(node ast.Node) error {
	// TODO: Implement immutability checking
	return nil
}

// typesCompatible checks if two types are compatible for assignment
func (a *Analyzer) typesCompatible(expected, actual *ast.Type) bool {
	if expected == nil || actual == nil {
		return false
	}

	// For now, simple base type comparison
	if expected.BaseType != actual.BaseType {
		return false
	}

	// Special handling for function types
	if expected.IsFunctionType() && actual.IsFunctionType() {
		return a.functionSignaturesCompatible(
			expected.GetFunctionSignature(),
			actual.GetFunctionSignature(),
		)
	}

	return true
}

// functionSignaturesCompatible checks if two function signatures are compatible
func (a *Analyzer) functionSignaturesCompatible(expected, actual *ast.FunctionSignature) bool {
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
		if !a.typesCompatible(expectedParam.Type, actualParam.Type) {
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

	return a.typesCompatible(expected.ReturnType, actual.ReturnType)
}

func (a *Analyzer) inferType(expr ast.Expression) *ast.Type {
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
