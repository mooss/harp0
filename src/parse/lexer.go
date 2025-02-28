package parse

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

var ErrLexing = errors.New("lexing error")

// Lexer performs lexical analysis for Harp source code, that is to say it turns input text into tokens.
type Lexer struct {
	// input is the source code being lexically analyzed.
	input string

	// currentPosition is the position of the current character.
	currentPosition int

	// nextPosition is the position of the next character.
	nextPosition int

	// current is the character under examination.
	current rune

	// line is the current line number in the input.
	line int

	// column is the current column number in the input.
	column int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.forward()
	return l
}

// forward moves the lexer to the forward position.
func (l *Lexer) forward() {
	l.currentPosition = l.nextPosition

	if l.nextPosition >= len(l.input) { // EOF.
		l.current = 0
		return
	}

	r, size := utf8.DecodeRuneInString(l.input[l.nextPosition:])
	l.current = r
	l.nextPosition += size
	l.column += size
}

// peekChar return the rune of *the next byte* (not the next rune!!!).
func (l *Lexer) peekChar() rune {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return rune(l.input[l.nextPosition])
}

// NextToken produces the next token by moving the lexer forward.
func (l *Lexer) NextToken() (Token, error) {
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
		tok = l.monotok(TOKEN_DOT)
	case ':':
		tok = l.monotok(TOKEN_COLON)
	case ';':
		tok = l.monotok(TOKEN_SEMICOLON)
	case '|':
		tok = l.monotok(TOKEN_PIPE)
	case '\'':
		tok = l.monotok(TOKEN_QUOTE)
	case '_':
		tok = l.monotok(TOKEN_UNDER)
	case '"':
		var err error
		tok.Line = l.line
		tok.Column = l.column
		tok.Type = TOKEN_STRING
		tok.Literal, err = l.readString()

		if err != nil {
			return tok, err
		}
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if canStartSymbol(l.current) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readSymbol()
			tok.Type = lookupSymbol(tok.Literal)
			return tok, nil
		} else if isDigit(l.current) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Type, tok.Literal = l.readNumber()
			return tok, nil
		} else {
			tok = l.monotok(TOKEN_ILLEGAL)
		}
	}

	l.forward()
	return tok, nil
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

//////////////
// Skippers //

func (l *Lexer) skipWhitespace() {
	for {
		switch l.current {
		case ' ', '\t', '\r':
			l.forward()
		case '\n':
			l.line++
			l.column = 0
			l.forward()
		case ';':
			l.skipComment()
		default:
			return
		}
	}
}

func (l *Lexer) skipComment() {
	for l.current != '\n' && l.current != 0 {
		l.forward()
	}
}

/////////////
// Readers //

func (l *Lexer) readSymbol() string {
	position := l.currentPosition

	for canStartSymbol(l.current) || isDigit(l.current) {
		l.forward()
	}

	return l.input[position:l.currentPosition]
}

func (l *Lexer) readNumber() (TokenType, string) {
	tokenType := TOKEN_INT
	position := l.currentPosition

	for isDigit(l.current) {
		l.forward()
	}

	if l.current == '.' && isDigit(l.peekChar()) {
		tokenType = TOKEN_FLOAT
		l.forward() // consume the dot
		for isDigit(l.current) {
			l.forward()
		}
	}

	return tokenType, l.input[position:l.currentPosition]
}

func (l *Lexer) readString() (string, error) {
	l.forward() // Consume opening ".
	position := l.currentPosition

	for {
		if l.current == 0 {
			return "", fmt.Errorf("%w: met EOF while reading string", ErrLexing)
		}

		if l.current == '"' {
			break
		}

		if l.current == '\\' { // Handle escape sequences.
			l.forward()
		}

		l.forward()
	}

	str := l.input[position:l.currentPosition]
	l.forward() // Consume closing ".

	return str, nil
}

/////////////////////
// Rune predicates //

// canStartSymbol returns true if the given rune can start a valid symbol
// (unicode letter, _, -, +, / or *).
func canStartSymbol(run rune) bool {
	return unicode.IsLetter(run) || run == '_' || run == '-' || run == '+' || run == '/' || run == '*'
}

// isDigit returns true if run is an ASCII digit.
func isDigit(run rune) bool {
	return '0' <= run && run <= '9'
}
