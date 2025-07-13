# Mars Compiler Quick Reference

## Common Patterns

### 1. Adding Position Tracking to AST Nodes

```go
// When creating any AST node, always set the position
node := &ast.VarDecl{
    Name: &ast.Identifier{
        Name:     "x",
        Position: ast.Position{Line: 1, Column: 1},
    },
    Type: &ast.Type{BaseType: "int"},
    Position: ast.Position{Line: 1, Column: 1},
}
```

### 2. Error Reporting Pattern

```go
// Always use the MarsReporter for consistent error handling
a.errors.AddErrorWithHelp(
    node.Pos(),
    errors.ErrCodeTypeError,
    "type mismatch: expected int, got string",
    "cast the value or change the variable type",
)
```

### 3. Type Compatibility Checking

```go
// Use the typesCompatible method for all type comparisons
if !a.typesCompatible(expectedType, actualType) {
    // Handle type mismatch
}

// For function types, use functionSignaturesCompatible
if expected.IsFunctionType() && actual.IsFunctionType() {
    if !a.functionSignaturesCompatible(
        expected.GetFunctionSignature(),
        actual.GetFunctionSignature(),
    ) {
        // Handle function type mismatch
    }
}
```

### 4. Symbol Table Usage

```go
// Define a symbol
err := a.symbols.Define(name, type, mutable, isFunction, declaredAt)

// Resolve a symbol
symbol, exists := a.symbols.Resolve(name)

// Enter/exit scopes
a.symbols.EnterScope()
defer a.symbols.ExitScope()
```

## Code Examples

### 1. Parsing a Variable Declaration

```go
func (p *Parser) parseVarDecl() *ast.VarDecl {
    pos := p.curToken.Position
    
    // Check for 'mut' keyword
    mutable := false
    if p.curToken.Type == lexer.MUT {
        mutable = true
        p.nextToken()
    }
    
    // Parse identifier
    name := &ast.Identifier{
        Name:     p.curToken.Literal,
        Position: p.curToken.Position,
    }
    p.nextToken()
    
    // Parse type (if present)
    var varType *ast.Type
    if p.curToken.Type == lexer.COLON {
        p.nextToken()
        varType = p.parseType()
    }
    
    // Parse value (if present)
    var value ast.Expression
    if p.curToken.Type == lexer.ASSIGN {
        p.nextToken()
        value = p.parseExpression(LOWEST)
    }
    
    return &ast.VarDecl{
        Mutable:  mutable,
        Name:     name,
        Type:     varType,
        Value:    value,
        Position: pos,
    }
}
```

### 2. Type Checking a Function Call

```go
func (a *Analyzer) checkFunctionCall(call *ast.FunctionCall) *ast.Type {
    // Check function expression
    funcType := a.checkExpression(call.Function)
    if funcType == nil {
        return nil
    }
    
    // Ensure it's a function type
    if !funcType.IsFunctionType() {
        a.errors.AddError(
            call.Position,
            errors.ErrCodeTypeError,
            fmt.Sprintf("cannot call non-function type %s", funcType.String()),
        )
        return nil
    }
    
    signature := funcType.GetFunctionSignature()
    
    // Check argument count
    if len(call.Arguments) != len(signature.Parameters) {
        a.errors.AddError(
            call.Position,
            errors.ErrCodeTypeError,
            fmt.Sprintf("function expects %d arguments, got %d",
                len(signature.Parameters), len(call.Arguments)),
        )
        return nil
    }
    
    // Check argument types
    for i, arg := range call.Arguments {
        argType := a.checkExpression(arg)
        if argType == nil {
            continue
        }
        
        expectedType := signature.Parameters[i].Type
        if !a.typesCompatible(expectedType, argType) {
            a.errors.AddError(
                arg.Pos(),
                errors.ErrCodeTypeError,
                fmt.Sprintf("argument %d: expected %s, got %s",
                    i+1, expectedType.String(), argType.String()),
            )
        }
    }
    
    return signature.ReturnType
}
```

### 3. Creating Function Types

```go
// Create a function type from parameters and return type
func createFunctionType(params []*ast.Parameter, returnType *ast.Type) *ast.Type {
    signature := &ast.FunctionSignature{
        Parameters: params,
        ReturnType: returnType,
    }
    return ast.NewFunctionType(signature)
}

// Example usage
funcType := createFunctionType(
    []*ast.Parameter{
        {Name: &ast.Identifier{Name: "x"}, Type: ast.NewBaseType("int")},
        {Name: &ast.Identifier{Name: "y"}, Type: ast.NewBaseType("int")},
    },
    ast.NewBaseType("int"),
)
```

## Troubleshooting

### 1. Common Parser Errors

