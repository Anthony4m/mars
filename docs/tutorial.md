# Mars Language Tutorial

Welcome to Mars! This tutorial will guide you through the basics of the Mars programming language.

## Getting Started

### Your First Program

Create a file named `hello.mars`:

```mars
func main() {
    log("Hello, Mars!");
}
```

**Note**: Currently, Mars has a runtime evaluator but no code generation. You can test programs using the test runner in `cmd/test_errors/`.

## Language Basics

### Variables

Variables in Mars are immutable by default. Use `mut` to declare mutable variables:

```mars
// Immutable variable with explicit type
x : int = 42;

// Immutable variable with type inference
name := "Mars";

// Mutable variable
mut y : int = 10;
y = 20;  // OK - y is mutable

// x = 30;  // Error: cannot assign to immutable variable
```

### Type Inference vs Explicit Types

```mars
// Type inference (recommended when obvious)
age := 25;
pi := 3.14159;
message := "Hello, World!";
isActive := true;

// Explicit types (required in some contexts)
count : int = 0;
temperature : float = 98.6;
username : string = "admin";
enabled : bool = false;
```

### Functions

Functions are declared with the `func` keyword:

```mars
func add(a: int, b: int) -> int {
    return a + b;
}

func greet(name: string) {
    log("Hello, " + name + "!");
}

// Function with no parameters
func getCurrentTime() -> string {
    return "2024-01-01";
}
```

### Control Flow

#### If Statements

```mars
func max(a: int, b: int) -> int {
    if a > b {
        return a;
    } else {
        return b;
    }
}

// If-else chains
func getGrade(score: int) -> string {
    if score >= 90 {
        return "A";
    } else if score >= 80 {
        return "B";
    } else if score >= 70 {
        return "C";
    } else {
        return "F";
    }
}
```

#### For Loops

```mars
func sum(n: int) -> int {
    mut total := 0;
    for i := 0; i < n; i = i + 1 {
        total = total + i;
    }
    return total;
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
    return Point{x: x, y: y};
}

func createPerson(name: string, age: int) -> Person {
    return Person{
        name: name,
        age: age,
        height: 5.8
    };
}
```

**Note**: Struct literals are parsed but not yet evaluated at runtime.

### Arrays and Slices

Mars supports both fixed-size arrays and dynamic slices:

```mars
// Fixed-size array
numbers : [5]int = [1, 2, 3, 4, 5];

// Dynamic slice
dynamicNumbers := [1, 2, 3, 4, 5];

// Array operations
firstElement := numbers[0];
slice := numbers[1:4];  // Slicing: elements 1, 2, 3

// Empty arrays
emptyFixed : [10]int = [];
emptyDynamic := [];
```

**Note**: Array literals and indexing are parsed but not yet evaluated at runtime.

### Member Access and Function Calls

```mars
struct Calculator {
    value: int;
}

calc := Calculator{value: 10};
result := calc.value;

// Function calls
sum := add(5, 3);
greeting := greet("Alice");

// Chained operations
point := Point{x: 1, y: 2};
xCoord := point.x;
```

**Note**: Member access is parsed but not yet evaluated at runtime.

### Unsafe Blocks

For low-level operations, use unsafe blocks:

```mars
unsafe {
    // Pointer operations
    ptr : *int = alloc(int);
    *ptr = 42;
    value := *ptr;
    free(ptr);
}
```

**Note**: Unsafe blocks are parsed but not yet evaluated at runtime.

## Built-in Functions

### Output Functions

#### `log(value)`, `print(value)`, `println(value)`, `printf(format, ...)`

Output functions in Mars:

```mars
log("Hello, World!");           // String output
print("No newline");            // Print without newline
println("With newline");        // Print with newline
printf("Value: %s", "test");    // Formatted printing

log(42);                        // Integer output
log(3.14);                      // Float output
log(true);                      // Boolean output
log(add(5, 3));                 // Function result output
```

### Array Functions

#### `len(value)`, `append(array, value)`

Array manipulation functions:

```mars
let arr := [1, 2, 3];
let length := len(arr);         // Get array length
let newArr := append(arr, 4);   // Append to array

let str := "hello";
let strLen := len(str);         // Get string length
```

### Math Functions

