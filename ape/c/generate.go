package c

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
	"github.com/pcen/ape/ape/types"
)

var (
	tab = []byte("\t")
)

const (
	builtins = `int printf(const char*, ...);
void* malloc(unsigned long);
void* realloc(void*, unsigned long);
double pow(double, double);
void println(int i){printf("%d\n",i);}
int ipow(int x, int y){return (int)pow((double)x, (double)y);}
double dpow(double x, double y){return pow(x, y);}
`
)

var (
	listImplementations = map[types.Type]string{
		types.IntList: "ape_ivec",
	}
)

var (
	indexFunctions = map[types.Type]string{
		types.IntList: "ape_ivec_get",
	}
)

func GenerateCode(decls []ast.Declaration, env types.Environment) *codegen {
	cg := newCodegen(env)
	cg.write(builtins)
	cg.write(implementVector(listImplementations[types.IntList], "int"))
	cg.program(decls)
	return cg
}

type codegen struct {
	Code  strings.Builder
	Env   types.Environment
	level int
}

func newCodegen(env types.Environment) *codegen {
	return &codegen{Env: env}
}

func (cg *codegen) write(s string) {
	cg.Code.WriteString(s)
}

// start indented line
func (cg *codegen) sil(s string) {
	cg.indent()
	cg.write(s)
}

func (cg *codegen) indent() {
	cg.Code.Write(bytes.Repeat(tab, cg.level))
}

func (cg *codegen) indented(f func()) {
	cg.level++
	f()
	cg.level--
}

func (cg *codegen) receiverArgs(receiver ast.Expression, exprs ...ast.Expression) {
	cg.write("(&")
	cg.expr(receiver)
	if len(exprs) > 0 {
		cg.write(", ")
	}
	for i, e := range exprs {
		cg.gen(e)
		if i != len(exprs)-1 {
			cg.write(", ")
		}
	}
	cg.write(")")
}

func (cg *codegen) method(dot *ast.DotExpr, call *ast.CallExpr) {
	// assume method for now
	if dot.Field.Ident.Lexeme == "push" {
		dot.Field.Ident.Lexeme = "ape_ivec_push"
	}
	cg.gen(dot.Field)
	cg.receiverArgs(dot.Expr, call.Args...)
}

func (cg *codegen) index(receiver ast.Expression, index ast.Expression) {
	// TODO: the function used here depends on the type of the expression
	// - for lists use ape_<ctype>vec_get
	// - for maps use ape_<ctype>map_get
	// etc...
	fmt.Println("cg index of type ", cg.Env.Expressions[receiver])
	fmt.Println(receiver.ExprStr(), reflect.TypeOf(receiver))
	cg.write("ape_ivec_get")
	cg.receiverArgs(receiver, index)
}

func (cg *codegen) expr(expr ast.Expression) {
	sepWithOpLiteral := func(lhs ast.Expression, op token.Kind, rhs ast.Expression) {
		cg.gen(lhs)
		cg.write(" " + op.String() + " ")
		cg.gen(rhs)
	}

	sepWithString := func(lhs ast.Expression, op string, rhs ast.Expression) {
		cg.gen(lhs)
		cg.write(" " + op + " ")
		cg.gen(rhs)
	}

	switch e := expr.(type) {

	case *ast.LiteralExpr:
		switch e.Kind {
		case token.Integer:
			cg.write(e.Lexeme)
		case token.True:
			cg.write("1")
		case token.False:
			cg.write("0")
		case token.String:
			cg.write(`"`)
			cg.write(e.Lexeme)
			cg.write(`"`)
		default:
			panic("cannot codegen for literal expr of type " + e.Kind.String())
		}

	case *ast.IdentExpr:
		cg.write(e.Ident.Lexeme)

	case *ast.BinaryOp:
		switch e.Op {
		case token.Plus, token.Minus, token.Star, token.Divide:
			sepWithOpLiteral(e.Lhs, e.Op, e.Rhs)

		case token.Equal, token.NotEqual, token.Greater, token.GreaterEq, token.Less, token.LessEq:
			sepWithOpLiteral(e.Lhs, e.Op, e.Rhs)

		case token.Power:
			cg.write("ipow(")
			cg.expr(e.Lhs)
			cg.write(", ")
			cg.expr(e.Rhs)
			cg.write(")")

		case token.And:
			sepWithString(e.Lhs, "&&", e.Rhs)

		case token.Or:
			sepWithString(e.Lhs, "||", e.Rhs)

		case token.ShiftLeft, token.ShiftRight:
			sepWithOpLiteral(e.Lhs, e.Op, e.Rhs)

		case token.Ampersand:
			// wrap in parenthesis since & is higher precidence than in c
			cg.write("(")
			sepWithOpLiteral(e.Lhs, e.Op, e.Rhs)
			cg.write(")")

		case token.Mod:
			// TODO: calculate actual modulus
			sepWithString(e.Lhs, "%", e.Rhs)

		default:
			panic("invalid binary op: " + e.Op.String())
		}

	case *ast.UnaryOp:
		switch e.Op {
		case token.Minus, token.Bang, token.Tilde:
			cg.write(e.Op.String())
		}
		cg.expr(e.Expr)

	case *ast.CallExpr:
		// check for method call
		if dot, ok := e.Callee.(*ast.DotExpr); ok {
			fmt.Printf("callee is dot field from type %v\n", cg.Env.Expressions[dot.Expr])
			cg.method(dot, e)
		} else {
			cg.expr(e.Callee)
			cg.args(e.Args)
		}

	case *ast.DotExpr:
		cg.expr(e.Expr)
		cg.write(".")
		cg.expr(e.Field)

	case *ast.IndexExpr:
		cg.index(e.Expr, e.Index)

	case *ast.TypeExpr:
		if t, ok := cg.Env.Expressions[expr].(types.List); ok {
			cg.write(listImplementations[t])
			break
		}
		switch e.Name {
		case types.String.String():
			cg.write("char*")
		default:
			cg.write(e.Name)
		}

	default:
		panic("cannot gen expr of type " + reflect.TypeOf(expr).String())
	}
}

