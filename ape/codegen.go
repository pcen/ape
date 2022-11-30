package ape

import (
	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/op"
	"github.com/pcen/ape/ape/token"
)

func GenerateCode(decls []ast.Declaration) []op.Code {
	// assume that each node is at least 2 opcodes
	cg := &codegen{opcodes: make([]op.Code, 2*len(decls))}
	cg.gen(decls)
	return cg.opcodes
}

type codegen struct {
	opcodes []op.Code
}

func (cg *codegen) code(c op.Code) {
	cg.opcodes = append(cg.opcodes, c)
}

func (cg *codegen) lit(expr *ast.LiteralExpr) {
	switch expr.Kind {
	case token.Integer:
	case token.Rational:
	case token.String:
	default:
		panic("codegen: unknown LiteralExpr TokenType")
	}
}

func (cg *codegen) ident(expr *ast.IdentExpr) {}

func (cg *codegen) unary(expr *ast.UnaryOp) {}

func (cg *codegen) binary(expr *ast.BinaryOp) {}

func (cg *codegen) decl(decl ast.Declaration) {}

func (cg *codegen) gen(decls []ast.Declaration) {
	for _, d := range decls {
		cg.decl(d)
	}
}
