// Package lexer implements lexical analysis for the Sludge programming language.
package lexer

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/Toolnado/sludge/token"
)

// operators defines valid two-character operator combinations.
// The first character is the main operator, the second is a possible continuation.
var operators = map[rune]map[rune]struct{}{
	'=': {'=': {}, '>': {}},
	'!': {'=': {}},
	'<': {'=': {}},
	'>': {'=': {}},
	'&': {'&': {}},
	'|': {'|': {}},
	'+': {'=': {}},
	'-': {'=': {}},
	'*': {'=': {}},
	'/': {'=': {}},
	'%': {'=': {}},
	'(': {},
	')': {},
	'{': {},
	'}': {},
	'[': {},
	']': {},
	',': {},
	'.': {},
	';': {},
	':': {},
}

// operatorsMap maps string representations of operators to their token types.
var operatorsMap = map[string]token.TokenType{
	"||": token.OR,
	"&&": token.AND,
	"==": token.EQUAL_EQUAL,
	"!=": token.BANG_EQUAL,
	"<=": token.LESS_EQUAL,
	">=": token.GREATER_EQUAL,
	"<":  token.LESS,
	">":  token.GREATER,
	"+":  token.PLUS,
	"-":  token.MINUS,
	"*":  token.STAR,
	"/":  token.SLASH,
	"%":  token.PERCENT,
	"(":  token.LEFT_PAREN,
	")":  token.RIGHT_PAREN,
	"{":  token.LEFT_BRACE,
	"}":  token.RIGHT_BRACE,
	"[":  token.LEFT_BRACKET,
	"]":  token.RIGHT_BRACKET,
	",":  token.COMMA,
	".":  token.DOT,
	";":  token.SEMICOLON,
	":":  token.COLON,
	"=>": token.ARROW,
	"+=": token.PLUS_EQUAL,
	"-=": token.MINUS_EQUAL,
	"*=": token.STAR_EQUAL,
	"/=": token.SLASH_EQUAL,
	"%=": token.PERCENT_EQUAL,
	"!":  token.BANG,
	"=":  token.EQUAL,
}

// scanOperator processes operators and returns the corresponding token.
// Supports both single and compound operators.
func (l *Lexer) scanOperator(ch rune, pos token.Position) token.Token {
	ttype := token.ILLEGAL
	text := string(ch)
	if op, ok := operators[ch]; ok {
		char := l.peek()
		if _, ok := op[char]; ok {
			l.next()
			text += string(char)
		}
	}

	if typ, ok := operatorsMap[text]; ok {
		ttype = typ
	} else {
		l.addError(fmt.Sprintf("Unexpected character sequence: %s", text))
	}

	return token.New(pos, ttype, text)
}

// scanFloat processes floating-point numbers and returns a FLOAT token.
func (l *Lexer) scanFloat() token.Token {
	return token.New(l.position(), token.FLOAT, l.text())
}

// scanInteger processes integer numbers and returns an INTEGER token.
func (l *Lexer) scanInteger() token.Token {
	return token.New(l.position(), token.INTEGER, l.text())
}

// scanString processes string literals and returns a STRING token.
func (l *Lexer) scanString() token.Token {
	return token.New(l.position(), token.STRING, l.substring())
}

// scanRawString processes raw strings with support for interpolation (${}) and templates (@{}).
// Returns the last processed token from the sequence.
func (l *Lexer) scanRawString() token.Token {
	text := l.text()
	if !l.hasInterpolation(text) {
		return token.New(l.position(), token.RAW_STRING, text)
	}

	tokens := l.processInterpolations(text)
	for _, t := range tokens {
		l.addToken(t)
	}
	return tokens[len(tokens)-1]
}

// hasInterpolation checks if the given text contains any interpolation expressions (${})
// or template expressions (@{}). Returns true if either type of expression is found.
func (l *Lexer) hasInterpolation(text string) bool {
	return strings.Contains(text, "${") || strings.Contains(text, "@{")
}

// processInterpolations breaks down a raw string containing interpolations into a sequence of tokens.
// It handles both string interpolation (${}) and template (@{}) expressions.
// The input text should be a raw string including quotes. The function:
// - Strips the surrounding quotes
// - Identifies interpolation/template expressions
// - Creates appropriate tokens for the literal text segments and expressions
// - Returns the complete sequence of tokens representing the string
func (l *Lexer) processInterpolations(text string) []token.Token {
	var tokens []token.Token
	currentText := text[1 : len(text)-1] // Trim opening and closing quotes

	for {
		interpIndex := strings.Index(currentText, "${")
		macroIndex := strings.Index(currentText, "@{")

		if interpIndex == -1 && macroIndex == -1 {
			if currentText != "" {
				tokens = append(tokens, token.New(l.position(), token.RAW_STRING, currentText))
			}
			break
		}

		startIndex, tokenType := l.findNextInterpolation(interpIndex, macroIndex)

		if startIndex > 0 {
			tokens = append(tokens, token.New(l.position(), token.RAW_STRING, currentText[:startIndex]))
		}

		exprStart := startIndex + 2
		exprEnd := strings.Index(currentText[exprStart:], "}")
		if exprEnd == -1 {
			l.addError("Unclosed interpolation or template")
			break
		}
		exprEnd += exprStart

		tokens = append(tokens, token.New(l.position(), tokenType, currentText[exprStart:exprEnd]))
		currentText = currentText[exprEnd+1:]
	}

	return tokens
}

// findNextInterpolation determines which type of expression (interpolation or template)
// appears first in the text based on their indices. Returns:
// - The starting index of the first expression
// - The appropriate token type (INTERPOLATION for ${} or TEMPLATE for @{})
func (l *Lexer) findNextInterpolation(interpIndex, macroIndex int) (int, token.TokenType) {
	startIndex := -1
	tokenType := token.INTERPOLATION

	if interpIndex == -1 {
		startIndex = macroIndex
		tokenType = token.TEMPLATE
	} else if macroIndex == -1 {
		startIndex = interpIndex
	} else if macroIndex < interpIndex {
		startIndex = macroIndex
		tokenType = token.TEMPLATE
	} else {
		startIndex = interpIndex
	}

	return startIndex, tokenType
}

// scanIdentifier processes identifiers and keywords.
// Returns an IDENTIFIER token or corresponding keyword token.
func (l *Lexer) scanIdentifier() token.Token {
	ttype := token.IDENTIFIER
	if keyword, isKeyword := token.IsKeyword(l.text()); isKeyword {
		ttype = keyword
	}
	return token.New(l.position(), ttype, l.text())
}

// scan performs scanning of the next token from the input stream.
// Determines the token type and calls the appropriate handler.
func (l *Lexer) scan() token.Token {
	ch := l.advance()
	pos := l.position()

	if ch == scanner.EOF {
		return token.New(pos, token.EOF, "")
	}

	switch ch {
	case scanner.Float:
		return l.scanFloat()
	case scanner.Int:
		return l.scanInteger()
	case scanner.String, scanner.Char:
		return l.scanString()
	case scanner.RawString:
		return l.scanRawString()
	case scanner.Ident:
		return l.scanIdentifier()
	default:
		return l.scanOperator(ch, pos)
	}
}
