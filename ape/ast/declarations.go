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
	Ident *IdentExpr
	Type  *TypeExpr
}

func (d *ParamDecl) DeclStr() string {
	return fmt.Sprintf("%v %v", d.Ident.ExprStr(), d.Type)
}

type VarDecl struct {
	Mutable bool
	Ident   token.Token // TODO: are not sure if IdentExpr is a good idea, but should be consistent
	Type    *TypeExpr
	Value   Expression
}

func (d *VarDecl) DeclStr() string {
	mut := "<IMMUTABLE>"
	if d.Mutable {
		mut = "<MUTABLE>"
	}
	val := "<DEFAULT VALUE>"
	if d.Value != nil {
		val = d.Value.ExprStr()
	}
	typ := "<IMPLICIT TYPE>"
	if d.Type != nil {
		typ = d.Type.Name
	}
	return fmt.Sprintf("(%v %v %v %v)", mut, d.Ident, typ, val)
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
