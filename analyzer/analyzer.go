package analyzer

import (
	"fmt"
	"mars/ast"
	"mars/errors"
	"os"
)

func debugLog(msg string) {
	f, err := os.OpenFile("analyzer_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(msg + "\n")
	}
}

// Analyzer performs semantic analysis on the AST
type Analyzer struct {
	errors          *errors.MarsReporter
	symbols         *SymbolTable
	types           *TypeChecker
	immutable       *ImmutabilityChecker
	sourceCode      string
	filename        string
	currentFunction *ast.FuncDecl
	inUnsafeContext bool
	inLoopContext   bool
}

// New creates a new analyzer instance
func New(sourceCode, filename string) *Analyzer {
	return &Analyzer{
		errors:          errors.NewMarsReporter(sourceCode, filename),
		symbols:         NewSymbolTable(),
		types:           NewTypeChecker(),
		immutable:       NewImmutabilityChecker(),
		sourceCode:      sourceCode,
		filename:        filename,
		currentFunction: nil,
		inUnsafeContext: false,
		inLoopContext:   false,
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
	if err := a.CheckTypes(node); err != nil {
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
	case *ast.BlockStatement:
		// Only collect declarations from all statements, do not enter/exit scope in first pass
		for _, stmt := range n.Statements {
			if err := a.collectDeclarations(stmt); err != nil {
				return err
			}
		}
	case *ast.IfStatement:
		// Collect declarations from if statement blocks
		if n.Consequence != nil {
			if err := a.collectDeclarations(n.Consequence); err != nil {
				return err
			}
		}
		if n.Alternative != nil {
			if err := a.collectDeclarations(n.Alternative); err != nil {
				return err
			}
		}
	case *ast.ForStatement:
		// Collect declarations from for statement blocks
		if n.Init != nil {
			if err := a.collectDeclarations(n.Init); err != nil {
				return err
			}
		}
		if n.Body != nil {
			if err := a.collectDeclarations(n.Body); err != nil {
				return err
			}
		}
		if n.Post != nil {
			if err := a.collectDeclarations(n.Post); err != nil {
				return err
			}
		}
	default:
		return nil // ✅ Ignore non-declaration nodes
	}

	return nil
}

func (a *Analyzer) CheckVarDecl(decl *ast.VarDecl) error {
	// DEBUG: Log variable being resolved and current scope pointer
	debugLog(fmt.Sprintf("[DEBUG] Resolving variable '%s' in scope %p", decl.Name.Name, a.symbols.CurrentScope))
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
		return nil
	}

	declared := sym.Type // ast.Type
	hasAnnot := decl.Type != nil
	hasInit := decl.Value != nil

	// Check the initializer expression for errors (e.g., struct literal errors)
	if hasInit {
		a.CheckTypes(decl.Value)
	}

	switch {
	// 1) explicit type + initializer → check compatibility
	case hasAnnot && hasInit:
		actual := a.inferExpressionType(decl.Value)
		if !a.types.typesCompatible(&declared, actual) {
			help := fmt.Sprintf("cast the value to %s or change the variable's type", declared.BaseType)
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
		actual := a.inferExpressionType(decl.Value)
		if actual.BaseType == "unknown" {
			a.errors.AddErrorWithHelp(
				decl.Name.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("cannot infer type for variable %q", decl.Name.Name),
				fmt.Sprintf("variable may not exist, declare %s before using", decl.Name.Name),
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
	return nil
}

func (a *Analyzer) collectVariableDeclaration(decl *ast.VarDecl) error {
	// DEBUG: Log variable being defined and current scope pointer
	debugLog(fmt.Sprintf("[DEBUG] Defining variable '%s' in scope %p", decl.Name.Name, a.symbols.CurrentScope))
	var varType ast.Type
	if decl.Type != nil {
		varType = *decl.Type
	} else if decl.Value != nil {
		inferredType := a.inferExpressionType(decl.Value)
		if inferredType != nil {
			varType = *inferredType
		} else {
			// If we can't infer the type, use unknown
			varType = ast.Type{BaseType: "unknown"}
		}
	} else {
		// No type annotation and no initializer - this should be caught by the parser
		varType = ast.Type{BaseType: "unknown"}
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

func (a *Analyzer) CheckLiteral(lit *ast.Literal) error {
	// inferType accepts an Expression, so pass the nod e itself
	inferred := a.types.inferType(lit)
	if inferred.BaseType == "unknown" {
		a.errors.AddError(
			lit.Position,
			errors.ErrCodeTypeError,
			fmt.Sprintf("cannot determine type for literal %v", lit.Value),
		)
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

	prevFunc := a.currentFunction
	a.currentFunction = decl
	defer func() { a.currentFunction = prevFunc }()

	// Add parameters to local scope
	for _, param := range decl.Signature.Parameters {
		err := a.symbols.Define(param.Name.Name, *param.Type, false, false, param.Name)
		if err != nil {
			return err
		}
	}

	// Now check the body - recursion is safe because
	// the function name is already in the symbol table
	return a.CheckTypes(decl.Body)
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

	//Traverse INTO the unsafe block to find declarations
	return a.collectDeclarations(block.Body)

	// Note: We don't create any symbol for the unsafe block itself
	// We just look for declarations inside it
}

// CheckTypes performs type checking on the AST
func (a *Analyzer) CheckTypes(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, decl := range n.Declarations {
			if err := a.CheckTypes(decl); err != nil {
				return err
			}
		}

	case *ast.FuncDecl:
		return a.checkFunctionBody(n)
	case *ast.VarDecl:
		return a.CheckVarDecl(n)
	case *ast.Literal:
		return a.CheckLiteral(n)
	case *ast.Identifier:
		return a.checkIdentifier(n)
	case *ast.BinaryExpression:
		return a.checkBinaryExpression(n)
	case *ast.ExpressionStatement:
		// Check the expression within the statement
		if n.Expression != nil {
			return a.CheckTypes(n.Expression)
		}
	case *ast.BlockStatement:
		a.symbols.EnterScope()
		defer a.symbols.ExitScope()
		// Check all statements in a block
		for _, stmt := range n.Statements {
			if err := a.CheckTypes(stmt); err != nil {
				return err
			}
		}
	case *ast.FunctionCall:
		return a.checkFunctionCall(n)

	case *ast.StructLiteral:
		return a.checkStructLiteral(n)

	case *ast.IfStatement:
		// Check condition
		if err := a.CheckTypes(n.Condition); err != nil {
			return err
		}
		condType := a.inferExpressionType(n.Condition)
		if condType.BaseType != "bool" {
			a.errors.AddError(
				n.Condition.Pos(),
				errors.ErrCodeTypeError,
				fmt.Sprintf("condition must be boolean, found '%s'", condType.BaseType),
			)
		}

		// Check consequence block
		if n.Consequence != nil {
			if err := a.CheckTypes(n.Consequence); err != nil {
				return err
			}
		}

		// Check alternative block (else)
		if n.Alternative != nil {
			if err := a.CheckTypes(n.Alternative); err != nil {
				return err
			}
		}

		return nil
	case *ast.ReturnStatement:
		if a.currentFunction == nil {
			a.errors.AddError(
				n.Position,
				errors.ErrCodeSyntaxError,
				"return statement outside function",
			)
			return nil
		}

		sig := a.currentFunction.Signature
		if n.Value == nil {
			// Empty return
			if sig.ReturnType != nil {
				a.errors.AddError(
					n.Position,
					errors.ErrCodeTypeError,
					fmt.Sprintf("function '%s' should return '%s'",
						a.currentFunction.Name.Name, sig.ReturnType.String()),
				)
			}
		} else {
			// Check return value
			if err := a.CheckTypes(n.Value); err != nil {
				return err
			}

			if sig.ReturnType == nil {
				a.errors.AddError(
					n.Position,
					errors.ErrCodeTypeError,
					"function has no return type but returns a value",
				)
			} else {
				returnType := a.inferExpressionType(n.Value)
				if !a.types.typesCompatible(sig.ReturnType, returnType) {
					a.errors.AddError(
						n.Position,
						errors.ErrCodeTypeError,
						fmt.Sprintf("cannot return '%s' from function with return type '%s'",
							returnType.String(), sig.ReturnType.String()),
					)
				}
			}
		}
		return nil
	case *ast.UnaryExpression:
		return a.CheckTypes(n.Right)
	case *ast.AssignmentStatement:
		return a.CheckAssignment(n)

	case *ast.ForStatement:
		// Enter loop context
		prevLoopContext := a.inLoopContext
		a.inLoopContext = true
		defer func() { a.inLoopContext = prevLoopContext }()

		// Check initialization, condition, and post statements
		if n.Init != nil {
			if err := a.CheckTypes(n.Init); err != nil {
				return err
			}
		}
		if n.Condition != nil {
			if err := a.CheckTypes(n.Condition); err != nil {
				return err
			}
			// Verify condition is boolean
			condType := a.inferExpressionType(n.Condition)
			if condType.BaseType != "bool" {
				a.errors.AddError(
					n.Condition.Pos(),
					errors.ErrCodeTypeError,
					"for loop condition must be boolean",
				)
			}
		}
		if n.Post != nil {
			if err := a.CheckTypes(n.Post); err != nil {
				return err
			}
		}
		// Check body
		return a.CheckTypes(n.Body)

	case *ast.PrintStatement:
		// Check the expression being printed
		if n.Expression != nil {
			return a.CheckTypes(n.Expression)
		}

	case *ast.BreakStatement:
		if !a.inLoopContext {
			a.errors.AddErrorWithHelp(
				n.Position,
				errors.ErrCodeSyntaxError,
				"break statement outside loop",
				"break can only be used inside for loops",
			)
		}
		return nil

	case *ast.ContinueStatement:
		if !a.inLoopContext {
			a.errors.AddErrorWithHelp(
				n.Position,
				errors.ErrCodeSyntaxError,
				"continue statement outside loop",
				"continue can only be used inside for loops",
			)
		}
		return nil

	case *ast.ArrayLiteral:
		if len(n.Elements) == 0 {
			return nil
		}
		// Check all elements have compatible types
		for _, elem := range n.Elements {
			if err := a.CheckTypes(elem); err != nil {
				return err
			}
		}
		expectedType := a.inferExpressionType(n.Elements[0])
		for i := 1; i < len(n.Elements); i++ {
			elem := n.Elements[i]
			actualType := a.inferExpressionType(elem)

			if !a.types.typesCompatible(expectedType, actualType) {
				a.errors.AddErrorWithHelp(
					elem.Pos(),
					errors.ErrCodeTypeError,
					fmt.Sprintf("mismatched types in array literal: found '%s', expected '%s'",
						actualType.String(), expectedType.String()),
					"all elements in an array literal must have the same type",
				)
			}
		}
		return nil
	case *ast.IndexExpression:
		// Check array and index
		if err := a.CheckTypes(n.Object); err != nil {
			return err
		}
		if err := a.CheckTypes(n.Index); err != nil {
			return err
		}
		// Verify index is integer
		indexType := a.inferExpressionType(n.Index)
		if indexType.BaseType != "int" {
			a.errors.AddError(
				n.Index.Pos(),
				errors.ErrCodeTypeError,
				"array index must be integer",
			)
		}
		return nil

	case *ast.MemberExpression:
		// Check struct field access
		if err := a.CheckTypes(n.Object); err != nil {
			return err
		}
		objectType := a.inferExpressionType(n.Object)

		// Check if the object's type is a struct.
		if objectType.StructName == "" { // If StructName is empty, it's not a struct type
			a.errors.AddErrorWithHelp(
				n.Object.Pos(),
				errors.ErrCodeTypeError,
				fmt.Sprintf("cannot access field '%s' on non-struct type '%s'",
					n.Property.Name, objectType.String()),
				"member access is only allowed on struct types",
			)
			return nil // Stop further checking for this expression as it's fundamentally wrong
		}
		// Resolve the struct type symbol to get its field definitions.
		structSym, err := a.symbols.Resolve(objectType.StructName)
		if err != nil {
			// This error means the struct type itself was not found in any scope.
			a.errors.AddError(
				n.Object.Pos(),
				errors.ErrCodeUndefinedType,
				fmt.Sprintf("internal error: struct type '%s' not found in symbol table", objectType.StructName),
			)
			return nil // Stop processing to prevent a panic.
		}

		// Search for the property (field) within the struct's fields.
		var foundField *ast.FieldDecl // Use a pointer, which is nil by default.
		for _, field := range structSym.Type.StructFields {
			if field.Name.Name == n.Property.Name {
				foundField = field // Assign the pointer directly.
				break
			}
		}
		// If the field was not found, report an error.
		if foundField == nil { // The check is now safe and clear.
			a.errors.AddErrorWithHelp(
				n.Property.Pos(),             // Point the error at the field name itself.
				errors.ErrCodeUndefinedField, // Use a more specific error code.
				fmt.Sprintf("field '%s' does not exist on type '%s'",
					n.Property.Name, objectType.String()),
				fmt.Sprintf("check the spelling or define '%s' in struct '%s'",
					n.Property.Name, objectType.StructName),
			)
			return nil
		}
		return nil
	case *ast.UnsafeBlock:
		a.symbols.EnterScope()
		defer a.symbols.ExitScope()

		//set the unsafe context flag
		prevUnsafeState := a.inUnsafeContext
		a.inUnsafeContext = true
		defer func() { a.inUnsafeContext = prevUnsafeState }()
		// Check all statements in a block
		return a.CheckTypes(n.Body)
		// Example: If I had a new AST node for pointer dereference
	//case *ast.PointerDereferencing: // Assuming I add this AST node later
	//	if !a.isInUnsafeBlock {
	//		a.errors.AddErrorWithHelp(
	//			n.Position,
	//			errors.ErrCodeUnsafeOperation, // I might need a new error code
	//			"pointer dereference is only allowed inside an 'unsafe' block",
	//			"wrap this operation in an 'unsafe { ... }' block",
	//		)
	//	}
	// ... then proceed with type checking the dereferenced expression ...
	//return a.CheckTypes(n.Expression) // Assuming n.Expression is the pointer
	default:
		return nil
	}
	return nil
}

func (a *Analyzer) checkIdentifier(ident *ast.Identifier) error {
	_, err := a.symbols.Resolve(ident.Name)
	if err != nil {
		a.errors.AddErrorWithHelp(ident.Position, errors.ErrCodeUndefinedVar, fmt.Sprintf("undefined  %q", ident.Name), "variable must be defined before use")
	}
	return nil
}

// checkImmutability verifies immutability rules
func (a *Analyzer) checkImmutability(node ast.Node) error {
	// TODO: Implement immutability checking
	return nil
}

func (a *Analyzer) CheckAssignment(stmt *ast.AssignmentStatement) error {
	// 1) resolve the variable
	sym, err := a.symbols.Resolve(stmt.Name.Name)
	if err != nil {
		// undefined‐variable error
		return err
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
	actual := a.inferExpressionType(stmt.Value)
	if !a.types.typesCompatible(&sym.Type, actual) {
		a.errors.AddError(
			stmt.Name.Position,
			errors.ErrCodeTypeError,
			fmt.Sprintf("cannot assign %s to %s", actual.BaseType, sym.Type.BaseType),
		)
	}
	return nil
}

func (a *Analyzer) checkBinaryExpression(expr *ast.BinaryExpression) error {
	if err := a.CheckTypes(expr.Left); err != nil {
		return err
	}
	if err := a.CheckTypes(expr.Right); err != nil {
		return err
	}

	leftType := a.inferExpressionType(expr.Left)
	rightType := a.inferExpressionType(expr.Right)
	switch expr.Operator {
	case "+", "-", "*", "/", "%":
		if !isNumericType(leftType) || !isNumericType(rightType) {
			a.errors.AddError(
				expr.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("invalid operation: %s %s %s (operator %s not defined on %s)",
					leftType.String(), expr.Operator, rightType.String(),
					expr.Operator, leftType.String()),
			)
		}
	case "==", "!=":
		if !a.types.typesCompatible(leftType, rightType) {
			a.errors.AddError(
				expr.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("invalid operation: %s %s %s (mismatched types)",
					leftType.String(), expr.Operator, rightType.String()),
			)
		}
	case "<", ">", "<=", ">=":
		// Comparison needs ordered types
		if !isOrderedType(leftType) || !isOrderedType(rightType) {
			a.errors.AddError(
				expr.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("invalid operation: %s %s %s (operator %s not defined on %s)",
					leftType.String(), expr.Operator, rightType.String(),
					expr.Operator, leftType.String()),
			)
		}
	case "&&":
		if !isBool(leftType) && !isBool(rightType) {
			a.errors.AddError(
				expr.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("types mismatch: %s %s %s (operator %s not defined on %s)",
					leftType.String(), expr.Operator, rightType.String(),
					expr.Operator, leftType.String()),
			)
		}

	case "||":
		if !isBool(leftType) || !isBool(rightType) {
			a.errors.AddError(
				expr.Position,
				errors.ErrCodeTypeError,
				fmt.Sprintf("logical operators require boolean operands"),
			)
		}

	}
	return nil
}

func (a *Analyzer) checkFunctionCall(call *ast.FunctionCall) error {
	if err := a.CheckTypes(call.Function); err != nil {
		return err
	}

	for _, arg := range call.Arguments {
		if err := a.CheckTypes(arg); err != nil {
			return err
		}
	}

	ident, ok := call.Function.(*ast.Identifier)
	if !ok {
		// TODO: Handle method calls, function expressions
		return nil
	}

	sym, err := a.symbols.Resolve(ident.Name)
	if err != nil {
		return err
	}

	if !sym.IsFunction {
		a.errors.AddError(
			call.Position,
			errors.ErrCodeTypeError,
			fmt.Sprintf("'%s' is not a function", ident.Name),
		)
		return nil
	}

	funcSig := sym.Type.GetFunctionSignature()
	if funcSig == nil {
		return nil
	}

	//check argument count
	if len(call.Arguments) != len(funcSig.Parameters) {
		a.errors.AddErrorWithHelp(
			call.Position,
			errors.ErrCodeTypeError,
			fmt.Sprintf("wrong number of arguments in call to '%s'", ident.Name),
			fmt.Sprintf("expected %d arguments, got %d",
				len(funcSig.Parameters), len(call.Arguments)),
		)
		return nil
	}

	//chack argument types
	for i, arg := range call.Arguments {
		argType := a.inferExpressionType(arg)
		paramType := funcSig.Parameters[i].Type

		if !a.types.typesCompatible(paramType, argType) {
			paramName := funcSig.Parameters[i].Name.Name
			a.errors.AddErrorWithHelp(
				arg.Pos(),
				errors.ErrCodeTypeError,
				fmt.Sprintf("cannot use '%s' as type '%s' in argument to '%s'",
					argType.String(), paramType.String(), ident.Name),
				fmt.Sprintf("parameter '%s' expects type '%s'",
					paramName, paramType.String()),
			)
		}
	}
	return nil
}

func (a *Analyzer) checkStructLiteral(lit *ast.StructLiteral) error {
	// 1) Resolve the struct's type symbol.
	sym, err := a.symbols.Resolve(lit.Type.Name)
	if err != nil {
		a.errors.AddError(lit.Type.Position,
			errors.ErrCodeUndefinedType,
			fmt.Sprintf("unknown type %q", lit.Type.Name),
		)
		return nil
	}

	// 2) Ensure it really is a struct.
	if sym.Type.StructName != lit.Type.Name {
		a.errors.AddError(lit.Type.Position,
			errors.ErrCodeTypeError,
			fmt.Sprintf("%q is not a struct type", lit.Type.Name),
		)
		return nil
	}

	// Build a map of declared fields → types.
	declared := make(map[string]*ast.Type, len(sym.Type.StructFields))
	for _, f := range sym.Type.StructFields {
		declared[f.Name.Name] = f.Type
	}

	seen := map[string]bool{}
	for _, init := range lit.Fields {
		name := init.Name.Name

		// Duplicate in *this* literal?
		if seen[name] {
			a.errors.AddError(init.Position,
				errors.ErrCodeDuplicateDecl,
				fmt.Sprintf("duplicate field %q in literal", name),
			)
			continue
		}

		expected, ok := declared[name]
		if !ok {
			a.errors.AddError(init.Position,
				errors.ErrCodeUndefinedField,
				fmt.Sprintf("field %q does not exist on %s", name, lit.Type.Name),
			)
			continue
		}

		// 3) Type-check the initializer expression.
		actual := a.inferExpressionType(init.Value)
		if !a.types.typesCompatible(expected, actual) {
			a.errors.AddErrorWithHelp(
				init.Value.Pos(),
				errors.ErrCodeTypeError,
				fmt.Sprintf("cannot use %s to initialize field %q (type is %s)",
					actual.BaseType, name, expected.BaseType),
				fmt.Sprintf("field %q expects %s", name, expected.BaseType),
			)
		}

		seen[name] = true
	}

	return nil
}

func (a *Analyzer) inferExpressionType(expr ast.Expression) *ast.Type {
	if expr == nil {
		return &ast.Type{BaseType: "unknown"}
	}

	switch e := expr.(type) {
	case *ast.Literal:
		return a.types.inferType(expr)

	case *ast.Identifier:
		symbol, err := a.symbols.Resolve(e.Name)
		if err != nil {
			return &ast.Type{BaseType: "unknown"}
		}
		return &symbol.Type

	case *ast.FunctionCall:
		// Get return type of the function
		if ident, ok := e.Function.(*ast.Identifier); ok {
			symbol, err := a.symbols.Resolve(ident.Name)
			if err == nil && symbol.IsFunction {
				sig := symbol.Type.GetFunctionSignature()
				if sig != nil && sig.ReturnType != nil {
					return sig.ReturnType
				}
			}
		}
		return &ast.Type{BaseType: "void"} // No return type

	case *ast.BinaryExpression:
		leftType := a.inferExpressionType(e.Left)
		rightType := a.inferExpressionType(e.Right)

		switch e.Operator {
		case "+", "-", "*", "/", "%":
			// Arithmetic operators
			if leftType.BaseType == "float" || rightType.BaseType == "float" {
				return &ast.Type{BaseType: "float"}
			}
			// Both operands are integers
			if leftType.BaseType == "int" && rightType.BaseType == "int" {
				return &ast.Type{BaseType: "int"}
			}
			// If we get here, one or both operands are not numeric
			return &ast.Type{BaseType: "unknown"}

		case "==", "!=", "<", ">", "<=", ">=":
			// Comparison operators
			return &ast.Type{BaseType: "bool"}

		case "&&", "||":
			// Logical operators
			return &ast.Type{BaseType: "bool"}

		default:
			return &ast.Type{BaseType: "unknown"}
		}

	case *ast.UnaryExpression:
		switch e.Operator {
		case "!":
			return &ast.Type{BaseType: "bool"}
		case "-":
			return a.inferExpressionType(e.Right)
		default:
			return &ast.Type{BaseType: "unknown"}
		}

	case *ast.ArrayLiteral:
		// Infer array type from first element
		if len(e.Elements) > 0 {
			elemType := a.inferExpressionType(e.Elements[0])
			return &ast.Type{
				ArrayType: elemType,
				// ArraySize could be set to len(e.Elements) for fixed arrays
			}
		}
		return &ast.Type{BaseType: "unknown"}

	case *ast.StructLiteral:
		// For struct literals, return the struct type
		return &ast.Type{
			StructName: e.Type.Name,
		}

	default:
		return &ast.Type{BaseType: "unknown"}
	}
}

func isBool(t *ast.Type) bool {
	return t.BaseType == "bool"
}

// Helper functions
func isNumericType(t *ast.Type) bool {
	return t.BaseType == "int" || t.BaseType == "float"
}

func isOrderedType(t *ast.Type) bool {
	return isNumericType(t) || t.BaseType == "string"
}
