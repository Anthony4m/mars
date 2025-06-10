// parser/parser.go
package parser

import (
	"fmt"
	"mars/ast"
	"mars/lexer"
	"strconv"
)

type parser struct {
	lexer     *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	prevToken lexer.Token
	errors    []string
}

func NewParser(lexer *lexer.Lexer) *parser {
	p := &parser{lexer: lexer, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *parser) recordError(message string) {
	p.errors = append(p.errors, message)
}

func (p *parser) nextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *parser) previousToken() lexer.Token {
	return p.prevToken
}

func (p *parser) GetErrors() []string {
	return p.errors
}

// ParseProgram is the main entry point for parsing.
func (p *parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Declarations = []ast.Declaration{}

	for p.curToken.Type != lexer.EOF {
		decl := p.parseDeclaration()
		if decl != nil {
			program.Declarations = append(program.Declarations, decl)
		}
	}

	return program
}

// ===== TYPE PARSING =====

// parseType parses type expressions using your unified Type struct
func (p *parser) parseType() *ast.Type {
	switch p.curToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.STRING_KW, lexer.BOOL:
		return p.parseBaseType()
	case lexer.LBRACKET:
		return p.parseArrayOrSliceType()
	case lexer.ASTERISK:
		return p.parsePointerType()
	case lexer.IDENT:
		return p.parseStructTypeReference()
	default:
		p.recordError(fmt.Sprintf("expected type, got %s", p.curToken.Type))
		return nil
	}
}

func (p *parser) parseBaseType() *ast.Type {
	baseType := &ast.Type{BaseType: p.curToken.Literal}
	p.nextToken() // consume the type token
	return baseType
}

func (p *parser) parseArrayOrSliceType() *ast.Type {
	p.nextToken() // consume '['

	arrayType := &ast.Type{}

	// Check if it's a sized array [N] or dynamic slice []
	if p.curTokenIs(lexer.NUMBER) {
		// Parse array size
		if size, err := strconv.Atoi(p.curToken.Literal); err == nil {
			arrayType.ArraySize = &size
		} else {
			p.recordError("invalid array size")
			return nil
		}
		p.nextToken() // consume size
	}
	// else: dynamic slice (ArraySize remains nil)

	if !p.expectCurrent(lexer.RBRACKET) {
		return nil
	}

	// Parse element type
	arrayType.ArrayType = p.parseType()
	return arrayType
}

func (p *parser) parsePointerType() *ast.Type {
	p.nextToken() // consume '*'
	pointeeType := p.parseType()
	if pointeeType == nil {
		return nil
	}
	return &ast.Type{PointerType: pointeeType}
}

func (p *parser) parseStructTypeReference() *ast.Type {
	structType := &ast.Type{StructName: p.curToken.Literal}
	p.nextToken() // consume struct name
	return structType
}

// ===== DECLARATION PARSING =====

func (p *parser) parseDeclaration() ast.Declaration {
	switch p.curToken.Type {
	case lexer.FUNC:
		return p.parseFunctionDeclaration()
	case lexer.MUT:
		return p.parseVariableDeclaration()
	case lexer.STRUCT:
		return p.parseStructDeclaration()
	case lexer.ENUM:
		return p.parseEnumDeclaration()
	case lexer.TYPE:
		return p.parseTypeDeclaration()
	case lexer.UNSAFE:
		return p.parseUnsafeDeclaration()
	case lexer.IF, lexer.FOR, lexer.RETURN, lexer.LOG, lexer.BREAK, lexer.CONTINUE:
		return p.parseStatement()
	case lexer.IDENT:
		// Check for different identifier contexts
		if p.peekTokenIs(lexer.COLON) {
			return p.parseVariableDeclaration() // x : int = 5
		} else if p.peekTokenIs(lexer.COLONEQ) {
			return p.parseVariableDeclaration() // x := 5
		} else if p.peekTokenIs(lexer.EQ) {
			return p.parseAssignment() // x = 5
		}
		return p.parseExpressionStatement()
	default:
		p.recordError(fmt.Sprintf("unexpected token %s at top level", p.curToken.Type))
		p.synchronize()
		return nil
	}
}