#### `sin(angle)`, `cos(angle)`, `sqrt(value)`

Mathematical functions (angles in radians):

```mars
let sine := sin(0);             // Sine of 0
let cosine := cos(0);           // Cosine of 0
let root := sqrt(16);           // Square root of 16
```

### Utility Functions

#### `now()`

Get current time:

```mars
let currentTime := now();       // Current time as string
```

**Available**: ✅ **Fully implemented**

## Best Practices

### 1. Immutability First

```mars
// Prefer immutable variables
username := "admin";
maxRetries := 3;

// Use mut only when necessary
mut counter := 0;
mut attempts := 0;
```

### 2. Type Safety

```mars
// Good: Clear intent
func processUser(id: int, name: string) -> bool {
    // Process user
    return true;
}

// Good: Explicit when needed
userCount : int = 0;
```

### 3. Error Handling

```mars
func divide(a: int, b: int) -> int {
    if b == 0 {
        log("Error: Division by zero");
        return 0;
    }
    return a / b;
}
```

### 4. Function Design

```mars
// Good: Single responsibility
func calculateArea(width: int, height: int) -> int {
    return width * height;
}

func validateInput(input: string) -> bool {
    return input != "";
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
        };
    }
    return Result{
        value: a / b,
        isValid: true,
        error: ""
    };
}
```

## Advanced Features

### Built-in Functions

Mars comes with a comprehensive set of built-in functions:

#### Type Conversion and Checking

```mars
// Type conversion
let num := toInt("42");        // String to int
let float := toFloat("3.14");  // String to float
let str := toString(42);       // Any value to string
let type := getType("hello");  // Get type as string

// Type checking
let isInt := isInt(42);        // true
let isFloat := isFloat(3.14);  // true
let isString := isString("hello"); // true
let isArray := isArray([1, 2, 3]); // true
let isBool := isBool(true);    // true
```

#### Array Operations

```mars
let arr := [1, 2, 3];

// Basic operations
let length := len(arr);        // Get array length
let newArr := append(arr, 4);  // Create new array with element

// Advanced operations
push(arr, 5);                  // Add element to end (modifies array)
let popped := pop(arr);        // Remove and return last element
reverse(arr);                  // Reverse array in place
let joined := join(arr, ", "); // Join elements with separator
```

#### Math Functions

```mars
// Basic math
let power := pow(2, 3);        // 2^3 = 8
let floor := floor(3.7);       // 3.7 → 3
let ceiling := ceil(3.2);      // 3.2 → 4
let absolute := abs(-5);       // |-5| = 5
let minimum := min(3, 7);      // 3
let maximum := max(3, 7);      // 7

// Trigonometry
let sine := sin(0);            // 0
let cosine := cos(0);          // 1
let root := sqrt(16);          // 4
```

#### Output Functions

```mars
log("Hello, World!");          // Print with newline
print("No newline");           // Print without newline
println("With newline");       // Print with newline
printf("Value: %s", "test");   // Formatted printing
```

### String and Array Slicing

Mars supports Python-style slicing for both strings and arrays:

```mars
let str := "Hello, Mars!";

// String slicing
let slice1 := str[0:5];        // "Hello"
let slice2 := str[:5];         // "Hello" (from start)
let slice3 := str[7:];         // "Mars!" (to end)
let slice4 := str[-6:-1];      // "Mars" (negative indices)

let arr := [1, 2, 3, 4, 5];

// Array slicing
let arrSlice1 := arr[1:4];     // [2, 3, 4]
let arrSlice2 := arr[:3];      // [1, 2, 3]
let arrSlice3 := arr[2:];      // [3, 4, 5]
let arrSlice4 := arr[-3:-1];   // [3, 4]
```

### String and Array Indexing

```mars
let str := "Hello, Mars!";
let char := str[0];            // "H" (first character)

let arr := [1, 2, 3, 4, 5];
let elem := arr[2];            // 3 (third element)

// Array assignment
arr[0] = 10;                   // Modify array element
```

### Comments

Mars supports both single-line and multi-line comments:

```mars
// This is a single-line comment
x := 42; // Another single-line comment

/* This is a multi-line comment
   that spans multiple lines */

y := 10; /* Inline block comment */

/* Nested /* block */ comments work too */
```

