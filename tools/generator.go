package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	generateDefaultAST()
}

func generateAST(file string, exprs map[string][]string) {
	f, err := os.Create(file)
	if err != nil {
		log.Println("GenerageAST error:", err)
		return
	}

	fmt.Fprint(f, "package ast\n\n")
	fmt.Fprint(f, "import \"github.com/Toolnado/sludge/token\"\n\n")

	for expr, fields := range exprs {
		fmt.Fprint(f, "type ")
		fmt.Fprint(f, expr)
		fmt.Fprint(f, " struct {\n")
		for _, field := range fields {
			fmt.Fprint(f, "	")
			fmt.Fprint(f, field)
			fmt.Fprint(f, "\n")
		}
		fmt.Fprint(f, "}\n\n")
		generateConstructor(expr, fields, f)
		generateVisitor(expr, f)
	}
}

func generateVisitor(expr string, f io.Writer) {
	r := string(strings.ToLower(expr)[0])
	if r == "v" {
		r += string(strings.ToLower(expr)[1])
	}

	fmt.Fprintf(f, "func (%s *%s) Accept(v IASTVisitor) (any, error) { return v.Visit%s(%s)}\n\n", r, expr, expr, r)
}

func generateConstructor(expr string, fields []string, f io.Writer) {
	fmt.Fprint(f, "func New")
	fmt.Fprint(f, expr)
	fmt.Fprint(f, "(")
	for i, field := range fields {
		if i != 0 {
			fmt.Fprint(f, ", ")
		}
		fmt.Fprint(f, field)
	}
	fmt.Fprint(f, ") *")
	fmt.Fprint(f, expr)
	fmt.Fprint(f, " {\n")
	fmt.Fprint(f, "	return &")
	fmt.Fprint(f, expr)
	fmt.Fprint(f, "{\n")
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		fmt.Fprint(f, "		")
		fmt.Fprint(f, name)
		fmt.Fprint(f, ": ")
		fmt.Fprint(f, name)
		fmt.Fprint(f, ",\n")
	}
	fmt.Fprint(f, "	}\n}\n\n")
}

func generateDefaultAST() {
	generateAST("expr.go", map[string][]string{
		"BinaryExpr": {
			"Left Expr",
			"Operator token.Token",
			"Right Expr",
		},
		"UnaryExpr": {
			"Operator token.Token",
			"Right Expr",
		},
		"LiteralExpr":  {"Value any"},
		"GroupingExpr": {"Expession Expr"},
		"VariableExpr": {"Name token.Token"},
		"AssignExpr": {
			"Name token.Token",
			"Value Expr",
		},
		"LogicalExpr": {
			"Left Expr",
			"Operator token.Token",
			"Right Expr",
		},
		"CallExpr": {
			"Callee Expr",
			"Paren token.Token",
			"Arguments []Expr",
		},
	})

	generateAST("stmt.go", map[string][]string{
		"PrintStmt": {"Expession Expr"},
		"ExprStmt":  {"Expession Expr"},
		"VarStmt": {
			"Name token.Token",
			"Initializer Expr",
		},
		"BlockStmt": {"Statements []Stmt"},
		"IfStmt": {
			"Condition Expr",
			"ThenBranch Stmt",
			"ElseBranch Stmt",
		},
		"WhileStmt": {
			"Condition Expr",
			"Body Stmt",
		},
	})
}
