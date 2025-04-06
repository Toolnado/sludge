package token

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		pos       Position
		tokenType TokenType
		lexeme    string
		literal   interface{}
		want      Token
	}{
		{
			name: "create token with position",
			pos: Position{
				Filename: "test.go",
				Line:     1,
				Column:   1,
			},
			tokenType: IDENTIFIER,
			lexeme:    "test",
			literal:   "test",
			want: Token{
				Type:     IDENTIFIER,
				Lexeme:   "test",
				Literal:  "test",
				Position: Position{Filename: "test.go", Line: 1, Column: 1},
			},
		},
		{
			name:      "create token with empty position",
			pos:       Position{},
			tokenType: STRING,
			lexeme:    "\"hello\"",
			literal:   "hello",
			want: Token{
				Type:     STRING,
				Lexeme:   "\"hello\"",
				Literal:  "hello",
				Position: Position{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.pos, tt.tokenType, tt.lexeme, tt.literal)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
