package lex

type TokenType string

const (
	/////////////
	// Special //

	// End of file.
	TOKEN_EOF TokenType = "EOF"
	// Invalid rune identified at the start of a token.
	TOKEN_INVALID TokenType = "INVALID"
	// Comment that stretches to the end of the line (semicolon).
	TOKEN_COMMENT TokenType = "COMMENT" // ;

	///////////
	// Atoms //

	// Identifier mapped to a value.
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

	// Opening parenthesis.
	TOKEN_LPAREN TokenType = "LPAREN" // (
	// Closing parenthesis.
	TOKEN_RPAREN TokenType = "RPAREN" // )
	// Opening curly brace.
	TOKEN_LBRACE TokenType = "LBRACE" // {
	// Closing curly brace.
	TOKEN_RBRACE TokenType = "RBRACE" // }
	// Opening square bracket.
	TOKEN_LBRACKET TokenType = "LBRACKET" // [
	// Closing square bracket.
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
	TOKEN_UNDERSCORE TokenType = "UNDER" // _
	// Pipe.
	TOKEN_PIPE TokenType = "PIPE" // |
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
