package types

import (
	"fmt"

	"github.com/pcen/ape/ape/ast"
)

func (c *Checker) CheckDeclaration(decl ast.Declaration) (t Type) {
	switch d := decl.(type) {
	case *ast.TypedDecl:

		declType, ok := c.Scope.LookupType(d.Type)
		if !ok {
			c.err(d.Ident.Position, "undefined type %v for %v", d.Type, d.Ident.Lexeme)
		} else if err := c.Scope.DeclareSymbol(d.Ident.Lexeme, declType); err != nil {
			c.err(d.Ident.Position, err.Error())
		}
		valueType := c.CheckExpr(d.Value)
		if !Same(valueType, declType) {
			c.errTypeMissmatch(d.Ident.Position, d.Ident.Lexeme, declType.String(), valueType.String())
		}
		// DEBUG
		fmt.Printf("identifier %v has type %v\n", d.Ident.Lexeme, valueType)
		t = valueType

	case *ast.ClassDecl:
		return

	case *ast.FuncDecl:
		retType, ok := c.Scope.LookupType(d.ReturnType.Lexeme)
		if !ok {
			c.err(d.Name.Position, "undefined return type for %v: %v", d.Name.Lexeme, d.ReturnType.Lexeme)
		}
		fmt.Printf("%v returns type %v", d.Name, retType)
		t = retType

	default:
		panic("cannot type check declaration " + d.DeclStr())
	}
	return t
}
