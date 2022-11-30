package ape

import (
	"strconv"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/op"
	"github.com/pcen/ape/ape/token"
)

func GenerateCode(stmts []ast.Statement) *codegen {
	// assume that each node is at least 2 opcodes
	cg := newCodegen()
	cg.demo(stmts)
	return cg
}

type codegen struct {
	Code     []op.Code
	LitIdx   map[int32]uint8
	IdentIdx map[string]uint8
}

func newCodegen() *codegen {
	return &codegen{
		Code:     make([]op.Code, 0),
		LitIdx:   make(map[int32]uint8),
		IdentIdx: make(map[string]uint8),
	}
}

func (cg *codegen) lit(value int32) uint8 {
	if idx, ok := cg.LitIdx[value]; ok {
		return idx
	}
	idx := uint8(len(cg.LitIdx))
	cg.LitIdx[value] = idx
	return idx
}

func (cg *codegen) ident(name string) uint8 {
	if idx, ok := cg.IdentIdx[name]; ok {
		return idx
	}
	idx := uint8(len(cg.IdentIdx))
	cg.IdentIdx[name] = idx
	return idx
}

func (cg *codegen) op(c op.Code) {
	cg.Code = append(cg.Code, c)
}

func (cg *codegen) put(u uint8) {
	cg.Code = append(cg.Code, op.Code(u))
}

func (cg *codegen) expr(expr ast.Expression) {
	switch e := expr.(type) {

	case *ast.LiteralExpr:
		var idx uint8
		switch e.Kind {
		case token.Integer:
			i, err := strconv.ParseInt(e.Lexeme, 10, 32)
			if err != nil {
				panic(err)
			}
			idx = cg.lit(int32(i))
		case token.True:
			idx = cg.lit(1)
		case token.False:
			idx = cg.lit(0)
		default:
			panic("cannot codegen for literal expr of type " + e.Kind.String())
		}
		cg.op(op.Constant)
		cg.put(idx)

	case *ast.IdentExpr:
		idx := cg.ident(e.Ident.Lexeme)
		cg.op(op.Get)
		cg.put(idx)

	case *ast.BinaryOp:
		cg.gen(e.Lhs)
		cg.gen(e.Rhs)
		switch e.Op {
		case token.Plus:
			cg.op(op.Add)
		case token.Minus:
			cg.op(op.Subtract)
		case token.Star:
			cg.op(op.Multiply)
		case token.Divide:
			cg.op(op.Divide)
		default:
			panic("invalid binary op: " + e.Op.String())
		}
	}
}

func (cg *codegen) stmt(stmt ast.Statement) {
	switch t := stmt.(type) {
	case *ast.ExprStmt:
		cg.expr(t.Expr)

	case *ast.TypedDeclStmt:
		cg.decl(t.Decl)

	case *ast.AssignmentStmt:
		cg.gen(t.Rhs)
		cg.op(op.Set)
		if ie, ok := t.Lhs.(*ast.IdentExpr); ok {
			idx := cg.ident(ie.Ident.Lexeme)
			cg.put(idx)
		} else {
			panic("target of assignment must be identifier")
		}
	}
}

func (cg *codegen) decl(decl ast.Declaration) {
	switch d := decl.(type) {
	case *ast.TypedDecl:
		cg.gen(d.Value)
		cg.op(op.Set)
		idx := cg.ident(d.Ident.Lexeme)
		cg.put(idx)

	default:
		panic("cannot codegen decl: " + decl.DeclStr())
	}
}

func (cg *codegen) gen(node ast.Node) {
	switch n := node.(type) {
	case ast.Declaration:
		cg.decl(n)
	case ast.Statement:
		cg.stmt(n)
	case ast.Expression:
		cg.expr(n)
	default:
		panic("unknown interface type for ast.Node")
	}
}

func (cg *codegen) demo(stmts []ast.Statement) {
	for _, s := range stmts {
		cg.stmt(s)
	}
}
