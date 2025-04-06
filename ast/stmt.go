package ast

//go:generate go run ../tools/generator.go

type Expr interface {
	Accept(IASTVisitor) any
}

type Stmt interface{ Expr }

// expression     → literal
//                | unary
//                | binary
//                | grouping ;

// literal        → NUMBER | STRING | "true" | "false" | "nil" ;
// grouping       → "(" expression ")" ;
// unary          → ( "-" | "!" ) expression | expression ("++" | "--");
// binary         → expression operator expression ;
// operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
//                | "+"  | "-"  | "*" | "/" | "**" | "%" ;

// Name		Operators		Associates
// Equality		== !=		Left
// Comparison	> >= < <=	Left
// Term			- +			Left
// Factor		/ * %		Left
// Unary		! -			Right
