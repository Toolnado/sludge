package lexer

import (
	"strings"
	"testing"

	"github.com/Toolnado/sludge/token"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "empty input",
			input: "",
			expected: []token.Token{
				{Type: token.EOF},
			},
		},
		{
			name:  "identifier",
			input: "test",
			expected: []token.Token{
				{Type: token.IDENTIFIER, Literal: "test"},
				{Type: token.EOF},
			},
		},
		{
			name:  "string",
			input: `"hello"`,
			expected: []token.Token{
				{Type: token.STRING, Literal: "hello"},
				{Type: token.EOF},
			},
		},
		{
			name:  "string with single quotes",
			input: "'hello'",
			expected: []token.Token{
				{Type: token.STRING, Literal: "hello"},
				{Type: token.EOF},
			},
		},
		{
			name:  "string with backticks",
			input: "`hello`",
			expected: []token.Token{
				{Type: token.RAW_STRING, Literal: "hello"},
				{Type: token.EOF},
			},
		},
		{
			name:  "number",
			input: "123",
			expected: []token.Token{
				{Type: token.INTEGER, Literal: "123"},
				{Type: token.EOF},
			},
		},
		{
			name:  "operators",
			input: "+ - * /",
			expected: []token.Token{
				{Type: token.PLUS, Literal: "+"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.STAR, Literal: "*"},
				{Type: token.SLASH, Literal: "/"},
				{Type: token.EOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.input))
			tokens := l.ScanTokens()

			if len(tokens) != len(tt.expected) {
				t.Errorf("ScanTokens() returned %d tokens, expected %d", len(tokens), len(tt.expected))
				return
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: got type %v, expected %v", i, tok.Type, tt.expected[i].Type)
				}
				if tok.Literal != tt.expected[i].Literal {
					t.Errorf("token %d: got literal %q, expected %q", i, tok.Literal, tt.expected[i].Literal)
				}
			}
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "unclosed string",
			input:   `"unclosed string`,
			wantErr: true,
		},
		{
			name:    "valid string",
			input:   `"closed string"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.input))
			l.ScanTokens()
			errors := l.Errors()

			if tt.wantErr && len(errors) == 0 {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("expected no errors but got: %v", errors)
			}
		})
	}
}
