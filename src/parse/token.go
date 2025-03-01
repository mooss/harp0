package parse

type TokenType string

const (
	// Special.
	TOKEN_EOF     TokenType = "EOF"
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	TOKEN_COMMENT TokenType = "COMMENT"

	// Atoms.
	TOKEN_SYMBOL TokenType = "SYMBOL"
	TOKEN_INT    TokenType = "INT"
	TOKEN_FLOAT  TokenType = "FLOAT"
	TOKEN_STRING TokenType = "STRING"

	// Stoprunes, runes validly end any token and can appear right next to anything (whitespace is
	// also a stoprune but is only a delimiter).
	TOKEN_LPAREN   TokenType = "LPAREN"   // (
	TOKEN_RPAREN   TokenType = "RPAREN"   // )
	TOKEN_LBRACE   TokenType = "LBRACE"   // {
	TOKEN_RBRACE   TokenType = "RBRACE"   // }
	TOKEN_LBRACKET TokenType = "LBRACKET" // [
	TOKEN_RBRACKET TokenType = "RBRACKET" // ]

	// Other single characters with syntactic meaning in Harp.
	TOKEN_DOT   TokenType = "DOT"   // .
	TOKEN_COLON TokenType = "COLON" // :
	TOKEN_PIPE  TokenType = "PIPE"  // |
	TOKEN_QUOTE TokenType = "QUOTE" // '
	TOKEN_UNDER TokenType = "UNDER" // _
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
