package tests

import (
	"fmt"
	"testing"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/types"
)

const (
	// prog1 = `
	// 	module test_module
	// 	class Foo {}
	// 	class Bar {}
	// 	class Biz {}
	// `

	// prog2 = `
	// 	module test_module
	// 	var a int = 1
	// 	var b string = "foo"
	// 	var c bool = false
	// 	var d int = 2

	// 	func fizzle() {
	// 		var e int = a * (d - (5 + 1))
	// 	}

	// 	val pi float = 3.1415
	// `

	// prog3 = `
	// 	module test_module
	// 	var a int = 1
	// 	var b int = 2
	// 	var c int = a * b
	// 	var d float = c ** 2.5 ** 4

	// 	func foo(a float) int {
	// 		var b string = "jejeje"
	// 		return b
	// 	}

	// 	var e string = foo()
	// `

	prog4 = `
		module test
		func main() {
			SKIP {
				REVERSE 1
			} SEIZE {
				println("default seize")
			}
		}

		func notMain() {
			SKIP {
				REVERSE "string"
			} SEIZE ("string") {
				println("seize on a string")
			}
		}
	`

	badSeize = `
		module test
		func main() {
			SKIP {
				REVERSE 1.0
			} SEIZE (1) {
				println("bad")
			}
		}
	`
)

var (
	progs = []string{
		prog4, badSeize,
		// prog1, prog2, prog3,
	}
)

func TestChecker(t *testing.T) {
	for i, prog := range progs {
		t.Run(fmt.Sprintf("program %v", i), func(t *testing.T) {
			f, _ := Parse(prog)
			c := types.NewChecker(f)
			ast.PrettyPrint(f.Ast)

			c.GatherModuleScope()
			c.Scope.Print()

			c.Check()

			for _, err := range c.Errors {
				fmt.Printf("checker error: %v\n", err)
			}

		})
	}
}
