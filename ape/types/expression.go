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
		t = c.CheckExpr(e.Expr)

	case *ast.UnaryOp:
		t = c.CheckExpr(e.Expr)

	case *ast.BinaryOp:
		t1 := c.CheckExpr(e.Lhs)
		t2 := c.CheckExpr(e.Rhs)
		if !t1.Is(t2) {
			c.err(e.Op.Position, "invalid types for binary op: %v %v %v", t1, e.Op, t2)
			t = Invalid
			break
		}
		t = t1
		switch e.Op.Kind {
		case token.Equal, token.NotEqual, token.Greater, token.GreaterEq, token.Less, token.LessEq:
			t = Bool
		}

	case *ast.IdentExpr:
		typ, ok := c.Scope.LookupSymbol(e.Ident.Lexeme)
		if !ok {
			c.errUndefinedIdent(e)
		}
		t = typ

	case *ast.CallExpr:
		t = c.CheckExpr(e.Callee)
		for _, arg := range e.Args {
			c.CheckExpr(arg)
		}

	case *ast.DotExpr:
		et := c.CheckExpr(e.Expr)
		// the type of Field depends on the type of the receiver
		switch et.(type) {
		case List:
			if e.Field.Ident.Lexeme == "push" {
				t = NewFunction(nil, nil)
			}
		default:
			fmt.Println("WARNING: unknown receiver type in dot expression")
			t = c.CheckExpr(e.Field)
		}

	case *ast.IndexExpr:
		t = c.CheckExpr(e.Expr)
		c.CheckExpr(e.Index)
		if list, ok := t.(List); ok {
			t = list.Data
		} else if m, ok := t.(Map); ok {
			t = m.Value
		} else {
			panic("cannot index into non-list type")
		}

	case *ast.LitListExpr:
		t = c.CheckExpr(e.Elements[0])
		if len(e.Elements) >= 2 {
			for i := 1; i < len(e.Elements); i++ {
				te := c.CheckExpr(e.Elements[i])
				if !te.Is(t) {
					panic("inconsistent types in list literal")
				}
			}
		}
		t = NewList(t)

	case *ast.LitMapExpr:
		var kt Type = nil
		var vt Type = nil
		for k, v := range e.Elements {
			if kt == nil && vt == nil {
				kt = c.CheckExpr(k)
				vt = c.CheckExpr(v)
			} else {
				currentKt := c.CheckExpr(k)
				currentVt := c.CheckExpr(v)
				if currentKt != kt {
					panic("key type in map does not match first key type")
				} else if currentVt != vt {
					panic("value type in map does not match first value type")
				}
			}
		}
		t = NewMap(kt, vt)

	default:
		panic(fmt.Sprintf("cannot type check expressions of type %v", reflect.TypeOf(expr)))
	}

	c.Types[expr] = t
	return t
}
