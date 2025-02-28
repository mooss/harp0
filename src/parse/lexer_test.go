package parse

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
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
			name:  "Comments",
			input: "; This is a comment\n123",
			expected: []Token{
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
			name:  "Mixed symbols and numbers",
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
			name:  "Comments",
			input: "ignore the rest ; !!!!!@#",
			expected: []Token{
				{Type: TOKEN_SYMBOL, Literal: "ignore", Line: 1, Column: 0},
				{Type: TOKEN_SYMBOL, Literal: "the", Line: 1, Column: 7},
				{Type: TOKEN_SYMBOL, Literal: "rest", Line: 1, Column: 11},
				{Type: TOKEN_SEMICOLON, Literal: ";", Line: 1, Column: 13},
				{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 11},
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
		// {
		// 	name:  "Unclosed string",
		// 	input: `"hello`,
		// 	expected: []Token{
		// 		{Type: TOKEN_STRING, Literal: "hello", Line: 1, Column: 0},
		// 		{Type: TOKEN_EOF, Literal: "", Line: 1, Column: 6},
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
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
