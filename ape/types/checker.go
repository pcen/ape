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
	c.err(expr.Ident.Position, "undefined identifier %v", expr.Ident.Lexeme)
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
	Types      map[ast.Expression]Type
	File       *ast.File
	Errors     []checkerError
}

func NewChecker(File *ast.File) *Checker {
	moduleScope := NewScope(GlobalScope())
	return &Checker{
		Scope:      moduleScope,
		scopeStack: []*Scope{moduleScope},
		Types:      make(map[ast.Expression]Type),
		File:       File,
	}
}

func (c *Checker) pushScope() {
	top := NewScope(c.Scope)
	c.scopeStack = append(c.scopeStack, top)
	c.Scope = top
}

func (c *Checker) popScope() {
	if len(c.scopeStack) <= 1 {
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
		var returns Type
		var ok bool
		returns, ok = Void, true
		if d.ReturnType != nil {
			returns, ok = c.Scope.LookupType(d.ReturnType.Name)
		}
		if !ok {
			fmt.Println("unknown return type for ", d.Name.Lexeme)
		}
		if err := c.Scope.DeclareSymbol(d.Name.Lexeme, NewFunction(nil, []Type{returns})); err != nil {
			c.err(d.Name.Position, err.Error())
		}
	}

	for _, d := range filter[*ast.TypedDecl](c.File.Ast) {
		if typ, ok := c.Scope.LookupType(d.Type.Name); !ok {
			c.err(d.Ident.Position, "unknown type in declaration of %v, %v", d.Ident.Lexeme, d.Type)
		} else if err := c.Scope.DeclareSymbol(d.Ident.Lexeme, typ); err != nil {
			c.err(d.Ident.Position, err.Error())
		} else if exprType := c.CheckExpr(d.Value); !typ.Is(exprType) {
			c.errTypeMissmatch(d.Ident.Position, d.Ident.Lexeme, d.Type.Name, exprType.String())
		}
	}

}

func (c *Checker) Check() Environment {
	c.GatherModuleScope()
	for _, decl := range c.File.Ast {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			fmt.Println("type checking func", d.Name.Lexeme)
			c.CheckDeclaration(d)
		}
	}
	for _, e := range c.Errors {
		fmt.Println(e)
	}
	return Environment{Expressions: c.Types}
}
