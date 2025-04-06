package parser

import (
	"github.com/Toolnado/sludge/token"
)

// match checks if the current token matches any of the provided types.
// If a match is found, it advances the parser and returns true.
// Otherwise, it returns false.
func (p *Parser) match(types ...token.TokenType) bool {
	for _, _type := range types {
		if p.check(_type) {
			p.advance()
			return true
		}
	}
	return false
}

// check returns true if the current token matches the provided type.
// It returns false if the parser is at the end of the input or the types don't match.
func (p *Parser) check(_type token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == _type
}

// advance moves the parser to the next token if it's not at the end.
// It returns the previous token.
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd returns true if the parser has reached the end of the input.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// peek returns the current token without advancing the parser.
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

// previous returns the last token that was consumed by the parser.
func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

// consume checks if the current token matches the expected type.
// If so, it advances the parser and returns the token.
// If not, it logs an error with the current token and message, then returns the previous token.
func (p *Parser) consume(t token.TokenType, message string) (token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return token.Token{}, NewError(p.peek(), message)
}
