package tests

import (
	"fmt"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

func Parse(source string) (*ast.File, []ape.ParseError) {
	tokens := ape.NewLexer().LexString(source)
	fmt.Println(tokens)
	parser := ape.NewParser(tokens)
	node := parser.File()
	errors, _ := parser.Errors()
	return node, errors
}
