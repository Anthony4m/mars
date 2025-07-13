// lexer/token.go
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
	BREAK
	CONTINUE

	// Type keywords (needed for parser)
	INT       // int
	FLOAT     // float
	STRING_KW // string (keyword, not literal)
	BOOL      // bool

	// Declaration keywords (for future use)
	ENUM // enum
	TYPE // type

	// Operators
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	PERCENT  // %
	BANG     // !
	EQ       // =
	COLONEQ  // :=
	EQEQ     // ==
	BANGEQ   // !=
	LT       // <
	LTEQ     // <=
	GT       // >
	GTEQ     // >=
	AND      // &&
	OR       // ||

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
	"mut":      MUT,
	"func":     FUNC,
	"struct":   STRUCT,
	"unsafe":   UNSAFE,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"return":   RETURN,
	"log":      LOG,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"break":    BREAK,
	"continue": CONTINUE,

	// Type keywords (these are essential for the parser)
	"int":    INT,
	"float":  FLOAT,
	"string": STRING_KW,
	"bool":   BOOL,

	// Declaration keywords (for future expansion)
	"enum": ENUM,
	"type": TYPE,
}

// LookupIdent checks if the given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// String method for TokenType for better debugging
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case COMMENT:
		return "COMMENT"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case MUT:
		return "MUT"
	case FUNC:
		return "FUNC"
	case STRUCT:
		return "STRUCT"
	case UNSAFE:
		return "UNSAFE"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case RETURN:
		return "RETURN"
	case LOG:
		return "LOG"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NIL:
		return "NIL"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING_KW:
		return "STRING_KW"
	case BOOL:
		return "BOOL"
	case ENUM:
		return "ENUM"
	case TYPE:
		return "TYPE"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case PERCENT:
		return "PERCENT"
	case BANG:
		return "BANG"
	case EQ:
		return "EQ"
	case COLONEQ:
		return "COLONEQ"
	case EQEQ:
		return "EQEQ"
	case BANGEQ:
		return "BANGEQ"
	case LT:
		return "LT"
	case LTEQ:
		return "LTEQ"
	case GT:
		return "GT"
	case GTEQ:
		return "GTEQ"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case ARROW:
		return "ARROW"
	default:
		return "UNKNOWN"
	}
}
