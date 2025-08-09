# Contributing to Mars Programming Language

Thank you for your interest in contributing to Mars! This document provides guidelines and information for contributors.

## **What We're Building**

Mars is a modern, statically-typed programming language designed for:
- **Algorithmic Problem Solving**: Successfully solves LeetCode Hard problems
- **Clean Syntax**: Modern, expressive language design
- **Go Integration**: Seamless interoperability with Go ecosystem
- **Type Safety**: Strong static typing with inference

## **Development Setup**

### Prerequisites
- Go 1.21 or later
- Git

### Getting Started
```bash
# Clone the repository
git clone https://github.com/yourusername/mars.git
cd mars

# Run tests
go test ./...

# Test Mars programs
go run cmd/mars/*.go run examples/while_loop_test.mars
```

## **Current Status**

### **Completed Features**
- **Core Language**: Variables, functions, control flow
- **Data Structures**: Arrays, nested arrays, indexing
- **Control Flow**: If/else, for loops, while loops, break/continue
- **Type System**: Static typing with inference
- **Algorithmic Support**: 5+ LeetCode problems solved

### **Areas for Contribution**

#### **High Priority**
1. **String Operations**: Enhanced string manipulation
2. **Hash Maps**: Dictionary data structures
3. **Dynamic Arrays**: Append, remove operations
4. **Standard Library**: More built-in functions

#### **Medium Priority**
1. **Package System**: Module imports and dependencies
2. **Concurrency**: Goroutines and channels
3. **IDE Support**: Language server and tooling
4. **Performance**: Code optimization

#### **Low Priority**
1. **Documentation**: More examples and tutorials
2. **Testing**: Additional test cases
3. **Examples**: More algorithmic problems

## **How to Contribute**

### **1. Choose an Issue**
- Check existing issues or create a new one
- Comment on issues you'd like to work on
- We'll assign it to you

### **2. Development Workflow**
```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# Add tests for new functionality
# Update documentation

# Run tests
go test ./...

# Commit your changes
git commit -m "feat: add your feature description"

# Push and create a pull request
git push origin feature/your-feature-name
```

### **3. Code Standards**
- **Go Style**: Follow Go formatting and conventions
- **Comments**: Document public APIs and complex logic
- **Tests**: Add tests for new functionality
- **Documentation**: Update relevant docs

### **4. Pull Request Guidelines**
- **Title**: Clear, descriptive title
- **Description**: Explain what and why (not how)
- **Tests**: Include test cases
- **Documentation**: Update docs if needed

## **Testing**

### **Run All Tests**
```bash
go test ./...
```

### **Test Mars Programs**
```bash
# Test basic functionality
go run cmd/mars/*.go run examples/while_loop_test.mars

# Test algorithmic problems
go run cmd/mars/*.go run examples/two_sum_working_final.mars
go run cmd/mars/*.go run examples/binary_search_with_while.mars
```

### **Add New Tests**
- Add test files in appropriate directories
- Test both success and error cases
- Include edge cases and boundary conditions

## **Documentation**

### **Code Documentation**
- Document all public APIs
- Include usage examples
- Explain complex algorithms

### **User Documentation**
- Update README.md for new features
- Add examples to docs/
- Create tutorials for new capabilities

## **Bug Reports**

### **Before Reporting**
1. Check existing issues
2. Try to reproduce the bug
3. Test with latest version

### **Bug Report Template**
```
**Description**
Brief description of the bug

**Steps to Reproduce**
1. Step 1
2. Step 2
3. Step 3

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., macOS, Linux, Windows]
- Go version: [e.g., 1.21.0]
- Mars version: [if applicable]

**Additional Context**
Any other relevant information
```

## **Feature Requests**

### **Before Requesting**
1. Check if feature already exists
2. Consider if it fits Mars's goals
3. Think about implementation complexity

### **Feature Request Template**
```
**Problem**
Description of the problem this feature would solve

**Proposed Solution**
Description of the proposed feature

**Alternatives Considered**
Other approaches you've considered

**Additional Context**
Any other relevant information
```

## **Community Guidelines**

### **Be Respectful**
- Treat all contributors with respect
- Welcome newcomers
- Provide constructive feedback

### **Be Helpful**
- Answer questions when you can
- Share knowledge and experience
- Help review pull requests

### **Be Patient**
- Development takes time
- Complex features require careful design
- Quality is more important than speed

## **Getting Help**

### **Discussions**
- Use GitHub Discussions for questions
- Share ideas and proposals
- Get help with implementation

### **Issues**
- Use issues for bugs and feature requests
- Provide detailed information
- Be patient for responses

## **Recognition**

### **Contributors**
- All contributors will be listed in CONTRIBUTORS.md
- Significant contributions will be highlighted
- We appreciate all levels of contribution

### **Hall of Fame**
- Major feature implementations
- Critical bug fixes
- Outstanding documentation

## **License**

By contributing to Mars, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Mars! 