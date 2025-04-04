package parser

import (
	"log"

	"github.com/Toolnado/sludge/token"
)

func (p *Parser) match(types ...token.TokenType) bool {
	for _, _type := range types {
		if p.check(_type) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(_type token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == _type
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(_type token.TokenType, message string) token.Token {
	if p.check(_type) {
		return p.advance()
	}
	log.Println("[ERROR]:", p.peek(), message)
	return p.previous()
}
