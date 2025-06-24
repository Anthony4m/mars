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
	if err := a.CheckExpression(node); err != nil {
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

func (a *Analyzer) CheckVarDecl(decl *ast.VarDecl) {
	// 1) resolve symbol (must have been defined in pass 1)
	sym, err := a.symbols.Resolve(decl.Name.Name)
	if err != nil {
		help := fmt.Sprintf("declare variable %q before use", decl.Name.Name)
		a.errors.AddErrorWithHelp(
			decl.Name.Position,
			errors.ErrCodeUndefinedVar,
			fmt.Sprintf("variable %q not found", decl.Name.Name),
			help,
		)
		return
	}

	declared := sym.Type // ast.Type
	hasAnnot := decl.Type != nil
	hasInit := decl.Value != nil

	switch {
	// 1) explicit type + initializer → check compatibility
	case hasAnnot && hasInit:
		actual := a.types.inferType(decl.Value)
		if !a.types.typesCompatible(&declared, actual) {
			help := fmt.Sprintf("cast the value to %s or change the variable’s type", declared.BaseType)
			a.errors.AddErrorWithHelp(
				decl.Name.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("mismatched types: expected %s, found %s", declared.BaseType, actual.BaseType),
				help,
			)
		}

	// 2) explicit type only → nothing more to check
	case hasAnnot:
		// ok

	// 3) initializer only (x := e) → infer
	case hasInit:
		actual := a.types.inferType(decl.Value)
		if actual.BaseType == "unknown" {
			a.errors.AddError(
				decl.Name.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("cannot infer type for variable %q", decl.Name.Name),
			)
		}

	// 4) neither → error
	default:
		a.errors.AddErrorWithHelp(
			decl.Name.Position,
			errors.ErrCodeSyntaxError,
			fmt.Sprintf("variable %q needs a type annotation or an initializer", decl.Name.Name),
			"add `: T` or `:= value` to the declaration",
		)
	}
}

func (a *Analyzer) collectVariableDeclaration(decl *ast.VarDecl) error {
	var varType ast.Type
	if decl.Type != nil {
		varType = *decl.Type
	} else if decl.Value != nil {
		varType = *a.types.inferType(decl.Value)
	}

	if err := a.symbols.Define(decl.Name.Name, varType, decl.Mutable, false, decl); err != nil {
		a.errors.AddErrorWithHelp(
			decl.Name.Position,
			errors.ErrCodeDuplicateDecl,
			fmt.Sprintf("variable '%s' is already defined in this scope", decl.Name.Name),
			"give this variable a different name",
		)

		// Add note about original declaration
		if original, _ := a.symbols.Resolve(decl.Name.Name); original != nil && original.DeclaredAt != nil {
			//TODO: Think about what to do about original.DeclaredAt
			//origPos := (*original.DeclaredAt).Pos()
			// You could enhance this to show the original location
		}
	}

	return nil
}

func (a *Analyzer) collectFunctionDeclaration(decl *ast.FuncDecl) error {
	// 1. Create the function's type from the signature.
	// TODO: This needs to be enhanced to handle:
	// - Type cycles (e.g. mutually recursive function types)
	// - Future generic type parameters
	// - Return type inference
	// For now we assume the parser has validated the basic signature structure
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
	return a.CheckExpression(decl.Body)
}

func (a *Analyzer) collectStructDeclaration(decl *ast.StructDecl) error {
	if decl.Name == nil {
		return fmt.Errorf("struct declaration must have a name")
	}

	// 1. Create a proper struct type with the struct name and fields
	structType := ast.NewStructType(decl.Name.Name, decl.Fields)

	// 2. Add the struct type to the symbol table immediately (for cycle safety)
	err := a.symbols.Define(decl.Name.Name, *structType, false, false, decl)
	if err != nil {
		a.errors.AddErrorWithHelp(
			decl.Name.Position,
			errors.ErrCodeDuplicateDecl,
			fmt.Sprintf("struct '%s' is already defined in this scope", decl.Name.Name),
			"give this struct a different name",
		)
		return err
	}

	// 3. Process struct fields for validation
	if decl.Fields != nil {
		// Track field names for uniqueness checking
		fieldNames := make(map[string]bool)

		for _, field := range decl.Fields {
			if field.Name == nil || field.Type == nil {
				a.errors.AddError(
					field.Position,
					errors.ErrCodeSyntaxError,
					"struct field must have both a name and type",
				)
				continue
			}

			// Check for duplicate field names
			if fieldNames[field.Name.Name] {
				a.errors.AddErrorWithHelp(
					field.Name.Position,
					errors.ErrCodeDuplicateDecl,
					fmt.Sprintf("duplicate field name '%s' in struct '%s'", field.Name.Name, decl.Name.Name),
					"each field name must be unique within a struct",
				)
				continue
			}

			if field.Type.BaseType == "" && field.Type.StructName == "" &&
				field.Type.ArrayType == nil && field.Type.PointerType == nil {
				a.errors.AddError(
					field.Type.Position,
					errors.ErrCodeTypeError,
					fmt.Sprintf("invalid type for field '%s'", field.Name.Name),
				)
			}

			fieldNames[field.Name.Name] = true

			// Note: We don't define fields in the current scope here
			// Fields are only accessible through struct instances, not in the struct declaration scope
		}
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

// CheckExpression performs type checking on the AST
func (a *Analyzer) CheckExpression(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, decl := range n.Declarations {
			if err := a.CheckExpression(decl); err != nil {
				return err
			}
		}

	case *ast.FuncDecl:
		return a.checkFunctionBody(n)
	case *ast.VarDecl:
		a.CheckVarDecl(n)
		return nil
	case ast.Expression:
		return nil
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

// TODO: Verify this correctness
func (a *Analyzer) CheckAssignment(stmt *ast.AssignmentStatement) {
	// 1) resolve the variable
	sym, err := a.symbols.Resolve(stmt.Name.Name)
	if err != nil {
		// undefined‐variable error
		return
	}

	// 2) if the symbol was declared immutable, error
	if !sym.IsMutable {
		a.errors.AddError(
			stmt.Name.Position,
			errors.ErrCodeImmutable,
			fmt.Sprintf("cannot assign to immutable variable %q", stmt.Name.Name),
		)
		// still continue so we can report other errors
	}

	// 3) type‐check the right‐hand side
	actual := a.types.inferType(stmt.Value)
	if !a.types.typesCompatible(&sym.Type, actual) {
		a.errors.AddError(
			stmt.Name.Position,
			errors.ErrCodeTypeError,
			fmt.Sprintf("cannot assign %s to %s", actual.BaseType, sym.Type.BaseType),
		)
	}
}
