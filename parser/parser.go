// Package parser parser/parser.go
package parser

import (
	"fmt"
	"mars/ast"
	"mars/errors"
	"mars/lexer"
	"strconv"
)

type parser struct {
	lexer     *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	prevToken lexer.Token
	errors    *errors.ErrorList
}

func NewParser(lexer *lexer.Lexer) *parser {
	p := &parser{lexer: lexer, errors: errors.NewErrorList()}
	p.nextToken()
	p.nextToken()
	return p
}

// Helper function to convert token position to AST position
func tokenToPosition(token lexer.Token) ast.Position {
	return ast.Position{
		Line:   token.Line,
		Column: token.Column,
	}
}

// Helper function to get current token position
func (p *parser) currentPosition() ast.Position {
	return tokenToPosition(p.curToken)
}

// Helper function to get previous token position
func (p *parser) previousPosition() ast.Position {
	return tokenToPosition(p.prevToken)
}

func (p *parser) recordError(message string) {
	p.errors.AddError(message, p.curToken.Line, p.curToken.Column)
}

func (p *parser) recordSyntaxError(message string) {
	p.errors.Add(errors.NewSyntaxError(message, p.curToken.Line, p.curToken.Column))
}

func (p *parser) nextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *parser) previousToken() lexer.Token {
	return p.prevToken
}

func (p *parser) GetErrors() *errors.ErrorList {
	return p.errors
}

// ParseProgram is the main entry point for parsing.
func (p *parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Declarations = []ast.Declaration{}
	program.Position = p.currentPosition()

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
		p.recordSyntaxError(fmt.Sprintf("expected type, got %s", p.curToken.Type))
		return nil
	}
}

func (p *parser) parseBaseType() *ast.Type {
	baseType := &ast.Type{
		BaseType: p.curToken.Literal,
		Position: p.currentPosition(),
	}
	p.nextToken() // consume the type token
	return baseType
}

func (p *parser) parseArrayOrSliceType() *ast.Type {
	startPos := p.currentPosition()
	p.nextToken() // consume '['

	arrayType := &ast.Type{
		Position: startPos,
	}

	// Check if it's a sized array [N] or dynamic slice []
	if p.curTokenIs(lexer.NUMBER) {
		// Parse array size
		if size, err := strconv.Atoi(p.curToken.Literal); err == nil {
			arrayType.ArraySize = &size
		} else {
			p.recordSyntaxError("invalid array size")
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
	startPos := p.currentPosition()
	p.nextToken() // consume '*'
	pointeeType := p.parseType()
	if pointeeType == nil {
		return nil
	}
	return &ast.Type{
		PointerType: pointeeType,
		Position:    startPos,
	}
}

func (p *parser) parseStructTypeReference() *ast.Type {
	structType := &ast.Type{
		StructName: p.curToken.Literal,
		Position:   p.currentPosition(),
	}
	p.nextToken() // consume struct name
	return structType
}

// ===== DECLARATION PARSING =====

func (p *parser) parseDeclaration() ast.Declaration {
	switch p.curToken.Type {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!
	case lexer.COMMENT:
		// Skip comments by consuming the token and returning nil
		p.nextToken()
		return nil
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
	case lexer.LBRACE:
		return p.parseBlockStatement()
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
		p.recordSyntaxError(fmt.Sprintf("unexpected token %s at top level", p.curToken.Type))
		p.synchronize()
		return nil
	}
}

// parseFunctionDeclaration handles: "func" IDENT "(" [ Params ] ")" [ "->" Type ] Block
func (p *parser) parseFunctionDeclaration() ast.Declaration {
	startPos := p.currentPosition()
	funcDecl := &ast.FuncDecl{
		Position: startPos,
	}
	p.nextToken() // consume "func"

	if !p.curTokenIs(lexer.IDENT) {
		p.recordSyntaxError("expected function name")
		return nil
	}
	funcDecl.Name = &ast.Identifier{
		Name:     p.curToken.Literal,
		Position: p.currentPosition(),
	}
	p.nextToken() // consume function name

	if !p.expectCurrent(lexer.LPAREN) {
		return nil
	}

	// Create the signature node right here
	signature := &ast.FunctionSignature{
		Position: startPos,
	}

	// Parse parameters INTO the signature node
	if !p.curTokenIs(lexer.RPAREN) {
		signature.Parameters = p.parseParameters()
	}

	if !p.expectCurrent(lexer.RPAREN) {
		return nil
	}

	// Parse optional return type INTO the signature node
	if p.curTokenIs(lexer.ARROW) {
		p.nextToken() // consume "->"
		signature.ReturnType = p.parseType()
	}

	// Attach the completed signature to the function declaration
	funcDecl.Signature = signature

	// Parse function body
	funcDecl.Body = p.parseBlockStatement()

	// Optional semicolon
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

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
		p.recordSyntaxError("expected parameter name")
		return nil
	}

	param := &ast.Parameter{
		Name: &ast.Identifier{
			Name:     p.curToken.Literal,
			Position: p.currentPosition(),
		},
		Position: p.currentPosition(),
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
	startPos := p.currentPosition()
	varDecl := &ast.VarDecl{
		Position: startPos,
	}

	// Check for "mut" keyword
	if p.curTokenIs(lexer.MUT) {
		varDecl.Mutable = true
		p.nextToken() // consume "mut"
	}

	if !p.curTokenIs(lexer.IDENT) {
		p.recordSyntaxError("expected variable name")
		return nil
	}
	varDecl.Name = &ast.Identifier{
		Name:     p.curToken.Literal,
		Position: p.currentPosition(),
	}
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
		p.recordSyntaxError("expected ':' or ':=' in variable declaration")
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
	startPos := p.currentPosition()
	structDecl := &ast.StructDecl{
		Position: startPos,
	}
	p.nextToken() // consume "struct"

	if !p.curTokenIs(lexer.IDENT) {
		p.recordSyntaxError("expected struct name")
		return nil
	}
	structDecl.Name = &ast.Identifier{
		Name:     p.curToken.Literal,
		Position: p.currentPosition(),
	}
	p.nextToken() // consume struct name

	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}

	// Parse fields
	for !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
		field := p.parseFieldDeclaration()
		if field != nil {
			structDecl.Fields = append(structDecl.Fields, field)
		} else {
			// If field parsing failed, try to recover by skipping to next field or end
			if !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
				p.synchronize()
				// If we're still not at the end or closing brace, skip this field and continue
				if !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
					continue
				}
			}
			break
		}
	}

	if !p.expectCurrent(lexer.RBRACE) {
		return nil
	}

	return structDecl
}

