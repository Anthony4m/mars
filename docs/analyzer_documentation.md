# Mars Semantic Analyzer: A Deep Dive

This document explores the architecture and implementation of the Mars semantic analyzer. It's designed for new contributors to understand not just *what* the code does, but *why* it was designed this way, drawing heavily on the principles outlined in "The Road Not Taken: Why Mars Chose Manual Recursion Over the Visitor Pattern".

**Core Philosophy:** Pragmatism over purity, clarity over cleverness, and simplicity over sophistication. The analyzer prioritizes debuggability, performance, and maintainability.

## Core Components & Data Structures

### `Analyzer` Struct

```go
type Analyzer struct {
    errors          *errors.MarsReporter
    symbols         *SymbolTable
    types           *TypeChecker
    immutable       *ImmutabilityChecker
    sourceCode      string
    filename        string
    currentFunction *ast.FuncDecl
}
```

- **What it is:** This is the main state-bearing struct for the entire semantic analysis process. It holds references to all the necessary sub-components.
- **Why it exists:**
    - **Centralized State:** It acts as a "context" object, passed through the recursive analysis functions. This avoids global state and makes the data flow explicit.
    - **Sub-Component Registry:** It holds instances of helpers like the `SymbolTable` and `TypeChecker`, promoting separation of concerns. For example, all symbol lookups are delegated to `symbols`, keeping the `Analyzer` focused on traversal logic.
    - **`currentFunction`:** This field is crucial for context-sensitive checks, like validating `return` statements. It's a perfect example of the fine-grained state control that manual recursion provides.

### `New()` Function

```go
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
```

- **What it is:** A simple constructor for the `Analyzer`.
- **Why it exists:** It properly initializes the `Analyzer` and its sub-components. Encapsulating this setup ensures that an `Analyzer` is always in a valid state before analysis begins.

## The Two-Phase Analysis Pipeline

The analyzer's master plan is orchestrated by the `Analyze` function, which embodies the two-phase approach.

### `Analyze()` Function

```go
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
```

- **What it is:** The main entry point for semantic analysis. It executes the two main passes in sequence.
- **Why it exists:** This function is the heart of the two-phase strategy.
    1.  **Phase 1: `collectDeclarations`:** This pass walks the AST to find all function, struct, and variable declarations, populating the symbol table. This solves the "forward reference" problem, allowing functions to call other functions defined later in the source file (mutual recursion).
    2.  **Early Exit:** It checks for errors after the first pass. If declarations are invalid (e.g., duplicate names), there's no point proceeding to type-check the code that uses them. This makes error messages cleaner and more relevant.
    3.  **Phase 2: `CheckTypes` & `checkImmutability`:** With a complete symbol table, this pass can now validate expressions, function calls, and assignments, ensuring that types are used correctly and immutability rules are enforced.

## Phase 1: Declaration Collection

The first pass is a recursive descent that only cares about declarations.

### `collectDeclarations()` Function

```go
func (a *Analyzer) collectDeclarations(node ast.Node) error {
    switch n := node.(type) {
    case *ast.Program:
        for _, decl := range n.Declarations {
            if err := a.collectDeclarations(decl); err != nil {
                return err
            }
        }
    case *ast.VarDecl:
        return a.collectVariableDeclaration(n)
    case *ast.FuncDecl:
        return a.collectFunctionDeclaration(n)
    // ... other declaration types
    default:
        return nil // ✅ Ignore non-declaration nodes
    }
    return nil
}
```

- **What it is:** A classic example of manual recursion using a `switch` statement on the AST node type.
- **Why it's designed this way:**
    - **Direct and Readable:** The `switch` statement makes the control flow obvious. You can immediately see which node types are handled in this pass.
    - **Selective Traversal:** It only recurses into nodes that can contain declarations (like `*ast.Program`). It explicitly ignores everything else by having a `default` case that does nothing. This is far more direct than a Visitor, where you might have to implement empty `Visit...` methods for all the node types you don't care about.
    - **Delegation:** For each declaration type, it calls a specific, aptly-named helper function (e.g., `collectFunctionDeclaration`), keeping the logic for each declaration type clean and isolated.

### `collectFunctionDeclaration()` Function

