// Package lexer implements lexical analysis for the Sludge programming language.
package lexer

import (
	"fmt"
	"strconv"
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
		l.addError(fmt.Sprintf("[%d:%d] --> Unexpected character sequence: %s", pos.Line, pos.Column, text))
	}

	return l.createToken(ttype)
}

// scanFloat processes floating-point numbers and returns a FLOAT token.
func (l *Lexer) scanFloat() token.Token {
	return l.createToken(token.FLOAT)
}

func (l *Lexer) createToken(tpe token.TokenType) token.Token {
	return token.New(l.position(), tpe, l.text(), l.parseLiteral(tpe, l.text()))
}

func (l *Lexer) parseLiteral(t token.TokenType, text string) any {
	switch t {
	case token.INTEGER:
		if i, err := strconv.ParseInt(text, 10, 64); err == nil {
			return i
		}
	case token.FLOAT:
		if f, err := strconv.ParseFloat(text, 64); err == nil {
			return f
		}
	case token.NULL:
		return nil
	case token.STRING, token.RAW_STRING:
		unquoted, err := strconv.Unquote(text)
		if err != nil {
			return text
		}
		return unquoted
	}
	return text
}

// scanInteger processes integer numbers and returns an INTEGER token.
func (l *Lexer) scanInteger() token.Token {
	return l.createToken(token.INTEGER)
}

// scanString processes double-quoted string literals.
// String interpolation is not supported in double-quoted strings.
func (l *Lexer) scanString() token.Token {
	return l.createToken(token.STRING)
}

// scanRawString processes raw strings (enclosed in backticks) with support for interpolation.
func (l *Lexer) scanRawString() token.Token {
	text := l.text()
	if !l.hasInterpolation(text) {
		return l.createToken(token.RAW_STRING)
	}

	// Process interpolation
	tokens := l.processInterpolations(text)
	// Add all tokens except the last one
	for i := 0; i < len(tokens)-1; i++ {
		l.addToken(tokens[i])
	}
	// Return the last token
	return tokens[len(tokens)-1]
}

// hasInterpolation checks if the given text contains any interpolation expressions (${})
// or template expressions (@{}). Returns true if either type of expression is found.
func (l *Lexer) hasInterpolation(text string) bool {
	return strings.Contains(text, "${") || strings.Contains(text, "@{")
}

// processInterpolations breaks down a raw string containing interpolations into a sequence of tokens.
func (l *Lexer) processInterpolations(text string) []token.Token {
	var tokens []token.Token
	currentText := text[1 : len(text)-1] // Trim opening and closing quotes

	for {
		interpIndex := strings.Index(currentText, "${")
		macroIndex := strings.Index(currentText, "@{")

		if interpIndex == -1 && macroIndex == -1 {
			if currentText != "" {
				tokens = append(tokens, token.New(l.position(), token.RAW_STRING, currentText, currentText))
			}
			break
		}

		startIndex, tokenType := l.findNextInterpolation(interpIndex, macroIndex)

		if startIndex > 0 {
			tokens = append(tokens, token.New(l.position(), token.RAW_STRING, currentText[:startIndex], currentText[:startIndex]))
		}

		exprStart := startIndex + 2
		exprEnd := strings.Index(currentText[exprStart:], "}")
		if exprEnd == -1 {
			l.addError("Unclosed interpolation or template")
			break
		}
		exprEnd += exprStart

		tokens = append(tokens, token.New(l.position(), tokenType, currentText[exprStart:exprEnd], currentText[exprStart:exprEnd]))
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
	return l.createToken(ttype)
}

// scan performs scanning of the next token from the input stream.
func (l *Lexer) scan() token.Token {
	ch := l.advance()
	pos := l.position()

	if ch == scanner.EOF {
		return l.createToken(token.EOF)
	}

	switch ch {
	case scanner.Float:
		return l.scanFloat()
	case scanner.Int:
		return l.scanInteger()
	case '\'':
		return l.scanSingleQuotedString()
	case scanner.String:
		return l.scanString()
	case scanner.RawString:
		t := l.scanRawString()
		// Не добавляем токен здесь, так как он уже добавлен в scanRawString
		return t
	case scanner.Ident:
		return l.scanIdentifier()
	default:
		if ch < 0 {
			l.addError(fmt.Sprintf("Unexpected character: %v", ch))
			return token.New(pos, token.ILLEGAL, string(ch), string(ch))
		}
		return l.scanOperator(ch, pos)
	}
}

// scanSingleQuotedString processes single-quoted string literals.
func (l *Lexer) scanSingleQuotedString() token.Token {
	var builder strings.Builder
	pos := l.position()

	for !l.isAtEnd() {
		ch := l.next()
		if ch == '\'' {
			break
		}
		if ch == '\\' {
			// Handle escaped characters
			next := l.next()
			if next == '\'' {
				builder.WriteRune('\'')
			} else {
				builder.WriteRune('\\')
				builder.WriteRune(next)
			}
		} else {
			builder.WriteRune(ch)
		}
	}

	return token.New(pos, token.STRING, builder.String(), builder.String())
}