func (p *parser) parseFieldDeclaration() *ast.FieldDecl {
	if !p.curTokenIs(lexer.IDENT) {
		p.recordSyntaxError("expected field name")
		return nil
	}

	field := &ast.FieldDecl{
		Name: &ast.Identifier{
			Name:     p.curToken.Literal,
			Position: p.currentPosition(),
		},
		Position: p.currentPosition(),
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
	startPos := p.currentPosition()
	p.nextToken() // consume "unsafe"
	block := p.parseBlockStatement()
	return &ast.UnsafeBlock{
		Body:     block,
		Position: startPos,
	}
}

// parseAssignment handles: IDENT "=" Expression ";" or IDENT "[" Expression "]" "=" Expression ";"
func (p *parser) parseAssignment() ast.Declaration {
	if !p.curTokenIs(lexer.IDENT) {
		p.recordSyntaxError("expected identifier in assignment")
		return nil
	}

	// Parse the left-hand side (could be identifier or indexed expression)
	startPos := p.currentPosition()
	object := p.parseIdentifierOrStructLiteral()

	// Check if this is an array assignment (object[index] = value)
	if p.curTokenIs(lexer.LBRACKET) {
		// This is an array assignment: arr[index] = value
		p.nextToken() // consume "["

		index := p.parseExpression()

		if !p.expectCurrent(lexer.RBRACKET) {
			return nil
		}

		if !p.expectCurrent(lexer.EQ) {
			p.recordSyntaxError("expected '=' after array index")
			return nil
		}
		//!!!!!!!!!!!!!!!!!!!!!!!!!!
		value := p.parseExpression()

		// Optional semicolon
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}

		return &ast.IndexAssignmentStatement{
			Object:   object,
			Index:    index,
			Value:    value,
			Position: startPos,
		}
	}

	// This is a regular assignment: identifier = value
	if !p.expectCurrent(lexer.EQ) {
		return nil
	}

	value := p.parseExpression()

	// Optional semicolon
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	//!!!!!!!!!!!!!!!!!!!!!!!!!!

	// Convert object back to identifier for regular assignment
	if ident, ok := object.(*ast.Identifier); ok {
		return &ast.AssignmentStatement{
			Name:     ident,
			Value:    value,
			Position: startPos,
		}
	}

	p.recordSyntaxError("expected identifier for assignment")
	return nil
}

