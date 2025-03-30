// Package lexer implements lexical analysis for the Sludge programming language.
package lexer

import (
	"fmt"
	"io"
	"text/scanner"

	"github.com/Toolnado/sludge/token"
)

// Lexer performs lexical analysis of source code and breaks it into tokens.
// It uses the text/scanner package for basic scanning functionality.
type Lexer struct {
	scanner  scanner.Scanner // The underlying scanner
	hadError bool            // Indicates if any errors occurred during scanning
	tokens   []token.Token   // Slice of scanned tokens
	errors   []error         // Slice of errors encountered during scanning
}

// New creates and initializes a new Lexer instance.
// It takes an io.Reader as input which provides the source code to be scanned.
func New(r io.Reader) *Lexer {
	lexer := &Lexer{}
	lexer.scanner.Init(r)

	// Configure scanner mode
	lexer.scanner.Mode = scanner.ScanIdents |
		scanner.ScanFloats |
		scanner.ScanInts |
		scanner.ScanStrings |
		scanner.ScanRawStrings |
		scanner.ScanComments |
		scanner.SkipComments

	// Configure error handling
	lexer.scanner.Error = func(s *scanner.Scanner, msg string) {
		lexer.hadError = true
		lexer.errors = append(lexer.errors, fmt.Errorf("[%d:%d] --> scanner error: %s", s.Position.Line, s.Position.Column, msg))
	}
	return lexer
}

// Errors returns all errors encountered during lexical analysis.
func (l *Lexer) Errors() []error {
	return l.errors
}

// ScanTokens performs the lexical analysis of the entire input.
// It returns a slice of all tokens found in the source code.
func (l *Lexer) ScanTokens() []token.Token {
	for !l.isAtEnd() {
		t := l.scan()
		l.addToken(t)
	}

	// Add EOF token only if it doesn't exist yet
	if len(l.tokens) == 0 || l.tokens[len(l.tokens)-1].Type != token.EOF {
		l.addToken(token.Token{
			Type:     token.EOF,
			Position: l.position(),
			Literal:  "",
		})
	}

	return l.tokens
}
