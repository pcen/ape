package tests

import (
	"fmt"
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

func parse(source string) ([]ast.Node, []ape.ParseError) {
	tokens := ape.NewLexer().LexString(source)
	fmt.Println(tokens)
	parser := ape.NewParser(tokens)
	node := parser.Program()
	errors, _ := parser.Errors()
	return node, errors
}

func TestParsing(t *testing.T) {
	source := `
	func main(a int) {
		var b int = 10
		var c int = 20
		return b * c
	}
	`
	prog, errs := parse(source)
	ast.PrintSlice(prog)
	fmt.Println(errs)
}
