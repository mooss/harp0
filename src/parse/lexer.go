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
	Reason LexicalErrorReason
}

// lxr is a shortcut to build lexical errors.
func lxr(tok Token, kind LexicalErrorReason) *LexicalError {
	return &LexicalError{tok, kind}
}

func (le LexicalError) Error() string {
	return fmt.Sprintf(
		"lexical error at line %d column %d: %s",
		le.Line, le.Column, le.Reason,
	)
}

type LexicalErrorReason string

const (
	TwoDotsInFloat   LexicalErrorReason = "met a second dot while reading float"
	NonDigitInNumber LexicalErrorReason = "met non-digit while reading number"
	EofInString      LexicalErrorReason = "met EOF while reading string"
	NewlineInString  LexicalErrorReason = "met unescaped newline while reading string"
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
func (l *Lexer) forward() {
	if l.currentPosition >= len(l.input) { // Already at EOF.
		return
	}

	l.currentPosition += l.currentWidth
	l.column += 1
	if l.currentPosition >= len(l.input) { // Reached EOF.
		l.current = 0
		return
	}

	l.current, l.currentWidth = utf8.DecodeRuneInString(l.input[l.currentPosition:])
}

// nextLine registers that the input has moved to the next line (it does not change the position).
func (l *Lexer) nextLine() {
	l.line++
	l.column = -1 // -1 to ensure first column is 0.
}

// peekChar return the rune of *the next byte* (not exactly the next rune).
func (l *Lexer) peekChar() rune {
	npos := l.currentPosition + l.currentWidth
	if npos >= len(l.input) {
		return 0
	}

	return rune(l.input[npos])
}

// NextToken produces the next token by moving the lexer forward.
func (l *Lexer) NextToken() (Token, *LexicalError) {
	var tok Token

	l.skipWhitespace()

	switch l.current {
	case '(':
		tok = l.monotok(TOKEN_LPAREN)
	case ')':
		tok = l.monotok(TOKEN_RPAREN)
	case '{':
		tok = l.monotok(TOKEN_LBRACE)
	case '}':
		tok = l.monotok(TOKEN_RBRACE)
	case '[':
		tok = l.monotok(TOKEN_LBRACKET)
	case ']':
		tok = l.monotok(TOKEN_RBRACKET)
	case '.':
		if isDigit(l.peekChar()) { // Float < 1.
			return l.readNumber()
		}

		tok = l.monotok(TOKEN_DOT)
	case ':':
		tok = l.monotok(TOKEN_COLON)
	case '|':
		tok = l.monotok(TOKEN_PIPE)
	case '\'':
		tok = l.monotok(TOKEN_QUOTE)
	case '_':
		tok = l.monotok(TOKEN_UNDER)
	case '"':
		return l.readString()
	case ';':
		return l.readComment()
	case 0:
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if canStartSymbol(l.current) {
			return l.readSymbol()
		} else if isDigit(l.current) {
			return l.readNumber()
		} else {
			tok = l.monotok(TOKEN_ILLEGAL)
		}
	}

	l.forward()
	return tok, nil
}

/////////////
// Readers //

func (l *Lexer) readComment() (Token, *LexicalError) {
	tok, start := l.start(TOKEN_COMMENT)

	for l.current != '\n' && l.current != 0 {
		l.forward()
	}

	tok.Literal = l.from(start)

	if l.current == '\n' {
		l.nextLine()
		l.forward()
	}

	return tok, nil
}

func (l *Lexer) readSymbol() (Token, *LexicalError) {
	tok, start := l.start(TOKEN_SYMBOL)

	for canStartSymbol(l.current) || isDigit(l.current) {
		l.forward()
	}

	tok.Literal = l.from(start)
	return tok, nil
}

func (l *Lexer) readNumber() (Token, *LexicalError) {
	tok, start := l.start(TOKEN_INT)

	for {
		if l.current == '.' {
			// One dot is a float, two dots is an error.
			if tok.Type == TOKEN_FLOAT {
				tok.Literal = l.from(start)
				return Token{}, lxr(tok, TwoDotsInFloat)
			}

			tok.Type = TOKEN_FLOAT
		} else if isStoprune(l.current) {
			break
		} else if !isDigit(l.current) {
			tok.Literal = l.from(start)
			return Token{}, lxr(tok, NonDigitInNumber)
		}

		l.forward()
	}

	tok.Literal = l.from(start)
	return tok, nil
}

func (l *Lexer) readString() (Token, *LexicalError) {
	tok, start := l.start(TOKEN_STRING)
	l.forward() // Consume opening double quote.
	start = l.currentPosition

loop:
	for {
		switch l.current {
		case 0:
			tok.Literal = l.from(start)
			return Token{}, lxr(tok, EofInString)
		case '\n':
			tok.Literal = l.from(start)
			return Token{}, lxr(tok, NewlineInString)
		case '"':
			break loop
		case '\\': // Handle escape sequences.
			l.forward()
		}

		l.forward()
	}

	tok.Literal = l.from(start)
	l.forward() // Consume closing double quote.

	return tok, nil
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

// Start returns a token initialized at the current line, as well as the current position.
func (l *Lexer) start(tokenType TokenType) (Token, int) {
	return Token{
		Type:   tokenType,
		Line:   l.line,
		Column: l.column,
	}, l.currentPosition
}

// monotok is a shorcut that builds single-rune tokens.
func (l *Lexer) monotok(tokenType TokenType) Token {
	return Token{
		Type:    tokenType,
		Literal: string(l.current),
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.current {
		case ' ', '\t', '\r':
			l.forward()
		case '\n':
			l.nextLine()
			l.forward()
		default:
			return
		}
	}
}

// from returns the substring between start and the current position.
func (l *Lexer) from(start int) string { return l.input[start:l.currentPosition] }