// parseFunctionDeclaration handles: "func" IDENT "(" [ Params ] ")" [ "->" Type ] Block
func (p *parser) parseFunctionDeclaration() ast.Declaration {
	funcDecl := &ast.FuncDecl{}
	p.nextToken() // consume "func"

	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected function name")
		return nil
	}
	funcDecl.Name = &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken() // consume function name

	if !p.expectCurrent(lexer.LPAREN) {
		return nil
	}

	// Parse parameters
	if !p.curTokenIs(lexer.RPAREN) {
		funcDecl.Parameters = p.parseParameters()
	}

	if !p.expectCurrent(lexer.RPAREN) {
		return nil
	}

	// Parse optional return type: [ "->" Type ]
	if p.curTokenIs(lexer.ARROW) {
		p.nextToken() // consume "->"
		funcDecl.ReturnType = p.parseType()
	}

	// Parse function body
	funcDecl.Body = p.parseBlockStatement()
	return funcDecl
}

func (p *parser) parseParameters() []*ast.Parameter {
	var params []*ast.Parameter

	// Parse first parameter
	param := p.parseParameter()
	if param != nil {
		params = append(params, param)
	}

	// Parse additional parameters
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // consume ","
		param := p.parseParameter()
		if param != nil {
			params = append(params, param)
		}
	}

	return params
}

func (p *parser) parseParameter() *ast.Parameter {
	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected parameter name")
		return nil
	}

	param := &ast.Parameter{
		Name: &ast.Identifier{Name: p.curToken.Literal},
	}
	p.nextToken() // consume parameter name

	if !p.expectCurrent(lexer.COLON) {
		return nil
	}

	param.Type = p.parseType()
	return param
}

// parseVariableDeclaration handles both explicit and inferred declarations
func (p *parser) parseVariableDeclaration() ast.Declaration {
	varDecl := &ast.VarDecl{}

	// Check for "mut" keyword
	if p.curTokenIs(lexer.MUT) {
		varDecl.Mutable = true
		p.nextToken() // consume "mut"
	}

	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected variable name")
		return nil
	}
	varDecl.Name = &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken() // consume variable name

	if p.curTokenIs(lexer.COLON) {
		// Explicit type: x : int = 5
		p.nextToken() // consume ":"
		varDecl.Type = p.parseType()

		if p.curTokenIs(lexer.EQ) {
			p.nextToken() // consume "="
			varDecl.Value = p.parseExpression()
		}
	} else if p.curTokenIs(lexer.COLONEQ) {
		// Type inference: x := 5
		p.nextToken() // consume ":="
		varDecl.Value = p.parseExpression()
	} else {
		p.recordError("expected ':' or ':=' in variable declaration")
		return nil
	}

	// Optional semicolon
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return varDecl
}

// parseStructDeclaration handles: "struct" IDENT "{" { FieldDecl } "}"
func (p *parser) parseStructDeclaration() ast.Declaration {
	structDecl := &ast.StructDecl{}
	p.nextToken() // consume "struct"

	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected struct name")
		return nil
	}
	structDecl.Name = &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken() // consume struct name

	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}

	// Parse fields
	for !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
		field := p.parseFieldDeclaration()
		if field != nil {
			structDecl.Fields = append(structDecl.Fields, field)
		}
	}

	if !p.expectCurrent(lexer.RBRACE) {
		return nil
	}

	return structDecl
}

func (p *parser) parseFieldDeclaration() *ast.FieldDecl {
	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected field name")
		return nil
	}

	field := &ast.FieldDecl{
		Name: &ast.Identifier{Name: p.curToken.Literal},
	}
	p.nextToken() // consume field name

	if !p.expectCurrent(lexer.COLON) {
		return nil
	}

	field.Type = p.parseType()

	// Optional semicolon
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return field
}

// parseUnsafeDeclaration handles: "unsafe" Block
func (p *parser) parseUnsafeDeclaration() ast.Declaration {
	p.nextToken() // consume "unsafe"
	block := p.parseBlockStatement()
	return &ast.UnsafeBlock{Body: block}
}

