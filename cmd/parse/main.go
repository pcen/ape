package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("supply file to parse")
	}
	file := os.Args[1]
	lexer := ape.NewLexer()
	tokens := lexer.LexFile(file)
	for _, t := range tokens {
		fmt.Printf("%v: %v\n", t.Type, t.Lexeme)
	}

	parser := ape.NewParser(tokens)
	expr := parser.Program()
	fmt.Println("ast:")
	fmt.Println(expr.StmtStr())
}