// TODO: Implement these
// Placeholder implementations
func (p *parser) parseEnumDeclaration() ast.Declaration {
	p.recordSyntaxError("enum declarations not yet implemented")
	p.synchronize()
	return nil
}

func (p *parser) parseTypeDeclaration() ast.Declaration {
	p.recordSyntaxError("type declarations not yet implemented")
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
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseLogicalAnd()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
			Position: pos,
		}
	}
	return expr
}

func (p *parser) parseLogicalAnd() ast.Expression {
	expr := p.parseEquality()
	for p.curTokenIs(lexer.AND) {
		op := p.curToken.Literal
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseEquality()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
			Position: pos,
		}
	}
	return expr
}

func (p *parser) parseEquality() ast.Expression {
	expr := p.parseComparison()

	for p.curTokenIs(lexer.EQEQ) || p.curTokenIs(lexer.BANGEQ) {
		op := p.curToken.Literal
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseComparison()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
			Position: pos,
		}
	}
	return expr
}

func (p *parser) parseComparison() ast.Expression {
	expr := p.parseTerm()

	for p.curTokenIs(lexer.LT) || p.curTokenIs(lexer.GT) || p.curTokenIs(lexer.LTEQ) || p.curTokenIs(lexer.GTEQ) {
		op := p.curToken.Literal
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseTerm()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
			Position: pos,
		}
	}
	return expr
}

func (p *parser) parseTerm() ast.Expression {
	expr := p.parseFactor()

	for p.curTokenIs(lexer.PLUS) || p.curTokenIs(lexer.MINUS) {
		op := p.curToken.Literal
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseFactor()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
			Position: pos,
		}
	}
	return expr
}

func (p *parser) parseFactor() ast.Expression {
	expr := p.parseUnary()

	for p.curTokenIs(lexer.ASTERISK) || p.curTokenIs(lexer.SLASH) || p.curTokenIs(lexer.PERCENT) {
		op := p.curToken.Literal
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseUnary()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: op,
			Right:    right,
			Position: pos,
		}
	}
	return expr
}

func (p *parser) parseUnary() ast.Expression {
	if p.curTokenIs(lexer.BANG) || p.curTokenIs(lexer.MINUS) {
		op := p.curToken.Literal
		pos := p.currentPosition()
		p.nextToken()
		right := p.parseUnary()
		return &ast.UnaryExpression{
			Operator: op,
			Right:    right,
			Position: pos,
		}
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
		p.synchronize()
		p.recordSyntaxError(fmt.Sprintf("unexpected token %s in expression", p.curToken.Type))
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
	pos := p.currentPosition()
	p.nextToken() // consume identifier

	// Check if this is a struct literal: Point{x: 1, y: 2}
	if p.curTokenIs(lexer.LBRACE) {
		return p.parseStructLiteral(name)
	}

	return &ast.Identifier{
		Name:     name,
		Position: pos,
	}
}

// parseStructLiteral handles: IDENT "{" [ FieldInit ( "," FieldInit )* ] "}"
func (p *parser) parseStructLiteral(typeName string) ast.Expression {
	startPos := p.previousPosition() // Position of the type name
	structLit := &ast.StructLiteral{
		Type: &ast.Identifier{
			Name:     typeName,
			Position: startPos,
		},
		Position: startPos,
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
		p.recordSyntaxError("expected field name")
		return nil
	}

	field := &ast.FieldInit{
		Name: &ast.Identifier{
			Name:     p.curToken.Literal,
			Position: p.currentPosition(),
		},
		Position: p.currentPosition(),
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
	startPos := p.currentPosition()
	p.nextToken() // consume "["

	//!!!!!!!!!!!!!!!!!!!!!!!!!!
	// Check if it's slicing (has colon immediately or after start expression)
	if p.curTokenIs(lexer.COLON) {
		// Handle [:end] case
		p.nextToken() // consume ":"

		var endExpr ast.Expression
		if !p.curTokenIs(lexer.RBRACKET) {
			endExpr = p.parseExpression()
		}

		if !p.expectCurrent(lexer.RBRACKET) {
			return nil
		}

		// Create SliceExpression with nil start
		return &ast.SliceExpression{
			Object:   object,
			Start:    nil,
			End:      endExpr,
			Position: startPos,
		}
	}

	// Parse start expression
	startExpr := p.parseExpression()

	// Check if it's slicing (has colon after start expression)
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
			Object:   object,
			Start:    startExpr,
			End:      endExpr,
			Position: startPos,
		}
	}

	// Regular indexing
	if !p.expectCurrent(lexer.RBRACKET) {
		return nil
	}

	return &ast.IndexExpression{
		Object:   object,
		Index:    startExpr,
		Position: startPos,
	}
}

