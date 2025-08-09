// lexer/lexer.go
package lexer

import (
	"unicode"
	"unicode/utf8"
)

// Lexer represents a lexical scanner
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	prevChar := l.ch

	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
	}

	if prevChar == '\n' {
		l.line++
		l.column = 0
	}
	l.column++
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Store the start position for the token
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = EQEQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = EQ
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = BANGEQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = BANG
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = LTEQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = LT
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = GTEQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = GT
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok.Type = AND
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok.Type = OR
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '+':
		tok.Type = PLUS
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok.Type = ARROW
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = MINUS
			tok.Literal = string(l.ch)
		}
		l.readChar()
		return tok
	case '*':
		tok.Type = ASTERISK
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '/':
		if l.peekChar() == '*' {
			tok.Type = COMMENT
			tok.Literal = l.readBlockComment()
			return tok
		} else if l.peekChar() == '/' {
			tok.Type = COMMENT
			tok.Literal = l.readLineComment()
			return tok
		}
		tok.Type = SLASH
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '%':
		tok.Type = PERCENT
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '(':
		tok.Type = LPAREN
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case ')':
		tok.Type = RPAREN
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '{':
		tok.Type = LBRACE
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '}':
		tok.Type = RBRACE
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '[':
		tok.Type = LBRACKET
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case ']':
		tok.Type = RBRACKET
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case ',':
		tok.Type = COMMA
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '.':
		tok.Type = DOT
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case ':':
		if l.peekChar() == '=' {
			// consume '=' after ':'
			ch := l.ch
			l.readChar()
			tok.Type = COLONEQ
			tok.Literal = string(ch) + string(l.ch)
			tok.Line = l.line
			tok.Column = l.column
		} else {
			tok.Type = COLON
			tok.Literal = string(l.ch)
			tok.Line = l.line
			tok.Column = l.column
		}
		l.readChar()
		return tok
	case ';':
		tok.Type = SEMICOLON
		tok.Literal = string(l.ch)
		l.readChar()
		return tok
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		l.readChar()
		return tok
	case 0:
		tok.Type = EOF
		tok.Literal = ""
		return tok
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
			l.readChar()
			return tok
		}
	}
}

// PeekTokenN returns the N-th upcoming token without consuming any input.
// n = 1 returns the next token that would be produced by NextToken(),
// n = 2 returns the token after that, and so on.
func (l *Lexer) PeekTokenN(n int) Token {
	// Make a shallow copy of the lexer state so advancing doesn't affect the original
	clone := *l
	var tok Token
	for i := 0; i < n; i++ {
		tok = clone.NextToken()
	}
	return tok
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	var identifier []rune

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		identifier = append(identifier, l.ch)
		l.readChar()
	}

	return string(identifier)
}

// readNumber reads a number literal
func (l *Lexer) readNumber() string {
	var number []rune

	for isDigit(l.ch) {
		number = append(number, l.ch)
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		number = append(number, l.ch)
		l.readChar()
		for isDigit(l.ch) {
			number = append(number, l.ch)
			l.readChar()
		}
	}

	return string(number)
}

// readString reads a string literal
func (l *Lexer) readString() string {
	position := l.position + 1 // skip the opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter checks if a rune is a letter
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit checks if a rune is a digit
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// readBlockComment reads a block comment /* ... */
func (l *Lexer) readBlockComment() string {
	position := l.position
	nestLevel := 1

	// Skip the opening /*
	l.readChar() // Skip '/'
	l.readChar() // Skip '*'

	for nestLevel > 0 {
		if l.ch == 0 {
			// EOF reached before comment end
			break
		}

		if l.ch == '/' && l.peekChar() == '*' {
			// Found nested comment start
			nestLevel++
			l.readChar() // Skip '/'
			l.readChar() // Skip '*'
		} else if l.ch == '*' && l.peekChar() == '/' {
			// Found comment end
			nestLevel--
			l.readChar() // Skip '*'
			l.readChar() // Skip '/'
		} else {
			// Track newlines in comments
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

// readLineComment reads a single-line comment // ...
func (l *Lexer) readLineComment() string {
	position := l.position

	// Skip the opening //
	l.readChar() // Skip '/'
	l.readChar() // Skip '/'

	// Read until end of line or EOF
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	return l.input[position:l.position]
}

// Tokens returns all tokens from the input
func (l *Lexer) Tokens() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}
