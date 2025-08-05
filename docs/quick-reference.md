# Mars Language Quick Reference

A concise reference for the Mars programming language, showing what's implemented and working.

## Language Status: **Core Features Complete** ✅

Mars has a solid foundation with all core language features working. The language can execute basic programs with variables, functions, control flow, and output.

## Core Syntax

### Variables

```mars
// Type inference (recommended)
x := 42;
name := "Mars";
isActive := true;

// Explicit types
count : int = 0;
temperature : float = 98.6;

// Mutable variables
mut counter := 0;
counter = 1;  // OK - mutable
```

### Functions

```mars
// Function declaration
func add(a: int, b: int) -> int {
    return a + b;
}

// Function call
result := add(5, 3);
log(result);
```

### Control Flow

```mars
// If statements
if x > 10 {
    log("x is greater than 10");
} else {
    log("x is 10 or less");
}

// For loops
for i := 0; i < 5; i = i + 1 {
    log(i);
}

// Break and continue
for i := 0; i < 10; i = i + 1 {
    if i == 5 {
        break;
    }
    if i == 2 {
        continue;
    }
    log(i);
}
```

### Built-in Functions

```mars
// Output functions (fully implemented)
log("Hello, World!");
print("No newline");
println("With newline");
printf("Value: %s", "test");

// Array functions
let arr := [1, 2, 3];
let length := len(arr);
let newArr := append(arr, 4);

// Math functions
let sine := sin(0);
let cosine := cos(0);
let root := sqrt(16);

// Time function
let currentTime := now();
```

## Data Types

### Basic Types ✅ **Working**

```mars
// Integer
x : int = 42;
y := 100;

// Float
pi : float = 3.14159;
temp := 98.6;

// String
name : string = "Mars";
message := "Hello, World!";

// Boolean
isActive : bool = true;
enabled := false;

// Null
empty := null;
```

### Complex Types 🔄 **Parsed, Not Evaluated**

```mars
// Array types (parsed, not evaluated)
numbers : [5]int = [1, 2, 3, 4, 5];
dynamic := [1, 2, 3, 4, 5];

// Struct types (parsed, not evaluated)
struct Point {
    x: int;
    y: int;
}

point := Point{x: 1, y: 2};

// Pointer types (parsed, not evaluated)
ptr : *int;
```

## Operators

### Arithmetic ✅ **Working**

```mars
a := 10 + 5;    // 15
b := 10 - 5;    // 5
c := 10 * 5;    // 50
d := 10 / 5;    // 2
e := 10 % 3;    // 1
```

### Comparison ✅ **Working**

```mars
a := 5 == 5;    // true
b := 5 != 3;    // true
c := 5 > 3;     // true
d := 5 >= 5;    // true
e := 3 < 5;     // true
f := 5 <= 5;    // true
```

### Logical ✅ **Working**

```mars
a := true && true;   // true
b := true || false;  // true
c := !false;         // true
```

### Assignment ✅ **Working**

```mars
mut x := 10;
x = 20;  // Assignment to mutable variable
```

## Control Structures

### If Statements ✅ **Working**

```mars
if condition {
    // code
} else if other_condition {
    // code
} else {
    // code
}
```

### For Loops ✅ **Working**

```mars
// C-style for loop
for init; condition; post {
    // code
}

// Example
for i := 0; i < 5; i = i + 1 {
    log(i);
}
```

### Break/Continue ✅ **Working**

```mars
for i := 0; i < 10; i = i + 1 {
    if i == 5 {
        break;  // Exit loop
    }
    if i == 2 {
        continue;  // Skip iteration
    }
    log(i);
}
```

## Functions

### Function Declaration ✅ **Working**

```mars
func functionName(param1: type1, param2: type2) -> returnType {
    // function body
    return value;
}
```

### Function Calls ✅ **Working**

```mars
result := functionName(arg1, arg2);
```

### Closures ✅ **Working**

```mars
func createCounter() -> func() -> int {
    mut count := 0;
    return func() -> int {
        count = count + 1;
        return count;
    };
}
```

## Blocks and Scoping ✅ **Working**

```mars
{
    x := 10;
    {
        y := 20;
        log(x + y);  // 30
    }
    // y is not accessible here
}
```

## Error Handling ✅ **Working**

```mars
// Runtime errors are caught and reported
func divide(a: int, b: int) -> int {
    if b == 0 {
        log("Error: Division by zero");
        return 0;
    }
    return a / b;
}
```

## What's Not Working Yet

### Array Operations 🔄 **Parsed, Not Evaluated**

```mars
// These are parsed but not evaluated at runtime
numbers := [1, 2, 3, 4, 5];
first := numbers[0];
slice := numbers[1:3];
```

### Struct Operations 🔄 **Parsed, Not Evaluated**

```mars
// These are parsed but not evaluated at runtime
struct Person {
    name: string;
    age: int;
}

person := Person{name: "Alice", age: 30};
name := person.name;
```

### Unsafe Blocks 🔄 **Parsed, Not Evaluated**

