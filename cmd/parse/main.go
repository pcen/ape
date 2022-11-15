package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("supply file to parse")
	}
	file := os.Args[1]
	lexer := ape.NewLexer()
	tokens := lexer.LexFile(file)

	parser := ape.NewParser(tokens)
	prog := parser.Program()
	fmt.Println("ast:")
	ast.PrettyPrint(prog)

	if errors, ok := parser.Errors(); ok {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