### Practical Examples

#### Working with Arrays

```mars
func processArray(arr: []int) {
    // Add elements
    push(arr, 100);
    push(arr, 200);
    
    // Reverse the array
    reverse(arr);
    
    // Join elements for display
    let display := join(arr, " → ");
    log("Array: " + display);
    
    // Check types
    log("Is array: " + toString(isArray(arr)));
    log("Length: " + toString(len(arr)));
}

// Usage
let numbers := [1, 2, 3];
processArray(numbers);
```

#### Type-Safe Operations

```mars
func safeOperation(value) {
    if isInt(value) {
        log("Processing integer: " + toString(value));
    } else if isString(value) {
        log("Processing string: " + value);
    } else if isArray(value) {
        log("Processing array with " + toString(len(value)) + " elements");
    } else {
        log("Unknown type: " + getType(value));
    }
}

// Test with different types
safeOperation(42);
safeOperation("hello");
safeOperation([1, 2, 3]);
```

## Language Features Summary

### ✅ **Supported and Working:**

**Core Syntax:**
- ✅ `x := value` (type inference)
- ✅ `x : type = value` (explicit type)
- ✅ `mut x := value` (mutable with inference)
- ✅ `mut x : type = value` (mutable with explicit type)
- ✅ `func name(params) -> type { ... }`
- ✅ `struct Name { field: type; }`
- ✅ `if condition { ... } else { ... }`
- ✅ `for init; condition; post { ... }`
- ✅ Comments: `//` and `/* */`
- ✅ `unsafe { ... }` (parsing only)

**Operations:**
- ✅ Arithmetic: `+`, `-`, `*`, `/`, `%`
- ✅ Comparison: `==`, `!=`, `>`, `>=`, `<`, `<=`
- ✅ Logical: `&&`, `||`, `!`
- ✅ Assignment: `=`

**Types:**
- ✅ `int`, `float`, `string`, `bool`, `null`
- ✅ Array types: `[N]Type`, `[]Type`
- ✅ Pointer types: `*Type`
- ✅ Struct types: `struct Name`

**Built-in Functions (25+ functions):**
- ✅ Output: `log()`, `print()`, `println()`, `printf()`
- ✅ Type conversion: `toInt()`, `toFloat()`, `toString()`, `getType()`
- ✅ Type checking: `isInt()`, `isFloat()`, `isString()`, `isArray()`, `isBool()`
- ✅ Array operations: `len()`, `append()`, `push()`, `pop()`, `reverse()`, `join()`
- ✅ Math: `sin()`, `cos()`, `sqrt()`, `pow()`, `floor()`, `ceil()`, `abs()`, `min()`, `max()`
- ✅ Time: `now()`

**String and Array Operations:**
- ✅ Indexing: `str[0]`, `arr[2]`
- ✅ Slicing: `str[0:5]`, `str[:5]`, `str[7:]`, `str[-6:-1]`
- ✅ Assignment: `arr[0] = 10`

### 🔄 **Parsed but Not Yet Evaluated:**

- 🔄 `obj.field` (member access)
- 🔄 Struct literals: `Point{x: 1, y: 2}`
- 🔄 Unsafe block operations

### 📋 **Planned Features:**

- [ ] String functions (`substring`, `indexOf`, `split`, `toLowerCase`, etc.)
- [ ] File I/O (`readFile`, `writeFile`, `exists`)
- [ ] Compound assignments (`+=`, `-=`, `*=`, `/=`)
- [ ] Package system and imports
- [ ] Code generation to Go
- [ ] Concurrency support

## Testing Your Code

Currently, you can test Mars programs using the test runner:

```bash
# Test a simple program
go run cmd/test_errors/main.go

# Test with a file
go run cmd/test_errors/main.go your_program.mars
```

## Next Steps

1. **Try the examples** in this tutorial
2. **Experiment with functions** and control flow
3. **Check the test suite** in `evaluator/evaluator_test.go` for more examples
4. **Contribute** to implement missing features like array/struct runtime support

### Want to Contribute?

If you're interested in contributing to Mars, check out our [Contributing Guide](../CONTRIBUTING.md) for:

- Development setup instructions
- Coding standards and guidelines
- Current development priorities
- How to submit pull requests
- Good first issues for beginners