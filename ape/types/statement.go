package types

import "github.com/pcen/ape/ape/ast"

func (c *Checker) CheckStatement(stmt ast.Statement) {
	switch s := stmt.(type) {

	case *ast.BlockStmt:
		for _, s := range s.Content {
			c.CheckStatement(s)
		}

	case *ast.TypedDeclStmt:
		c.CheckDeclaration(s.Decl)

	// case *ast.AssignmentStmt:

	default:
		panic("cannot check statement " + s.StmtStr())
	}
}
