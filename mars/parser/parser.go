// parser/parser.go
package parser

import (
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

// ParseProgram is the main entry point for parsing.
// It creates the root Program node and parses all statements until EOF.
func (p *parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Declarations = append(program.Declarations, stmt)
		}
		p.nextToken() // Move to next token after statement
	}

	return program
}

// parseExpression is the single entry‐point for any expression,
// matching the "Expression = Assignment" production.
func (p *parser) parseExpression() ast.Expression {
	return p.parseLogicalOr()
}

// Handles "a || b || c"
func (p *parser) parseLogicalOr() ast.Expression {
	expr := p.parseLogicalAnd()

	// we loop until we find a token that is not OR
	for p.curTokenIs(lexer.OR) {
		op := p.curToken.Literal
		p.nextToken() // consume ||
		right := p.parseLogicalAnd()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// Handles "a && b && c"
func (p *parser) parseLogicalAnd() ast.Expression {
	expr := p.parseEquality()
	// we loop until we find a token that is not AND
	for p.curTokenIs(lexer.AND) {
		op := p.curToken.Literal
		p.nextToken() // consume &&
		right := p.parseEquality()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// Handles == and !=
func (p *parser) parseEquality() ast.Expression {
	expr := p.parseComparison()

	// we loop until we find a token that is not EQ or NEQ
	for p.curTokenIs(lexer.EQEQ) || p.curTokenIs(lexer.BANGEQ) {
		op := p.curToken.Literal
		p.nextToken() // consume EQ or NEQ
		right := p.parseComparison()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// Handles <, >, <=, >=
func (p *parser) parseComparison() ast.Expression {
	expr := p.parseTerm()

	// we loop until we find a token that is not LT, GT, LE, or GE
	for p.curTokenIs(lexer.LT) || p.curTokenIs(lexer.GT) || p.curTokenIs(lexer.LTEQ) || p.curTokenIs(lexer.GTEQ) {
		op := p.curToken.Literal
		p.nextToken() // consume LT, GT, LE, or GE
		right := p.parseTerm()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// Handles + and -
func (p *parser) parseTerm() ast.Expression {
	expr := p.parseFactor()

	// we loop until we find plus or minus
	for p.curTokenIs(lexer.PLUS) || p.curTokenIs(lexer.MINUS) {
		op := p.curToken.Literal
		p.nextToken() // consume plus or minus
		right := p.parseFactor()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// Handles *, /, %
func (p *parser) parseFactor() ast.Expression {
	expr := p.parseUnary()

	// we loop until we find times or divide
	for p.curTokenIs(lexer.ASTERISK) || p.curTokenIs(lexer.SLASH) || p.curTokenIs(lexer.PERCENT) {
		op := p.curToken.Literal
		p.nextToken() // consume times or divide
		right := p.parseUnary()
		expr = &ast.BinaryExpression{Left: expr, Operator: op, Right: right}
	}
	return expr
}

// Handles unary ! and -
func (p *parser) parseUnary() ast.Expression {
	if p.curTokenIs(lexer.BANG) || p.curTokenIs(lexer.MINUS) {
		op := p.curToken.Literal
		p.nextToken()
		right := p.parseUnary()
		return &ast.UnaryExpression{Operator: op, Right: right}
	}
	return p.parsePrimary()
}

// parseIdentifier parses an identifier expression
func (p *parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken() // consume the identifier
	return ident
}

// parseNumberLiteral parses a number literal
func (p *parser) parseNumberLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.recordError("failed to parse number: " + err.Error())
		return nil
	}
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: val,
	}
	p.nextToken() // consume the number
	return lit
}

// parseStringLiteral parses a string literal
func (p *parser) parseStringLiteral() ast.Expression {
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: p.curToken.Literal,
	}
	p.nextToken() // consume the string
	return lit
}

// parseBooleanLiteral parses a boolean literal
func (p *parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: p.curToken.Type == lexer.TRUE,
	}
	p.nextToken() // consume the boolean
	return lit
}

// parseNilLiteral parses a nil literal
func (p *parser) parseNilLiteral() ast.Expression {
	lit := &ast.Literal{
		Token: p.curToken.Literal,
		Value: nil,
	}
	p.nextToken() // consume the nil
	return lit
}

// parseArrayLiteral parses an array literal
func (p *parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{}
	p.nextToken() // consume [

	// Handle empty array
	if p.curTokenIs(lexer.RBRACKET) {
		p.nextToken() // consume ]
		return array
	}

	// Parse array elements
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

// parseCallExpression parses a function call expression
func (p *parser) parseCallExpression(function ast.Expression) ast.Expression {
	// 1. Create the AST node for the entire call, storing the 'function' part.
	call := &ast.FunctionCall{
		Function: function,
	}
	p.nextToken() // 2. Consume the opening parenthesis '('.

	// 3. Handle the case of a call with no arguments, like myFunction().
	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken() // Consume the closing parenthesis ')'.
		return call   // The call expression is complete.
	}

	// 4. If there are arguments, parse the first one.
	//    It calls p.parseExpression(), so an argument can be any valid expression.
	call.Arguments = append(call.Arguments, p.parseExpression())

	// 5. Loop as long as you see commas, parsing each subsequent argument.
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken() // Consume the comma.
		call.Arguments = append(call.Arguments, p.parseExpression())
	}

	// 6. Ensure the argument list is correctly closed with a ')'.
	if !p.expectCurrent(lexer.RPAREN) {
		return nil // Error
	}

	// 7. Return the completed FunctionCall AST node.
	return call
}

// parseIndexExpression parses an array indexing expression
func (p *parser) parseIndexExpression(array ast.Expression) ast.Expression {
	// 1. Create the AST node for the index operation, storing the 'array' part.
	index := &ast.IndexExpression{
		Object: array,
	}
	p.nextToken() // 2. Consume the opening bracket '['.

	// 3. Parse the expression *inside* the brackets. Because this calls
	//    p.parseExpression(), the index can be a complex expression itself.
	index.Index = p.parseExpression()

	// 4. Ensure the index is correctly closed with a ']'.
	if !p.expectCurrent(lexer.RBRACKET) {
		return nil // Error
	}

	// 5. Return the completed IndexExpression AST node.
	return index
}

// parseMemberExpression parses a member access expression
func (p *parser) parseMemberExpression(object ast.Expression) ast.Expression {
	p.nextToken() // 1. Consume the dot '.'.

	// 2. Check that what follows the dot is a valid identifier.
	if !p.curTokenIs(lexer.IDENT) {
		p.recordError("expected identifier after '.'")
		return nil // Error
	}

	// 3. Create the AST node for the member access expression.
	return &ast.MemberExpression{
		Object: object, // The expression on the left of the dot.
		// Parse the identifier on the right and assert it to the correct AST type.
		Property: p.parseIdentifier().(*ast.Identifier),
	}
}

// parsePrimary parses primary expressions and their suffixes
func (p *parser) parsePrimary() ast.Expression {
	var expr ast.Expression

	switch p.curToken.Type {
	case lexer.IDENT:
		expr = p.parseIdentifier()
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
		return nil // or an error node
	}

	// Handle suffixes (function calls, indexing, member access)
	for {
		switch p.curToken.Type {
		case lexer.LPAREN:
			expr = p.parseCallExpression(expr)
		case lexer.LBRACKET:
			expr = p.parseIndexExpression(expr)
		case lexer.DOT:
			expr = p.parseMemberExpression(expr)
		default:
			return expr
		}
	}
}

func (p *parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *parser) isAtEnd() bool {
	return p.peekToken.Type == lexer.EOF
}

// expectPeek checks if the next token is of the expected type
func (p *parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.recordError("expected next token to be " + string(t) + ", got " + string(p.peekToken.Type))
	p.synchronize() // Start fast-forwarding!
	return false
}

// Add this new helper function to your parser
func (p *parser) expectCurrent(t lexer.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken() // It's what we expect, so we consume it and move on
		return true
	}
	// If it's not what we expect, record an error
	p.recordError("expected current token to be " + string(t) + ", got " + string(p.curToken.Type))
	p.synchronize() // Start fast-forwarding!
	return false
}

// synchronize implements panic-mode error recovery:
// it skips tokens until it finds a synchronization point to resume parsing
func (p *parser) synchronize() {
	p.nextToken() // Advance past the erroneous token

	for !p.isAtEnd() {
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken() // Consume the semicolon
			return
		}

		switch p.curToken.Type {
		case lexer.FUNC, lexer.MUT, lexer.RBRACE, lexer.RBRACKET,
			lexer.FOR, lexer.IF, lexer.RETURN, lexer.UNSAFE:
			return // These tokens likely start a new statement or declaration
		}

		p.nextToken()
	}
}

// parseStatement parses a single statement
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
	default:
		return p.parseExpressionStatement()
	}
}

// parseIfStatement parses an if statement
func (p *parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{}
	p.nextToken() // consume 'if'

	stmt.Condition = p.parseExpression()
	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}

	stmt.Consequence = &ast.BlockStatement{}
	for !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
		stmt.Consequence.Statements = append(stmt.Consequence.Statements, p.parseStatement())
		p.nextToken()
	}

	return stmt
}

// parseForStatement parses a for statement
func (p *parser) parseForStatement() ast.Statement {
	stmt := &ast.ForStatement{}
	p.nextToken() // consume 'for'

	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}

	stmt.Body = &ast.BlockStatement{}
	for !p.curTokenIs(lexer.RBRACE) && !p.isAtEnd() {
		stmt.Body.Statements = append(stmt.Body.Statements, p.parseStatement())
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{}
	p.nextToken() // consume 'return'

	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Value = p.parseExpression()
	}

	return stmt
}

// parsePrintStatement parses a print/log statement
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

	return stmt
}

// parseExpressionStatement parses an expression statement
func (p *parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Expression: p.parseExpression(),
	}
	return stmt
}
