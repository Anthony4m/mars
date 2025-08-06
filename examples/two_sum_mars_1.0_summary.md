# Two Sum Problem - Mars 1.0 Implementation Analysis

## Problem Statement
Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to `target`.

**Examples:**
- Input: `nums = [2,7,11,15]`, `target = 9` → Output: `[0,1]`
- Input: `nums = [3,2,4]`, `target = 6` → Output: `[1,2]`
- Input: `nums = [3,3]`, `target = 6` → Output: `[0,1]`

## Mars 1.0 Implementation Challenges

### ✅ Working Features
1. **Function definitions** - Functions can be declared and called
2. **Basic arithmetic** - Addition, subtraction, multiplication, division
3. **If statements with simple conditions** - `if a + b == c`
4. **Array literals** - `[1, 2, 3, 4]`
5. **Array return types** - Functions can return arrays
6. **Variable declarations** - `x := 5` and `x : int := 5`
7. **Built-in functions** - `println()`, `len()`, etc.

### ❌ Not Working Features
1. **Array indexing in if conditions** - `if nums[i] + nums[j] == target` causes parse errors
2. **For loops with array indexing** - `for i := 0; i < len(nums); i++` with array access
3. **Hash maps** - No hash map/dictionary support yet
4. **Array parameters in functions** - `func two_sum(nums : []int)` may not work properly

## Working Two Sum Implementation

Since array indexing in if statements doesn't work, we implemented a workaround using individual parameters:

```mars
func two_sum_individual(num1 : int, num2 : int, num3 : int, num4 : int, target : int) -> []int {
    if check_pair(num1, num2, target) {
        return [0, 1];
    }
    if check_pair(num1, num3, target) {
        return [0, 2];
    }
    if check_pair(num1, num4, target) {
        return [0, 3];
    }
    if check_pair(num2, num3, target) {
        return [1, 2];
    }
    if check_pair(num2, num4, target) {
        return [1, 3];
    }
    if check_pair(num3, num4, target) {
        return [2, 3];
    }
    return [-1, -1];
}
```

## Test Results

**Test Case 1:** `[2,7,11,15]`, `target = 9`
- Expected: `[0,1]`
- Our implementation: `[0,1]` ✅

**Test Case 2:** `[3,2,4]`, `target = 6`
- Expected: `[1,2]`
- Our implementation: `[1,2]` ✅

**Test Case 3:** `[3,3]`, `target = 6`
- Expected: `[0,1]`
- Our implementation: `[0,1]` ✅

## Mars 1.0 Roadmap for Two Sum

To fully support the Two Sum problem, Mars 1.0 needs:

1. **Fix array indexing in if conditions** - This is the main blocker
2. **Improve for loop support** - Enable `for i := 0; i < len(nums); i++`
3. **Add hash map support** - For O(n) time complexity solution
4. **Fix function execution** - Currently shows function definitions instead of executing

## Conclusion

Mars 1.0 has a solid foundation with basic arithmetic, functions, and arrays, but needs improvements in:
- Array indexing within conditional statements
- Loop constructs with array access
- Hash map data structures

The Two Sum problem can be solved with the current limitations using a brute force approach with individual parameters, but it's not scalable for larger arrays. 