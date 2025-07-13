# Mars Language Grammar

This document defines the formal grammar of the Mars programming language in Extended Backus-Naur Form (EBNF).

## Grammar Rules

```
Program       = { Declaration } EOF ;

Declaration   = VarDecl
              | FuncDecl
              | StructDecl
              | UnsafeBlock
              | Statement ;

VarDecl       = [ "mut" ] IDENT ":" Type [ ":=" Expression ] ";" ;

FuncDecl      = "func" IDENT "(" [ Params ] ")" [ "->" Type ] Block ;
Params        = Param ( "," Param )* ;
Param         = IDENT ":" Type ;

StructDecl    = "struct" IDENT "{" { FieldDecl } "}" ;
FieldDecl     = IDENT ":" Type ";" ;

UnsafeBlock   = "unsafe" Block ;

Statement     = AssignmentStmt
              | ExprStmt
              | IfStmt
              | ForStmt
              | PrintStmt
              | ReturnStmt
              | Block ;

AssignmentStmt= IDENT "=" Expression ";" ;
ExprStmt      = Expression ";" ;

IfStmt        = "if" Expression Block [ "else" ( IfStmt | Block ) ] ;
ForStmt       = "for" [ Init ] ";" [ Condition ] ";" [ Post ] Block ;
Init          = VarDecl | ExprStmt ;
PrintStmt     = "log" "(" Expression ")" ";" ;
ReturnStmt    = "return" [ Expression ] ";" ;
Block         = "{" { Declaration } "}" ;

Expression    = LogicalOr ;
LogicalOr     = LogicalAnd { "||" LogicalAnd } ;
LogicalAnd    = Equality   { "&&" Equality } ;
Equality      = Comparison { ( "==" | "!=" ) Comparison } ;
Comparison    = Term       { ( ">" | ">=" | "<" | "<=" ) Term } ;
Term          = Factor     { ( "+" | "-" ) Factor } ;
Factor        = Unary      { ( "*" | "/" | "%" ) Unary } ;
Unary         = ( "!" | "-" ) Unary | Primary ;
Primary       = Literal
              | IDENT "(" [ Args ] ")"
              | IDENT
              | "(" Expression ")"
              | ArrayLit
              | StructLit ;

ArrayLit      = "[" [ Expression ( "," Expression )* ] "]" ;
StructLit     = IDENT "{" [ FieldInit ( "," FieldInit )* ] "}" ;
FieldInit     = IDENT ":" Expression ;
Args          = Expression ( "," Expression )* ;

Type          = BaseType
              | ArrayType
              | StructType
              | PointerType ;

BaseType      = "int" | "float" | "string" | "bool" ;
ArrayType     = ( "[" [ INTEGER ] "]" | "[]" ) Type ;
StructType    = "struct" IDENT ;
PointerType   = "*" Type ;

Literal       = NUMBER | STRING | BOOLEAN | "nil" ;
BOOLEAN       = "true" | "false" ;
IDENT         = LETTER ( LETTER | DIGIT | "_" )* ;
NUMBER        = INTEGER | FLOAT ;
INTEGER       = DIGIT+ ;
FLOAT         = DIGIT+ "." DIGIT+ ;
STRING        = "\"" ( CHAR | ESCAPE )* "\"" ;
NILL          = "nil" ;
```

## Lexical Elements

### Keywords
- `mut`: Declares a mutable variable
- `func`: Function declaration
- `struct`: Structure declaration
- `unsafe`: Unsafe code block
- `if`: Conditional statement
- `else`: Alternative branch
- `for`: Loop statement
- `return`: Function return
- `log`: Print statement

### Types
- `int`: Integer type
- `float`: Floating-point type
- `string`: String type
- `bool`: Boolean type

### Operators
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `==`, `!=`, `>`, `>=`, `<`, `<=`
- Logical: `&&`, `||`, `!`
- Assignment: `=`, `:=`

### Special Tokens
- `nil`: Null value
- `true`, `false`: Boolean literals

## Notes
1. All statements must end with a semicolon (`;`)
2. Block statements are enclosed in curly braces (`{` `}`)
3. Function parameters and return types are separated by `->`
4. Type declarations use a colon (`:`) syntax
5. Array types can be fixed-size `[N]Type` or dynamic `[]Type` 