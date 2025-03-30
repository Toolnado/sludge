package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Toolnado/sludge/lexer"
	"github.com/Toolnado/sludge/token"
)

func main() {
	var source string
	var filename string

	flag.StringVar(&source, "code", "", "Исходный код для лексического анализа")
	flag.StringVar(&filename, "file", "", "Файл с исходным кодом")
	flag.Parse()

	var reader io.Reader

	switch {
	case filename != "":
		// Чтение из файла
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Ошибка при открытии файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file

	case source != "":
		// Чтение из аргумента командной строки
		reader = strings.NewReader(source)

	default:
		// Чтение из stdin
		fmt.Println("Введите код (Ctrl+D для завершения):")
		reader = bufio.NewReader(os.Stdin)
	}

	l := lexer.New(reader)
	tokens := l.ScanTokens()

	printTokens(tokens)
	printErrors(l.Errors())
}

func printTokens(tokens []token.Token) {
	for _, token := range tokens {
		fmt.Printf("Тип: %-15s Значение: %s\n", token.Type, token.Literal)
	}
}

func printErrors(errors []error) {
	if len(errors) > 0 {
		fmt.Println("\nОшибки:")
		for _, err := range errors {
			fmt.Printf("- %v\n", err)
		}
	}
}
