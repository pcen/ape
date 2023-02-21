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

/** === Expression Code Begins ==== */

func (twi *TWI) evaluate(expr ast.Expression) value {
	switch t := expr.(type) {
	case *ast.LiteralExpr:
		twi.visitLiteralExpr(t)
	}

	panic("Expression type cannot be evaluated")
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

// TODO: Handle the OP= expressions (+=, -=, etc...)
func (twi *TWI) visitBinaryExpr(bin ast.BinaryOp) value {
	lv := twi.evaluate(bin.Lhs)
	rv := twi.evaluate(bin.Rhs)

	switch bin.Op.Kind {
	case token.Plus:
		switch lv.(type) {
		case val_int:
			l, r := cast_value[val_int](lv, rv)
			return val_int{l.Value + r.Value}
		case val_rational:
			l, r := cast_value[val_rational](lv, rv)
			return val_rational{l.Value + r.Value}
		case val_str:
			l, r := cast_value[val_str](lv, rv)
			return val_str{l.Value + r.Value}
		}

	case token.Minus:
		switch lv.(type) {
		case val_int:
			l, r := cast_value[val_int](lv, rv)
			return val_int{l.Value + r.Value}
		case val_rational:
			l, r := cast_value[val_rational](lv, rv)
			return val_rational{l.Value + r.Value}
		}
	case token.Star:
	case token.Divide:
	case token.Power:
	case token.Equal:
	case token.NotEqual:
	case token.Less:
	case token.LessEq:
	case token.Greater:
	case token.GreaterEq:
	}

	panic(fmt.Sprintf("Unknown binary operation: %s", bin.Op.Kind))
}

/** == Expression Utilities == */
func cast_value[T value](a value, b value) (T, T) {
	return a.(T), b.(T)
}

/** === Expression Code Ends ==== */
