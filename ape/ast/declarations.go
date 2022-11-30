package ast

import (
	"fmt"

	"github.com/pcen/ape/ape/token"
)

type Declaration interface {
	DeclStr() string
}

// function parameters
type ParamDecl struct {
	Ident token.Token
	Type  *TypeExpr
}

func (d *ParamDecl) DeclStr() string {
	return fmt.Sprintf("%v %v", d.Ident.Lexeme, d.Type)
}

type TypedDecl struct {
	Kind  token.Kind // val | var
	Ident token.Token
	Type  *TypeExpr
	Value Expression
}

func (d *TypedDecl) DeclStr() string {
	return fmt.Sprintf("(%v %v %v %v)", d.Kind, d.Ident, d.Type.ExprStr(), d.Value.ExprStr())
}

type FuncDecl struct {
	Name       token.Token
	Params     []*ParamDecl
	ReturnType *TypeExpr
	Body       *BlockStmt
}

func (d *FuncDecl) DeclStr() string {
	return fmt.Sprintf("(decl func %v)", d.Name.Lexeme)
}

type ClassDecl struct {
	Name token.Token
	Body []Declaration
}

func (d *ClassDecl) DeclStr() string {
	return fmt.Sprintf("(decl class %v)", d.Name.Lexeme)
}

type MemberDecl struct {
	Name token.Token
	Type *TypeExpr
}

func (d *MemberDecl) DeclStr() string {
	return fmt.Sprintf("(decl member %v %v)", d.Name, d.Type.ExprStr())
}

type ErrDecl struct{}

func (d *ErrDecl) DeclStr() string {
	return "err"
}
