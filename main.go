package main

import (
	"log"
	"strings"

	"github.com/Toolnado/sludge/interpreter"
	"github.com/Toolnado/sludge/lexer"
	"github.com/Toolnado/sludge/parser"
)

func main() {
	l := lexer.New(strings.NewReader("2+3*10/8-2-4+\"5\""))
	t := l.ScanTokens()
	p := parser.New(t)
	expr, err := p.Parse()
	if err != nil {
		log.Println(err)
		return
	}
	i := interpreter.New()
	v, err := i.Interpret(expr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("result:", v)
}
