package tests

import (
	"fmt"
	"testing"

	"github.com/pcen/ape/ape/types"
)

const (
	prog = `
		module test_module
		class Foo {}
		class Bar {}
		class Biz {}
	`

	prog2 = `
		module test_module
		var a int = 1
		var b string = "foo"
		var c bool = false
		var d int = 2

		func fizzle() {
			var e int = a * (d - (5 + 1))
		}

		val pi float = 3.1415
	`
)

var (
	progs = []string{
		prog, prog2,
	}
)

func TestChecker(t *testing.T) {
	for i, prog := range progs {
		t.Run(fmt.Sprintf("program %v", i), func(t *testing.T) {
			f, _ := Parse(prog)
			c := types.NewChecker(f)

			c.GatherModuleScope()
			c.Scope.Print()

			c.Check()

			for _, err := range c.Errors {
				fmt.Printf("checker error: %v\n", err)
			}

		})
	}
}
