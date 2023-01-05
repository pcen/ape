package types

import (
	"errors"
	"fmt"

	"github.com/pcen/ape/ape/ast"
)

func (c *Checker) ResolveTypeNode(n *ast.TypeExpr) (Type, error) {
	if n == nil {
		// TODO: this makes ast simpler, but might make parser bugs harder to find
		return Void, nil
	}
	typ, ok := c.Scope.LookupType(n.Name)
	if !ok {
		return Invalid, errors.New("undefined type in current scope")
	}
	if n.List {
		typ = NewList(typ)
	}
	return typ, nil
}

func (c *Checker) CheckDeclaration(decl ast.Declaration) {
	switch d := decl.(type) {
	case *ast.VarDecl:
		dtyp, err := c.ResolveTypeNode(d.Type)
		if err != nil {
			c.err(d.Ident.Position, "undefined type %v for %v", d.Type, d.Ident.Lexeme)
		}
		if err := c.Scope.DeclareSymbol(d.Ident.Lexeme, dtyp); err != nil {
			c.err(d.Ident.Position, err.Error())
		}
		if d.Value != nil {
			valueType := c.CheckExpr(d.Value)
			if !valueType.Is(dtyp) {
				c.errTypeMissmatch(d.Ident.Position, d.Ident.Lexeme, dtyp.String(), valueType.String())
			}
		}
		c.Types[d.Type] = dtyp

	case *ast.ClassDecl:
		return

	case *ast.FuncDecl:
		retType, err := c.ResolveTypeNode(d.ReturnType)
		if err != nil {
			c.err(d.Name.Position, "undefined return type for %v: %v", d.Name.Lexeme, d.ReturnType.Name)
		}
		paramSignature := make([]Type, 0, len(d.Params))
		for _, p := range d.Params {
			c.CheckDeclaration(p)
			paramSignature = append(paramSignature, c.Types[p.Type])
		}
		c.CheckStatement(d.Body)
		c.Scope.DeclareSymbol(d.Name.Lexeme, NewFunction(paramSignature, []Type{retType}))
		fmt.Printf("%v signature: %v -> %v\n", d.Name, paramSignature, retType)

	case *ast.ParamDecl:
		dtyp, err := c.ResolveTypeNode(d.Type)
		if err != nil {
			c.err(d.Ident.Ident.Position, "undefined type %v for %v", d.Type, d.Ident.Ident.Lexeme)
		}
		if err := c.Scope.DeclareSymbol(d.Ident.Ident.Lexeme, dtyp); err != nil {
			c.err(d.Ident.Ident.Position, err.Error())
		}
		fmt.Printf("param decl %v has type %v\n", d.Ident, dtyp)

		c.CheckExpr(d.Ident) // set expr type
		c.Types[d.Type] = dtyp

	default:
		panic("cannot type check declaration " + d.DeclStr())
	}
}
