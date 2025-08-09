# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2025-08-09

### Added
- Struct literals and member access (parser lookahead, analyzer validation, evaluator runtime `StructValue`).
- While loops end-to-end (lexer, AST, parser, evaluator).
- Modulo operator `%` with integer semantics.
- Clearer, symbol-based error messages with source context.
- Minimal example set (two-sum, three-sum, trapping-rain-water, maximum-subarray, best-time-to-buy-sell-stock-III).

### Changed
- Parser lookahead via non-consuming `PeekTokenN`; stabilized ambiguous `IDENT {` disambiguation.
- Analyzer error messaging for variable declarations, struct fields, and control flow.

### Known limitations
- Strings: char literals, escapes, indexing/slicing not fully supported.
- `for` condition-only loops not yet supported (use `while`).
- Multi-arg `println` not supported; use single-arg prints.

