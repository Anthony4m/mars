// ast/ast.go
package ast

import (
	"fmt"
)

// Node represents a node in the AST
// Pos returns line and column for error reporting (optional)
type Node interface {
	TokenLiteral() string
}

// Declaration represents a declaration or top-level statement
// (grammar: declaration = varDecl | funcDecl | unsafeBlock | statement)
type Declaration interface {
	Node
	declarationNode()
}

// Statement represents a statement node
// Statement nodes also satisfy Declaration at top level
type Statement interface {
	Node
	statementNode()
	// implement declarationNode to allow top-level statements
	declarationNode()
	// String returns a string representation of the statement
	String() string
}

// Expression represents an expression node in the AST
type Expression interface {
	Node
	expressionNode()
	// String returns a string representation of the expression
	String() string
}

// Program is the root node of the AST
// It contains both declarations and top-level statements
type Program struct {
	Declarations []Declaration
}

// VarDecl represents a variable declaration
type VarDecl struct {
	Mutable bool
	Name    *Identifier
	Type    *Type
	Value   Expression
}

// AssignmentStatement represents mutable variable assignment
type AssignmentStatement struct {
	Name  *Identifier
	Value Expression
}

// FuncDecl represents a function declaration
type FuncDecl struct {
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *Type
	Body       *BlockStatement
}

// StructDecl represents a struct declaration
type StructDecl struct {
	Name   *Identifier
	Fields []*FieldDecl
}

// UnsafeBlock represents an unsafe block
type UnsafeBlock struct {
	Body *BlockStatement
}

// Parameter represents a function parameter
type Parameter struct {
	Name *Identifier
	Type *Type
}

// FieldDecl represents a struct field declaration
type FieldDecl struct {
	Name *Identifier
	Type *Type
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Statements []Statement
}

// IfStatement represents an if statement
type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// ForStatement represents a for statement
type ForStatement struct {
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

// PrintStatement represents a print/log statement
type PrintStatement struct {
	Expression Expression
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value Expression
}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	Expression Expression
}

// Identifier represents an identifier
type Identifier struct {
	Name string
}

// Type represents a type
// Supports primitives, arrays, pointers, and structs
type Type struct {
	BaseType    string // "int", "float", "string", "bool"
	ArrayType   *Type  // For []T or [N]T
	ArraySize   *int   // nil for dynamic slices, value for fixed arrays
	PointerType *Type  // For *T
	StructName  string // For struct references
	MapType     *Type  // For map[K]V
}

// ArrayLiteral represents an array or slice literal
type ArrayLiteral struct {
	Elements []Expression
}

// StructLiteral represents a struct literal
type StructLiteral struct {
	Type   *Identifier
	Fields []*FieldInit
}

// FieldInit represents a struct field initialization
type FieldInit struct {
	Name  *Identifier
	Value Expression
}

// FunctionCall represents a function or method call
type FunctionCall struct {
	Function  Expression
	Arguments []Expression
}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

// UnaryExpression represents a unary operation
type UnaryExpression struct {
	Operator string
	Right    Expression
}

// Literal represents a literal value (number, string, boolean, nil)
type Literal struct {
	Token string
	Value interface{}
}

// MemberExpression represents object.member access
type MemberExpression struct {
	Object   Expression
	Property *Identifier
}

// BreakStatement represents a break within loops
type BreakStatement struct{}

// ContinueStatement represents a continue within loops
type ContinueStatement struct{}

// IndexExpression represents array indexing (a[i])
type IndexExpression struct {
	Object Expression
	Index  Expression
}

type SliceExpression struct {
	Object Expression
	Start  Expression // Can be nil for [:end]
	End    Expression // Can be nil for [start:]
}

// MapLiteral represents a map literal
type MapLiteral struct {
	KeyType   *Type
	ValueType *Type
	Elements  []Expression
}

