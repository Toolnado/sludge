package token

// Keywords maps string literals to their corresponding TokenType for all keywords
// in the Sludge programming language.
var Keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"const":    CONST,
	"var":      VAR,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"return":   RETURN,
	"break":    BREAK,
	"continue": CONTINUE,
	"null":     NULL,
	"import":   IMPORT,
}

// IsKeyword checks if a given string is a keyword in the language.
// Returns the corresponding TokenType and true if it is a keyword,
// or an empty TokenType and false if it is not.
func IsKeyword(s string) (TokenType, bool) {
	typ, ok := Keywords[s]
	return typ, ok
}