// parseAssignment handles: IDENT "=" Expression ";"
func (p *parser) parseAssignment() ast.Declaration {
	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected identifier in assignment")
		return nil
	}

	assignment := &ast.AssignmentStatement{
		Name: &ast.Identifier{Name: p.curToken.Literal},
	}
	p.nextToken() // consume identifier

	if !p.expectCurrent(lexer.EQ) {
		return nil
	}

	assignment.Value = p.parseExpression()

	// Optional semicolon
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return assignment
}

// TODO: Implement these
// Placeholder implementations
func (p *parser) parseEnumDeclaration() ast.Declaration {
	p.recordError("enum declarations not yet implemented")
	p.synchronize()
	return nil
}

func (p *parser) parseTypeDeclaration() ast.Declaration {
	p.recordError("type declarations not yet implemented")
	p.synchronize()
	return nil
}

// ===== ENHANCED EXPRESSION PARSING =====

func (p *parser) parseExpression() ast.Expression {
	return p.parseLogicalOr()
}

// Keep your existing operator precedence methods
func (p *parser) parseLogicalOr() ast.Expression {
	expr := p.parseLogicalAnd()

	for p.curTokenIs(lexer.OR) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseLogicalAnd()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *parser) parseLogicalAnd() ast.Expression {
	expr := p.parseEquality()
	for p.curTokenIs(lexer.AND) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseEquality()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *parser) parseEquality() ast.Expression {
	expr := p.parseComparison()

	for p.curTokenIs(lexer.EQEQ) || p.curTokenIs(lexer.BANGEQ) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseComparison()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *parser) parseComparison() ast.Expression {
	expr := p.parseTerm()

	for p.curTokenIs(lexer.LT) || p.curTokenIs(lexer.GT) || p.curTokenIs(lexer.LTEQ) || p.curTokenIs(lexer.GTEQ) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseTerm()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *parser) parseTerm() ast.Expression {
	expr := p.parseFactor()

	for p.curTokenIs(lexer.PLUS) || p.curTokenIs(lexer.MINUS) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseFactor()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *parser) parseFactor() ast.Expression {
	expr := p.parseUnary()

	for p.curTokenIs(lexer.ASTERISK) || p.curTokenIs(lexer.SLASH) || p.curTokenIs(lexer.PERCENT) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseUnary()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

func (p *parser) parseUnary() ast.Expression {
	if p.curTokenIs(lexer.BANG) || p.curTokenIs(lexer.MINUS) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseUnary()
		return &ast.UnaryExpression{Operator: op, Right: right}
	}
	return p.parsePrimary()
}

// Enhanced parsePrimary with struct literal support
func (p *parser) parsePrimary() ast.Expression {
	var expr ast.Expression

	switch p.curToken.Type {
	case lexer.IDENT:
		expr = p.parseIdentifierOrStructLiteral()
	case lexer.NUMBER:
		expr = p.parseNumberLiteral()
	case lexer.STRING:
		expr = p.parseStringLiteral()
	case lexer.TRUE, lexer.FALSE:
		expr = p.parseBooleanLiteral()
	case lexer.NIL:
		expr = p.parseNilLiteral()
	case lexer.LPAREN:
		p.nextToken()
		expr = p.parseExpression()
		if !p.expectCurrent(lexer.RPAREN) {
			return nil
		}
	case lexer.LBRACKET:
		expr = p.parseArrayLiteral()
	default:
		p.recordError(fmt.Sprintf("unexpected token %s in expression", p.curToken.Type))
		return nil
	}

	// Handle suffixes (function calls, indexing, member access, slicing)
	for {
		switch p.curToken.Type {
		case lexer.LPAREN:
			expr = p.parseCallExpression(expr)
		case lexer.LBRACKET:
			expr = p.parseIndexOrSliceExpression(expr)
		case lexer.DOT:
			expr = p.parseMemberExpression(expr)
		default:
			return expr
		}
	}
}