// TokenLiteral implementations
func (p *Program) TokenLiteral() string {
	if len(p.Declarations) > 0 {
		return p.Declarations[0].TokenLiteral()
	}
	return ""
}
func (vd *VarDecl) TokenLiteral() string             { return vd.Name.TokenLiteral() }
func (as *AssignmentStatement) TokenLiteral() string { return "=" }
func (fd *FuncDecl) TokenLiteral() string            { return fd.Name.TokenLiteral() }
func (sd *StructDecl) TokenLiteral() string          { return sd.Name.TokenLiteral() }
func (ub *UnsafeBlock) TokenLiteral() string         { return "unsafe" }
func (bs *BlockStatement) TokenLiteral() string      { return "{" }
func (is *IfStatement) TokenLiteral() string         { return "if" }
func (fs *ForStatement) TokenLiteral() string        { return "for" }
func (ps *PrintStatement) TokenLiteral() string      { return "log" }
func (rs *ReturnStatement) TokenLiteral() string     { return "return" }
func (es *ExpressionStatement) TokenLiteral() string { return es.Expression.TokenLiteral() }
func (i *Identifier) TokenLiteral() string           { return i.Name }
func (al *ArrayLiteral) TokenLiteral() string        { return "[" }
func (sl *StructLiteral) TokenLiteral() string       { return sl.Type.TokenLiteral() }
func (fc *FunctionCall) TokenLiteral() string        { return fc.Function.TokenLiteral() }
func (be *BinaryExpression) TokenLiteral() string    { return be.Operator }
func (ue *UnaryExpression) TokenLiteral() string     { return ue.Operator }
func (l *Literal) TokenLiteral() string              { return l.Token }
func (me *MemberExpression) TokenLiteral() string    { return me.Object.TokenLiteral() }
func (bs *BreakStatement) TokenLiteral() string      { return "break" }
func (cs *ContinueStatement) TokenLiteral() string   { return "continue" }
func (ie *IndexExpression) TokenLiteral() string     { return "[" }
func (se *SliceExpression) TokenLiteral() string     { return "[" }
func (ml *MapLiteral) TokenLiteral() string          { return "map" }

// Node type implementations
func (vd *VarDecl) declarationNode()             {}
func (vd *VarDecl) statementNode()               {}
func (as *AssignmentStatement) statementNode()   {}
func (as *AssignmentStatement) declarationNode() {}
func (fd *FuncDecl) declarationNode()            {}
func (sd *StructDecl) declarationNode()          {}
func (ub *UnsafeBlock) declarationNode()         {}
func (bs *BlockStatement) statementNode()        {}
func (bs *BlockStatement) declarationNode()      {}
func (is *IfStatement) statementNode()           {}
func (is *IfStatement) declarationNode()         {}
func (fs *ForStatement) statementNode()          {}
func (fs *ForStatement) declarationNode()        {}
func (ps *PrintStatement) statementNode()        {}
func (ps *PrintStatement) declarationNode()      {}
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) declarationNode()     {}
func (es *ExpressionStatement) statementNode()   {}
func (es *ExpressionStatement) declarationNode() {}
func (i *Identifier) expressionNode()            {}
func (al *ArrayLiteral) expressionNode()         {}
func (sl *StructLiteral) expressionNode()        {}
func (fc *FunctionCall) expressionNode()         {}
func (be *BinaryExpression) expressionNode()     {}
func (ue *UnaryExpression) expressionNode()      {}
func (l *Literal) expressionNode()               {}
func (me *MemberExpression) expressionNode()     {}
func (bs *BreakStatement) statementNode()        {}
func (bs *BreakStatement) declarationNode()      {}
func (cs *ContinueStatement) statementNode()     {}
func (cs *ContinueStatement) declarationNode()   {}
func (ie *IndexExpression) expressionNode()      {}
func (se *SliceExpression) expressionNode()      {}
func (ml *MapLiteral) expressionNode()           {}

// method to check if type is a slice vs fixed array
func (t *Type) IsSlice() bool {
	return t.ArrayType != nil && t.ArraySize == nil
}

func (t *Type) IsFixedArray() bool {
	return t.ArrayType != nil && t.ArraySize != nil
}

// Helper constructor functions for common types
func NewBaseType(name string) *Type {
	return &Type{BaseType: name}
}

func NewSliceType(elementType *Type) *Type {
	return &Type{ArrayType: elementType}
}

func NewArrayType(elementType *Type, size int) *Type {
	return &Type{ArrayType: elementType, ArraySize: &size}
}

func NewPointerType(pointeeType *Type) *Type {
	return &Type{PointerType: pointeeType}
}

func NewStructType(name string) *Type {
	return &Type{StructName: name}
}

func NewMapType(keyType *Type, valueType *Type) *Type {
	return &Type{MapType: &Type{BaseType: "map", ArrayType: keyType}}
}

// String method for better debugging
func (t *Type) String() string {
	if t.BaseType != "" {
		return t.BaseType
	}
	if t.ArrayType != nil {
		if t.ArraySize != nil {
			return fmt.Sprintf("[%d]%s", *t.ArraySize, t.ArrayType.String())
		}
		return fmt.Sprintf("[]%s", t.ArrayType.String())
	}
	if t.PointerType != nil {
		return fmt.Sprintf("*%s", t.PointerType.String())
	}
	if t.StructName != "" {
		return fmt.Sprintf("struct %s", t.StructName)
	}
	return "unknown"
}

