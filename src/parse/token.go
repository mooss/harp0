package parse

type TokenType string

const (
	// Special.
	TOKEN_EOF     TokenType = "EOF"
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	TOKEN_COMMENT TokenType = "COMMENT"

	// Atoms.
	TOKEN_SYMBOL     TokenType = "SYMBOL"
	TOKEN_PREDEFINED TokenType = "PREDEFINED" // A symbol predefined by the language.
	TOKEN_INT        TokenType = "INT"
	TOKEN_FLOAT      TokenType = "FLOAT"
	TOKEN_STRING     TokenType = "STRING"

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

// keywords is the lookup table of predefined symbols.
var keywords = map[string]struct{}{
	"def":    {},
	"let":    {},
	"fun":    {},
	"struct": {},
	"lambda": {},
}

// lookupSymbol looks up whether symbol is predefined and returns the appropriate token type.
func lookupSymbol(symbol string) TokenType {
	if _, ok := keywords[symbol]; ok {
		return TOKEN_PREDEFINED
	}

	return TOKEN_SYMBOL
}
