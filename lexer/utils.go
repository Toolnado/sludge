// Package lexer implements lexical analysis utilities for the Sludge programming language.
package lexer

import (
	"fmt"
	"text/scanner"

	"github.com/Toolnado/sludge/token"
)

// stringPosition returns a formatted string containing the current position in the source code
// in the format "filename:line:column". If no filename is set, returns "<input>:line:column".
func (l *Lexer) stringPosition() string {
	filename := l.scanner.Position.Filename
	if filename == "" {
		filename = "<input>"
	}
	return fmt.Sprintf("%s:%d:%d", filename, l.scanner.Position.Line, l.scanner.Position.Column)
}

// peek returns the next rune in the input without advancing the scanner.
func (l *Lexer) peek() rune {
	return l.scanner.Peek()
}

// next advances the scanner and returns the next rune.
func (l *Lexer) next() rune {
	return l.scanner.Next()
}

// advance scans and returns the next token from the input.
func (l *Lexer) advance() rune {
	return l.scanner.Scan()
}

// isAtEnd checks if the scanner has reached the end of input.
func (l *Lexer) isAtEnd() bool {
	return l.peek() == scanner.EOF
}

// position returns a token.Position struct containing the current scanner position information.
func (l *Lexer) position() token.Position {
	return token.Position{
		Filename: l.scanner.Position.Filename,
		Line:     l.scanner.Position.Line,
		Column:   l.scanner.Position.Column,
	}
}

// text returns the string value of the current token.
func (l *Lexer) text() string {
	return l.scanner.TokenText()
}

// addToken appends a new token to the lexer's token list.
func (l *Lexer) addToken(t token.Token) {
	l.tokens = append(l.tokens, t)
}

// addError adds a new error to the lexer's error list and sets the error flag.
// The error message includes the current position in the source code.
func (l *Lexer) addError(msg string) {
	if !l.hadError {
		l.hadError = true
	}
	l.errors = append(l.errors, fmt.Errorf("[%s] --> lexer error: %s", l.stringPosition(), msg))
}

// substring returns the content of a string literal without the surrounding quotes.
// Returns the original text if it's empty or too short to trim.
func (l *Lexer) substring() string {
	text := l.text()
	if text == "" || len(text) < 2 {
		return text
	}
	return text[1 : len(text)-1]
}
