package token

// Position represents a location in a source file, including the filename,
// byte offset from start of file, line number, and column number.
type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

// Token represents a lexical token with its position in source code,
// type classification, and literal string value.
type Token struct {
	Position Position
	Type     TokenType
	Literal  any
	Lexeme   string
}

// New creates a new Token with the given position, type and literal value.
// Returns a Token struct initialized with the provided values.
func New(
	position Position,
	t TokenType,
	lexeme string,
	literal any,
) Token {
	return Token{
		Position: position,
		Type:     t,
		Lexeme:   lexeme,
		Literal:  literal,
	}
}
