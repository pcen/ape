package interpreter

import (
	"fmt"
	"strconv"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

type Interpreter interface {
	Interpret(ast.Node)
}

type TWI struct {
	GlobalScope Scope
}

// TODO: Temp for testing
func (twi *TWI) Interpret(expr ast.Expression) {
	res := twi.evaluate(expr)
	fmt.Println(res)
}

/** === Expression Code Begins ==== */
func (twi *TWI) evaluate(expr ast.Expression) value {
	switch t := expr.(type) {
	case *ast.LiteralExpr:
		return twi.visitLiteralExpr(t)
	case *ast.BinaryOp:
		return twi.visitBinaryExpr(t)
	default:
		panic(fmt.Sprintf("Expression type cannot be evaluated: %s", t))
	}
}

func (twi *TWI) visitLiteralExpr(literal *ast.LiteralExpr) value {
	switch literal.Kind {
	case token.String:
		return val_str{literal.Lexeme}
	case token.Integer:
		val, _ := strconv.Atoi(literal.Lexeme)
		return val_int{val}
	case token.Rational:
		val, _ := strconv.ParseFloat(literal.Lexeme, 64)
		return val_rational{val}
	case token.True:
		return val_bool{true}
	case token.False:
		return val_bool{false}
	default:
		panic(fmt.Sprintf("Unknown literal expression kind: %s", literal.Kind))
	}
}

// TODO: Handle the equal expressions (+=, -=, etc...)
func (twi *TWI) visitBinaryExpr(bin *ast.BinaryOp) value {
	lv := twi.evaluate(bin.Lhs)
	rv := twi.evaluate(bin.Rhs)

	switch bin.Op.Kind {
	case token.Plus:
		switch lv.(type) {
		case val_str:
			return val_str{lv.(val_str).Value + rv.(val_str).Value}
		default:
			return lv.(number).Add(rv.(number)).(value) // Typechecked: we know this is a number
		}
	case token.Minus:
		return lv.(number).Subtract(rv.(number)).(value)
	case token.Star:
		return lv.(number).Multiply(rv.(number)).(value)
	case token.Divide:
		return lv.(number).Divide(rv.(number)).(value)
	case token.Power:
		return lv.(number).Power(rv.(number)).(value)
	case token.Less:
		return lv.(number).LessThan(rv.(number))
	case token.LessEq:
		return lv.(number).LessThanEq(rv.(number))
	case token.Greater:
		return lv.(number).GreaterThan(rv.(number))
	case token.GreaterEq:
		return lv.(number).GreaterThanEq(rv.(number))
	case token.Equal:
		return val_bool{lv.Equals(rv)}
	case token.NotEqual:
		return val_bool{!lv.Equals(rv)}
	}

	panic(fmt.Sprintf("Unknown binary operation: %s", bin.Op.Kind))
}

/** === Expression Code Ends ==== */
