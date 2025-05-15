// ast/ast.go
package ast

import "bytes"

// Node represents a node in the Abstract Syntax Tree.
type Node interface {
	TokenLiteral() string
	String() string
}

// TypeNode represents a type node in the AST.
// It's an interface to allow for different kinds of type representations (e.g., simple identifiers, slice types, pointer types).
type TypeNode interface {
	Node
	typeNode() // marker method
}

// Declaration represents a declaration or top-level statement
// (grammar: declaration = varDecl | funcDecl | unsafeBlock | statement)
type Declaration interface {
	Node
	declarationNode()
}

// Statement represents a statement node in the AST.
// Statement nodes also satisfy Declaration at top level
type Statement interface {
	Node
	statementNode()
	// implement declarationNode to allow top-level statements
	declarationNode()
}

// Expression represents an expression node in the AST
type Expression interface {
	Node
	expressionNode()
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
	BaseType    string
	ArrayType   *Type
	StructType  *Identifier
	PointerType *Type
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

// NumberLiteral represents a numeric literal
type NumberLiteral struct {
	Value string
}

// TokenLiteral implementations
func (p *Program) TokenLiteral() string {
	if len(p.Declarations) > 0 {
		return p.Declarations[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Declarations {
		out.WriteString(s.String())
	}
	return out.String()
}

// In your ast package:
type ParameterNode struct {
	// NodeImpl // If you have a base AST node interface/struct
	Name *Identifier
	Type TypeNode // ast.TypeNode could be an interface, or *Identifier if types are simple
}

func (vd *VarDecl) TokenLiteral() string { return vd.Name.TokenLiteral() }
func (vd *VarDecl) String() string {
	var out bytes.Buffer
	if vd.Mutable {
		out.WriteString("mut ")
	} else {
		out.WriteString("var ")
	}
	out.WriteString(vd.Name.String())
	if vd.Type != nil {
		out.WriteString(" : ")
		out.WriteString(vd.Type.String())
	}
	if vd.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vd.Value.String())
	}
	return out.String()
}

func (as *AssignmentStatement) TokenLiteral() string { return "=" }
func (as *AssignmentStatement) String() string {
	// Placeholder, can be more detailed
	return as.Name.String() + " = ..."
}

func (fd *FuncDecl) TokenLiteral() string { return fd.Name.TokenLiteral() }
func (fd *FuncDecl) String() string {
	// Placeholder, can be more detailed
	return "func " + fd.Name.String() + "()"
}

func (sd *StructDecl) TokenLiteral() string { return sd.Name.TokenLiteral() }
func (sd *StructDecl) String() string {
	// Placeholder
	return "struct " + sd.Name.String() + " {}"
}

func (ub *UnsafeBlock) TokenLiteral() string { return "unsafe" }
func (ub *UnsafeBlock) String() string {
	// Placeholder
	return "unsafe {}"
}

func (bs *BlockStatement) TokenLiteral() string { return "{" }
func (bs *BlockStatement) String() string {
	// Placeholder
	return "{ ... }"
}

func (is *IfStatement) TokenLiteral() string { return "if" }
func (is *IfStatement) String() string {
	// Placeholder
	return "if (...) { ... }"
}

func (fs *ForStatement) TokenLiteral() string { return "for" }
func (fs *ForStatement) String() string {
	// Placeholder
	return "for { ... }"
}

func (ps *PrintStatement) TokenLiteral() string { return "log" }
func (ps *PrintStatement) String() string {
	// Placeholder
	return "log(...)"
}

func (rs *ReturnStatement) TokenLiteral() string { return "return" }
func (rs *ReturnStatement) String() string {
	// Placeholder
	return "return ..."
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Expression.TokenLiteral() }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier) TokenLiteral() string    { return i.Name }
func (i *Identifier) String() string          { return i.Name }
func (i *Identifier) typeNode()               {}
func (al *ArrayLiteral) TokenLiteral() string { return "[" }
func (al *ArrayLiteral) String() string {
	// Placeholder
	return "[...]"
}

func (sl *StructLiteral) TokenLiteral() string { return sl.Type.TokenLiteral() }
func (sl *StructLiteral) String() string {
	// Placeholder
	return sl.Type.String() + "{...}"
}

func (fc *FunctionCall) TokenLiteral() string { return fc.Function.TokenLiteral() }
func (fc *FunctionCall) String() string {
	// Placeholder
	return fc.Function.String() + "(...)"
}

func (be *BinaryExpression) TokenLiteral() string { return be.Operator }
func (be *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.Operator + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")
	return out.String()
}

func (ue *UnaryExpression) TokenLiteral() string { return ue.Operator }
func (ue *UnaryExpression) String() string {
	// Placeholder
	return ue.Operator + "(...)"
}

func (l *Literal) TokenLiteral() string { return l.Token }
func (l *Literal) String() string       { return l.Token }

func (me *MemberExpression) TokenLiteral() string { return me.Object.TokenLiteral() }
func (me *MemberExpression) String() string {
	// Placeholder
	return me.Object.String() + "." + me.Property.String()
}

func (bs *BreakStatement) TokenLiteral() string { return "break" }
func (bs *BreakStatement) String() string       { return "break" }

func (cs *ContinueStatement) TokenLiteral() string { return "continue" }
func (cs *ContinueStatement) String() string       { return "continue" }

func (ie *IndexExpression) TokenLiteral() string { return "[" }
func (ie *IndexExpression) String() string {
	// Placeholder
	return "(...)[...]"
}

func (t *Type) TokenLiteral() string {
	if t.ArrayType != nil {
		return "[]" + t.ArrayType.TokenLiteral()
	}
	if t.PointerType != nil {
		return "*" + t.PointerType.TokenLiteral()
	}
	if t.StructType != nil {
		return t.StructType.TokenLiteral()
	}
	return t.BaseType
}

func (t *Type) String() string {
	// This makes Type satisfy the Node interface if needed as part of TypeNode.
	return t.TokenLiteral()
}

func (n *NumberLiteral) TokenLiteral() string { return n.Value }
func (n *NumberLiteral) String() string       { return n.Value }

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
func (n *NumberLiteral) expressionNode()         {}
