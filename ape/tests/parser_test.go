package tests

import (
	"fmt"
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

func parse(source string) ([]ast.Node, []ape.ParseError) {
	tokens := ape.NewLexer().LexString(source)
	parser := ape.NewParser(tokens)
	node := parser.Program()
	errors, _ := parser.Errors()
	return node, errors
}

func TestParsing(t *testing.T) {
	source := `
	class foobar {

	}

	func main(a int) {
		var b int = 10
		var c int = 20
		a.b().c.d()() *= 2
		a[foo.bar]++
		if a > b {
			b += 20
		} elif a == b {
			b = 0
		} else {
			b **= 2
		}
		var d int = foobar()
		return a + b * c - d
	}

	func loopy(word string) {
		var a int = 0
		while a < 10 {
			a += 1
		}
		for var i int = 1; i < 20; i++ {
			a *= i
		}
	}
	`

	prog, errs := parse(source)
	fmt.Println("ast:")
	ast.PrettyPrint(prog)

	if len(errs) > 0 {
		fmt.Println("\nerrors:")
		for _, err := range errs {
			fmt.Println(err)
		}
	}
}
