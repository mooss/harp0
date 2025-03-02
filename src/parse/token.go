package parse

type TokenType string

const (
	/////////////
	// Special //

	// End of file.
	TOKEN_EOF TokenType = "EOF"
	// Illegal rune that has been immediately identified as being incorrect.
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	// Comment that stretches to the end of the line (semicolon).
	TOKEN_COMMENT TokenType = "COMMENT" // ;

	///////////
	// Atoms //

	// Identifier supposed to be mapped to a value.
	TOKEN_SYMBOL TokenType = "SYMBOL"
	// Integer.
	TOKEN_INT TokenType = "INT"
	// Floating point number.
	TOKEN_FLOAT TokenType = "FLOAT"
	// Double quoted string.
	TOKEN_DQSTRING TokenType = "STRING"

	///////////////
	// Stoprunes //
	// Runes that validly end any token and can appear right next to anything.
	// Whitespace is also a stoprune but is only a delimiter, so it's not represented here.

	// Left parenthesis.
	TOKEN_LPAREN TokenType = "LPAREN" // (
	// RIght parenthesis.
	TOKEN_RPAREN TokenType = "RPAREN" // )
	// Left curly brace.
	TOKEN_LBRACE TokenType = "LBRACE" // {
	// Right curly brace.
	TOKEN_RBRACE TokenType = "RBRACE" // }
	// Left square bracket.
	TOKEN_LBRACKET TokenType = "LBRACKET" // [
	// Right square bracket.
	TOKEN_RBRACKET TokenType = "RBRACKET" // ]

	/////////////////
	// Other runes //

	// Dot, meant to be followed by a symbol (method call or field access).
	TOKEN_DOT TokenType = "DOT"
	// Colon, meant to be followed by a symbol (method of a type).
	TOKEN_COLON TokenType = "COLON"
	// Single quote.
	TOKEN_QUOTE TokenType = "QUOTE" // '
	// Underscore.
	TOKEN_UNDER TokenType = "UNDER" // _
	// Pipe.
	TOKEN_PIPE TokenType = "PIPE" // |
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
