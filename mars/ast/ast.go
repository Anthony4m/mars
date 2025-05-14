// ast/ast.go
package ast

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
func (n *NumberLiteral) TokenLiteral() string { return n.Value }

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