```go
func (a *Analyzer) collectFunctionDeclaration(decl *ast.FuncDecl) error {
    // ... create function type ...
    if err := a.symbols.Define(decl.Name.Name, *funcType, false, true, decl); err != nil {
        // ... handle duplicate declaration error ...
    }
    // 3. Defer checking the body until the second pass.
    return nil
}
```

- **What it is:** This function adds a function's signature to the symbol table.
- **Why it's designed this way:**
    - **Signature Only:** Crucially, it *only* records the function's name and type signature. It **does not** analyze the function's body.
    - **Enabling Mutual Recursion:** By defining the function symbol before analyzing any function bodies, it ensures that when `CheckTypes` is called later, a call to `isEven()` inside `isOdd()` (or vice-versa) will find a valid symbol in the table.

## Phase 2: Type Checking

After all symbols are known, the second pass validates their usage.

### `CheckTypes()` Function

```go
func (a *Analyzer) CheckTypes(node ast.Node) error {
    switch n := node.(type) {
    case *ast.Program:
        // ... recurse ...
    case *ast.FuncDecl:
        return a.checkFunctionBody(n)
    case *ast.VarDecl:
        return a.CheckVarDecl(n)
    case *ast.BinaryExpression:
        return a.checkBinaryExpression(n)
    // ... many other expression and statement types
    }
    return nil
}
```

- **What it is:** Another manually recursive function that forms the backbone of the type-checking pass.
- **Why it's designed this way:**
    - **Comprehensive:** Unlike `collectDeclarations`, this function has cases for almost every type of statement and expression, because types need to be checked everywhere.
    - **Explicit Control Flow:** The order of operations is explicit. For an `*ast.IfStatement`, for example, the code first checks the condition, then the consequence, then the alternative. This direct control is a major advantage over the Visitor pattern's fixed traversal order.

### `checkFunctionBody()` Function

```go
func (a *Analyzer) checkFunctionBody(decl *ast.FuncDecl) error {
    // Enter function scope
    a.symbols.EnterScope()
    defer a.symbols.ExitScope()

    // ... set currentFunction ...

    // Add parameters to local scope
    for _, param := range decl.Signature.Parameters {
        // ... define parameter in symbol table ...
    }

    // Now check the body
    return a.CheckTypes(decl.Body)
}
```

- **What it is:** This function handles the type checking for the *inside* of a function.
- **Why it's designed this way:**
    - **Bulletproof Scoping:** It demonstrates the power of combining manual recursion with Go's `defer` statement. `a.symbols.EnterScope()` is called, and `a.symbols.ExitScope()` is guaranteed to run when the function exits, no matter how. This is simpler and less error-prone than manually managing scope entry/exit in a Visitor.
    - **Context Management:** It sets `a.currentFunction`, providing the necessary context for checking `return` statements within the body.

---

## Detailed Function Analysis

This section provides a deep dive into the core logic of the type-checking pass. These functions are called from the main `CheckTypes` recursive descent and are where the rules of the Mars language are enforced.

### `CheckVarDecl`

This function validates a variable declaration like `var x: int = 10` or `mut y := "hello"`.

```go
func (a *Analyzer) CheckVarDecl(decl *ast.VarDecl) error {
    // ...
    switch {
    // 1) explicit type + initializer → check compatibility
    case hasAnnot && hasInit:
        actual := a.types.inferType(decl.Value)
        if !a.types.typesCompatible(&declared, actual) {
            // ... report error ...
        }

    // 2) explicit type only → nothing more to check
    case hasAnnot:
        // ok

    // 3) initializer only (x := e) → infer
    case hasInit:
        actual := a.types.inferType(decl.Value)
        if actual.BaseType == "unknown" {
            // ... report error ...
        }

    // 4) neither → error
    default:
        // ... report error ...
    }
    return nil
}
```

-   **Purpose:** To ensure a variable declaration is well-formed and that the types are consistent.
-   **Student Focus: State Machine for Declarations:** This function is a great example of implementing a state machine for language rules. It exhaustively checks the four possible states of a variable declaration:
    1.  **`var x: int = 10` (Type and Initializer):** The most explicit form. The job here is to verify that the type of the initializer (`10` is `int`) is compatible with the declared type (`int`). It calls `inferType` on the right-hand side and `typesCompatible` to perform the check.
    2.  **`var x: int` (Type Only):** This is perfectly valid. The variable is declared and has a type, but no initial value. No further type checking is needed at this point.
    3.  **`x := 10` (Initializer Only):** This is type inference. The compiler must determine the type of the variable from the initializer. It calls `inferType` and if the type cannot be determined (e.g., from an empty array literal `[]`), it's an error.
    4.  **`var x` (Neither):** This is a syntax error in Mars. A variable needs either an explicit type or an initializer to infer one from.
