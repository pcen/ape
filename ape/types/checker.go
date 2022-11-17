package types

import (
	"fmt"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

type Checker interface {
	Check([]ast.Declaration) (map[ast.Expression]Type, map[string]Type)
}

type checker struct {
	types   map[ast.Expression]Type
	symbols map[string]Type
}

func NewChecker() Checker {
	return &checker{
		types:   make(map[ast.Expression]Type),
		symbols: make(map[string]Type),
	}
}

func (c *checker) Check(decls []ast.Declaration) (map[ast.Expression]Type, map[string]Type) {
	for _, decl := range decls {
		c.checkDeclaration(decl)
	}
	return c.types, c.symbols
}

func (c *checker) checkDeclaration(decl ast.Declaration) {
	switch d := decl.(type) {
	case *ast.TypedDecl:
		fmt.Println(d)
	}
}

func (c *checker) checkExpression(expr ast.Expression) (t Type) {
	switch e := expr.(type) {
	// case *ast.GroupExpr:
	// 	t = c.checkExpression(e.Expr)
	// 	c.types[e] = t
	case *ast.LiteralExpr:
		t = c.checkLiteral(e)
		c.types[e] = t
	case *ast.IdentExpr:
		t = c.checkIdent(e)
		c.types[e] = t
	case *ast.UnaryOp:
		t = c.checkUnary(e)
		c.types[e] = t
	case *ast.BinaryOp:
		t = c.checkBinary(e)
		c.types[e] = t
	default:
		panic("checker: unknown expression ast node type")
	}
	return t
}

func (c *checker) checkStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		c.checkExpression(s.Expr)
	default:
		panic("unknown statement ast node")
	}
}

func (c *checker) checkTypedDecl(decl *ast.TypedDecl) {
	switch decl.Kind {
	case token.Val, token.Var:
		typeKind := LookupPrimitive(decl.Type)
		exprType := c.checkExpression(decl.Value)
		if typeKind != exprType.Kind {
			panic("declaration gets wrong type")
		}
		c.symbols[decl.Ident.Lexeme] = exprType
	default:
		panic("unknown decl kind")
	}
}

func (c *checker) checkLiteral(lit *ast.LiteralExpr) Type {
	switch lit.Kind {
	case token.Number:
		return NewType(Int)
	case token.String:
		return NewType(String)
	case token.True, token.False:
		return NewType(Bool)
	default:
		panic("checker: unknown literal ast type")
	}
}

func (c *checker) checkIdent(ident *ast.IdentExpr) Type {
	t, known := c.symbols[ident.Lexeme]
	if !known {
		panic(fmt.Sprintf("checker: %v not in symbol table", ident))
	}
	c.types[ident] = t
	return t
}

func (c *checker) checkUnary(unary *ast.UnaryOp) (t Type) {
	t = c.checkExpression(unary.Expr)
	c.types[unary] = t
	return t
}

func (c *checker) checkBinary(binary *ast.BinaryOp) Type {
	lhsType := c.checkExpression(binary.Lhs)
	rhsType := c.checkExpression(binary.Rhs)
	if lhsType.Kind != rhsType.Kind {
		panic("binary operands are not equal")
	}
	c.types[binary] = lhsType
	return lhsType
}