// String implementations for statements
func (vd *VarDecl) String() string {
	var s string
	if vd.Mutable {
		s += "mut "
	}
	s += vd.Name.Name + " : " + vd.Type.String()
	if vd.Value != nil {
		s += " := " + vd.Value.String()
	}
	return s + ";"
}

func (as *AssignmentStatement) String() string {
	return as.Name.Name + " = " + as.Value.String() + ";"
}

func (fd *FuncDecl) String() string {
	var s string
	s += "func " + fd.Name.Name + "("
	for i, param := range fd.Parameters {
		if i > 0 {
			s += ", "
		}
		s += param.Name.Name + " : " + param.Type.String()
	}
	s += ")"
	if fd.ReturnType != nil {
		s += " -> " + fd.ReturnType.String()
	}
	s += " " + fd.Body.String()
	return s
}

func (sd *StructDecl) String() string {
	var s string
	s += "struct " + sd.Name.Name + " {"
	for _, field := range sd.Fields {
		s += "\n\t" + field.Name.Name + " : " + field.Type.String() + ";"
	}
	s += "\n}"
	return s
}

func (ub *UnsafeBlock) String() string {
	return "unsafe " + ub.Body.String()
}

func (bs *BlockStatement) String() string {
	var s string
	s += "{"
	for _, stmt := range bs.Statements {
		s += "\n\t" + stmt.String()
	}
	s += "\n}"
	return s
}

func (is *IfStatement) String() string {
	var s string
	s += "if " + is.Condition.String() + " " + is.Consequence.String()
	if is.Alternative != nil {
		s += " else " + is.Alternative.String()
	}
	return s
}

func (fs *ForStatement) String() string {
	var s string
	s += "for "
	if fs.Init != nil {
		s += fs.Init.String()
	}
	s += "; "
	if fs.Condition != nil {
		s += fs.Condition.String()
	}
	s += "; "
	if fs.Post != nil {
		s += fs.Post.String()
	}
	s += " " + fs.Body.String()
	return s
}

func (ps *PrintStatement) String() string {
	return "log(" + ps.Expression.String() + ");"
}

func (rs *ReturnStatement) String() string {
	if rs.Value == nil {
		return "return;"
	}
	return "return " + rs.Value.String() + ";"
}

func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return "<nil expression>;"
	}
	return es.Expression.String() + ";"
}
func (bs *BreakStatement) String() string {
	return "break;"
}

func (cs *ContinueStatement) String() string {
	return "continue;"
}

// String implementations for expressions
func (i *Identifier) String() string {
	return i.Name
}

func (al *ArrayLiteral) String() string {
	var s string
	s += "["
	for i, elem := range al.Elements {
		if i > 0 {
			s += ", "
		}
		s += elem.String()
	}
	s += "]"
	return s
}

func (sl *StructLiteral) String() string {
	var s string
	s += sl.Type.Name + "{"
	for i, field := range sl.Fields {
		if i > 0 {
			s += ", "
		}
		s += field.Name.Name + ": " + field.Value.String()
	}
	s += "}"
	return s
}

func (fc *FunctionCall) String() string {
	var s string
	s += fc.Function.String() + "("
	for i, arg := range fc.Arguments {
		if i > 0 {
			s += ", "
		}
		s += arg.String()
	}
	s += ")"
	return s
}

func (be *BinaryExpression) String() string {
	return "(" + be.Left.String() + " " + be.Operator + " " + be.Right.String() + ")"
}

func (ue *UnaryExpression) String() string {
	return "(" + ue.Operator + ue.Right.String() + ")"
}

func (l *Literal) String() string {
	switch v := l.Value.(type) {
	case string:
		return "\"" + v + "\""
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (me *MemberExpression) String() string {
	return me.Object.String() + "." + me.Property.Name
}

func (ie *IndexExpression) String() string {
	return ie.Object.String() + "[" + ie.Index.String() + "]"
}

func (se *SliceExpression) String() string {
	var s string
	s += se.Object.String() + "["
	if se.Start != nil {
		s += se.Start.String()
	}
	s += ":"
	if se.End != nil {
		s += se.End.String()
	}
	s += "]"
	return s
}

func (ml *MapLiteral) String() string {
	var s string
	s += "map[" + ml.KeyType.String() + "]" + ml.ValueType.String() + "{"
	for i, elem := range ml.Elements {
		if i > 0 {
			s += ", "
		}
		s += elem.String()
	}
	s += "}"
	return s
}
