package types

import (
	"errors"
	"fmt"

	"github.com/pcen/ape/ape/ast"
)

func (c *Checker) ResolveTypeNode(n *ast.TypeExpr) (Type, error) {
	if n == nil {
		return Invalid, errNotTyped
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

func (c *Checker) varDeclWithValue(d *ast.VarDecl) Type {
	dtyp, err := c.ResolveTypeNode(d.Type)
	etyp := c.CheckExpr(d.Value)
	if err == errNotTyped {
		dtyp = etyp // inferred
	} else if err != nil {
		c.err(d.Ident.Position, "invalid type %v for %v", d.Type, d.Ident.Lexeme)
		return Invalid
	}
	if err := c.Scope.DeclareSymbol(d.Ident.Lexeme, dtyp); err != nil {
		c.err(d.Ident.Position, err.Error())
		return Invalid
	}
	if !etyp.Is(dtyp) {
		c.errTypeMissmatch(d.Ident.Position, d.Ident.Lexeme, dtyp.String(), etyp.String())
	}
	return dtyp
}

func (c *Checker) varDeclWithoutValue(d *ast.VarDecl) Type {
	dtyp, err := c.ResolveTypeNode(d.Type)
	if err == errNotTyped {
		c.err(d.Ident.Position, "%v cannot have implicit type in declaration without value", d.Ident.Lexeme)
		return Invalid
	} else if err != nil {
		c.err(d.Ident.Position, "invalid type %v for %v", d.Type, d.Ident.Lexeme)
		return Invalid
	}
	if err := c.Scope.DeclareSymbol(d.Ident.Lexeme, dtyp); err != nil {
		c.err(d.Ident.Position, err.Error())
		return Invalid
	}
	return dtyp
}

func (c *Checker) CheckDeclaration(decl ast.Declaration) {
	switch d := decl.(type) {
	case *ast.VarDecl:
		var dtyp Type
		if d.Value != nil {
			dtyp = c.varDeclWithValue(d)
		} else {
			dtyp = c.varDeclWithoutValue(d)
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
		c.CheckExpr(d.Ident)
		c.Types[d.Type] = dtyp

	default:
		panic("cannot type check declaration " + d.DeclStr())
	}
}
