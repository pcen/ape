package ape

import (
	"reflect"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

func GenerateCode(node ast.Node) []Opcode {
	cg := &codegen{opcodes: make([]Opcode, 0)}
	cg.gen(node)
	return cg.opcodes
}

type codegen struct {
	opcodes []Opcode
}

func (cg *codegen) code(oc Opcode) {
	cg.opcodes = append(cg.opcodes, oc)
}

func (cg *codegen) lit(expr *ast.LiteralExpr) {
	switch expr.Kind {
	case token.Number:
		// deduce size
	case token.True:
		cg.code(OpTrue)
	case token.False:
		cg.code(OpFalse)
	case token.String:
	default:
		panic("codegen: unknown LiteralExpr TokenType")
	}
}

func (cg *codegen) ident(expr *ast.IdentExpr) {}

func (cg *codegen) unary(expr *ast.UnaryOp) {}

func (cg *codegen) binary(expr *ast.BinaryOp) {}

func (cg *codegen) typedDecl(stmt *ast.TypedDecl) {}

func (cg *codegen) gen(node ast.Node) {
	switch e := node.(type) {
	case *ast.GroupExpr:
		cg.gen(e.Expr)
	case *ast.LiteralExpr:
		cg.lit(e)
	case *ast.IdentExpr:
		cg.ident(e)
	case *ast.UnaryOp:
		cg.unary(e)
	case *ast.BinaryOp:
		cg.binary(e)
	case *ast.TypedDecl:
		cg.typedDecl(e)
	default:
		panic("unknown AST node type: " + reflect.TypeOf(node).Name())
	}
}