-   **Design Choice:** Using a `switch` on boolean conditions (`hasAnnot`, `hasInit`) makes the logic clear and directly maps to the language's grammar rules.

### `CheckAssignment`

This function validates an assignment to an existing variable, like `x = 20`.

```go
func (a *Analyzer) CheckAssignment(stmt *ast.AssignmentStatement) error {
    // 1) resolve the variable
    sym, err := a.symbols.Resolve(stmt.Name.Name)
    if err != nil { /* ... error ... */ }

    // 2) if the symbol was declared immutable, error
    if !sym.IsMutable { /* ... error ... */ }

    // 3) type‐check the right‐hand side
    actual := a.types.inferType(stmt.Value)
    if !a.types.typesCompatible(&sym.Type, actual) { /* ... error ... */ }

    return nil
}
```

-   **Purpose:** To enforce the rules of assignment: the variable must exist, it must be mutable, and the new value's type must be compatible.
-   **Student Focus: Enforcing Language Semantics:**
    1.  **Symbol Resolution:** The first step in any operation involving a variable is to look it up in the symbol table (`a.symbols.Resolve`). This confirms the variable has been declared in the current scope.
    2.  **Mutability Check:** This is where a core feature of Mars—immutability by default—is enforced. The symbol table stores whether a variable was declared with `mut`. If `sym.IsMutable` is false, an error is reported. This demonstrates how semantic flags stored during declaration are used later.
    3.  **Type Compatibility:** Just like in a variable declaration, the analyzer infers the type of the right-hand side value and checks if it's compatible with the variable's declared type (stored in `sym.Type`).

### `checkFunctionCall`

This function validates a function call, like `myFunc(10, "hello")`.

```go
func (a *Analyzer) checkFunctionCall(call *ast.FunctionCall) error {
    // ...
    // 1. Resolve function name
    sym, err := a.symbols.Resolve(ident.Name)
    if err != nil { /* ... */ }

    // 2. Check if it's a function
    if !sym.IsFunction { /* ... */ }

    funcSig := sym.Type.GetFunctionSignature()

    // 3. Check argument count (arity)
    if len(call.Arguments) != len(funcSig.Parameters) { /* ... */ }

    // 4. Check argument types
    for i, arg := range call.Arguments {
        argType := a.types.inferType(arg)
        paramType := funcSig.Parameters[i].Type
        if !a.types.typesCompatible(paramType, argType) { /* ... */ }
    }
    return nil
}
```

-   **Purpose:** To ensure a function call is valid: the function exists, the number of arguments is correct, and the type of each argument matches the function's signature.
-   **Student Focus: Validating Signatures:**
    1.  **Symbol Resolution:** It starts by resolving the function's name.
    2.  **Type Kind Check:** It then checks `sym.IsFunction`. This is important because in Mars, variables and functions share the same namespace, so a user could try to "call" a variable. The symbol table must distinguish between different kinds of symbols.
    3.  **Arity Check:** A simple but critical check comparing the number of arguments provided to the number of parameters expected. This catches many common errors.
    4.  **Per-Argument Type Check:** It iterates through the arguments and parameters in lock-step, comparing the type of each argument (`inferType(arg)`) with the expected type from the signature (`paramType`). This demonstrates a practical application of type checking lists of related items.

### `checkStructLiteral`

This function validates the instantiation of a struct, like `MyStruct{field1: 10, field2: "hi"}`.

```go
func (a *Analyzer) checkStructLiteral(lit *ast.StructLiteral) error {
    // 1) Resolve the struct’s type symbol.
    sym, err := a.symbols.Resolve(lit.Type.Name)
    // ...

    // 2) Ensure it really is a struct.
    if sym.Type.StructName != lit.Type.Name { /* ... */ }

    // Build a map of declared fields → types.
    declared := make(map[string]*ast.Type, len(sym.Type.StructFields))
    for _, f := range sym.Type.StructFields {
        declared[f.Name.Name] = f.Type
    }

    seen := map[string]bool{}
    for _, init := range lit.Fields {
        // 3. Check for duplicate fields in the literal
        if seen[name] { /* ... */ }

        // 4. Check that the field exists on the struct
        expected, ok := declared[name]
        if !ok { /* ... */ }

        // 5. Type-check the initializer expression.
        actual := a.types.inferType(init.Value)
        if !a.types.typesCompatible(expected, actual) { /* ... */ }

        seen[name] = true
    }
    return nil
}
```

