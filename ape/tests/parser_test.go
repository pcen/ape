package tests

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
)

var (
	parser_test1 = `
	module test
	class foobar {
		Fiz string
		Buz int
		func Woof() {
			println(this.Fiz, " says woof woof!")
		}
	}

	var abc int = 100

	func main(a int) {
		var b int = 10
		var c int = 20
		a.b().c.d()() *= 2
		a[foo.bar]++
		if zzzzzzz > b {
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
	}`

	parser_test2 = `
	module test
	func main(a int) {
		b: int = 10
		SKIP {
			1 + 2
			REVERSE x
		} SEIZE(3) {
			b = 3
		}
		return b
	}`
)

func TestParsing(t *testing.T) {
	prog, errs := Parse(parser_test2)
	fmt.Println("module:", prog.Module)
	fmt.Println("ast:")
	ast.PrettyPrint(prog.Ast)
	if len(errs) > 0 {
		fmt.Println("\nerrors:")
		for _, err := range errs {
			fmt.Println(err)
		}
	}
}

// go test -run TestRandomParse ./ape/tests/*** -v
func TestRandomParse(t *testing.T) {
	fuzzDir := "../../tests/fuzz/"
	filepath.WalkDir(fuzzDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if err != nil {
			panic(err)
		}
		t.Run(d.Name(), func(t *testing.T) {
			tokens := ape.NewLexer().LexFile(path)
			parser := ape.NewParser(tokens)
			f := parser.File()
			fmt.Println("ast:")
			ast.PrettyPrint(f.Ast)
			errs, hasErrors := parser.Errors()
			if hasErrors {
				fmt.Println("\nerrors:")
				for _, err := range errs {
					fmt.Println(err)
				}
				t.Fatal("parser errors")
			}
		})
		return nil
	})
}
