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

func generateAST(exprs map[string][]string) {
	f, err := os.Create("ast.go")
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
	r := strings.ToLower(expr)[0]
	fmt.Fprintf(f, "func (%c *%s) Accept(v IASTVisitor) {v.Visit(%c)}\n\n", r, expr, r)
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
	generateAST(map[string][]string{
		"Binary": {
			"Left Expr",
			"Operator token.Token",
			"Right Expr",
		},
		"Unary": {
			"Operator token.Token",
			"Right Expr",
		},
		"Literal": {
			"Value string",
		},
		"Grouping": {
			"Expession Expr",
		},
	})
}
