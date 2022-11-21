package types

import (
	"fmt"
	"reflect"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

func (c *Checker) CheckExpr(expr ast.Expression) (t Type) {
	switch e := expr.(type) {

	case *ast.LiteralExpr:
		switch e.Kind {
		case token.String:
			t = String
		case token.Integer:
			t = Int
		case token.Rational:
			t = Float
		case token.True, token.False:
			t = Bool
		default:
			// panic on upstream parser error
			panic(fmt.Sprintf("invalid token kind for LiteralExpr: %v", e.Kind))
		}

	case *ast.GroupExpr:
		return c.CheckExpr(e.Expr)

	case *ast.UnaryOp:
		return c.CheckExpr(e.Expr)

	case *ast.BinaryOp:
		t1 := c.CheckExpr(e.Lhs)
		t2 := c.CheckExpr(e.Rhs)
		if !Same(t1, t2) {
			c.err(token.Position{}, "invalid types for binary op: %v %v %v", t1, e.Op, t2)
		}
		t = t1

	case *ast.IdentExpr:
		typ, ok := c.Scope.LookupSymbol(e.Lexeme)
		if !ok {
			c.errUndefinedIdent(e)
		}
		t = typ

	default:
		panic(fmt.Sprintf("cannot type check expressions of type %v", reflect.TypeOf(expr)))
	}

	fmt.Printf("%v has type %v\n", expr.ExprStr(), t)
	return t
}