**Problem**: "unexpected token" errors
```go
// Solution: Check token precedence and lookahead
if p.peekToken.Type == lexer.SEMICOLON {
    p.nextToken() // Consume semicolon
}
```

**Problem**: Missing position information
```go
// Solution: Always set position when creating nodes
node.Position = p.curToken.Position
```

### 2. Type Checking Issues

**Problem**: Type inference not working
```go
// Solution: Check literal type mapping
func (a *Analyzer) inferType(expr ast.Expression) *ast.Type {
    switch e := expr.(type) {
    case *ast.Literal:
        switch e.Value.(type) {
        case int, int64:
            return ast.NewBaseType("int")
        case float64:
            return ast.NewBaseType("float")
        case string:
            return ast.NewBaseType("string")
        case bool:
            return ast.NewBaseType("bool")
        }
    }
    return ast.NewBaseType("unknown")
}
```

**Problem**: Function type compatibility failing
```go
// Solution: Ensure function signatures are properly compared
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
```

### 3. Symbol Table Issues

**Problem**: Variables not found in scope
```go
// Solution: Check scope management
a.symbols.EnterScope()
// ... add symbols ...
defer a.symbols.ExitScope() // Always exit scope
```

**Problem**: Duplicate declaration errors
```go
// Solution: Check if symbol already exists
if existing, _ := a.symbols.Resolve(name); existing != nil {
    a.errors.AddError(
        pos,
        errors.ErrCodeDuplicateDecl,
        fmt.Sprintf("'%s' is already declared in this scope", name),
    )
}
```

## Testing Patterns

### 1. Parser Test Structure

```go
func TestVariableDeclaration(t *testing.T) {
    tests := []struct {
        input              string
        expectedIdentifier string
        expectedType       string
        expectedValue      interface{}
    }{
        {"x := 5;", "x", "", 5},
        {"x : int = 5;", "x", "int", 5},
        {"mut y := 10.5;", "y", "", 10.5},
    }
    
    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := NewParser(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)
        
        if len(program.Declarations) != 1 {
            t.Fatalf("expected 1 declaration, got %d", len(program.Declarations))
        }
        
        stmt := program.Declarations[0]
        if !testVariableDeclaration(t, stmt, tt.expectedIdentifier, tt.expectedType, tt.expectedValue) {
            return
        }
    }
}
```

### 2. Type Checker Test Structure

```go
func TestTypeCompatibility(t *testing.T) {
    tests := []struct {
        expected *ast.Type
        actual   *ast.Type
        shouldBe bool
    }{
        {ast.NewBaseType("int"), ast.NewBaseType("int"), true},
        {ast.NewBaseType("int"), ast.NewBaseType("string"), false},
    }
    
    analyzer := New("", "")
    
    for _, tt := range tests {
        result := analyzer.typesCompatible(tt.expected, tt.actual)
        if result != tt.shouldBe {
            t.Errorf("typesCompatible(%s, %s) = %t, want %t",
                tt.expected.String(), tt.actual.String(), result, tt.shouldBe)
        }
    }
}
```

## Performance Tips

### 1. AST Construction
- Reuse position structs when possible
- Use string interning for identifiers
- Minimize allocations in hot paths

### 2. Type Checking
- Cache type inference results
- Use efficient type comparison
- Avoid repeated symbol lookups

### 3. Error Reporting
- Collect errors in batches
- Use efficient string formatting
- Minimize position struct copying

## Debugging

### 1. Enable Debug Logging
```go
// Add debug prints to track execution
fmt.Printf("Parsing declaration at line %d\n", p.curToken.Position.Line)
```

### 2. AST Visualization
```go
// Print AST structure for debugging
func printAST(node ast.Node, indent int) {
    prefix := strings.Repeat("  ", indent)
    fmt.Printf("%s%T: %s\n", prefix, node, node.String())
    
    // Recursively print children
    switch n := node.(type) {
    case *ast.Program:
        for _, decl := range n.Declarations {
            printAST(decl, indent+1)
        }
    case *ast.BlockStatement:
        for _, stmt := range n.Statements {
            printAST(stmt, indent+1)
        }
    }
}
```

### 3. Token Stream Debugging
```go
// Print token stream for debugging
func debugTokens(input string) {
    l := lexer.New(input)
    for {
        tok := l.NextToken()
        fmt.Printf("Token: %s '%s' at %d:%d\n",
            tok.Type.String(), tok.Literal, tok.Position.Line, tok.Position.Column)
        if tok.Type == lexer.EOF {
            break
        }
    }
}
```

---

*This quick reference provides common patterns and solutions for Mars compiler development. For detailed architecture information, see [technical-architecture.md](./technical-architecture.md).* 