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
	GlobalScope  Scope
	CurrentScope Scope
}

// TODO: Temp for testing
func (twi *TWI) Interpret(expr ast.Expression) {
	res := twi.evaluateExpr(expr)
	fmt.Println(res)
}

/** === Expression Code Begins ==== */
func (twi *TWI) evaluateExpr(expr ast.Expression) value {
	switch t := expr.(type) {
	case *ast.LiteralExpr:
		return twi.visitLiteralExpr(t)
	case *ast.BinaryOp:
		return twi.visitBinaryExpr(t)
	case *ast.GroupExpr:
		return twi.visitGroupExpr(t)
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
	lv := twi.evaluateExpr(bin.Lhs)
	rv := twi.evaluateExpr(bin.Rhs)

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
	case token.Mod:
		return lv.(val_int).Mod(rv.(val_int))
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
	case token.And:
		return val_bool{lv.(val_bool).Value && rv.(val_bool).Value}
	case token.Or:
		return val_bool{lv.(val_bool).Value || rv.(val_bool).Value}
	}

	panic(fmt.Sprintf("Unknown binary operation: %s", bin.Op.Kind))
}

func (twi *TWI) visitUnaryExpr(unary *ast.UnaryOp) value {
	val := twi.evaluateExpr(unary.Expr)

	switch unary.Op {
	case token.Bang:
		return val_bool{!val.(val_bool).Value}
	default:
		panic("Unknown unary token")
	}
}

func (twi *TWI) visitGroupExpr(group *ast.GroupExpr) value {
	return twi.evaluateExpr(group.Expr)
}

/** === Expression Code Ends ==== */
