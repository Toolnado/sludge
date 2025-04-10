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
		var i = 0;
		while (i < 10){
			print i
			i = i + 1
		};
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
