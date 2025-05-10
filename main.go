package main

import (
	"log"
	"strings"

	"github.com/Toolnado/sludge/interpreter"
	"github.com/Toolnado/sludge/lexer"
	"github.com/Toolnado/sludge/parser"
)

func main() {
	l := lexer.New(strings.NewReader(`
		var a = 0;
		var temp;

		for (var b = 1; a < 10000; b = temp + b) {
			print a;
			temp = a;
			a = b;
		}

		var time = clock()
		print(time)
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
		log.Println(err)
	}
}
