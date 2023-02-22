package main

import (
	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/interpreter"
	"github.com/pcen/ape/ape/token"
)

func main() {

	lhs := ast.NewLiteralExpr(token.NewLexeme(token.Integer, "13", token.Position{1, 1}))
	rhs := ast.NewLiteralExpr(token.NewLexeme(token.Integer, "26", token.Position{1, 1}))

	op := ast.NewBinaryOp(lhs, token.New(token.Star, token.Position{2, 2}), rhs)

	another := ast.NewLiteralExpr(token.NewLexeme(token.Rational, "17.5", token.Position{1, 1}))

	addition := ast.NewBinaryOp(another, token.New(token.Plus, token.Position{2, 2}), op)

	i := interpreter.TWI{}

	i.Interpret(addition)
}
