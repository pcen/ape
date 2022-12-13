package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/c"
)

const (
	output = "./ape.i"
)

/*
 compiles demo scripts to bytecode
*/

func printMap[T comparable](m map[T]struct{}) {
	for k := range m {
		fmt.Printf("%v\n", k)
	}
}

func compile(script string) (string, error) {
	lexer := ape.NewLexer()
	fmt.Println("lexing...")
	tokens := lexer.LexFile(script)
	fmt.Println("parsing...")
	parser := ape.NewParser(tokens)
	stmts := parser.Demo()
	if errs, hasErrs := parser.Errors(); hasErrs {
		return "", fmt.Errorf("parser error(s): %v", errs)
	}
	fmt.Println("generating code...")
	code := c.GenerateCode(stmts)

	fmt.Println("identifiers:")
	printMap(code.Idents)

	src := code.Code.String()
	fmt.Printf("src:\n%v\n", src)

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.WriteString(f, src)
	return output, nil
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
	exec.Command("gcc", "ape.i").Run()
}