// Keep your existing literal parsing methods
func (p *parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{
		Name:     p.curToken.Literal,
		Position: p.currentPosition(),
	}
	p.nextToken()
	return ident
}

func (p *parser) parseNumberLiteral() ast.Expression {
	// First try to parse as an integer
	if intVal, err := strconv.ParseInt(p.curToken.Literal, 10, 64); err == nil {
		lit := &ast.Literal{
			Token:    p.curToken.Literal,
			Value:    int(intVal), // Convert to int to match the type checker's expectations
			Position: p.currentPosition(),
		}
		p.nextToken()
		return lit
	}

	// If integer parsing fails, try as float
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.recordSyntaxError("failed to parse number: " + err.Error())
		p.synchronize()
		return nil
	}
	lit := &ast.Literal{
		Token:    p.curToken.Literal,
		Value:    val,
		Position: p.currentPosition(),
	}
	p.nextToken()
	return lit
}

func (p *parser) parseStringLiteral() ast.Expression {
	lit := &ast.Literal{
		Token:    p.curToken.Literal,
		Value:    p.curToken.Literal,
		Position: p.currentPosition(),
	}
	p.nextToken()
	return lit
}

func (p *parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.Literal{
		Token:    p.curToken.Literal,
		Value:    p.curToken.Type == lexer.TRUE,
		Position: p.currentPosition(),
	}
	p.nextToken()
	return lit
}

func (p *parser) parseNilLiteral() ast.Expression {
	lit := &ast.Literal{
		Token:    p.curToken.Literal,
		Value:    nil,
		Position: p.currentPosition(),
	}
	p.nextToken()
	return lit
}

func (p *parser) parseArrayLiteral() ast.Expression {
	startPos := p.currentPosition()
	array := &ast.ArrayLiteral{
		Position: startPos,
	}
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
	startPos := p.currentPosition()
	call := &ast.FunctionCall{
		Function: function,
		Position: startPos,
	}
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
	startPos := p.currentPosition()
	index := &ast.IndexExpression{
		Object:   array,
		Position: startPos,
	}
	p.nextToken() // consume [

	index.Index = p.parseExpression()

	if !p.expectCurrent(lexer.RBRACKET) {
		return nil
	}

	return index
}

func (p *parser) parseMemberExpression(object ast.Expression) ast.Expression {
	startPos := p.currentPosition()
	p.nextToken() // consume .

	if !p.curTokenIs(lexer.IDENT) {
		p.recordSyntaxError("expected identifier after '.'")
		return nil
	}

	// Safe parsing without dangerous type assertion
	property := &ast.Identifier{
		Name:     p.curToken.Literal,
		Position: p.currentPosition(),
	}
	p.nextToken() // consume identifier

	return &ast.MemberExpression{
		Object:   object,
		Property: property,
		Position: startPos,
	}
}

// ===== STATEMENT PARSING =====
func (p *parser) parseStatement() ast.Statement {
	// Guard: if we're at the end of a block or file, do not parse a statement
	if p.curTokenIs(lexer.RBRACE) || p.curTokenIs(lexer.EOF) {
		return nil
	}

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
	startPos := p.currentPosition()
	block := &ast.BlockStatement{
		Position: startPos,
	}

	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}

	for !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
		// Guard: if we're at RBRACE or EOF, break
		if p.curTokenIs(lexer.RBRACE) || p.curTokenIs(lexer.EOF) {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		// Ensure we advance past the statement if we're not at the end
		if !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
			p.nextToken()
		}
	}

	if !p.expectCurrent(lexer.RBRACE) {
		return nil
	}

	return block
}

