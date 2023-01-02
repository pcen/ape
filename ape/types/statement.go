package types

import (
	"reflect"

	"github.com/pcen/ape/ape/ast"
)

func (c *Checker) CheckStatement(stmt ast.Statement) {
	switch s := stmt.(type) {

	case *ast.BlockStmt:
		for _, s := range s.Content {
			c.CheckStatement(s)
		}

	case *ast.TypedDeclStmt:
		c.CheckDeclaration(s.Decl)

	case *ast.ExprStmt:
		c.CheckExpr(s.Expr)

	case *ast.ForStmt:
		c.pushScope()
		c.CheckDeclaration(s.Init)
		c.CheckExpr(s.Cond)
		c.CheckStatement(s.Incr)
		c.CheckStatement(s.Body)
		c.popScope()

	case *ast.IncStmt:
		typ := c.CheckExpr(s.Expr)
		if !typ.Is(Int) {
			// TODO: all integer types can be incremented, not just int
			// c allows ++/-- on floats too
			panic("cannot increment non-integer type")
		}

	case *ast.AssignmentStmt:
		l := c.CheckExpr(s.Lhs)
		r := c.CheckExpr(s.Rhs)
		if !r.Is(l) {
			panic("type missmatch in assignment statement")
		}

	default:
		panic("cannot check statement " + s.StmtStr() + ", " + reflect.TypeOf(stmt).String())
	}
}
