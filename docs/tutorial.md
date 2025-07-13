# Mars Language Tutorial

Welcome to Mars! This tutorial will guide you through the basics of the Mars programming language.

## Getting Started

### Your First Program

Create a file named `hello.mars`:

```mars
func main() {
    log("Hello, Mars!")
}
```

## Language Basics

### Variables

Variables in Mars are immutable by default. Use `mut` to declare mutable variables:

```mars
// Immutable variable with explicit type
x : int = 42

// Immutable variable with type inference
name := "Mars"

// Mutable variable
mut y : int = 10
y = 20  // OK - y is mutable

// x = 30  // Error: cannot assign to immutable variable
```

### Type Inference vs Explicit Types

```mars
// Type inference (recommended when obvious)
age := 25
pi := 3.14159
message := "Hello, World!"
isActive := true

// Explicit types (required in some contexts)
count : int = 0
temperature : float = 98.6
username : string = "admin"
enabled : bool = false
```

### Functions

Functions are declared with the `func` keyword:

```mars
func add(a: int, b: int) -> int {
    return a + b
}

func greet(name: string) {
    log("Hello, " + name + "!")
}

// Function with no parameters
func getCurrentTime() -> string {
    return "2024-01-01"
}
```

### Control Flow

#### If Statements

```mars
func max(a: int, b: int) -> int {
    if a > b {
        return a
    } else {
        return b
    }
}

// If-else chains
func getGrade(score: int) -> string {
    if score >= 90 {
        return "A"
    } else if score >= 80 {
        return "B"
    } else if score >= 70 {
        return "C"
    } else {
        return "F"
    }
}
```

#### For Loops

```mars
func sum(n: int) -> int {
    mut total := 0
    for i := 0; i < n; i = i + 1 {
        total = total + i
    }
    return total
}

// Infinite loop (use with break)
func waitForInput() {
    for {
        // Loop body
        // Use break to exit
    }
}
```

### Structs

Define custom types using structs:

```mars
struct Point {
    x: int;
    y: int;
}

// Struct with mixed types
struct Person {
    name: string;
    age: int;
    height: float;
}

func createPoint(x: int, y: int) -> Point {
    return Point{x: x, y: y}
}

func createPerson(name: string, age: int) -> Person {
    return Person{
        name: name,
        age: age,
        height: 5.8
    }
}
```

### Arrays and Slices

Mars supports both fixed-size arrays and dynamic slices:

```mars
// Fixed-size array
numbers : [5]int = [1, 2, 3, 4, 5]

// Dynamic slice
dynamicNumbers := [1, 2, 3, 4, 5]

// Array operations
firstElement := numbers[0]
slice := numbers[1:4]  // Slicing: elements 1, 2, 3

// Empty arrays
emptyFixed : [10]int = []
emptyDynamic := []
```

### Member Access and Function Calls

```mars
struct Calculator {
    value: int;
}

calc := Calculator{value: 10}
result := calc.value

// Function calls
sum := add(5, 3)
greeting := greet("Alice")

// Chained operations
point := Point{x: 1, y: 2}
xCoord := point.x
```

### Unsafe Blocks

For low-level operations, use unsafe blocks:

```mars
unsafe {
    // Pointer operations
    ptr : *int = alloc(int)
    *ptr = 42
    value := *ptr
    free(ptr)
}
```

## Best Practices

### 1. Immutability First

```mars
// Prefer immutable variables
username := "admin"
maxRetries := 3

// Use mut only when necessary
mut counter := 0
mut attempts := 0
```

### 2. Type Safety

```mars
// Good: Clear intent
func processUser(id: int, name: string) -> bool {
    // Process user
    return true
}

// Good: Explicit when needed
userCount : int = 0
```

### 3. Error Handling

```mars
func divide(a: int, b: int) -> int {
    if b == 0 {
        log("Error: Division by zero")
        return 0
    }
    return a / b
}
```

### 4. Function Design

```mars
// Good: Single responsibility
func calculateArea(width: int, height: int) -> int {
    return width * height
}

func validateInput(input: string) -> bool {
    return input != ""
}
```

## Common Patterns

### 1. Data Structures

```mars
struct Node {
    value: int;
    next: *Node;  // Pointer to next node (in unsafe blocks)
}

struct Config {
    host: string;
    port: int;
    timeout: float;
}
```

### 2. Option-like Pattern

```mars
struct Result {
    value: int;
    isValid: bool;
    error: string;
}

func safeDiv(a: int, b: int) -> Result {
    if b == 0 {
        return Result{
            value: 0,
            isValid: false,
            error: "Division by zero"
        }
    }
    return Result{
        value: a / b,
        isValid: true,
        error: ""
    }
}
```

## Language Features Summary

### Supported Syntax:
- ✅ `x := value` (type inference)
- ✅ `x : type = value` (explicit type)
- ✅ `mut x := value` (mutable with inference)
- ✅ `mut x : type = value` (mutable with explicit type)
- ✅ `func name(params) -> type { ... }`
- ✅ `struct Name { field: type; }`
- ✅ `if condition { ... } else { ... }`
- ✅ `for init; condition; post { ... }`
- ✅ `array[index]` and `array[start:end]`
- ✅ `obj.field` and `func(args)`
- ✅ `unsafe { ... }`

### Not Supported:
- ❌ `:=` for explicit type declarations
- ❌ Method declarations on structs
- ❌ Generics
- ❌ Modules/imports