// parseIdentifierOrStructLiteral handles both identifiers and struct literals
func (p *parser) parseIdentifierOrStructLiteral() ast.Expression {
	name := p.curToken.Literal
	p.nextToken() // consume identifier

	// Check if this is a struct literal: Point{x: 1, y: 2}
	if p.curTokenIs(lexer.LBRACE) {
		return p.parseStructLiteral(name)
	}

	return &ast.Identifier{Name: name}
}

// parseStructLiteral handles: IDENT "{" [ FieldInit ( "," FieldInit )* ] "}"
func (p *parser) parseStructLiteral(typeName string) ast.Expression {
	structLit := &ast.StructLiteral{
		Type: &ast.Identifier{Name: typeName},
	}
	p.nextToken() // consume "{"

	if p.curTokenIs(lexer.RBRACE) {
		p.nextToken() // consume "}"
		return structLit
	}

	// Parse field initializers
	field := p.parseFieldInit()
	if field != nil {
		structLit.Fields = append(structLit.Fields, field)
	}

	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // consume ","
		field := p.parseFieldInit()
		if field != nil {
			structLit.Fields = append(structLit.Fields, field)
		}
	}

	if !p.expectCurrent(lexer.RBRACE) {
		return nil
	}

	return structLit
}

func (p *parser) parseFieldInit() *ast.FieldInit {
	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected field name")
		return nil
	}

	field := &ast.FieldInit{
		Name: &ast.Identifier{Name: p.curToken.Literal},
	}
	p.nextToken() // consume field name

	if !p.expectCurrent(lexer.COLON) {
		return nil
	}

	field.Value = p.parseExpression()
	return field
}

// parseIndexOrSliceExpression handles both indexing and slicing
func (p *parser) parseIndexOrSliceExpression(object ast.Expression) ast.Expression {
	p.nextToken() // consume "["

	startExpr := p.parseExpression()

	// Check if it's slicing (has colon)
	if p.curTokenIs(lexer.COLON) {
		p.nextToken() // consume ":"

		var endExpr ast.Expression
		if !p.curTokenIs(lexer.RBRACKET) {
			endExpr = p.parseExpression()
		}

		if !p.expectCurrent(lexer.RBRACKET) {
			return nil
		}

		// Create SliceExpression
		return &ast.SliceExpression{
			Object: object,
			Start:  startExpr,
			End:    endExpr,
		}
	}

	// Regular indexing
	if !p.expectCurrent(lexer.RBRACKET) {
		return nil
	}

	return &ast.IndexExpression{
		Object: object,
		Index:  startExpr,
	}
}

// Keep your existing literal parsing methods
func (p *parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken()
	return ident
}

func (p *parser) parseNumberLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.recordError("failed to parse number: " + err.Error())
		p.synchronize()
		return nil
	}
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: val,
	}
	p.nextToken()
	return lit
}

func (p *parser) parseStringLiteral() ast.Expression {
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: p.curToken.Literal,
	}
	p.nextToken()
	return lit
}

func (p *parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: p.curToken.Type == lexer.TRUE,
	}
	p.nextToken()
	return lit
}

func (p *parser) parseNilLiteral() ast.Expression {
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: nil,
	}
	p.nextToken()
	return lit
}

func (p *parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{}
	p.nextToken() // consume [

	if p.curTokenIs(lexer.RBRACKET) {
		p.nextToken() // consume ]
		return array
	}

	array.Elements = append(array.Elements, p.parseExpression())
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		array.Elements = append(array.Elements, p.parseExpression())
	}

	if !p.expectCurrent(lexer.RBRACKET) {
		return nil
	}

	return array
}

func (p *parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.FunctionCall{Function: function}
	p.nextToken() // consume (

	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken() // consume )
		return call
	}

	call.Arguments = append(call.Arguments, p.parseExpression())

	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		call.Arguments = append(call.Arguments, p.parseExpression())
	}

	if !p.expectCurrent(lexer.RPAREN) {
		return nil
	}

	return call
}

