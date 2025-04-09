package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Toolnado/sludge/interpreter"
	"github.com/Toolnado/sludge/lexer"
	"github.com/Toolnado/sludge/parser"
)

func main() {
	l := lexer.New(strings.NewReader(`
		var a = 1
		if (a == 1){
			a = 10 
			print a
		}
	`))
	t := l.ScanTokens()
	p := parser.New(t)
	stmts, err := p.Parse()
	if err != nil {
		log.Println(err)
		return
	}
	i := interpreter.New()
	_, err = i.Interpret(stmts)
	if err != nil {
		fmt.Println(err)
	}
}
