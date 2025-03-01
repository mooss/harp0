package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

////////////
// Errors //
////////////

type LexicalError struct {
	// Token is the token that was being parsed when the error occured.
	Token

	// Reason explains what triggered the error.
	Reason LexicalFailure
}

func (le LexicalError) Error() string {
	return fmt.Sprintf(
		"lexical error at line %d column %d: %s",
		le.Line, le.Column, le.Reason,
	)
}

// LexicalFailure describes what caused the lexer to fail.
type LexicalFailure string

const (
	TwoDotsInFloat   LexicalFailure = "met a second dot while reading float"
	NonDigitInNumber LexicalFailure = "met non-digit while reading number"
	EofInString      LexicalFailure = "met EOF while reading string"
	NewlineInString  LexicalFailure = "met unescaped newline while reading string"
)

///////////
// Lexer //
///////////

// Lexer performs lexical analysis for Harp source code, that is to say it turns input text into tokens.
type Lexer struct {
	// input is the source code being lexically analyzed.
	input string

	// currentPosition is the position of the current character.
	currentPosition int

	// current is the character under examination.
	current rune

	// currentWidth is the currentWidth of the current rune (the number of bytes used to represent it).
	currentWidth int

	// line is the current line number in the input.
	line int

	// column is the current column number in the input.
	column int
}

func NewLexer(input string) *Lexer {
	if len(input) == 0 {
		return &Lexer{line: 1}
	}

	l := &Lexer{input: input, line: 1, column: -1} // -1 to ensure first column is 0.
	l.forward()
	return l
}

// forward moves the lexer to the forward position.
func (lex *Lexer) forward() {
	if lex.currentPosition >= len(lex.input) { // Already at EOF.
		return
	}

	lex.currentPosition += lex.currentWidth
	lex.column += 1
	if lex.currentPosition >= len(lex.input) { // Reached EOF.
		lex.current = 0
		return
	}

	lex.current, lex.currentWidth = utf8.DecodeRuneInString(lex.input[lex.currentPosition:])
}

// nextLine registers that the input has moved to the next line (it does not change the position).
func (lex *Lexer) nextLine() {
	lex.line++
	lex.column = -1 // -1 to ensure first column is 0.
}

// peekChar return the rune of *the next byte* (not exactly the next rune).
func (lex *Lexer) peekChar() rune {
	npos := lex.currentPosition + lex.currentWidth
	if npos >= len(lex.input) {
		return 0
	}

	return rune(lex.input[npos])
}

// NextToken produces the next token by moving the lexer forward.
func (lex *Lexer) NextToken() (Token, *LexicalError) {
	var tok Token

	lex.skipWhitespace()

	switch lex.current {
	case '(':
		tok = lex.monotok(TOKEN_LPAREN)
	case ')':
		tok = lex.monotok(TOKEN_RPAREN)
	case '{':
		tok = lex.monotok(TOKEN_LBRACE)
	case '}':
		tok = lex.monotok(TOKEN_RBRACE)
	case '[':
		tok = lex.monotok(TOKEN_LBRACKET)
	case ']':
		tok = lex.monotok(TOKEN_RBRACKET)
	case '.':
		if isDigit(lex.peekChar()) { // Float < 1.
			return lex.read(readNumber, TOKEN_INT)
		}

		tok = lex.monotok(TOKEN_DOT)
	case ':':
		tok = lex.monotok(TOKEN_COLON)
	case '|':
		tok = lex.monotok(TOKEN_PIPE)
	case '\'':
		tok = lex.monotok(TOKEN_QUOTE)
	case '_':
		tok = lex.monotok(TOKEN_UNDER)
	case '"':
		return lex.read(readString, TOKEN_STRING)
	case ';':
		return lex.read(readComment, TOKEN_COMMENT)
	case 0:
		tok.Line = lex.line
		tok.Column = lex.column
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if canStartSymbol(lex.current) {
			return lex.read(readSymbol, TOKEN_SYMBOL)
		} else if isDigit(lex.current) {
			return lex.read(readNumber, TOKEN_INT)
		} else {
			tok = lex.monotok(TOKEN_ILLEGAL)
		}
	}

	// The current character is a part of the returned token, so it must be skipped.
	lex.forward()

	return tok, nil
}

/////////////
// Readers //

// reader defines a function iterating forward in a lexer to build a token.
type reader func(*Lexer, *Token) LexicalFailure

// read takes a reader, does boilerplate pre and post processing and builds a token.
func (lex *Lexer) read(
	fun reader, typ TokenType,
) (Token, *LexicalError) {
	tok := Token{
		Type:   typ,
		Line:   lex.line,
		Column: lex.column,
	}
	start := lex.currentPosition

	fail := fun(lex, &tok)
	tok.Literal = lex.input[start:lex.currentPosition]

	if fail != "" {
		return Token{}, &LexicalError{tok, fail}
	}

	return tok, nil
}

func readComment(lex *Lexer, tok *Token) LexicalFailure {
	for lex.current != '\n' && lex.current != 0 {
		lex.forward()
	}

	return ""
}

func readNumber(lex *Lexer, tok *Token) LexicalFailure {
	for {
		switch run := lex.current; {
		case run == '.':
			// One dot is a float, two dots is an error.
			if tok.Type == TOKEN_FLOAT {
				return TwoDotsInFloat
			}

			tok.Type = TOKEN_FLOAT
		case isStoprune(run):
			return ""
		case !isDigit(run):
			return NonDigitInNumber
		}

		lex.forward()
	}
}

func readString(lex *Lexer, tok *Token) LexicalFailure {
	lex.forward() // Consume opening double quote.

	for {
		switch lex.current {
		case 0:
			return EofInString
		case '\n':
			return NewlineInString
		case '"':
			lex.forward()
			return ""
		case '\\': // Handle escape sequences.
			lex.forward()
		}

		lex.forward()
	}
}

func readSymbol(lex *Lexer, tok *Token) LexicalFailure {
	for canStartSymbol(lex.current) || isDigit(lex.current) {
		lex.forward()
	}

	return ""
}

/////////////////////
// Rune predicates //

// canStartSymbol returns true if the given rune can start a valid symbol
// (unicode letter, _, -, +, / or *).
func canStartSymbol(run rune) bool {
	return unicode.IsLetter(run) || strings.ContainsRune("_-+/*", run)
}

// isDigit returns true if run is an ASCII digit.
func isDigit(run rune) bool {
	return '0' <= run && run <= '9'
}

// isStoprune returns true when given a stoprune, that is to say a rune that can validly end any
// token and can appear right next to anything.
// For instance, `(` is a stoprune, but `:` is not (it cannot end an int).
func isStoprune(run rune) bool {
	return strings.ContainsRune("()[]{} \t\r\n\000", run)
}

///////////////////////
// Utility functions //

// monotok is a shorcut that builds single-rune tokens.
func (lex *Lexer) monotok(tokenType TokenType) Token {
	return Token{
		Type:    tokenType,
		Literal: string(lex.current),
		Line:    lex.line,
		Column:  lex.column,
	}
}

func (lex *Lexer) skipWhitespace() {
	for {
		switch lex.current {
		case ' ', '\t', '\r':
			lex.forward()
		case '\n':
			lex.nextLine()
			lex.forward()
		default:
			return
		}
	}
}
