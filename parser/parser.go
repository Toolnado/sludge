package parser

import (
	"github.com/Toolnado/sludge/ast"
	"github.com/Toolnado/sludge/token"
)

// Parser parses a sequence of tokens into an abstract syntax tree (AST).
// It implements a recursive descent parser based on the following grammar:
//
// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → remainder ( ( "-" | "+" ) remainder )* ;
// remainder      → factor ( "%" factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary  | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
type Parser struct {
	tokens   []token.Token // List of tokens to parse
	hadError bool          // Indicates if a parsing error has occurred
	current  int           // Index of the current token
}

// New creates a new parser from a slice of tokens.
func New(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// HadError returns true if the parser has encountered a syntax error.
func (p *Parser) HadError() bool {
	return p.hadError
}

// Parse starts the parsing process and returns the root expression node or an error.
func (p *Parser) Parse() (ast.Expr, error) {
	expr, err := p.expression()
	if err != nil {
		p.hadError = true
		p.synchronize()
		return nil, err
	}
	return expr, nil
}

// synchronize attempts to recover from a parsing error by advancing
// through tokens until it finds a likely statement boundary.
func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.FUNCTION, token.VAR, token.FOR, token.IF,
			token.WHILE, token.LET, token.CONST, token.RETURN:
			return
		}
		p.advance()
	}
}

// expression → equality ;
func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

// term → remainder ( ( "-" | "+" ) remainder )* ;
func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.remainder()
	if err != nil {
		return nil, err
	}
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.remainder()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

// remainder → factor ( "%" factor )* ;
func (p *Parser) remainder() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.PERCENT) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}
	return expr, nil
}

// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnary(operator, right), nil
	}
	return p.primary()
}

// primary → literal | "(" expression ")" ;
func (p *Parser) primary() (ast.Expr, error) {
	switch {
	case p.match(token.FALSE):
		return ast.NewLiteral(p.previous().Literal), nil
	case p.match(token.TRUE):
		return ast.NewLiteral(p.previous().Literal), nil
	case p.match(token.NULL):
		return ast.NewLiteral("null"), nil
	case p.match(token.STRING, token.RAW_STRING, token.INTEGER):
		return ast.NewLiteral(p.previous().Literal), nil
	case p.match(token.LEFT_PAREN):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression"); err != nil {
			return nil, err
		}
		return ast.NewGrouping(expr), nil
	default:
		return nil, NewError(p.peek(), "Expect expression.")
	}
}
