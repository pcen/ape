package types

import (
	"reflect"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

func (c *Checker) CheckStatement(stmt ast.Statement) Type {
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
			c.err(token.Position{}, "type missmatch in assignment statement: %v is not %v", l, r)
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

	case *ast.SwitchStmt:
		t := c.CheckExpr(s.Expr)
		if _, ok := t.(Primitive); !ok {
			c.err(s.Token.Position, "invalid type for switch value: %v", t)
		}
		for _, caseStmt := range s.Cases {
			c.CheckStatement(caseStmt)
		}

	case *ast.CaseStmt:
		c.CheckStatement(s.Body)

	case *ast.FallthroughtStmt:
		break

	case *ast.SkipStmt:
		var reverseType Type = nil
		for _, bodyStmt := range s.Body.Content {
			switch reverseStmt := bodyStmt.(type) {
			case *ast.ReverseStmt:
				// make sure the type reversed on is consistent throughout the current skip
				// statement block
				nextReverseType := c.CheckExpr(reverseStmt.Expr)
				if reverseType != nil && !nextReverseType.Is(Void) && !nextReverseType.Is(reverseType) {
					c.err(reverseStmt.Token.Position, "inconsistent reverse types in skip block")
				}
			default:
				c.CheckStatement(bodyStmt)
			}
		}
		for _, seize := range s.Seizes {
			// make sure that each seize statement seizes the same type as each reverse statement
			// in the preceding skip statement block
			accepts := c.CheckStatement(seize)
			if reverseType != nil && !accepts.Is(Void) && !accepts.Is(reverseType) {
				c.err(seize.Token.Position, "seize expr type does not match reverse expr type in skip block")
			}
		}

	case *ast.ReverseStmt:
		if s.Expr == nil {
			return Void
		}
		return c.CheckExpr(s.Expr)

	case *ast.SeizeStmt:
		var accepts Type = Void
		if s.Expr != nil {
			accepts = c.CheckExpr(s.Expr)
		}
		c.CheckStatement(s.Body)
		return accepts

	default:
		panic("cannot check statement " + s.StmtStr() + ", " + reflect.TypeOf(stmt).String())
	}
	return Void
}