```mars
// These are parsed but not evaluated at runtime
unsafe {
    ptr := alloc(int);
    *ptr = 42;
    value := *ptr;
    free(ptr);
}
```

### Built-in Functions ✅ **Working**

```mars
// Output functions
log("Hello");           // ✅ Working
print("No newline");    // ✅ Working
println("With newline"); // ✅ Working
printf("Value: %s", "test"); // ✅ Working

// Type conversion functions
let num := toInt("42");     // ✅ Working
let float := toFloat("3.14"); // ✅ Working
let str := toString(42);    // ✅ Working
let type := getType("hello"); // ✅ Working

// Type checking functions
let isInt := isInt(42);     // ✅ Working
let isFloat := isFloat(3.14); // ✅ Working
let isString := isString("hello"); // ✅ Working
let isArray := isArray([1, 2, 3]); // ✅ Working
let isBool := isBool(true); // ✅ Working

// Array functions
let arr := [1, 2, 3];
let length := len(arr);     // ✅ Working
let newArr := append(arr, 4); // ✅ Working
push(arr, 5);              // ✅ Working
let popped := pop(arr);     // ✅ Working
reverse(arr);              // ✅ Working
let joined := join(arr, ", "); // ✅ Working

// Math functions
let sine := sin(0);         // ✅ Working
let cosine := cos(0);       // ✅ Working
let root := sqrt(16);       // ✅ Working
let power := pow(2, 3);     // ✅ Working (2^3 = 8)
let floor := floor(3.7);    // ✅ Working (3.7 → 3)
let ceiling := ceil(3.2);   // ✅ Working (3.2 → 4)
let absolute := abs(-5);    // ✅ Working
let minimum := min(3, 7);   // ✅ Working
let maximum := max(3, 7);   // ✅ Working

// Time function
let currentTime := now();   // ✅ Working

// String and array slicing
let str := "Hello, Mars!";
let slice1 := str[0:5];     // ✅ Working ("Hello")
let slice2 := str[:5];      // ✅ Working ("Hello")
let slice3 := str[7:];      // ✅ Working ("Mars!")
let slice4 := str[-6:-1];   // ✅ Working ("Mars")

let arr := [1, 2, 3, 4, 5];
let arrSlice1 := arr[1:4];  // ✅ Working ([2, 3, 4])
let arrSlice2 := arr[:3];   // ✅ Working ([1, 2, 3])
let arrSlice3 := arr[2:];   // ✅ Working ([3, 4, 5])
let arrSlice4 := arr[-3:-1]; // ✅ Working ([3, 4])

// String indexing
let char := str[0];         // ✅ Working ("H")
let arrElem := arr[2];      // ✅ Working (3)

// Array assignment
arr[0] = 10;               // ✅ Working
```

## Testing Your Code

### Using the Test Runner

```bash
# Test a simple program
go run cmd/test_errors/main.go

# Test with a file
go run cmd/test_errors/main.go your_program.mars
```

### Example Working Program

```mars
func main() {
    mut sum := 0;
    for i := 1; i <= 10; i = i + 1 {
        sum = sum + i;
    }
    log("Sum of 1 to 10 is:");
    log(sum);
}
```

## Error Messages

Mars provides clear error messages with line and column information:

```
error[E0001]: unexpected token EOF in expression
          --> line 1, column 6
```

## Next Steps

1. **Try the examples** in this reference
2. **Check the test suite** in `evaluator/evaluator_test.go` for more examples
3. **Contribute** to implement missing features like array/struct runtime support

### Want to Contribute?

If you're interested in contributing to Mars, check out our [Contributing Guide](../CONTRIBUTING.md) for:

- Development setup instructions
- Coding standards and guidelines
- Current development priorities
- How to submit pull requests
- Good first issues for beginners

## Implementation Status Summary

- ✅ **Core Language**: Variables, functions, control flow, operators
- ✅ **Runtime**: AST evaluation, environment management, error handling
- ✅ **Built-ins**: Comprehensive set of 25+ functions including:
  - Output: `log()`, `print()`, `println()`, `printf()`
  - Type conversion: `toInt()`, `toFloat()`, `toString()`, `getType()`
  - Type checking: `isInt()`, `isFloat()`, `isString()`, `isArray()`, `isBool()`
  - Array operations: `len()`, `append()`, `push()`, `pop()`, `reverse()`, `join()`
  - Math: `sin()`, `cos()`, `sqrt()`, `pow()`, `floor()`, `ceil()`, `abs()`, `min()`, `max()`
  - Time: `now()`
- ✅ **String & Array Operations**: Indexing, slicing, assignment
- ✅ **Comments**: Single-line (`//`) and multi-line (`/* */`)
- 🔄 **Data Structures**: Parsed but not evaluated
- 🔄 **Advanced Features**: Unsafe blocks, member access
- 📋 **Tooling**: CLI compiler, REPL, code generation

**Current Status**: Feature-rich programming language ready for real-world tasks with comprehensive built-in library. 