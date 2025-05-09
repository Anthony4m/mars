package lexer

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT // For block comments /* ... */

	// Identifiers and literals
	IDENT  // variable names, function names, etc.
	NUMBER // integers and floats
	STRING // string literals

	// Keywords
	MUT
	FUNC
	STRUCT
	UNSAFE
	IF
	ELSE
	FOR
	RETURN
	LOG
	TRUE
	FALSE
	NIL

	// Operators
	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %
	BANG    // !
	EQ      // =
	COLONEQ // :=
	EQEQ    // ==
	BANGEQ  // !=
	LT      // <
	LTEQ    // <=
	GT      // >
	GTEQ    // >=
	AND     // &&
	OR      // ||

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	DOT       // .
	COLON     // :
	SEMICOLON // ;
	ARROW     // ->
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// keywords maps string keywords to their corresponding token types
var keywords = map[string]TokenType{
	"mut":    MUT,
	"func":   FUNC,
	"struct": STRUCT,
	"unsafe": UNSAFE,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"return": RETURN,
	"log":    LOG,
	"true":   TRUE,
	"false":  FALSE,
	"nil":    NIL,
}

// LookupIdent checks if the given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
