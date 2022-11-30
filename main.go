package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
)

/*
 compiles demo scripts to bytecode
*/

func compile(script string) (string, error) {
	lexer := ape.NewLexer()
	tokens := lexer.LexFile(script)
	parser := ape.NewParser(tokens)
	decls := parser.Demo()
	if errs, hasErrs := parser.Errors(); hasErrs {
		return "", fmt.Errorf("parser error(s): %v\n", errs)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("supply name of ape script as argument")
		os.Exit(1)
	}
	script := os.Args[1]
	compiled, err := compile(script)
	if err != nil {
		fmt.Printf("error compiling %v: %v\n", script, err.Error())
		os.Exit(1)
	}
	fmt.Printf("compiled %v to %v\n", script, compiled)

}
