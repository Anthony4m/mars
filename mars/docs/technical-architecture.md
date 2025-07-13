# Mars Compiler Technical Architecture

## Core Architecture Overview

The Mars compiler follows a traditional multi-stage pipeline architecture:

```
Source Code → Lexer → Parser → AST → Analyzer → Transpiler → Go Code
```

## Component Deep Dive

### 1. Lexer (`lexer/`)

**Purpose**: Convert source code into a stream of tokens with position information.

**Key Design Decisions**:
- **Position Tracking**: Every token includes line and column information
- **Lookahead**: Single token lookahead for operator disambiguation
- **Error Recovery**: Continues lexing after encountering invalid characters

**Core Types**:
```go
type Token struct {
    Type     TokenType
    Literal  string
    Position Position
}

type Lexer struct {
    input        string
    position     int
    readPosition int
    ch           byte
    line         int
    column       int
}
```

**Token Categories**:
- **Keywords**: `mut`, `func`, `struct`, `unsafe`, `if`, `for`, `return`, `log`
- **Literals**: Numbers, strings, booleans, `nil`
- **Operators**: Arithmetic, comparison, logical, assignment
- **Delimiters**: Parentheses, braces, semicolons, commas
- **Identifiers**: Variable names, function names, type names

### 2. Parser (`parser/`)

**Purpose**: Transform token stream into Abstract Syntax Tree (AST).

**Architecture**: Recursive descent parser with operator precedence parsing.

**Key Methods**:
```go
type Parser struct {
    l      *lexer.Lexer
    errors *errors.MarsReporter
    curToken  token.Token
    peekToken token.Token
}

// Main entry points
func (p *Parser) ParseProgram() *ast.Program
func (p *Parser) parseDeclaration() ast.Declaration
func (p *Parser) parseStatement() ast.Statement
func (p *Parser) parseExpression(precedence int) ast.Expression
```

**Expression Parsing Strategy**:
- **Pratt Parser**: Operator precedence parsing for expressions
- **Precedence Levels**: 
  - `LOWEST` (0)
  - `EQUALS` (1) - `==`, `!=`
  - `LESSGREATER` (2) - `<`, `>`, `<=`, `>=`
  - `SUM` (3) - `+`, `-`
  - `PRODUCT` (4) - `*`, `/`, `%`
  - `PREFIX` (5) - `!`, `-`
  - `CALL` (6) - Function calls, array indexing
  - `INDEX` (7) - Array/slice access

**Error Handling**:
- **Synchronization**: Skip tokens until safe recovery point
- **Error Collection**: Continue parsing, collect all errors
- **Context Preservation**: Maintain source position for error reporting

### 3. AST (`ast/`)

**Purpose**: Represent parsed code structure with position information.

**Design Principles**:
- **Position Tracking**: Every node has source position
- **Interface Hierarchy**: Clear separation of concerns
- **Immutability**: AST nodes are immutable after creation

**Core Interfaces**:
```go
type Node interface {
    TokenLiteral() string
    Pos() Position
}

type Declaration interface {
    Node
    declarationNode()
}

type Statement interface {
    Node
    statementNode()
    declarationNode() // Allow statements at top level
    String() string
}

type Expression interface {
    Node
    expressionNode()
    String() string
}
```

**Key Node Types**:

#### Declaration Nodes
- `VarDecl`: Variable declarations with mutability
- `FuncDecl`: Function declarations with signatures
- `StructDecl`: Struct type definitions
- `UnsafeBlock`: Unsafe memory operations

#### Statement Nodes
- `AssignmentStatement`: Variable assignment
- `IfStatement`: Conditional execution
- `ForStatement`: Loop constructs
- `PrintStatement`: Output operations
- `ReturnStatement`: Function returns
- `BlockStatement`: Statement sequences

#### Expression Nodes
- `Literal`: Constants (numbers, strings, booleans)
- `Identifier`: Variable/function names
- `BinaryExpression`: Binary operations
- `UnaryExpression`: Prefix operations
- `FunctionCall`: Function invocation
- `ArrayLiteral`: Array/slice literals
- `StructLiteral`: Struct instantiation

### 4. Analyzer (`analyzer/`)

**Purpose**: Perform semantic analysis, type checking, and error detection.

**Two-Pass Architecture**:

#### Pass 1: Declaration Collection
```go
func (a *Analyzer) collectDeclarations(node ast.Node) error
```
- Build symbol table
- Check for duplicate declarations
- Validate declaration syntax
- Collect function signatures

#### Pass 2: Type Checking
```go
func (a *Analyzer) checkTypes(node ast.Node) error
```
- Verify type compatibility
- Check function bodies
- Validate expressions
- Enforce scoping rules

