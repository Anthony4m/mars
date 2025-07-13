// ast/ast.go
package ast

import (
	"fmt"
)

// Position represents a position in the source code
type Position struct {
	Line   int
	Column int
}

// Node represents a node in the AST
// Pos returns line and column for error reporting (optional)
type Node interface {
	TokenLiteral() string
	Pos() Position
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
	Position     Position
}

// VarDecl represents a variable declaration
type VarDecl struct {
	Mutable  bool
	Name     *Identifier
	Type     *Type
	Value    Expression
	Position Position
}

// AssignmentStatement represents mutable variable assignment
type AssignmentStatement struct {
	Name     *Identifier
	Value    Expression
	Position Position
}

// FuncDecl represents a function declaration
type FuncDecl struct {
	Name      *Identifier
	Signature *FunctionSignature
	Body      *BlockStatement
	Position  Position
}

// StructDecl represents a struct declaration
type StructDecl struct {
	Name     *Identifier
	Fields   []*FieldDecl
	Position Position
}

// UnsafeBlock represents an unsafe block
type UnsafeBlock struct {
	Body     *BlockStatement
	Position Position
}

// Parameter represents a function parameter
type Parameter struct {
	Name     *Identifier
	Type     *Type
	Position Position
}

// FieldDecl represents a struct field declaration
type FieldDecl struct {
	Name     *Identifier
	Type     *Type
	Position Position
}

// BlockStatement represents a block of statements
type BlockStatement struct {
	Statements []Statement
	Position   Position
}

// IfStatement represents an if statement
type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Position    Position
}

// ForStatement represents a for statement
type ForStatement struct {
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
	Position  Position
}

// PrintStatement represents a print/log statement
type PrintStatement struct {
	Expression Expression
	Position   Position
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value    Expression
	Position Position
}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	Expression Expression
	Position   Position
}

// Identifier represents an identifier
type Identifier struct {
	Name     string
	Position Position
}

// FunctionSignature represents a function's type signature
type FunctionSignature struct {
	Parameters []*Parameter
	ReturnType *Type
	Position   Position
}

// Type represents a type
// Supports primitives, arrays, pointers, and structs
type Type struct {
	BaseType     string       // "int", "float", "string", "bool"
	ArrayType    *Type        // For []T or [N]T
	ArraySize    *int         // nil for dynamic slices, value for fixed arrays
	PointerType  *Type        // For *T
	StructName   string       // For struct references
	StructFields []*FieldDecl // For struct types - stores the field declarations
	MapType      *Type        // For map[K]V
	Position     Position
	// Function signature for function types
	FunctionSignature *FunctionSignature
}

// ArrayLiteral represents an array or slice literal
type ArrayLiteral struct {
	Elements []Expression
	Position Position
}

// StructLiteral represents a struct literal
type StructLiteral struct {
	Type     *Identifier
	Fields   []*FieldInit
	Position Position
}

// FieldInit represents a struct field initialization
type FieldInit struct {
	Name     *Identifier
	Value    Expression
	Position Position
}

// FunctionCall represents a function or method call
type FunctionCall struct {
	Function  Expression
	Arguments []Expression
	Position  Position
}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	Position Position
}

// UnaryExpression represents a unary operation
type UnaryExpression struct {
	Operator string
	Right    Expression
	Position Position
}

// Literal represents a literal value (number, string, boolean, nil)
type Literal struct {
	Token    string
	Value    interface{}
	Position Position
}

// MemberExpression represents object.member access
type MemberExpression struct {
	Object   Expression
	Property *Identifier
	Position Position
}

// BreakStatement represents a break within loops
type BreakStatement struct {
	Position Position
}

// ContinueStatement represents a continue within loops
type ContinueStatement struct {
	Position Position
}

// IndexExpression represents array indexing (a[i])
type IndexExpression struct {
	Object   Expression
	Index    Expression
	Position Position
}

type SliceExpression struct {
	Object   Expression
	Start    Expression // Can be nil for [:end]
	End      Expression // Can be nil for [start:]
	Position Position
}

