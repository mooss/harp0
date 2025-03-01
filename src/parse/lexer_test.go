package parse

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
		err      *LexicalError
	}{
		{
			name:  "Basic symbols",
			input: "def let fun struct lambda",
			expected: []Token{
				{Type: TOKEN_PREDEFINED, Literal: "def", Line: 1, Column: 0},
				{Type: TOKEN_PREDEFINED, Literal: "let", Line: 1, Column: 4},
				{Type: TOKEN_PREDEFINED, Literal: "fun", Line: 1, Column: 8},
				{Type: TOKEN_PREDEFINED, Literal: "struct", Line: 1, Column: 12},
				{Type: TOKEN_PREDEFINED, Literal: "lambda", Line: 1, Column: 19},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 25},
			},
		},
		{
			name:  "Numbers",
			input: "123 45.67 89.0",
			expected: []Token{
				{Type: TOKEN_INT, Literal: "123", Line: 1, Column: 0},
				{Type: TOKEN_FLOAT, Literal: "45.67", Line: 1, Column: 4},
				{Type: TOKEN_FLOAT, Literal: "89.0", Line: 1, Column: 10},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 14},
			},
		},
		{
			name:  "Strings",
			input: `"hello" "world"`,
			expected: []Token{
				{Type: TOKEN_STRING, Literal: "hello", Line: 1, Column: 0},
				{Type: TOKEN_STRING, Literal: "world", Line: 1, Column: 8},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 15},
			},
		},
		{
			name:  "Comment at start of line",
			input: "; This is a comment\n123",
			expected: []Token{
				{Type: TOKEN_COMMENT, Literal: "; This is a comment", Line: 1, Column: 0},
				{Type: TOKEN_INT, Literal: "123", Line: 2, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 2, Column: 3},
			},
		},
		{
			name:  "Symbols with special characters",
			input: "a-b_c/d*e",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Literal: "a-b_c/d*e", Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 9},
			},
		},
		{
			name:  "Mixed symbols and numbers", // Broken, DOT needs to be replaced with DOTSYMBOL.
			input: "a123 b45.67",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Literal: "a123", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "b45", Line: 1, Column: 5},
				{Type: TOKEN_DOT, Literal: ".", Line: 1, Column: 8},
				{Type: TOKEN_INT, Literal: "67", Line: 1, Column: 9},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 11},
			},
		},
		{
			name:  "Parens and braces",
			input: "(a [b] {c})",
			expected: []Token{
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
			expected: []Token{
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
			expected: []Token{
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
			expected: []Token{
				{Type: TOKEN_ILLEGAL, Literal: "!", Line: 1, Column: 0},
				{Type: TOKEN_ILLEGAL, Literal: "@", Line: 1, Column: 1},
				{Type: TOKEN_ILLEGAL, Literal: "#", Line: 1, Column: 2},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 3},
			},
		},
		{
			name:  "Empty input",
			input: "",
			expected: []Token{
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 0},
			},
		},
		{
			name:  "Whitespace",
			input: " \t\n ",
			expected: []Token{
				{Type: TOKEN_EOF, Literal: "", Line: 2, Column: 1},
			},
		},
		{
			name:  "Unterminated string",
			input: `"hello`,
			err: &LexicalError{
				Token: Token{Type: TOKEN_STRING, Literal: `hello`, Line: 1, Column: 0},
				Kind:  UnterminatedString,
			},
		},
		{
			name:  "String with escaped characters",
			input: `"hello\nworld\t\"quoted\"\\escaped\\"`,
			expected: []Token{
				{Type: TOKEN_STRING, Literal: `hello\nworld\t\"quoted\"\\escaped\\`, Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 37},
			},
		},
		{
			name:  "Unicode characters",
			input: "你好世界 ; This is a comment with Unicode: こんにちは",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Literal: "你好世界", Line: 1, Column: 0},
				{Type: TOKEN_COMMENT, Literal: "; This is a comment with Unicode: こんにちは", Line: 1, Column: 5},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 44},
			},
		},
		{
			name:  "Long numbers",
			input: "12345678901234567890 1234567890.1234567890",
			expected: []Token{
				{Type: TOKEN_INT, Literal: "12345678901234567890", Line: 1, Column: 0},
				{Type: TOKEN_FLOAT, Literal: "1234567890.1234567890", Line: 1, Column: 21},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 42},
			},
		},
		{
			name:  "Float without leading zero",
			input: ".123",
			expected: []Token{
				{Type: TOKEN_FLOAT, Literal: ".123", Line: 1, Column: 0},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "Multiple dots in float (invalid)",
			input: "1.2.3",
			err: &LexicalError{
				Token: Token{Type: TOKEN_FLOAT, Literal: "1.2", Line: 1, Column: 0},
				Kind:  BadNumber,
			},
		},
		{
			name:  "Int then non-digit",
			input: "1abc",
			err: &LexicalError{
				Token: Token{Type: TOKEN_INT, Literal: "1", Line: 1, Column: 0},
				Kind:  BadNumber,
			},
		},
		{
			name:  "Floats x.y float then non-digit",
			input: "1.0abc",
			err: &LexicalError{
				Token: Token{Type: TOKEN_FLOAT, Literal: "1.0", Line: 1, Column: 0},
				Kind:  BadNumber,
			},
		},
		{
			name:  "Zero width characters",
			input: "a\u200b\u200cb",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Literal: "a", Line: 1, Column: 0},
				{Type: TOKEN_ILLEGAL, Literal: "\u200b", Line: 1, Column: 1},
				{Type: TOKEN_ILLEGAL, Literal: "\u200c", Line: 1, Column: 2},
				{Type: TOKEN_SYMBOL, Literal: "b", Line: 1, Column: 3},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "BOM character", // Byte order mark, weird unicode thingie.
			input: "\ufeffabc",
			expected: []Token{
				{Type: TOKEN_ILLEGAL, Literal: "\ufeff", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "abc", Line: 1, Column: 1},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			if tt.err != nil {
				var (
					err *LexicalError
					tok Token
				)

				for {
					tok, err = lexer.NextToken()
					if err != nil || tok.Type == TOKEN_EOF {
						break
					}
				}

				if err == nil || *err != *tt.err {
					t.Errorf("expected %+v[%s], got %+v[%s]",
						tt.err.Token, tt.err.Kind, err.Token, err.Kind)
				}

				return
			}

			for _, expected := range tt.expected {
				tok, err := lexer.NextToken()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if tok != expected {
					t.Errorf("expected %+v, got %+v", expected, tok)
				}
			}
		})
	}
}
