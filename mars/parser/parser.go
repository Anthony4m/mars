// parser/parser.go
package parser

import (
	"fmt"
	"mars/ast"
	"mars/lexer"
)

type parser struct {
	lexer     *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
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
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != lexer.EOF {
		decl := p.parseDeclaration()
		if decl != nil {
			program.Declarations = append(program.Declarations, decl)
		}
	}
	return program
}

func (p *parser) parseDeclaration() ast.Declaration {
	switch p.curToken.Type {
	case lexer.MUT:
		return p.parseVarDecl()
	case lexer.FUNC:
		return p.parseFuncDecl()
	case lexer.STRUCT:
		return p.parseStructDecl()
	case lexer.UNSAFE:
		return p.parseUnsafeBlock()
	default:
		return p.parseStatement()
	}
}

// mut var_name := expression || var_name := expression
func (p *parser) parseVarDecl() *ast.VarDecl {
	isMutable := false
	if p.curToken.Type == lexer.MUT {
		isMutable = true
		p.nextToken() // Skip 'mut'
	}

	name := p.parseIdentifier()
	p.expect(lexer.COLONEQ) // Expect :=

	varDecl := &ast.VarDecl{
		Mutable: isMutable,
		Name:    name,
		Value:   p.parseExpression(),
	}

	return varDecl
}

func (p *parser) parseIdentifier() *ast.Identifier {
	ident := &ast.Identifier{Name: p.curToken.Literal}
	p.nextToken() // Skip identifier
	return ident
}

func (p *parser) expect(t lexer.TokenType) {
	if p.curToken.Type != t {
		msg := fmt.Sprintf("Expected token type %v but got %v at line %d, column %d",
			t, p.curToken.Type, p.curToken.Line, p.curToken.Column)
		p.recordError(msg)
		return // Return early since we encountered an error
	}
	p.nextToken()
}

func (p *parser) currentTokenIs(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.curToken.Type == t {
			return true
		}
	}
	return false
}

// Example: Parse additive expressions (+, -)
// The main trick here is that functions handling lower precedence operators (like +, -)
// will call functions that handle higher precedence operators (like *, /)
// to get their operands. This naturally groups operations according to their precedence.
func (p *parser) parseExpression() ast.Expression {
	// Logic: In this version, it simply delegates to parseAdditive().
	// This means that the lowest level of precedence it directly considers is addition/subtraction.
	// If you had even lower precedence operators (e.g., logical OR, assignments if they were expressions),
	// parseExpression might start there, or parseAdditive would be the entry if +/- were the lowest.
	return p.parseAdditive()
}

// Logic Breakdown:
//
//	Rule 1: expr := p.parseMultiplicative(): Crucially, to get the first part (left operand) of an additive expression (e.g., the A in A + B), it calls parseMultiplicative(). This means anything that parseMultiplicative can parse (like X * Y or just a number Z) can be the left operand of a + or -.
//	Rule 2 (The Loop): for p.currentTokenIs(lexer.PLUS, lexer.MINUS):
//	    It then checks if the current token is a + or -.
//	    If it is, it means we have an ongoing additive expression (e.g., we've parsed A and now see + B).
//	    The loop continues as long as it finds + or - operators. This handles chains like A + B - C.
//	Inside the Loop (Building the AST):
//	    operator := p.curToken.Literal: Stores the operator (+ or -).
//	    p.nextToken(): Consumes the operator token.
//	    expr = &ast.BinaryExpression{ Left: expr, ... }: This is where the Abstract Syntax Tree (AST) node is built. The expr that was parsed so far becomes the Left child of the new BinaryExpression.
//	    Right: p.parseMultiplicative(): Rule 3: To get the right operand (e.g., the B in A + B), it again calls parseMultiplicative(). This ensures that if you have A + B * C, B * C will be parsed by parseMultiplicative first and become a single unit (an AST node) before being made the right child of the + operation.
//
// Associativity: Because expr (the result of the previous operation) becomes the Left operand of the new operation, this structure naturally handles left-associativity for + and -. For example, A + B - C is parsed as (A + B) - C.
func (p *parser) parseAdditive() ast.Expression {
	expr := p.parseMultiplicative() // Rule 1: Get the left operand

	// Rule 2: Loop for subsequent '+' or '-' operators
	for p.currentTokenIs(lexer.PLUS, lexer.MINUS) {
		operator := p.curToken.Literal // Get the operator
		p.nextToken()                  // Consume the operator token
		// The current 'expr' is the left side of the new BinaryExpression
		// The right side is whatever parseMultiplicative gives us next
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    p.parseMultiplicative(), // Rule 3: Get the right operand
		}
	}
	return expr
}

