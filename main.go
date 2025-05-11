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
		function sayHi(first, last) {
			print "Hi, " + first + " " + last + "!";
		}

		sayHi("Dear", "Reader");
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
