package parse

import (
	"testing"
)

type expected struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	Reason  LexicalFailure
}

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []expected // Empty reason means no error and the token fields are used instead.
	}{
		{
			name:  "Basic symbols",
			input: "def let fun struct lambda",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "def", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "let", Line: 1, Column: 4},
				{Type: TOKEN_SYMBOL, Literal: "fun", Line: 1, Column: 8},
				{Type: TOKEN_SYMBOL, Literal: "struct", Line: 1, Column: 12},
				{Type: TOKEN_SYMBOL, Literal: "lambda", Line: 1, Column: 19},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 25},
			},
		},
		{
			name:  "Numbers",
			input: "123 45.67 89.0",
			expected: []expected{
				{Type: TOKEN_INT, Literal: "123", Line: 1, Column: 0},
				{Type: TOKEN_FLOAT, Literal: "45.67", Line: 1, Column: 4},
				{Type: TOKEN_FLOAT, Literal: "89.0", Line: 1, Column: 10},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 14},
			},
		},
		{
			name:  "Strings",
			input: `"hello" "world"`,
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `"hello"`, Line: 1, Column: 0},
				{Type: TOKEN_DQSTRING, Literal: `"world"`, Line: 1, Column: 8},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 15},
			},
		},
		{
			name:  "Comment at start of line",
			input: "; This is a comment\n123",
			expected: []expected{
				{Type: TOKEN_COMMENT, Literal: "; This is a comment", Line: 1, Column: 0},
				{Type: TOKEN_INT, Literal: "123", Line: 2, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 2, Column: 3},
			},
		},
		{
			name:  "Symbols with special characters",
			input: "a-b_c/d*e",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "a-b_c/d*e", Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 9},
			},
		},
		{
			name:  "Method call",
			input: "obj.method(arg)",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "obj", Line: 1, Column: 0},
				{Type: TOKEN_DOT, Literal: ".", Line: 1, Column: 3},
				{Type: TOKEN_SYMBOL, Literal: "method", Line: 1, Column: 4},
				{Type: TOKEN_LPAREN, Literal: "(", Line: 1, Column: 10},
				{Type: TOKEN_SYMBOL, Literal: "arg", Line: 1, Column: 11},
				{Type: TOKEN_RPAREN, Literal: ")", Line: 1, Column: 14},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 15},
			},
		},
		{
			name:  "Mixed symbols and numbers",
			input: "a123 b45.67",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "a123", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "b45", Line: 1, Column: 5,
					Reason: InvalidAfterSymbol.WithStrhex(".6")},
				{Type: TOKEN_FLOAT, Literal: ".67", Line: 1, Column: 8},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 11},
			},
		},
		{
			name:  "Parens and braces",
			input: "(a [b] {c})",
			expected: []expected{
				{Type: TOKEN_LPAREN, Literal: "(", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "a", Line: 1, Column: 1},
				{Type: TOKEN_LBRACKET, Literal: "[", Line: 1, Column: 3},
				{Type: TOKEN_SYMBOL, Literal: "b", Line: 1, Column: 4},
				{Type: TOKEN_RBRACKET, Literal: "]", Line: 1, Column: 5},
				{Type: TOKEN_LBRACE, Literal: "{", Line: 1, Column: 7},
				{Type: TOKEN_SYMBOL, Literal: "c", Line: 1, Column: 8},
				{Type: TOKEN_RBRACE, Literal: "}", Line: 1, Column: 9},
				{Type: TOKEN_RPAREN, Literal: ")", Line: 1, Column: 10},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 11},
			},
		},
		{
			name:  "Special characters",
			input: ". : | ' _",
			expected: []expected{
				{Type: TOKEN_DOT, Literal: ".", Line: 1, Column: 0},
				{Type: TOKEN_COLON, Literal: ":", Line: 1, Column: 2},
				{Type: TOKEN_PIPE, Literal: "|", Line: 1, Column: 4},
				{Type: TOKEN_QUOTE, Literal: "'", Line: 1, Column: 6},
				{Type: TOKEN_UNDER, Literal: "_", Line: 1, Column: 8},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 9},
			},
		},
		{
			name:  "Comment at end of line",
			input: "ignore the rest ; !!!!!@#",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "ignore", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "the", Line: 1, Column: 7},
				{Type: TOKEN_SYMBOL, Literal: "rest", Line: 1, Column: 11},
				{Type: TOKEN_COMMENT, Literal: "; !!!!!@#", Line: 1, Column: 16},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 25},
			},
		},
		{
			name:  "Illegal characters",
			input: "!@#",
			expected: []expected{
				{Type: TOKEN_ILLEGAL, Literal: "!", Line: 1, Column: 0},
				{Type: TOKEN_ILLEGAL, Literal: "@", Line: 1, Column: 1},
				{Type: TOKEN_ILLEGAL, Literal: "#", Line: 1, Column: 2},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 3},
			},
		},
		{
			name:  "Empty input",
			input: "",
			expected: []expected{
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 0},
			},
		},
		{
			name:  "Whitespace",
			input: " \t\n ",
			expected: []expected{
				{Type: TOKEN_EOF, Literal: "", Line: 2, Column: 1},
			},
		},
		{
			name:  "Unterminated string",
			input: `"hello`,
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `"hello`, Line: 1, Column: 0,
					Reason: EofInString},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "Unescaped newline in string",
			input: "\"hello\n\"",
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `"hello`, Line: 1, Column: 0,
					Reason: NewlineInString},
				{Type: TOKEN_DQSTRING, Literal: `"`, Line: 2, Column: 0,
					Reason: EofInString},
				{Type: TOKEN_EOF, Literal: "", Line: 2, Column: 1},
			},
		},
		{
			name:  "String with escaped characters",
			input: `"hello\nworld\t\"quoted\"\\escaped\\"`,
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `"hello\nworld\t\"quoted\"\\escaped\\"`, Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 37},
			},
		},
		{
			name:  "Unicode characters",
			input: "你好世界 ; This is a comment with Unicode: こんにちは",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "你好世界", Line: 1, Column: 0},
				{Type: TOKEN_COMMENT, Literal: "; This is a comment with Unicode: こんにちは", Line: 1, Column: 5},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 44},
			},
		},
		{
			name:  "Long numbers",
			input: "12345678901234567890 1234567890.1234567890",
			expected: []expected{
				{Type: TOKEN_INT, Literal: "12345678901234567890", Line: 1, Column: 0},
				{Type: TOKEN_FLOAT, Literal: "1234567890.1234567890", Line: 1, Column: 21},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 42},
			},
		},
		{
			name:  "Float without leading zero",
			input: ".123",
			expected: []expected{
				{Type: TOKEN_FLOAT, Literal: ".123", Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "Multiple dots in float (invalid)",
			input: "1.2.3",
			expected: []expected{
				{Type: TOKEN_FLOAT, Literal: "1.2", Line: 1, Column: 0,
					Reason: TwoDotsInFloat},
				{Type: TOKEN_FLOAT, Literal: ".3", Line: 1, Column: 3},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 5},
			},
		},
		{
			name:  "Int then non-digit",
			input: "1abc",
			expected: []expected{
				{Type: TOKEN_INT, Literal: "1", Line: 1, Column: 0,
					Reason: NonDigitInNumber},
				{Type: TOKEN_SYMBOL, Literal: "abc", Line: 1, Column: 1},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "x.y float then non-digit",
			input: "1.0abc",
			expected: []expected{
				{Type: TOKEN_FLOAT, Literal: "1.0", Line: 1, Column: 0,
					Reason: NonDigitInNumber},
				{Type: TOKEN_SYMBOL, Literal: "abc", Line: 1, Column: 3},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "Zero width characters",
			input: "a \u200b\u200cb",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "a", Line: 1, Column: 0},
				{Type: TOKEN_ILLEGAL, Literal: "\u200b", Line: 1, Column: 2},
				{Type: TOKEN_ILLEGAL, Literal: "\u200c", Line: 1, Column: 3},
				{Type: TOKEN_SYMBOL, Literal: "b", Line: 1, Column: 4},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 5},
			},
		},
		{
			name:  "BOM character", // Byte order mark, weird unicode thingie.
			input: "\ufeffabc",
			expected: []expected{
				{Type: TOKEN_ILLEGAL, Literal: "\ufeff", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "abc", Line: 1, Column: 1},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "Symbol followed by |",
			input: "lost|",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "lost", Line: 1, Column: 0,
					Reason: InvalidAfterSymbol.WithStrhex("|")},
				{Type: TOKEN_PIPE, Literal: "|", Line: 1, Column: 4},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 5},
			},
		},
		{
			name:  "Symbol followed by .1",
			input: "lost.1",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "lost", Line: 1, Column: 0,
					Reason: InvalidAfterSymbol.WithStrhex(".1")},
				{Type: TOKEN_FLOAT, Literal: ".1", Line: 1, Column: 4},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "Empty string",
			input: `""`,
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `""`, Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 2},
			},
		},
		{
			name:  "String with only whitespace",
			input: `" \t\r\n "`,
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `" \t\r\n "`, Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 10},
			},
		},
		{
			name:  "Dot followed by non-digit, non-symbol start",
			input: ".:",
			expected: []expected{
				{Type: TOKEN_DOT, Literal: ".", Line: 1, Column: 0},
				{Type: TOKEN_COLON, Literal: ":", Line: 1, Column: 1},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 2},
			},
		},
		{
			name:  "More illegal characters",
			input: "§±~`°•",
			expected: []expected{
				{Type: TOKEN_ILLEGAL, Literal: "§", Line: 1, Column: 0},
				{Type: TOKEN_ILLEGAL, Literal: "±", Line: 1, Column: 1},
				{Type: TOKEN_ILLEGAL, Literal: "~", Line: 1, Column: 2},
				{Type: TOKEN_ILLEGAL, Literal: "`", Line: 1, Column: 3},
				{Type: TOKEN_ILLEGAL, Literal: "°", Line: 1, Column: 4},
				{Type: TOKEN_ILLEGAL, Literal: "•", Line: 1, Column: 5},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "CRLF Line Endings",
			input: "abc\r\ndef",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "abc", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "def", Line: 2, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 2, Column: 3},
			},
		},
		{
			name:  "Tab character",
			input: "abc\tdef",
			expected: []expected{
				{Type: TOKEN_SYMBOL, Literal: "abc", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "def", Line: 1, Column: 4},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 7},
			},
		},
		{
			name:  "More escaped characters in string",
			input: `"\\ \\\r \b \f \v \040 \x41"`,
			expected: []expected{
				{Type: TOKEN_DQSTRING, Literal: `"\\ \\\r \b \f \v \040 \x41"`, Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 28},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			for _, exp := range tt.expected {
				expTok := Token{exp.Type, exp.Literal, exp.Line, exp.Column}
				expFail := exp.Reason
				gotFail := LexicalFailure("")
				gotTok, err := lexer.NextToken()

				if expFail == "" {
					expFail = "<nil>"
				}
				if err == nil {
					gotFail = "<nil>"
				} else {
					gotFail = err.Reason
					gotTok = err.Token
				}

				if expFail != gotFail {
					t.Errorf("expected failure:\n> %s\ngot:\n> %s", expFail, gotFail)
				}
				if expTok != gotTok {
					t.Errorf("expected %+v, got: %+v", expTok, gotTok)
				}
			}

			// Last expected token must be EOF (to be sure that the whole sentence is tested).
			lastExp := tt.expected[len(tt.expected)-1]
			if lastExp.Type != TOKEN_EOF {
				t.Errorf("last expected token type must be EOF")
			}
		})
	}
}
