# Mars Compiler Architecture

This document outlines the architecture and design decisions of the Mars programming language compiler.

## Compiler Pipeline

The Mars compiler follows a traditional multi-stage pipeline:

1. **Lexical Analysis** (`lexer/`)
   - Converts source code into tokens
   - Handles whitespace, comments, and basic syntax validation
   - Provides detailed error reporting with line/column information

2. **Parsing** (`parser/`)
   - Recursive descent parser implementation
   - Builds Abstract Syntax Tree (AST)
   - Validates syntax against formal grammar
   - Handles operator precedence and associativity

3. **Static Analysis** (`analyzer/`)
   - Type checking and inference
   - Scope resolution
   - Immutability checking
   - Safety analysis for unsafe blocks

4. **Code Generation** (`transpiler/`)
   - AST to Go code transformation
   - Handles memory management
   - Generates idiomatic Go code
   - Applies Go formatting

5. **Runtime** (`runtime/`)
   - Garbage collection
   - Memory management
   - Unsafe operations support

## Package Structure

```
mars/
├── lexer/          # Token definitions and lexical analysis
├── parser/         # Syntax analysis and AST construction
├── analyzer/       # Static analysis and type checking
├── transpiler/     # Go code generation
├── runtime/        # Runtime support and GC
├── cmd/
│   ├── zcc/       # Compiler CLI
│   └── repl/      # Interactive REPL
└── docs/          # Documentation
```

## Key Design Decisions

### 1. Go as Implementation Language
- Leverages Go's strong type system
- Benefits from Go's garbage collector
- Easy integration with Go ecosystem
- Good performance characteristics

### 2. Recursive Descent Parser
- Handwritten parser for better error messages
- Direct mapping to grammar rules
- Easy to maintain and extend
- Good performance for our use case

### 3. Memory Management
- Default garbage collection
- Optional manual memory management in unsafe blocks
- Zero-cost abstractions where possible
- Safe by default, unsafe by opt-in

### 4. Type System
- Static typing with type inference
- Support for basic types and user-defined types
- Pointer and array types
- Immutable by default, mutable by opt-in

### 5. Error Handling
- Detailed error messages with source locations
- Multiple error collection
- Clear error recovery
- User-friendly error formatting

## Implementation Guidelines

### Code Style
- Follow Go's standard formatting
- Use `go fmt` on all code
- Document all exported symbols
- Write tests for all new features

### Testing Strategy
- Unit tests for each package
- Integration tests for compiler pipeline
- Property-based tests for parser
- Fuzzing for error handling

### Error Handling
- Use custom error types
- Include source location in errors
- Provide helpful error messages
- Support error recovery where possible

### Performance Considerations
- Minimize allocations
- Use efficient data structures
- Profile critical paths
- Optimize for common cases

## Future Considerations

### Planned Features
1. Concurrency support
2. Package system
3. Standard library
4. Build system integration

### Potential Optimizations
1. Parallel compilation
2. Incremental compilation
3. Better error recovery
4. More aggressive optimizations

## Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Write tests
4. Implement feature
5. Submit pull request

### Code Review Process
1. Automated tests must pass
2. Code must be formatted
3. Documentation must be updated
4. At least one reviewer must approve

### Release Process
1. Update version numbers
2. Update changelog
3. Create release tag
4. Build and test release
5. Publish release 