# Contributing to Mars

Thank you for your interest in contributing to the Mars programming language! This document provides guidelines and information for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Contribution Workflow](#contribution-workflow)
- [Current Development Priorities](#current-development-priorities)
- [Reporting Issues](#reporting-issues)
- [Code of Conduct](#code-of-conduct)

## Getting Started

### Prerequisites

- **Go 1.21 or later** - [Download here](https://golang.org/dl/)
- **Git** - [Download here](https://git-scm.com/)
- **Basic understanding of compiler design** - Familiarity with lexers, parsers, and ASTs is helpful

### Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/Anthony4m/mars.git
   cd mars
   ```
3. **Set up the upstream remote**:
   ```bash
   git remote add upstream https://github.com/Anthony4m/mars.git
   ```
4. **Build the project**:
   ```bash
   go build ./...
   ```
5. **Run tests**:
   ```bash
   go test ./...
   ```

## Development Setup

### Environment Setup

1. **Install Go** and ensure it's in your PATH
2. **Set up your IDE** - Recommended: VS Code with Go extension
3. **Install development tools**:
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   go install golang.org/x/lint/golint@latest
   ```

### Building and Testing

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./evaluator

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Run linter
golint ./...
```

### Testing Your Changes

```bash
# Test the lexer
go test ./lexer -v

# Test the parser
go test ./parser -v

# Test the analyzer
go test ./analyzer -v

# Test the evaluator
go test ./evaluator -v

# Test error handling
go test ./errors -v
```

## Project Structure

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

### Key Components

- **`lexer/`**: Converts source code into tokens
- **`parser/`**: Builds Abstract Syntax Tree from tokens
- **`analyzer/`**: Performs type checking and semantic analysis
- **`evaluator/`**: Executes AST nodes at runtime
- **`ast/`**: Defines AST node types and interfaces
- **`errors/`**: Provides error reporting and handling

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `go fmt` to format your code
- Use `goimports` to organize imports
- Write clear, descriptive comments for exported functions and types

### Code Organization

- **Package structure**: Keep related functionality together
- **File naming**: Use descriptive names that reflect the content
- **Function length**: Keep functions focused and reasonably sized
- **Error handling**: Always check and handle errors appropriately

### Naming Conventions

```go
// Use PascalCase for exported names
type FunctionValue struct {
    Name string
}

// Use camelCase for unexported names
func (e *Evaluator) evalFunctionCall(n *ast.FunctionCall) Value {
    // implementation
}

// Use UPPER_CASE for constants
const (
    INTEGER_TYPE = "INTEGER"
    STRING_TYPE  = "STRING"
)
```

### Error Handling

```go
// Always check errors
result, err := someOperation()
if err != nil {
    return e.newError(pos, ErrRuntimeError, "operation failed: %v", err)
}

// Use structured error reporting
return e.newError(position, ErrTypeMismatch, 
    "cannot assign %s to %s", sourceType, targetType)
```

## Testing

### Writing Tests

- **Test coverage**: Aim for >95% coverage for new code
- **Test organization**: Group related tests together
- **Test names**: Use descriptive names that explain what's being tested
- **Test data**: Use realistic test cases

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected interface{}
        shouldError bool
    }{
        {
            name:     "basic case",
            input:    "x := 42;",
            expected: 42,
            shouldError: false,
        },
        {
            name:     "error case",
            input:    "x := ;",
            expected: nil,
            shouldError: true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./evaluator

# Run tests with verbose output
go test -v ./evaluator

# Run specific test
go test -run TestFunctionName ./evaluator
```

## Contribution Workflow

### 1. Choose an Issue

- Check the [Issues](https://github.com/Anthony4m/mars/issues) page
- Look for issues labeled `good first issue` for beginners
- Comment on the issue to let others know you're working on it

### 2. Create a Feature Branch

```bash
# Update your fork
git fetch upstream
git checkout main
git merge upstream/main

# Create a new branch
git checkout -b feature/your-feature-name
```

### 3. Make Your Changes

- Write your code following the coding standards
- Add tests for new functionality
- Update documentation if needed
- Ensure all tests pass

### 4. Commit Your Changes

```bash
# Stage your changes
git add .

# Commit with a descriptive message
git commit -m "feat: add array literal evaluation

- Implement ArrayLiteral evaluation in evaluator
- Add tests for array literal operations
- Update documentation for array support"
```

### 5. Push and Create a Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:
- Clear description of changes
- Reference to related issues
- Screenshots if UI changes
- Test results

### 6. Code Review

- Address review comments
- Make requested changes
- Ensure CI checks pass
- Wait for maintainer approval

## Current Development Priorities

### High Priority 🔥

1. **Array/Struct Runtime Support**
   - Implement `ArrayLiteral` evaluation
   - Implement `StructLiteral` evaluation
   - Add array indexing support
   - Add member access support

2. **Built-in Functions**
   - Add `len()` function for arrays and strings
   - Add `append()` function for slices
   - Add basic math functions (`abs()`, `min()`, `max()`)
   - Add string functions (`substring()`, `contains()`)

3. **CLI Tools**
   - Create full `zcc` compiler interface
   - Add `zcc run` command
   - Add `zcc build` command
   - Add interactive REPL

### Medium Priority 🟡

4. **Code Generation**
   - Create transpiler to Go
   - Generate idiomatic Go code
   - Handle memory management
   - Integrate with Go toolchain

5. **Standard Library**
   - Expand built-in function library
   - Add utility functions
   - Add data structure operations
   - Add I/O functions

### Low Priority 🟢

6. **Advanced Features**
   - Package system and imports
   - Concurrency support
   - Generics
   - Macros

## Good First Issues

If you're new to the project, here are some good starting points:

- **Documentation**: Improve existing docs or add examples
- **Tests**: Add test cases for edge cases
- **Error Messages**: Improve error reporting
- **Code Cleanup**: Refactor existing code
- **Performance**: Optimize existing implementations

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Clear description** of the problem
2. **Steps to reproduce** the issue
3. **Expected behavior** vs actual behavior
4. **Environment details** (OS, Go version, etc.)
5. **Code example** that demonstrates the issue
6. **Error messages** if any

### Feature Requests

When requesting features, please include:

1. **Clear description** of the feature
2. **Use case** and motivation
3. **Proposed syntax** or interface
4. **Examples** of how it would be used
5. **Consideration** of implementation complexity

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Use welcoming and inclusive language
- Be collaborative and constructive
- Focus on what is best for the community
- Show empathy towards other community members

### Enforcement

- Unacceptable behavior will not be tolerated
- Violations may result in temporary or permanent ban
- Maintainers have the right to remove, edit, or reject contributions

## Getting Help

### Resources

- **Documentation**: Check the [docs/](docs/) directory
- **Issues**: Search existing [issues](https://github.com/Anthony4m/mars/issues)
- **Discussions**: Use GitHub Discussions for questions
- **Code**: Review existing code for examples

### Communication

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Pull Requests**: For code contributions
- **Email**: For private or sensitive matters

## Recognition

Contributors will be recognized in:

- **README.md**: List of contributors
- **Release notes**: Credit for significant contributions
- **Documentation**: Attribution for major features
- **Community**: Recognition in discussions and events

## License

By contributing to Mars, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Mars! Your contributions help make this language better for everyone. 