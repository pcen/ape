// package main

// import (
// 	"github.com/pcen/ape/ape/ast"
// 	"github.com/pcen/ape/ape/interpreter"
// 	"github.com/pcen/ape/ape/token"
// )

// func main() {

// 	lhs := ast.NewLiteralExpr(token.NewLexeme(token.Integer, "13", token.Position{1, 1}))
// 	rhs := ast.NewLiteralExpr(token.NewLexeme(token.Integer, "26", token.Position{1, 1}))

// 	op := ast.NewBinaryOp(lhs, token.New(token.Plus, token.Position{2, 2}), rhs)

// 	another := ast.NewLiteralExpr(token.NewLexeme(token.Rational, "17.5", token.Position{1, 1}))

// 	addition := ast.NewBinaryOp(another, token.New(token.Star, token.Position{2, 2}), op)

// 	i := interpreter.TWI{}

// 	i.Interpret(addition)
// }

package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/interpreter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("supply file to parse")
	}
	file := os.Args[1]
	lexer := ape.NewLexer()
	tokens := lexer.LexFile(file)

	parser := ape.NewParser(tokens)
	prog := parser.Program()
	fmt.Println("ast:")
	ast.PrettyPrint(prog)

	twi := interpreter.NewTWI()

	for _, decl := range prog {
		twi.Interpret(decl)
	}

	twi.RunMain()

	if errors, ok := parser.Errors(); ok {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
