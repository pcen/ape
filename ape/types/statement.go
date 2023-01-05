package types

import (
	"reflect"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
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
		// Init is nil for while loops
		if s.Init != nil {
			c.CheckDeclaration(s.Init)
		}
		c.CheckExpr(s.Cond)
		// Incr is nil for while loops
		if s.Incr != nil {
			c.CheckStatement(s.Incr)
		}
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

	case *ast.IfStmt:
		if !c.CheckExpr(s.If.Cond).Is(Bool) {
			c.err(token.Position{}, "if condition must have boolean type")
		}
		c.CheckStatement(s.If.Body)
		for _, elif := range s.Elifs {
			c.CheckStatement(elif)
		}
		if s.Else != nil {
			c.CheckStatement(s.Else)
		}

	case *ast.CondBlockStmt:
		if !c.CheckExpr(s.Cond).Is(Bool) {
			c.err(token.Position{}, "elif condition must have boolean type")
		}
		c.CheckStatement(s.Body)

	case *ast.BreakStmt:
		break

	default:
		panic("cannot check statement " + s.StmtStr() + ", " + reflect.TypeOf(stmt).String())
	}
}
