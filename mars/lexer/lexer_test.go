package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		mut x := 42;
		/* This is a block comment */
		func add(a: int, b: int) -> int {
			return a + b;
		}
		/* This is a
		   multi-line
		   block comment */
		struct Point {
			x: int;
			y: int;
		}
		/* This is a /* nested */ block comment */
		unsafe {
			var p: *Point;
		}
		if x > 0 {
			log("positive");
		} else {
			log("negative");
		}
		for i := 0; i < 10; i = i + 1 {
			log(i);
		}
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{MUT, "mut", 2, 3},
		{IDENT, "x", 2, 7},
		{COLONEQ, ":=", 2, 10},
		{NUMBER, "42", 2, 12},
		{SEMICOLON, ";", 2, 14},
		{COMMENT, "/* This is a block comment */", 3, 3},
		{FUNC, "func", 4, 3},
		{IDENT, "add", 4, 8},
		{LPAREN, "(", 4, 11},
		{IDENT, "a", 4, 12},
		{COLON, ":", 4, 13},
		{IDENT, "int", 4, 15},
		{COMMA, ",", 4, 18},
		{IDENT, "b", 4, 20},
		{COLON, ":", 4, 21},
		{IDENT, "int", 4, 23},
		{RPAREN, ")", 4, 26},
		{ARROW, "->", 4, 28},
		{IDENT, "int", 4, 31},
		{LBRACE, "{", 4, 35},
		{RETURN, "return", 5, 4},
		{IDENT, "a", 5, 11},
		{PLUS, "+", 5, 13},
		{IDENT, "b", 5, 15},
		{SEMICOLON, ";", 5, 16},
		{RBRACE, "}", 6, 3},
		{COMMENT, "/* This is a\n\t\t   multi-line\n\t\t   block comment */", 7, 3},
		{STRUCT, "struct", 12, 3},
		{IDENT, "Point", 12, 10},
		{LBRACE, "{", 12, 16},
		{IDENT, "x", 13, 4},
		{COLON, ":", 13, 5},
		{IDENT, "int", 13, 7},
		{SEMICOLON, ";", 13, 10},
		{IDENT, "y", 14, 4},
		{COLON, ":", 14, 5},
		{IDENT, "int", 14, 7},
		{SEMICOLON, ";", 14, 10},
		{RBRACE, "}", 15, 3},
		{COMMENT, "/* This is a /* nested */ block comment */", 16, 3},
		{UNSAFE, "unsafe", 17, 3},
		{LBRACE, "{", 17, 10},
		{IDENT, "var", 18, 4},
		{IDENT, "p", 18, 8},
		{COLON, ":", 18, 9},
		{ASTERISK, "*", 18, 11},
		{IDENT, "Point", 18, 12},
		{SEMICOLON, ";", 18, 17},
		{RBRACE, "}", 19, 3},
		{IF, "if", 20, 3},
		{IDENT, "x", 20, 6},
		{GT, ">", 20, 8},
		{NUMBER, "0", 20, 10},
		{LBRACE, "{", 20, 12},
		{LOG, "log", 21, 4},
		{LPAREN, "(", 21, 7},
		{STRING, "positive", 21, 8},
		{RPAREN, ")", 21, 18},
		{SEMICOLON, ";", 21, 19},
		{RBRACE, "}", 22, 3},
		{ELSE, "else", 22, 5},
		{LBRACE, "{", 22, 10},
		{LOG, "log", 23, 4},
		{LPAREN, "(", 23, 7},
		{STRING, "negative", 23, 8},
		{RPAREN, ")", 23, 18},
		{SEMICOLON, ";", 23, 19},
		{RBRACE, "}", 24, 3},
		{FOR, "for", 25, 3},
		{IDENT, "i", 25, 7},
		{COLONEQ, ":=", 25, 10},
		{NUMBER, "0", 25, 12},
		{SEMICOLON, ";", 25, 13},
		{IDENT, "i", 25, 15},
		{LT, "<", 25, 17},
		{NUMBER, "10", 25, 19},
		{SEMICOLON, ";", 25, 21},
		{IDENT, "i", 25, 23},
		{EQ, "=", 25, 25},
		{IDENT, "i", 25, 27},
		{PLUS, "+", 25, 29},
		{NUMBER, "1", 25, 31},
		{LBRACE, "{", 25, 33},
		{LOG, "log", 26, 4},
		{LPAREN, "(", 26, 7},
		{IDENT, "i", 26, 8},
		{RPAREN, ")", 26, 9},
		{SEMICOLON, ";", 26, 10},
		{RBRACE, "}", 27, 3},
		{EOF, "", 28, 2},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}
