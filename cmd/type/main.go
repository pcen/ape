package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/types"
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

	checker := types.NewChecker()
	typeMap, symbolTable := checker.Check(prog)
	fmt.Println("expression node types:")
	for key, val := range typeMap {
		fmt.Printf("%+v: %v\n", key, val)
	}
	fmt.Println("\nsymbol table:")
	for key, val := range symbolTable {
		fmt.Printf("%v: %v\n", key, val)
	}
}
