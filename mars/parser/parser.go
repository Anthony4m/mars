// parser/parser.go
package parser

import (
	"mars/ast"
	"mars/lexer"
)

type parser struct {
	lexer     *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func NewParser(lexer *lexer.Lexer) *parser {
	p := &parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
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
		// TODO: Add proper error handling
		panic("unexpected token")
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
func (p *parser) parseExpression() ast.Expression {
	return p.parseAdditive()
}

func (p *parser) parseAdditive() ast.Expression {
	expr := p.parseMultiplicative()

	for p.currentTokenIs(lexer.PLUS, lexer.MINUS) {
		operator := p.curToken.Literal
		p.nextToken()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    p.parseMultiplicative(),
		}
	}

	return expr
}

func (p *parser) parseMultiplicative() ast.Expression {
	expr := p.parsePrimary()

	for p.currentTokenIs(lexer.ASTERISK, lexer.SLASH) {
		operator := p.curToken.Literal
		p.nextToken()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: operator,
			Right:    p.parsePrimary(),
		}
	}

	return expr
}

func (p *parser) parsePrimary() ast.Expression {
	switch p.curToken.Type {
	case lexer.NUMBER:
		lit := &ast.NumberLiteral{Value: p.curToken.Literal}
		p.nextToken()
		return lit
	case lexer.IDENT:
		return p.parseIdentifier()
	case lexer.LPAREN:
		p.nextToken()
		expr := p.parseExpression()
		p.expect(lexer.RPAREN)
		return expr
	default:
		// TODO: Add proper error handling
		panic("unexpected token in primary expression")
	}
}

func (p *parser) parseFuncDecl() ast.Declaration {
	p.nextToken() // skip 'func'
	name := p.parseIdentifier()

	p.expect(lexer.LPAREN)
	// TODO: Parse parameters
	p.expect(lexer.RPAREN)

	// TODO: Parse return type and body
	return &ast.FuncDecl{Name: name}
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
