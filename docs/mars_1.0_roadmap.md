# Mars 1.0 Roadmap - Based on Two Sum Implementation

## 🎯 **Goal: Complete Two Sum Implementation**

The Two Sum problem reveals critical gaps in Mars that need to be addressed for 1.0 release.

## ❌ **Current Issues Identified**

### 1. **Variable Declarations Inside Functions**
**Error**: `unexpected token COLONEQ at top level`
**Problem**: Can't declare variables inside function bodies
**Example**:
```mars
func two_sum(nums : []int, target : int) -> []int {
    i := 0;  // ❌ This fails
    // ...
}
```

### 2. **For Loop Syntax**
**Error**: `expected ';' after for loop init`
**Problem**: For loop parsing expects C-style syntax but implementation is incomplete
**Example**:
```mars
for i := 0; i < len(nums); i = i + 1 {  // ❌ This fails
    // ...
}
```

### 3. **Array Indexing**
**Error**: `unexpected token RBRACKET in expression`
**Problem**: Can't access array elements with `nums[i]`
**Example**:
```mars
nums[i]  // ❌ This fails
```

### 4. **Array Return Types**
**Error**: `expected current token to be RBRACE, got RETURN`
**Problem**: Can't return arrays from functions
**Example**:
```mars
func two_sum(...) -> []int {  // ❌ Return type parsing fails
    return [i, j];  // ❌ Array literal fails
}
```

### 5. **Hash Map Support**
**Error**: `unexpected token RBRACE at top level`
**Problem**: No hash map/dictionary support
**Example**:
```mars
seen := {};  // ❌ Hash map syntax not supported
seen[key] = value;  // ❌ Hash map operations not supported
```

### 6. **If Statements Inside Functions**
**Error**: `unexpected token IF in expression`
**Problem**: If statements don't work inside function bodies
**Example**:
```mars
func example() {
    if condition {  // ❌ This fails
        // ...
    }
}
```

## ✅ **What Currently Works**

1. **Function Declarations**: `func name(param : type) -> returnType`
2. **Basic Arithmetic**: `+`, `-`, `*`, `/`
3. **Function Calls**: `functionName(arg1, arg2)`
4. **Return Statements**: `return expression;` (simple types only)
5. **Print Statements**: `log("message");`
6. **Boolean Operations**: `==`, `!=`, `>`, `<`, etc.
7. **Array Literals**: `[1, 2, 3]` (at top level only)

## 🚀 **Mars 1.0 Implementation Plan**

### **Phase 1: Core Language Features (Priority: HIGH)**

#### 1.1 Variable Declarations Inside Functions
- [ ] Support `:=` syntax inside function bodies
- [ ] Support `let` declarations inside function bodies
- [ ] Implement variable scoping rules

#### 1.2 For Loop Implementation
- [ ] Complete for loop parsing: `for init; condition; post { body }`
- [ ] Support range-based loops: `for i in 0..len(array)`
- [ ] Implement loop variable scoping

#### 1.3 If Statement Support
- [ ] Enable if statements inside function bodies
- [ ] Support else clauses
- [ ] Support nested if statements

### **Phase 2: Data Structures (Priority: HIGH)**

#### 2.1 Array Operations
- [ ] Array indexing: `array[index]`
- [ ] Array assignment: `array[index] = value`
- [ ] Array length: `len(array)`
- [ ] Array slicing: `array[start:end]`

#### 2.2 Array Return Types
- [ ] Support `[]type` return types
- [ ] Array literal returns: `return [a, b, c];`
- [ ] Array type checking

#### 2.3 Hash Map Support
- [ ] Hash map literals: `{}`
- [ ] Hash map operations: `map[key]`, `map[key] = value`
- [ ] Hash map type: `map<keyType, valueType>`

### **Phase 3: Advanced Features (Priority: MEDIUM)**

#### 3.1 Control Flow
- [ ] While loops: `while condition { body }`
- [ ] Break and continue statements
- [ ] Switch/match statements

#### 3.2 Error Handling
- [ ] Try-catch blocks
- [ ] Error types and propagation
- [ ] Panic and recover

#### 3.3 Modules and Imports
- [ ] Module system
- [ ] Import statements
- [ ] Package management

## 🧪 **Testing Strategy**

### **Two Sum Test Cases**
```mars
// Test Case 1: [2,7,11,15], target = 9 → [0,1]
// Test Case 2: [3,2,4], target = 6 → [1,2]  
// Test Case 3: [3,3], target = 6 → [0,1]
```

### **Implementation Milestones**
1. **Milestone 1**: Basic variable declarations and for loops
2. **Milestone 2**: Array indexing and operations
3. **Milestone 3**: Complete Two Sum implementation
4. **Milestone 4**: Hash map optimization
5. **Milestone 5**: Performance testing and optimization

## 📊 **Success Criteria**

Mars 1.0 will be considered complete when:

1. ✅ **Two Sum can be fully implemented** with both brute force and hash map approaches
2. ✅ **All basic data structures** (arrays, maps) are supported
3. ✅ **Control flow** (loops, conditionals) works inside functions
4. ✅ **Variable declarations** work in all contexts
5. ✅ **Type system** is complete and consistent
6. ✅ **Error handling** is robust and user-friendly

## 🎯 **Next Steps**

1. **Fix variable declarations** inside function bodies
2. **Complete for loop implementation**
3. **Enable if statements** inside functions
4. **Implement array indexing**
5. **Add hash map support**
6. **Test with Two Sum implementation**

This roadmap ensures Mars 1.0 will be a fully working programming language capable of solving real algorithmic problems like Two Sum! 