func (cg *codegen) stmt(stmt ast.Statement) {
	switch t := stmt.(type) {
	case *ast.ExprStmt:
		cg.expr(t.Expr)

	case *ast.TypedDeclStmt:
		cg.decl(t.Decl)

	case *ast.AssignmentStmt:
		cg.gen(t.Lhs)
		cg.write(" = ")
		cg.gen(t.Rhs)
		if _, ok := t.Lhs.(*ast.IdentExpr); !ok {
			panic("target of assignment must be identifier")
		}

	case *ast.IfStmt:
		cg.write("if (")
		cg.expr(t.If.Cond)
		cg.write(") {\n")
		cg.indented(func() {
			cg.stmt(t.If.Body)
		})
		cg.sil("}")
		if t.Else != nil {
			cg.write(" else {\n")
			cg.indented(func() {
				cg.stmt(t.Else)
			})
			cg.indent()
			cg.write("}")
		}

	case *ast.BlockStmt:
		for _, s := range t.Content {
			cg.indent()
			cg.stmt(s)
			cg.write(";\n")
		}

	case *ast.ForStmt:
		cg.write("for (")
		if t.Init != nil {
			cg.decl(t.Init)
		}
		cg.write("; ")
		cg.expr(t.Cond)
		cg.write("; ")
		if t.Incr != nil {
			cg.stmt(t.Incr)
		}
		cg.write(") {\n")
		cg.indented(func() {
			cg.stmt(t.Body)
		})
		cg.sil("}")

	case *ast.IncStmt:
		cg.expr(t.Expr)
		if t.Op.Kind == token.Increment {
			cg.write("++")
		} else if t.Op.Kind == token.Decrement {
			cg.write("--")
		}

	case *ast.BreakStmt:
		cg.write("break")

	default:
		panic("cannot gen for statement " + reflect.TypeOf(stmt).String())
	}
}

func (cg *codegen) args(exprs []ast.Expression) {
	cg.write("(")
	for i, e := range exprs {
		cg.gen(e)
		if i != len(exprs)-1 {
			cg.write(", ")
		}
	}
	cg.write(")")
}

func (cg *codegen) params(decls []*ast.ParamDecl) {
	cg.write("(")
	for i, pd := range decls {
		cg.write(pd.Type.Name + " " + pd.Ident.Lexeme)
		if i != len(decls)-1 {
			cg.write(", ")
		}
	}
	cg.write(")")
}

func (cg *codegen) decl(decl ast.Declaration) {
	switch d := decl.(type) {
	case *ast.TypedDecl:
		cg.expr(d.Type)
		cg.write(" " + d.Ident.Lexeme)
		if d.Value != nil {
			cg.write(" = ")
			cg.gen(d.Value)
		}

	case *ast.FuncDecl:
		cg.write("\n")
		if d.Name.Lexeme == "main" {
			cg.write("int main(int argc, char* argv[]) {\n")
			cg.level++
			cg.stmt(d.Body)
			cg.level--
			cg.write("}\n")
		} else {
			typ := "void"
			if d.ReturnType != nil {
				typ = d.ReturnType.Name
			}
			cg.write(typ + " " + d.Name.Lexeme)
			cg.params(d.Params)
			cg.write(" {\n")
			cg.level++
			cg.gen(d.Body)
			cg.level--
			cg.write("}\n")
		}

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
		panic("unknown interface type for ast.Node" + reflect.TypeOf(node).String())
	}
}

func (cg *codegen) program(decls []ast.Declaration) {
	for _, d := range decls {
		cg.decl(d)
	}
}
