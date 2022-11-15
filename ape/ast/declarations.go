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
	Type  string
}

func (d *ParamDecl) DeclStr() string {
	return fmt.Sprintf("%v %v", d.Ident.Lexeme, d.Type)
}

type TypedDecl struct {
	Kind  token.Kind // val | var
	Ident token.Token
	Type  string
	Value Expression
}

func (d *TypedDecl) DeclStr() string {
	return fmt.Sprintf("(%v %v %v %v)", d.Kind, d.Ident, d.Type, d.Value.ExprStr())
}

type FuncDecl struct {
	Name   token.Token
	Params []*ParamDecl
	Body   *BlockStmt
}

func (d *FuncDecl) DeclStr() string {
	// TODO: nil unsafe
	return fmt.Sprintf("func %v\n%v", d.Name.Lexeme, d.Body.StmtStr())
}

type ErrDecl struct {
}

func (d *ErrDecl) DeclStr() string {
	return "err"
}
