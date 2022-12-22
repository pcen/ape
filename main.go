package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/c"
)

/*
 compiles demo scripts to bytecode
*/

func printMap[T comparable](m map[T]struct{}) {
	for k := range m {
		fmt.Printf("%v\n", k)
	}
}

func writeCode(path string, sb *strings.Builder) (string, error) {
	base := filepath.Base(path)
	suffix := filepath.Ext(base)
	base = strings.TrimSuffix(base, suffix)
	output := fmt.Sprintf("./out/%v.i", base)
	return output, os.WriteFile(output, []byte(sb.String()), 0664)
}

func compile(path string) (string, error) {
	lexer := ape.NewLexer()
	fmt.Println("lexing...")
	tokens := lexer.LexFile(path)
	fmt.Println("parsing...")
	parser := ape.NewParser(tokens)
	file := parser.File()
	if errs, hasErrs := parser.Errors(); hasErrs {
		return "", fmt.Errorf("parser error(s): %v", errs)
	}
	fmt.Println("generating code...")
	code := c.GenerateCode(file.Ast)

	fmt.Println("identifiers:")
	printMap(code.Idents)

	return writeCode(path, &code.Code)
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

	_, err = exec.Command("gcc", compiled, "-o", "bin").CombinedOutput()
	if err != nil {
		fmt.Printf("error compiling: %v\n", err.Error())
	}
}