func (p *parser) parseIfStatement() ast.Statement {
	startPos := p.currentPosition()
	stmt := &ast.IfStatement{
		Position: startPos,
	}
	p.nextToken() // consume 'if'

	stmt.Condition = p.parseExpression()

	// Ensure we're on LBRACE before parsing the block
	if !p.curTokenIs(lexer.LBRACE) {
		p.recordSyntaxError("expected '{' after if condition")
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()

	// Handle optional else clause
	if p.curTokenIs(lexer.ELSE) {
		p.nextToken() // consume 'else'
		if p.curTokenIs(lexer.IF) {
			// else if - parse as another if statement
			elseIf := p.parseIfStatement()
			stmt.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{elseIf},
				Position:   p.currentPosition(),
			}
		} else {
			// else block
			stmt.Alternative = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *parser) parseForStatement() ast.Statement {
	startPos := p.currentPosition()
	stmt := &ast.ForStatement{
		Position: startPos,
	}
	p.nextToken() // consume 'for'

	// Parse init (optional)
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken() // consume first semicolon
	} else if p.curTokenIs(lexer.LBRACE) {
		// No init, no condition, no post
	} else {
		stmt.Init = p.parseStatement()
		if !p.curTokenIs(lexer.SEMICOLON) {
			p.recordSyntaxError("expected ';' after for loop init")
			return nil
		}
		p.nextToken() // consume first semicolon
	}

	// Parse condition (optional)
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken() // consume second semicolon
	} else if p.curTokenIs(lexer.LBRACE) {
		// No condition, no post
	} else {
		stmt.Condition = p.parseExpression()
		if !p.curTokenIs(lexer.SEMICOLON) {
			p.recordSyntaxError("expected ';' after for loop condition")
			return nil
		}
		p.nextToken() // consume second semicolon
	}

	// Parse post (optional)
	if p.curTokenIs(lexer.LBRACE) {
		// No post
	} else if p.curTokenIs(lexer.SEMICOLON) {
		// Extra semicolon, skip
		p.nextToken()
	} else {
		stmt.Post = p.parseStatement()
	}

	// Ensure we're on LBRACE before parsing the block
	if !p.curTokenIs(lexer.LBRACE) {
		p.recordSyntaxError("expected '{' after for loop header")
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *parser) parseReturnStatement() ast.Statement {
	startPos := p.currentPosition()
	stmt := &ast.ReturnStatement{
		Position: startPos,
	}
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
	startPos := p.currentPosition()
	stmt := &ast.PrintStatement{
		Position: startPos,
	}
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
	startPos := p.currentPosition()
	p.nextToken() // consume 'break'
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	return &ast.BreakStatement{
		Position: startPos,
	}
}

func (p *parser) parseContinueStatement() ast.Statement {
	startPos := p.currentPosition()
	p.nextToken() // consume 'continue'
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	return &ast.ContinueStatement{
		Position: startPos,
	}
}

func (p *parser) parseExpressionStatement() ast.Statement {
	startPos := p.currentPosition()
	//!!!!!!!!!!!!!!!!!!!!!!!!!!
	// Parse the left-hand side of what might be an assignment
	leftExpr := p.parseExpression()

	// Check if this is an assignment (leftExpr = rightExpr)
	if p.curTokenIs(lexer.EQ) {
		p.nextToken() // consume "="
		rightExpr := p.parseExpression()

		// Optional semicolon
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}

		// Check if leftExpr is an IndexExpression (array assignment)
		if indexExpr, ok := leftExpr.(*ast.IndexExpression); ok {
			return &ast.IndexAssignmentStatement{
				Object:   indexExpr.Object,
				Index:    indexExpr.Index,
				Value:    rightExpr,
				Position: startPos,
			}
		}

		// Check if leftExpr is an Identifier (regular assignment)
		if ident, ok := leftExpr.(*ast.Identifier); ok {
			return &ast.AssignmentStatement{
				Name:     ident,
				Value:    rightExpr,
				Position: startPos,
			}
		}

		// For other types, treat as regular expression statement
		// (the assignment will be part of the expression)
	}

	// Regular expression statement
	stmt := &ast.ExpressionStatement{
		Expression: leftExpr,
		Position:   startPos,
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
	p.recordSyntaxError(fmt.Sprintf("expected next token to be %s, got %s", t.String(), p.peekToken.Type.String()))
	p.synchronize()
	return false
}

func (p *parser) expectCurrent(t lexer.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.recordSyntaxError(fmt.Sprintf("expected current token to be %s, got %s", t.String(), p.curToken.Type.String()))
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