func (p *parser) parseMultiplicative() ast.Expression {
	expr := p.parsePrimary() // Rule 1: Get the left operand

	// Rule 2: Loop for subsequent '*' or '/' operators
	for p.currentTokenIs(lexer.ASTERISK, lexer.SLASH) {
		operator := p.curToken.Literal
		p.nextToken()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    p.parsePrimary(), // Rule 3: Get the right operand
		}
	}
	return expr
}

// Logic Breakdown:
// lexer.NUMBER: If it's a number, it creates an ast.NumberLiteral node and advances the token.
// lexer.IDENT: If it's an identifier (like a variable name), it calls p.parseIdentifier() (which should return an ast.Identifier node or similar ast.Expression).
// lexer.LPAREN: This handles parenthesized expressions like (A + B).
//
//	It consumes the (.
//	Crucially, it calls p.parseExpression() recursively. This allows an entire new expression (with its own precedence rules) to be parsed within the parentheses.
//	It then expects a ).
//	The AST of the inner expression is returned. This is how parentheses override the default operator precedence.
//
// default: If the token is none of the above, it's an unexpected token in this context. It records an error (good!) and returns nil.
// How Precedence is Achieved (Example: A + B * C)
//
//	parseExpression() calls parseAdditive().
//	parseAdditive() calls expr_left := p.parseMultiplicative() to get its first operand.
//	parseMultiplicative() (for A): calls p.parsePrimary(). parsePrimary returns A. parseMultiplicative sees no * or / after A, so it returns A.
//	So, expr_left in parseAdditive() is A.
//	parseAdditive() sees the + token.
//	It stores +.
//	It calls expr_right := p.parseMultiplicative() to get the right operand of +.
//	parseMultiplicative() (for B * C):
//		Calls sub_expr_left := p.parsePrimary(). parsePrimary returns B.
//		parseMultiplicative() sees *. Stores *.
//		Calls sub_expr_right := p.parsePrimary(). parsePrimary returns C.
//		parseMultiplicative() builds BinaryExpression{Left: B, Operator: "*", Right: C}.
//		It sees no more * or / tokens. It returns the (B * C) AST node.
//	Back in parseAdditive(), expr_right is now the AST node for (B * C).
//	parseAdditive() builds BinaryExpression{Left: A, Operator: "+", Right: (B*C)}.
//	No more + or - tokens. parseAdditive() returns the final AST for (A + (B * C)).
//
// This shows how B * C is grouped together because parseAdditive calls parseMultiplicative to get its operands, and parseMultiplicative will fully resolve its higher-precedence operations before returning.
func (p *parser) parsePrimary() ast.Expression {
	switch p.curToken.Type {
	case lexer.NUMBER:
		lit := &ast.NumberLiteral{Value: p.curToken.Literal}
		p.nextToken()
		return lit
	case lexer.IDENT:
		return p.parseIdentifier() // Assumes parseIdentifier() returns an ast.Expression (e.g., ast.Identifier)
	case lexer.LPAREN: // '('
		p.nextToken()               // Consume '('
		expr := p.parseExpression() // Recursively parse the inner expression
		p.expect(lexer.RPAREN)      // Consume ')' - expect will handle error if not found
		return expr                 // Return the AST of the inner expression
	default:
		// This error handling is good!
		p.recordError(fmt.Sprintf("unexpected token in primary expression: %v", p.curToken.Type))
		return nil
	}
}

func (p *parser) parseFuncDecl() ast.Declaration {
	p.nextToken() // skip 'func'
	name := p.parseIdentifier()

	p.expect(lexer.LPAREN)
	p.parseParameterList()
	p.expect(lexer.RPAREN)

	// TODO: Parse return type and body
	return &ast.FuncDecl{Name: name}
}

