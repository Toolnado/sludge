package token

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		pos       Position
		tokenType TokenType
		literal   string
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
			literal:   "test",
			want: Token{
				Type:     IDENTIFIER,
				Position: Position{Filename: "test.go", Line: 1, Column: 1},
				Literal:  "test",
			},
		},
		{
			name:      "create token with empty position",
			pos:       Position{},
			tokenType: STRING,
			literal:   "hello",
			want: Token{
				Type:     STRING,
				Position: Position{},
				Literal:  "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.pos, tt.tokenType, tt.literal)
			if got != tt.want {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
