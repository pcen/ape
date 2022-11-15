package tests

import (
	"fmt"
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

func parseHelper(source string) ([]ast.Node, []ape.ParseError) {
	tokens := ape.NewLexer().LexString(source)
	fmt.Println(tokens)
	parser := ape.NewParser(tokens)
	node := parser.Program()
	errors, _ := parser.Errors()
	return node, errors
}

// func TestBadTypedDeclarations(t *testing.T) {
// 	source := `
// 	val a = 5
// 	val b int = 6
// 	`
// 	node, errors := parse(source)
// 	ast.PrintSlice(node)
// 	fmt.Println("errors:")
// 	for _, err := range errors {
// 		fmt.Println(err)
// 	}
// }

func TestBadGroupExpr(t *testing.T) {
	source := `
	val a int = 1 + (2 * 3
	var b int = 1 +
	int c = 2
	`
	node, errors := parseHelper(source)
	ast.PrintSlice(node)
	fmt.Println("errors:")
	for _, err := range errors {
		fmt.Println(err)
	}
}