// func anotherFunc(a :int, b :int, c :int)
func (p *parser) parseParameterList() []*ast.ParameterNode {
	parameters := []*ast.ParameterNode{}

	// Handle empty parameter list: if current token is RPAREN, e.g. foo()
	if p.currentTokenIs(lexer.RPAREN) {
		return parameters
	}

	// Parse the first parameter
	var paramNameNode *ast.Identifier
	var paramTypeNode ast.TypeNode

	// 1. Parse the parameter name (IDENT)
	if !p.currentTokenIs(lexer.IDENT) {
		p.recordError(fmt.Sprintf("SyntaxError: Expected parameter name (identifier), got %v at line %d, column %d. Violated rule: <FUNC_PARAM_NAME>", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return parameters // Return to allow parser to potentially recover or hit RPAREN expectation
	}
	paramNameNode = p.parseIdentifier() // Consumes IDENT

	// 2. Expect and consume COLON (:)
	if !p.currentTokenIs(lexer.COLON) {
		p.recordError(fmt.Sprintf("SyntaxError: Expected ':' after parameter name '%s', got %v at line %d, column %d. Violated rule: <FUNC_PARAM_COLON>", paramNameNode.Name, p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return parameters
	}
	p.nextToken() // Consume COLON

	// 3. Parse the Type (currently an IDENT, e.g., 'int')
	if !p.currentTokenIs(lexer.IDENT) {
		p.recordError(fmt.Sprintf("SyntaxError: Expected parameter type (identifier) after ':', got %v at line %d, column %d. Violated rule: <FUNC_PARAM_TYPE>", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return parameters
	}
	tempTypeIdent := p.parseIdentifier() // Consumes type IDENT
	paramTypeNode = tempTypeIdent        // ast.Identifier implements ast.TypeNode

	parameters = append(parameters, &ast.ParameterNode{Name: paramNameNode, Type: paramTypeNode})

	// Loop for subsequent parameters, e.g. (a:int, b:string)
	for p.currentTokenIs(lexer.COMMA) {
		p.nextToken() // Consume COMMA

		// After a comma, another full parameter declaration is expected.
		if !p.currentTokenIs(lexer.IDENT) {
			p.recordError(fmt.Sprintf("SyntaxError: Expected parameter name (identifier) after comma, got %v at line %d, column %d. Violated rule: <FUNC_PARAM_NAME_AFTER_COMMA>", p.curToken.Type, p.curToken.Line, p.curToken.Column))
			return parameters // Syntax error
		}
		paramNameNode = p.parseIdentifier()

		if !p.currentTokenIs(lexer.COLON) {
			p.recordError(fmt.Sprintf("SyntaxError: Expected ':' after parameter name '%s', got %v at line %d, column %d. Violated rule: <FUNC_PARAM_COLON>", paramNameNode.Name, p.curToken.Type, p.curToken.Line, p.curToken.Column))
			return parameters
		}
		p.nextToken() // Consume COLON

		if !p.currentTokenIs(lexer.IDENT) {
			p.recordError(fmt.Sprintf("SyntaxError: Expected parameter type (identifier) after ':', got %v at line %d, column %d. Violated rule: <FUNC_PARAM_TYPE>", p.curToken.Type, p.curToken.Line, p.curToken.Column))
			return parameters
		}
		tempTypeIdent = p.parseIdentifier()
		paramTypeNode = tempTypeIdent

		parameters = append(parameters, &ast.ParameterNode{Name: paramNameNode, Type: paramTypeNode})
	}

	return parameters
}

func (p *parser) parseStructDecl() ast.Declaration {
	p.nextToken() // skip 'struct'
	name := p.parseIdentifier()

	p.expect(lexer.LBRACE)
	// TODO: Parse fields
	p.expect(lexer.RBRACE)

	return &ast.StructDecl{Name: name}
}

func (p *parser) parseUnsafeBlock() ast.Declaration {
	p.nextToken() // skip 'unsafe'
	p.expect(lexer.LBRACE)
	// TODO: Parse unsafe block contents
	p.expect(lexer.RBRACE)

	return &ast.UnsafeBlock{}
}

func (p *parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.LOG:
		return p.parsePrintStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression()
	return &ast.ExpressionStatement{Expression: expr}
}

func (p *parser) parseIfStatement() ast.Statement {
	p.nextToken() // skip 'if'

	condition := p.parseExpression()

	p.expect(lexer.LBRACE)
	consequence := &ast.BlockStatement{}
	// TODO: Parse consequence statements
	p.expect(lexer.RBRACE)

	return &ast.IfStatement{
		Condition:   condition,
		Consequence: consequence,
	}
}

func (p *parser) parseForStatement() ast.Statement {
	p.nextToken() // skip 'for'

	// TODO: Parse init, condition, post
	p.expect(lexer.LBRACE)
	body := &ast.BlockStatement{}
	// TODO: Parse body statements
	p.expect(lexer.RBRACE)

	return &ast.ForStatement{Body: body}
}

func (p *parser) parsePrintStatement() ast.Statement {
	p.nextToken() // skip 'log'
	p.expect(lexer.LPAREN)
	expr := p.parseExpression()
	p.expect(lexer.RPAREN)

	return &ast.PrintStatement{Expression: expr}
}
