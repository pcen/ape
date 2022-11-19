package types

import (
	"fmt"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

type checkerError struct {
	pos token.Position
	msg string
}

func (e *checkerError) String() string {
	return fmt.Sprintf("%v: %v", e.pos, e.msg)
}

func (c *Checker) errTypeMissmatch(pos token.Position, ident string, expected, got string) {
	c.err(pos, "type missmatch for %v: expected %v, got %v", ident, expected, got)
}

func (c *Checker) errUndefinedIdent(expr *ast.IdentExpr) {
	c.err(expr.Token.Position, "undefined identifier %v", expr.Token.Lexeme)
}

func (c *Checker) err(pos token.Position, format string, a ...interface{}) {
	c.Errors = append(c.Errors, checkerError{
		pos: pos,
		msg: fmt.Sprintf(format, a...),
	})
}

type Checker struct {
	Scope      *Scope
	scopeStack []*Scope
	File       *ast.File
	Errors     []checkerError
}

func NewChecker(File *ast.File) *Checker {
	moduleScope := NewScope(GlobalScope())
	return &Checker{
		Scope:      moduleScope,
		scopeStack: []*Scope{moduleScope},
		File:       File,
	}
}

func (c *Checker) pushScope() {
	top := NewScope(c.Scope)
	c.scopeStack = append(c.scopeStack, top)
	c.Scope = top
}

func (c *Checker) popScope() {
	if len(c.scopeStack) <= 2 {
		panic("cannot pop module scope from scope stack")
	}
	c.scopeStack = c.scopeStack[:len(c.scopeStack)-1]
	c.Scope = c.scopeStack[len(c.scopeStack)-1]
}

// TODO put this somewhere in ast package
func filter[T ast.Declaration](decls []ast.Declaration) (filtered []T) {
	for _, decl := range decls {
		if d, ok := decl.(T); ok {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// GatherPackageScope "forward declares" package level type and function declarations
func (c *Checker) GatherModuleScope() {
	for _, d := range filter[*ast.ClassDecl](c.File.Ast) {
		if err := c.Scope.DeclareType(d.Name.Lexeme); err != nil {
			c.err(d.Name.Position, err.Error())
		}
	}

	for _, d := range filter[*ast.FuncDecl](c.File.Ast) {
		if err := c.Scope.DeclareSymbol(d.Name.Lexeme, Func); err != nil {
			c.err(d.Name.Position, err.Error())
		}
	}

	for _, d := range filter[*ast.TypedDecl](c.File.Ast) {
		if typ, ok := c.Scope.LookupType(d.Type); !ok {
			c.err(d.Ident.Position, "unknown type in declaration of %v, %v", d.Ident.Lexeme, d.Type)
		} else if err := c.Scope.DeclareSymbol(d.Ident.Lexeme, typ); err != nil {
			c.err(d.Ident.Position, err.Error())
		} else if exprType := c.CheckExpr(d.Value); !Same(typ, exprType) {
			c.errTypeMissmatch(d.Ident.Position, d.Ident.Lexeme, d.Type, exprType.String())
		}
	}
}

func (c *Checker) Check() {
	// c.GatherPackageScope()
	for _, decl := range c.File.Ast {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			fmt.Println("type checking func", d.Name.Lexeme)
			c.CheckStatement(d.Body)
		}
	}
}
