package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("supply file to lex")
	}
	file := os.Args[1]
	lexer := ape.NewLexer()
	tokens := lexer.LexFile(file)
	fmt.Printf("tokens in %v\n", file)
	for _, tok := range tokens {
		fmt.Printf("%v\n", tok)
	}
}