// MapLiteral represents a map literal
type MapLiteral struct {
	KeyType   *Type
	ValueType *Type
	Elements  []Expression
	Position  Position
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

// Position implementations
func (p *Program) Pos() Position              { return p.Position }
func (vd *VarDecl) Pos() Position             { return vd.Position }
func (as *AssignmentStatement) Pos() Position { return as.Position }
func (fd *FuncDecl) Pos() Position            { return fd.Position }
func (sd *StructDecl) Pos() Position          { return sd.Position }
func (ub *UnsafeBlock) Pos() Position         { return ub.Position }
func (bs *BlockStatement) Pos() Position      { return bs.Position }
func (is *IfStatement) Pos() Position         { return is.Position }
func (fs *ForStatement) Pos() Position        { return fs.Position }
func (ps *PrintStatement) Pos() Position      { return ps.Position }
func (rs *ReturnStatement) Pos() Position     { return rs.Position }
func (es *ExpressionStatement) Pos() Position { return es.Position }
func (i *Identifier) Pos() Position           { return i.Position }
func (al *ArrayLiteral) Pos() Position        { return al.Position }
func (sl *StructLiteral) Pos() Position       { return sl.Position }
func (fc *FunctionCall) Pos() Position        { return fc.Position }
func (be *BinaryExpression) Pos() Position    { return be.Position }
func (ue *UnaryExpression) Pos() Position     { return ue.Position }
func (l *Literal) Pos() Position              { return l.Position }
func (me *MemberExpression) Pos() Position    { return me.Position }
func (bs *BreakStatement) Pos() Position      { return bs.Position }
func (cs *ContinueStatement) Pos() Position   { return cs.Position }
func (ie *IndexExpression) Pos() Position     { return ie.Position }
func (se *SliceExpression) Pos() Position     { return se.Position }
func (ml *MapLiteral) Pos() Position          { return ml.Position }

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

func NewStructType(name string, fields []*FieldDecl) *Type {
	return &Type{
		StructName:   name,
		StructFields: fields,
	}
}

func NewMapType(keyType *Type, valueType *Type) *Type {
	return &Type{MapType: &Type{BaseType: "map", ArrayType: keyType}}
}

// NewFunctionType creates a function type with the given signature
func NewFunctionType(signature *FunctionSignature) *Type {
	return &Type{
		BaseType:          "function",
		FunctionSignature: signature,
	}
}

// IsFunctionType checks if the type represents a function
func (t *Type) IsFunctionType() bool {
	return t.BaseType == "function" && t.FunctionSignature != nil
}

// GetFunctionSignature returns the function signature if this is a function type
func (t *Type) GetFunctionSignature() *FunctionSignature {
	if t.IsFunctionType() {
		return t.FunctionSignature
	}
	return nil
}

// IsStructType checks if the type represents a struct
func (t *Type) IsStructType() bool {
	return t.BaseType != "struct"
}

// GetStructFields returns the struct fields if this is a struct type
func (t *Type) GetStructFields() []*FieldDecl {
	if t.IsStructType() {
		return t.StructFields
	}
	return nil
}

// String method for better debugging
func (t *Type) String() string {
	if t.BaseType != "" {
		if t.IsFunctionType() {
			return "func" + t.FunctionSignature.String()
		}
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
		if len(t.StructFields) > 0 {
			var s string
			s += fmt.Sprintf("struct %s {", t.StructName)
			for _, field := range t.StructFields {
				s += fmt.Sprintf("\n\t%s : %s;", field.Name.Name, field.Type.String())
			}
			s += "\n}"
			return s
		}
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
	for i, param := range fd.Signature.Parameters {
		if i > 0 {
			s += ", "
		}
		s += param.Name.Name + " : " + param.Type.String()
	}
	s += ")"
	if fd.Signature.ReturnType != nil {
		s += " -> " + fd.Signature.ReturnType.String()
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

// String returns a string representation of the function signature
func (fs *FunctionSignature) String() string {
	var s string
	s += "("
	for i, param := range fs.Parameters {
		if i > 0 {
			s += ", "
		}
		s += param.Name.Name + " : " + param.Type.String()
	}
	s += ")"
	if fs.ReturnType != nil {
		s += " -> " + fs.ReturnType.String()
	}
	return s
}