func (p *parser) parseIndexExpression(array ast.Expression) ast.Expression {
	index := &ast.IndexExpression{Object: array}
	p.nextToken() // consume [

	index.Index = p.parseExpression()

	if !p.expectCurrent(lexer.RBRACKET) {
		return nil
	}

	return index
}

func (p *parser) parseMemberExpression(object ast.Expression) ast.Expression {
	p.nextToken() // consume .

	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected identifier after '.'")
		return nil
	}

	// Safe parsing without dangerous type assertion
	property := &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken() // consume identifier

	return &ast.MemberExpression{
		Object:   object,
		Property: property,
	}
}

// ===== STATEMENT PARSING =====
func (p *parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.LOG:
		return p.parsePrintStatement()
	case lexer.BREAK:
		return p.parseBreakStatement()
	case lexer.CONTINUE:
		return p.parseContinueStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}

	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}

	for !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	if !p.expectCurrent(lexer.RBRACE) {
		return nil
	}

	return block
}

func (p *parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{}
	p.nextToken() // consume 'if'

	stmt.Condition = p.parseExpression()
	stmt.Consequence = p.parseBlockStatement()

	// Handle optional else clause
	if p.curTokenIs(lexer.ELSE) {
		p.nextToken() // consume 'else'
		if p.curTokenIs(lexer.IF) {
			// else if - parse as another if statement
			elseIf := p.parseIfStatement()
			stmt.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{elseIf},
			}
		} else {
			// else block
			stmt.Alternative = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *parser) parseForStatement() ast.Statement {
	stmt := &ast.ForStatement{}
	p.nextToken() // consume 'for'

	// Parse optional init; condition; post parts
	if !p.curTokenIs(lexer.LBRACE) {
		// Parse init (optional)
		if !p.curTokenIs(lexer.SEMICOLON) {
			stmt.Init = p.parseStatement()
		}

		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken() // consume ;

			// Parse condition (optional)
			if !p.curTokenIs(lexer.SEMICOLON) {
				stmt.Condition = p.parseExpression()
			}

			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken() // consume ;

				// Parse post (optional)
				if !p.curTokenIs(lexer.LBRACE) {
					stmt.Post = p.parseStatement()
				}
			}
		}
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{}
	p.nextToken() // consume 'return'

	if !p.curTokenIs(lexer.SEMICOLON) && !p.curTokenIs(lexer.RBRACE) {
		stmt.Value = p.parseExpression()
	}

	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *parser) parsePrintStatement() ast.Statement {
	stmt := &ast.PrintStatement{}
	p.nextToken() // consume 'log'

	if !p.expectCurrent(lexer.LPAREN) {
		return nil
	}

	stmt.Expression = p.parseExpression()

	if !p.expectCurrent(lexer.RPAREN) {
		return nil
	}

	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *parser) parseBreakStatement() ast.Statement {
	p.nextToken() // consume 'break'
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	return &ast.BreakStatement{}
}

func (p *parser) parseContinueStatement() ast.Statement {
	p.nextToken() // consume 'continue'
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	return &ast.ContinueStatement{}
}

func (p *parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Expression: p.parseExpression(),
	}

	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ===== UTILITY METHODS =====
func (p *parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *parser) isAtEnd() bool {
	return p.curToken.Type == lexer.EOF
}

func (p *parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.recordError(fmt.Sprintf("expected next token to be %s, got %s", t.String(), p.peekToken.Type.String()))
	p.synchronize()
	return false
}

func (p *parser) expectCurrent(t lexer.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.recordError(fmt.Sprintf("expected current token to be %s, got %s", t.String(), p.curToken.Type.String()))
	p.synchronize()
	return false
}

func (p *parser) synchronize() {
	p.nextToken()

	for !p.isAtEnd() {
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
			return
		}

		switch p.curToken.Type {
		case lexer.FUNC, lexer.MUT, lexer.STRUCT, lexer.ENUM, lexer.TYPE,
			lexer.RBRACE, lexer.RBRACKET, lexer.FOR, lexer.IF,
			lexer.RETURN, lexer.UNSAFE:
			return
		}

		p.nextToken()
	}
}