-   **Purpose:** To validate that a struct literal correctly matches its type definition.
-   **Student Focus: User-Defined Types:** This is a great example of how a compiler handles user-defined aggregate types.
    1.  **Type Resolution:** First, it resolves the struct's name (`MyStruct`) to get its definition from the symbol table.
    2.  **Field Map:** It then creates a hash map (`declared`) of the expected fields and their types from the struct's definition. This is a preparatory step for efficient lookups.
    3.  **Literal Field Iteration:** The function then iterates over the fields provided in the literal. For each field, it performs several checks:
        -   **Duplicate Check:** Ensures a field isn't provided twice in the same literal (`seen` map).
        -   **Existence Check:** Ensures the provided field actually exists in the `declared` map.
        -   **Type Check:** If the field exists, it checks that the type of the value provided in the literal is compatible with the field's expected type.

### `inferExpressionType`

This is a key utility function that determines the type of any given expression node.

```go
func (a *Analyzer) inferExpressionType(expr ast.Expression) *ast.Type {
    switch e := expr.(type) {
    case *ast.Literal:
        return a.types.inferType(expr) // Delegates to a helper

    case *ast.Identifier:
        symbol, err := a.symbols.Resolve(e.Name)
        if err != nil { return &ast.Type{BaseType: "unknown"} }
        return &symbol.Type

    case *ast.FunctionCall:
        // ... resolve function symbol ...
        // ... get signature ...
        return sig.ReturnType // The type of a call is its return type

    case *ast.BinaryExpression:
        // ...
        switch e.Operator {
        case "+", "-", "*", "/":
            return &ast.Type{BaseType: "float"} // Or int, simplified here
        case "==", "!=", "<", ">", "<=", ">=":
            return &ast.Type{BaseType: "bool"}
        }
    // ...
    }
}
```

-   **Purpose:** To compute and return the type of an expression without reporting errors itself. It's a query, not a validator.
-   **Student Focus: Recursive Type Inference:** This function is the embodiment of a recursive type inference algorithm.
    -   **Base Cases:** The recursion bottoms out at simple nodes. For a `Literal` (`10`, `"hi"`), the type is self-evident. For an `Identifier` (`x`), the type is found by looking it up in the symbol table.
    -   **Recursive Step:** For complex expressions, it determines its type based on the types of its sub-expressions.
        -   **`BinaryExpression`:** The type of `1 + 2` is `int`. The type of `1 < 2` is `bool`. This function implements those rules. It recursively calls `inferExpressionType` on its left and right children and then uses the operator to determine the resulting type.
        -   **`FunctionCall`:** The type of a function call expression is the return type of the function being called. This shows how different parts of the system connect: the result of a `checkFunctionCall` (the function's signature from the symbol table) is used here to determine the expression's type.
-   **Design Choice:** This function returns an `"unknown"` type when it can't figure something out. It does not add to `a.errors`. This is a good separation of concerns. The *calling* function (e.g., `CheckVarDecl`) is responsible for deciding if an unknown type is an error in its context and reporting it.

## Conclusion: The Pragmatic Path

The Mars analyzer is a case study in choosing the right tool for the job. By eschewing the "standard" Visitor pattern in favor of manual recursion with a two-phase analysis, the design achieves:

-   **Simplicity & Readability:** The control flow is immediately obvious from reading the code. There is no "magic" dispatch.
-   **Debuggability:** Setting a breakpoint in a `case` block shows you exactly what code is running for a given node type, and the call stack clearly shows how you got there.
-   **Flexibility:** Fine-grained control over traversal and state (like `currentFunction` and scopes) is trivial to implement.
-   **Performance:** Direct function calls and `switch` statements are more efficient than the interface-based calls required by the Visitor pattern.

This approach is a direct reflection of the project's values: build software that is easy to understand, maintain, and extend, even if it means taking the road less traveled by textbooks.