package c

import (
	"reflect"
	"strings"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

func GenerateCode(stmts []ast.Statement) *codegen {
	// assume that each node is at least 2 opcodes
	cg := newCodegen()
	// forward declare printf
	cg.write("int\tprintf(const char*, ...);\n")
	cg.write("int main(int argc, char* argv[]) {\n")
	cg.demo(stmts)
	cg.write("\n}\n")
	return cg
}

type codegen struct {
	Code   strings.Builder
	Idents map[string]struct{}
}

func newCodegen() *codegen {
	return &codegen{
		Idents: make(map[string]struct{}),
	}
}

func (cg *codegen) write(s string) {
	cg.Code.WriteString(s)
}

func (cg *codegen) expr(expr ast.Expression) {
	switch e := expr.(type) {

	case *ast.LiteralExpr:
		switch e.Kind {
		case token.Integer:
			cg.write(e.Lexeme)
		case token.True:
			cg.write("1")
		case token.False:
			cg.write("0")
		default:
			panic("cannot codegen for literal expr of type " + e.Kind.String())
		}

	case *ast.IdentExpr:
		cg.write(e.Ident.Lexeme)

	case *ast.BinaryOp:
		cg.gen(e.Lhs)
		switch e.Op {
		case token.Plus, token.Minus, token.Star, token.Divide:
			cg.write(" " + e.Op.String() + " ")
		case token.Equal, token.NotEqual, token.Greater, token.GreaterEq, token.Less, token.LessEq:
			cg.write(" " + e.Op.String() + " ")
		default:
			panic("invalid binary op: " + e.Op.String())
		}
		cg.gen(e.Rhs)

	case *ast.CallExpr:
		cg.write(`printf("%d\n",`)
		cg.gen(e.Args[0])
		if ie, ok := e.Callee.(*ast.IdentExpr); ok && strings.Contains(ie.Ident.Lexeme, "print") {
			cg.write(")")
		} else {
			panic("call expr is not print")
		}
	}
}

func (cg *codegen) stmt(stmt ast.Statement) {
	switch t := stmt.(type) {
	case *ast.ExprStmt:
		cg.expr(t.Expr)
		cg.write(";\n")

	case *ast.TypedDeclStmt:
		cg.decl(t.Decl)
		cg.write(";\n")

	case *ast.AssignmentStmt:
		cg.gen(t.Lhs)
		cg.write(" = ")
		cg.gen(t.Rhs)
		if _, ok := t.Lhs.(*ast.IdentExpr); !ok {
			panic("target of assignment must be identifier")
		}
		cg.write(";\n")

	case *ast.IfStmt:
		cg.write("if (")
		cg.expr(t.If.Cond)
		cg.write(") {\n")
		cg.stmt(t.If.Body)
		cg.write("}")
		if t.Else != nil {
			cg.write(" else {\n")
			cg.stmt(t.Else)
			cg.write("}")
		}
		cg.write("\n")

	case *ast.BlockStmt:
		for _, s := range t.Content {
			cg.stmt(s)
		}

	default:
		panic("cannot gen for statement " + reflect.TypeOf(stmt).String())
	}
}

func (cg *codegen) decl(decl ast.Declaration) {
	switch d := decl.(type) {
	case *ast.TypedDecl:
		cg.write("int " + d.Ident.Lexeme + " = ")
		cg.gen(d.Value)

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
