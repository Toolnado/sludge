package parser

import (
	"github.com/Toolnado/sludge/ast"
	"github.com/Toolnado/sludge/token"
)

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → remainder ( ( "-" | "+" ) remainder )* ;
// remainder      → factor ( % factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary  | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// expression     → equality ;
func (p *Parser) expression() ast.Expr {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() ast.Expr {
	expr := p.term()
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

// term           → remainder ( ( "-" | "+" ) remainder )* ;
func (p *Parser) term() ast.Expr {
	expr := p.remainder()
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.remainder()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

// remainder      → factor ( % factor )* ;
func (p *Parser) remainder() ast.Expr {
	expr := p.factor()
	for p.match(token.PERCENT) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() ast.Expr {
	expr := p.unary()
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr
}

// unary          → ( "!" | "-" ) unary  | primary ;
func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}

	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() ast.Expr {
	switch {
	case p.match(token.FALSE):
		return ast.NewLiteral(p.previous().Literal)
	case p.match(token.TRUE):
		return ast.NewLiteral(p.previous().Literal)
	case p.match(token.NULL):
		return ast.NewLiteral("null")
	case p.match(token.STRING, token.RAW_STRING, token.INTEGER):
		return ast.NewLiteral((p.previous().Literal))
	case p.match(token.LEFT_PAREN):
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.NewGrouping(expr)
	default:
		return &ast.Literal{}
	}
}
