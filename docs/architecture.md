# Mars Compiler Architecture

This document outlines the architecture and design decisions of the Mars programming language compiler.

## Compiler Pipeline

The Mars compiler follows a traditional multi-stage pipeline:

1. **Lexical Analysis** (`lexer/`) ✅ **COMPLETE**
   - Converts source code into tokens
   - Handles whitespace, comments, and basic syntax validation
   - Provides detailed error reporting with line/column information

2. **Parsing** (`parser/`) ✅ **COMPLETE**
   - Recursive descent parser implementation
   - Builds Abstract Syntax Tree (AST)
   - Validates syntax against formal grammar
   - Handles operator precedence and associativity

3. **Static Analysis** (`analyzer/`) ✅ **COMPLETE**
   - Type checking and inference
   - Scope resolution
   - Immutability checking
   - Safety analysis for unsafe blocks

4. **Runtime Evaluation** (`evaluator/`) ✅ **COMPLETE**
   - AST interpretation and execution
   - Environment management and scoping
   - Built-in function execution
   - Error handling and stack traces

5. **Code Generation** (`transpiler/`) ❌ **NOT IMPLEMENTED**
   - AST to Go code transformation
   - Handles memory management
   - Generates idiomatic Go code
   - Applies Go formatting

6. **Runtime** (`runtime/`) ❌ **NOT IMPLEMENTED**
   - Garbage collection
   - Memory management
   - Unsafe operations support

## Package Structure

```
mars/
├── lexer/          # Token definitions and lexical analysis ✅
├── parser/         # Syntax analysis and AST construction ✅
├── analyzer/       # Static analysis and type checking ✅
├── evaluator/      # Runtime evaluation and execution ✅
├── errors/         # Error handling and reporting ✅
├── ast/            # Abstract Syntax Tree definitions ✅
├── cmd/
│   └── test_errors/ # Simple test runner ✅
└── docs/           # Documentation ✅
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

## Implementation Status

### ✅ **Completed Components**

#### **Lexer** (`lexer/`)
- ✅ Complete token recognition for all language constructs
- ✅ Support for all operators, keywords, and literals
- ✅ Line/column position tracking
- ✅ Error recovery and reporting
- ✅ Comprehensive test coverage

#### **Parser** (`parser/`)
- ✅ Full recursive descent implementation
- ✅ Complete AST construction
- ✅ Operator precedence handling
- ✅ Error recovery and multiple error collection
- ✅ Support for all language constructs
- ✅ Comprehensive test coverage

#### **AST** (`ast/`)
- ✅ Complete node definitions
- ✅ Position tracking for all nodes
- ✅ String representation for debugging
- ✅ Type system integration
- ✅ Comprehensive test coverage

#### **Analyzer** (`analyzer/`)
- ✅ Type checking and inference
- ✅ Scope resolution and symbol tables
- ✅ Immutability checking
- ✅ Function signature validation
- ✅ Error reporting with context
- ✅ Comprehensive test coverage

#### **Evaluator** (`evaluator/`)
- ✅ Runtime execution of AST
- ✅ Environment management and scoping
- ✅ Variable declarations and assignments
- ✅ Function calls and closures
- ✅ Control flow (if/else, for loops)
- ✅ Built-in function support (`log()`)
- ✅ Error handling and stack traces
- ✅ Comprehensive test coverage

#### **Error System** (`errors/`)
- ✅ Structured error reporting
- ✅ Error codes and messages
- ✅ Context and help text
- ✅ Error chaining and severity levels
- ✅ Comprehensive test coverage

### 🔄 **In Progress**

#### **Runtime Features**
- 🔄 Array and struct literal evaluation
- 🔄 Member access and indexing
- 🔄 Unsafe block execution
- 🔄 More built-in functions

#### **CLI Tools**
- 🔄 Full compiler CLI (`zcc`)
- 🔄 Interactive REPL
- 🔄 Build system integration

### 📋 **Planned Components**

#### **Transpiler** (`transpiler/`)
- [ ] AST to Go code transformation
- [ ] Memory management code generation
- [ ] Go formatting and optimization
- [ ] Integration with Go toolchain

#### **Runtime** (`runtime/`)
- [ ] Garbage collection integration
- [ ] Unsafe memory management
- [ ] Performance profiling
- [ ] Debugging support

#### **Standard Library**
- [ ] Built-in function library
- [ ] Array and slice operations
- [ ] String manipulation functions
- [ ] Math and utility functions

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

## Current Limitations

### **Runtime Limitations**
- No code generation (interpreted execution only)
- Limited built-in functions (only `log()`)
- No array/struct runtime evaluation
- No unsafe block execution

### **Tooling Limitations**
- No full CLI compiler
- No REPL interface
- No build system integration
- No IDE support

### **Language Limitations**
- No package system
- No imports/exports
- No concurrency support
- No generics

## Future Considerations

### Planned Features
1. **Code Generation**: AST to Go transpiler
2. **Runtime Support**: Array/struct evaluation
3. **Standard Library**: Built-in functions
4. **Package System**: Module imports
5. **CLI Tools**: Full compiler and REPL
6. **Build System**: Dependency management

### Potential Optimizations
1. **Parallel Compilation**: Multi-threaded parsing
2. **Incremental Compilation**: Change detection
3. **Better Error Recovery**: More robust parsing
4. **Code Optimization**: AST transformations

## Contributing

We welcome contributions! Please see our [Contributing Guide](../CONTRIBUTING.md) for detailed information on:

- Development setup and environment
- Coding standards and guidelines
- Testing requirements
- Contribution workflow
- Current development priorities
- How to report issues and request features

### Quick Overview

**Development Workflow:**
1. Fork the repository
2. Create your feature branch
3. Write tests
4. Implement feature
5. Submit pull request

**Code Review Process:**
1. Automated tests must pass
2. Code must be formatted
3. Documentation must be updated
4. At least one reviewer must approve

**Release Process:**
1. Update version numbers
2. Update changelog
3. Create release tag
4. Build and test release
5. Publish release

## Current Development Focus

The current development priorities are:

1. **Array/Struct Runtime**: Implement evaluation of array and struct literals
2. **Built-in Functions**: Add standard library functions
3. **CLI Tools**: Build full compiler interface
4. **Code Generation**: Create transpiler to Go
5. **Documentation**: Complete user guides and examples

## Status Summary

Mars has a **solid foundation** with complete lexer, parser, analyzer, and evaluator implementations. The language can execute basic programs with variables, functions, control flow, and output. The next major milestone is implementing array/struct runtime support and building the transpiler to generate Go code. 