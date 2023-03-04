package tests

import (
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

const (
	prog = `{
	foo() @reverse bar()
	login() @reverse logout()
}`
)

func TestAnnotationParsing(t *testing.T) {
	tokens := ape.NewLexer().LexString(prog)
	block := ape.NewParser(tokens).BlockStmt()
	for _, stmt := range block.Content {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			t.Fatal("did not get expression statement")
		}
		if len(exprStmt.Annotations) == 0 {
			t.Fatal("did not get annotations")
		}
	}
}
