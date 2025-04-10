package parser

import (
	"log"

	"github.com/Toolnado/sludge/ast"
	"github.com/Toolnado/sludge/token"
)

// Parser parses a sequence of tokens into an abstract syntax tree (AST).
// It implements a recursive descent parser based on the following grammar:
// program        → declaration* EOF ;
// declaration    → varDecl
//                | statement ;

// statement      → exprStmt ;
// 				  | forStmt ;
// 				  | ifStmt ;
//                | printStmt ;
// 				  | whileStmt ;
//				  | block ;

// forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//                  expression? ";"
//                  expression? ")" statement ;

// whileStmt      → "while" "(" expression ")" statement ;

// ifStmt         → "if" "(" expression ")" statement
//                ( "else" statement )? ;

// block          → "{" declaration* "}" ;

// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment
//                | logic_or ;

// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
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
func (p *Parser) Parse() ([]ast.Stmt, error) {
	decls := []ast.Stmt{}
	for !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			log.Println(err)
			p.synchronize()
		} else {
			decls = append(decls, decl)
		}
	}
	return decls, nil
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "expect variable name.")
	if err != nil {
		return nil, NewError(p.previous(), err.Error())
	}
	var initializer ast.Expr
	if p.match(token.EQUAL) {
		expr, err := p.expression()
		if err != nil {
			return nil, NewError(p.previous(), err.Error())
		}
		initializer = expr
	}

	p.consume(token.SEMICOLON, "expect ';' after variable declaration.")
	return ast.NewVarStmt(name, initializer), nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, NewError(p.previous(), err.Error())
	}
	p.consume(token.SEMICOLON, "expect ';' after value.")
	return ast.NewPrintStmt(value), nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, NewError(p.previous(), err.Error())
	}
	p.consume(token.SEMICOLON, "expect ';' after value")
	return ast.NewExprStmt(value), nil
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

// expression → or;
func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, NewError(p.peek(), err.Error())
	}
	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, NewError(p.peek(), err.Error())
		}
		if v, ok := expr.(*ast.VariableExpr); ok {
			name := v.Name
			return ast.NewAssignExpr(name, value), nil
		}

		return nil, NewError(equals, "invalid assignment target")
	}

	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, NewError(p.peek(), err.Error())
	}
	if p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, NewError(p.peek(), err.Error())
		}
		return ast.NewLogicalExpr(expr, operator, right), nil
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, NewError(p.peek(), err.Error())
	}
	if p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, NewError(p.peek(), err.Error())
		}
		return ast.NewLogicalExpr(expr, operator, right), nil
	}

	return expr, nil
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
		expr = ast.NewBinaryExpr(expr, operator, right)
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
		expr = ast.NewBinaryExpr(expr, operator, right)
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
		expr = ast.NewBinaryExpr(expr, operator, right)
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
		expr = ast.NewBinaryExpr(expr, operator, right)
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
		expr = ast.NewBinaryExpr(expr, operator, right)
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
		return ast.NewUnaryExpr(operator, right), nil
	}
	return p.primary()
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	_, err := p.consume(token.RIGHT_BRACE, "expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "expect '(' after 'if'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "expect ')' after if condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		elseBranch = stmt
	}
	return ast.NewIfStmt(condition, thenBranch, elseBranch), nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	p.consume(token.LEFT_PAREN, "expect '(' after 'while'")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "expect ')' after condition")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return ast.NewWhileStmt(condition, body), nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR, token.LET) {
		i, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}
		initializer = i
	} else {
		i, err := p.expressionStatement()
		if err != nil {
			return nil, err
		}
		initializer = i
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		c, err := p.expression()
		if err != nil {
			return nil, err
		}
		condition = c
	}

	_, err = p.consume(token.SEMICOLON, "expect ';' after loop condition")
	if err != nil {
		return nil, err
	}
	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		i, err := p.expression()
		if err != nil {
			return nil, err
		}
		increment = i
	}
	_, err = p.consume(token.RIGHT_PAREN, "expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = ast.NewBlockStmt([]ast.Stmt{body, ast.NewExprStmt(increment)})
	}

	if condition == nil {
		condition = ast.NewLiteralExpr(true)
	}
	body = ast.NewWhileStmt(condition, body)

	if initializer != nil {
		body = ast.NewBlockStmt([]ast.Stmt{initializer, body})
	}

	return body, nil

}

// primary → literal | "(" expression ")" ;
func (p *Parser) primary() (ast.Expr, error) {
	switch {
	case p.match(token.FALSE):
		return ast.NewLiteralExpr(p.previous().Literal), nil
	case p.match(token.TRUE):
		return ast.NewLiteralExpr(p.previous().Literal), nil
	case p.match(token.NULL):
		return ast.NewLiteralExpr(nil), nil
	case p.match(token.STRING, token.RAW_STRING, token.INTEGER):
		return ast.NewLiteralExpr(p.previous().Literal), nil
	case p.match(token.IDENTIFIER):
		return ast.NewVariableExpr(p.previous()), nil
	case p.match(token.LEFT_BRACE):
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return ast.NewBlockStmt(stmts), nil
	case p.match(token.LEFT_PAREN):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(token.RIGHT_PAREN, "expect ')' after expression"); err != nil {
			return nil, err
		}
		return ast.NewGroupingExpr(expr), nil
	case p.match(token.WHILE):
		return p.whileStatement()
	case p.match(token.FOR):
		return p.forStatement()
	case p.match(token.IF):
		return p.ifStatement()
	default:
		return nil, NewError(p.peek(), "expect expression")
	}
}
