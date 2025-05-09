package ast

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
}

// Statement represents a statement node in the AST
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node in the AST
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST
type Program struct {
	Declarations []Declaration
}

// Declaration represents a declaration node in the AST
type Declaration interface {
	Node
	declarationNode()
}

// VarDecl represents a variable declaration
type VarDecl struct {
	IsMutable bool
	Name      *Identifier
	Type      *Type
	Value     Expression
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
	Condition Expression
	Then      *BlockStatement
	Else      *BlockStatement
}

// ForStatement represents a for statement
type ForStatement struct {
	Init      Declaration
	Condition Expression
	Post      Expression
	Body      *BlockStatement
}

// PrintStatement represents a print statement
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
type Type struct {
	BaseType    string
	ArrayType   *Type
	StructType  *Identifier
	PointerType *Type
}

// ArrayLiteral represents an array literal
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

// FunctionCall represents a function call
type FunctionCall struct {
	Function  Expression
	Arguments []Expression
}

// BinaryExpression represents a binary expression
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

// UnaryExpression represents a unary expression
type UnaryExpression struct {
	Operator string
	Right    Expression
}

// Literal represents a literal value
type Literal struct {
	Type  string // "number", "string", "boolean", "nil"
	Value interface{}
}

// TokenLiteral implementations
func (p *Program) TokenLiteral() string {
	if len(p.Declarations) > 0 {
		return p.Declarations[0].TokenLiteral()
	}
	return ""
}

func (vd *VarDecl) TokenLiteral() string             { return vd.Name.TokenLiteral() }
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
func (t *Type) TokenLiteral() string                 { return t.BaseType }
func (al *ArrayLiteral) TokenLiteral() string        { return "[" }
func (sl *StructLiteral) TokenLiteral() string       { return sl.Type.TokenLiteral() }
func (fc *FunctionCall) TokenLiteral() string        { return fc.Function.TokenLiteral() }
func (be *BinaryExpression) TokenLiteral() string    { return be.Left.TokenLiteral() }
func (ue *UnaryExpression) TokenLiteral() string     { return ue.Operator }
func (l *Literal) TokenLiteral() string              { return l.Type }

// Node type implementations
func (vd *VarDecl) declarationNode()           {}
func (fd *FuncDecl) declarationNode()          {}
func (sd *StructDecl) declarationNode()        {}
func (ub *UnsafeBlock) declarationNode()       {}
func (bs *BlockStatement) statementNode()      {}
func (is *IfStatement) statementNode()         {}
func (fs *ForStatement) statementNode()        {}
func (ps *PrintStatement) statementNode()      {}
func (rs *ReturnStatement) statementNode()     {}
func (es *ExpressionStatement) statementNode() {}
func (i *Identifier) expressionNode()          {}
func (al *ArrayLiteral) expressionNode()       {}
func (sl *StructLiteral) expressionNode()      {}
func (fc *FunctionCall) expressionNode()       {}
func (be *BinaryExpression) expressionNode()   {}
func (ue *UnaryExpression) expressionNode()    {}
func (l *Literal) expressionNode()             {}
