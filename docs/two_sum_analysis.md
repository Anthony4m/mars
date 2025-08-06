# Two Sum Implementation Analysis - Mars Language Gaps

## 🎯 **Objective**
Implement LeetCode 1: Two Sum in Mars to identify critical language features needed for 1.0 release.

## ❌ **Critical Issues Identified**

### **Issue 1: Variable Declarations Inside Functions**
```
Error: unexpected token COLONEQ at top level
Location: Line 20, column 12
```

**Problem**: Mars parser doesn't support variable declarations inside function bodies.
```mars
func example() {
    i := 0;  // ❌ Fails - parser expects top-level declaration
}
```

**Solution Needed**: 
- Extend `parseStatement()` to handle variable declarations
- Support both `:=` and `let` syntax inside functions
- Implement proper variable scoping

### **Issue 2: For Loop Implementation**
```
Error: expected ';' after for loop init
Location: Line 6, column 12
```

**Problem**: For loop parsing is incomplete - expects C-style syntax but doesn't handle variable declarations properly.
```mars
for i := 0; i < len(nums); i = i + 1 {  // ❌ Fails
    // ...
}
```

**Solution Needed**:
- Complete for loop parsing in `parseForStatement()`
- Support variable declarations in init clause
- Handle assignment expressions in post clause

### **Issue 3: Array Indexing**
```
Error: unexpected token RBRACKET in expression
Location: Line 9, column 17
```

**Problem**: Array indexing syntax `array[index]` is not supported.
```mars
nums[i]  // ❌ Fails - parser doesn't recognize bracket notation
```

**Solution Needed**:
- Extend `parsePrimary()` to handle array indexing
- Support `IndexExpression` in expression parsing
- Implement array bounds checking

### **Issue 4: Array Return Types**
```
Error: expected current token to be RBRACE, got RETURN
Location: Line 9, column 17
```

**Problem**: Function return types don't support arrays.
```mars
func two_sum(nums : []int, target : int) -> []int {  // ❌ Return type fails
    return [i, j];  // ❌ Array literal fails
}
```

**Solution Needed**:
- Extend type parsing to support `[]type` return types
- Support array literals in return statements
- Implement array type checking

### **Issue 5: Hash Map Support**
```
Error: unexpected token RBRACE at top level
Location: Line 18, column 14
```

**Problem**: No hash map/dictionary support in the language.
```mars
seen := {};  // ❌ Hash map literal not supported
seen[key] = value;  // ❌ Hash map operations not supported
```

**Solution Needed**:
- Add hash map literal parsing
- Implement hash map operations (get/set)
- Support hash map types in type system

### **Issue 6: If Statements Inside Functions**
```
Error: unexpected token IF in expression
Location: Line 8, column 13
```

**Problem**: If statements don't work inside function bodies.
```mars
func example() {
    if condition {  // ❌ Fails - parser doesn't expect if here
        // ...
    }
}
```

**Solution Needed**:
- Ensure `parseStatement()` includes if statement handling
- Support nested if statements
- Implement proper block scoping

## ✅ **What Currently Works**

1. **Function Declarations**: ✅ `func name(param : type) -> returnType`
2. **Basic Arithmetic**: ✅ `+`, `-`, `*`, `/`
3. **Function Calls**: ✅ `functionName(arg1, arg2)`
4. **Return Statements**: ✅ `return expression;` (simple types)
5. **Print Statements**: ✅ `log("message");`
6. **Boolean Operations**: ✅ `==`, `!=`, `>`, `<`
7. **Array Literals**: ✅ `[1, 2, 3]` (top level only)

## 🚀 **Implementation Priority**

### **Phase 1: Core Features (Week 1-2)**
1. **Variable Declarations Inside Functions** - HIGH PRIORITY
   - Extend statement parsing
   - Implement variable scoping
   - Support `:=` syntax

2. **If Statements Inside Functions** - HIGH PRIORITY
   - Fix statement parsing order
   - Support nested if statements
   - Implement proper block scoping

### **Phase 2: Control Flow (Week 3-4)**
3. **For Loop Completion** - HIGH PRIORITY
   - Complete for loop parsing
   - Support variable declarations in init
   - Handle assignment expressions

### **Phase 3: Data Structures (Week 5-6)**
4. **Array Indexing** - HIGH PRIORITY
   - Implement `IndexExpression` parsing
   - Support `array[index]` syntax
   - Add bounds checking

5. **Array Return Types** - MEDIUM PRIORITY
   - Extend type system for arrays
   - Support array literals in returns
   - Implement type checking

### **Phase 4: Advanced Features (Week 7-8)**
6. **Hash Map Support** - MEDIUM PRIORITY
   - Add hash map literals
   - Implement get/set operations
   - Support hash map types

## 🧪 **Testing Strategy**

### **Two Sum Test Cases**
```mars
// Test Case 1: [2,7,11,15], target = 9 → [0,1]
// Test Case 2: [3,2,4], target = 6 → [1,2]  
// Test Case 3: [3,3], target = 6 → [0,1]
```

### **Implementation Milestones**
1. **Milestone 1**: Basic variable declarations and if statements
2. **Milestone 2**: For loops and array indexing
3. **Milestone 3**: Array return types and literals
4. **Milestone 4**: Hash map support
5. **Milestone 5**: Complete Two Sum implementation

## 📊 **Success Metrics**

Mars 1.0 will be complete when:

- [ ] **Two Sum brute force** implementation works
- [ ] **Two Sum hash map** implementation works
- [ ] **All test cases** pass
- [ ] **Performance** is acceptable (O(n²) and O(n) solutions)
- [ ] **Error handling** is robust
- [ ] **Type safety** is maintained

## 🎯 **Next Steps**

1. **Start with variable declarations** - this is the foundation
2. **Fix if statements** - needed for control flow
3. **Complete for loops** - needed for iteration
4. **Add array indexing** - needed for data access
5. **Implement array returns** - needed for function outputs
6. **Add hash maps** - needed for optimization

This analysis provides a clear roadmap for making Mars 1.0 a fully functional programming language capable of solving real algorithmic problems! 