package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
)

/*
 compiles demo scripts to bytecode
*/

func printMap[T comparable](m map[T]uint8) {
	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}
}

func compile(script string) (string, error) {
	lexer := ape.NewLexer()
	tokens := lexer.LexFile(script)
	parser := ape.NewParser(tokens)
	stmts := parser.Demo()
	if errs, hasErrs := parser.Errors(); hasErrs {
		return "", fmt.Errorf("parser error(s): %v", errs)
	}
	code := ape.GenerateCode(stmts)

	fmt.Println("identifiers:")
	printMap(code.IdentIdx)

	fmt.Println("\nliterals:")
	printMap(code.LitIdx)

	fmt.Println("\nop codes:")

	for _, op := range code.Code {
		fmt.Println(op.String())
	}

	lits := make([]int32, len(code.LitIdx))
	for v, idx := range code.LitIdx {
		lits[idx] = v
	}

	f, err := os.Create("./out.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	numLits := int32(len(lits))

	binary.Write(f, binary.LittleEndian, numLits)
	for _, lit := range lits {
		binary.Write(f, binary.LittleEndian, lit)
	}
	for _, oc := range code.Code {
		binary.Write(f, binary.LittleEndian, byte(oc))
	}

	return "", nil
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