**Symbol Table Design**:
```go
type Symbol struct {
    Name       string
    Type       ast.Type
    Mutable    bool
    IsFunction bool
    DeclaredAt ast.Node
}

type SymbolTable struct {
    scopes []map[string]*Symbol
    scopeIndex int
}
```

**Type System Features**:
- **Type Inference**: Automatic type deduction from literals
- **Type Compatibility**: Strict type checking with conversions
- **Function Types**: Full function signature support
- **Array Types**: Slice and fixed array support

### 5. Error Reporting (`errors/`)

**Purpose**: Provide structured, user-friendly error messages.

**Error Structure**:
```go
type MarsError struct {
    Position  ast.Position
    Code      string
    Message   string
    Help      string
    Severity  Severity
}

type Severity int

const (
    Info Severity = iota
    Warning
    Error
)
```

**Error Categories**:
- `ErrCodeSyntaxError`: Grammar violations
- `ErrCodeTypeError`: Type mismatches
- `ErrCodeDuplicateDecl`: Redeclaration errors
- `ErrCodeUndefinedVar`: Undefined variable usage
- `ErrCodeImmutableAssign`: Assignment to immutable variable

**Features**:
- **Position Tracking**: Exact line/column location
- **Help Messages**: Suggested fixes and explanations
- **Error Codes**: Categorization for tooling
- **Severity Levels**: Info, Warning, Error
- **Source Context**: Preserve original source for display

## Data Flow

### 1. Parsing Flow
```
Source Code
    ↓
Token Stream (with positions)
    ↓
AST Construction (with position propagation)
    ↓
Program AST
```

### 2. Analysis Flow
```
Program AST
    ↓
Declaration Collection (Pass 1)
    ↓
Symbol Table Population
    ↓
Type Checking (Pass 2)
    ↓
Semantic Validation
    ↓
Error Collection
```

### 3. Error Propagation
```
Error Detection (any stage)
    ↓
Error Creation (with position and context)
    ↓
Error Collection (MarsReporter)
    ↓
Error Reporting (structured output)
```

## Memory Management

### AST Memory Model
- **Immutable Nodes**: AST nodes are created once and never modified
- **Position Sharing**: Position structs are small and can be shared
- **String Interning**: Consider string interning for identifiers

### Symbol Table Memory
- **Scope Stack**: Efficient scope management with slice-based stack
- **Symbol Reuse**: Reuse symbol structs when possible
- **Cleanup**: Automatic cleanup on scope exit

## Performance Considerations

### Parser Performance
- **Token Lookup**: O(1) token type checking
- **Expression Parsing**: O(n) for n tokens with precedence
- **Error Recovery**: Minimal overhead for error cases

### Analysis Performance
- **Symbol Lookup**: O(depth) for scope depth
- **Type Checking**: O(nodes) for AST traversal
- **Error Collection**: O(errors) for error storage

### Memory Usage
- **AST Size**: ~2-3x source code size
- **Symbol Table**: ~O(symbols) for active symbols
- **Error Storage**: ~O(errors) for collected errors

## Testing Strategy

### Unit Testing
- **Lexer Tests**: Token recognition and position tracking
- **Parser Tests**: Grammar compliance and error handling
- **AST Tests**: Node construction and string representation
- **Analyzer Tests**: Type checking and semantic validation

### Integration Testing
- **End-to-End**: Full pipeline validation
- **Error Scenarios**: Comprehensive error case coverage
- **Performance**: Parsing and analysis performance benchmarks

### Example Programs
- **Basic Constructs**: Variables, functions, control flow
- **Complex Types**: Arrays, structs, function types
- **Error Cases**: Invalid syntax, type errors, scope violations

## Extension Points

### Language Features
- **Generics**: Type parameter support
- **Modules**: Import/export system
- **Concurrency**: Goroutines and channels

### Tooling Integration
- **Language Server**: LSP protocol support
- **Debugger**: Source-level debugging
- **Profiler**: Performance analysis tools

### Optimization
- **Constant Folding**: Compile-time evaluation
- **Dead Code Elimination**: Unused code removal
- **Inlining**: Function inlining optimization

## Code Quality Metrics

### Test Coverage
- **Target**: ≥80% line coverage
- **Focus**: Critical paths and error handling
- **Automation**: CI/CD integration

### Code Complexity
- **Cyclomatic Complexity**: ≤10 per function
- **Cognitive Complexity**: ≤15 per function
- **Maintainability Index**: ≥65

### Documentation
- **GoDoc**: All exported symbols documented
- **Examples**: Usage examples for public APIs
- **Architecture**: Design decisions documented

---

*This document provides technical implementation details for the Mars compiler architecture. For high-level design decisions and progress tracking, see [mars-compiler-progress.md](./mars-compiler-progress.md).* 