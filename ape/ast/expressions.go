package ast

import (
	"fmt"
	"strings"

	"github.com/pcen/ape/ape/token"
)

type Expression interface {
	ExprStr() string
}

func exprListStr(exprs []Expression) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, expr := range exprs {
		sb.WriteString(expr.ExprStr())
		if i != len(exprs)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

type ErrExpr struct {
	What string
}

func (e *ErrExpr) ExprStr() string {
	return fmt.Sprintf("error: %v", e.What)
}

type GroupExpr struct {
	Expr Expression
}

func (e *GroupExpr) ExprStr() string {
	return fmt.Sprintf("(%v)", e.Expr.ExprStr())
}

type LiteralExpr struct {
	token.Token
}

func (e *LiteralExpr) ExprStr() string {
	return fmt.Sprintf("%v", e.Token)
}

func NewLiteralExpr(token token.Token) Expression {
	return &LiteralExpr{Token: token}
}

type IdentExpr struct {
	Ident token.Token
}

func (e *IdentExpr) ExprStr() string {
	return fmt.Sprintf("%v", e.Ident)
}

func NewIdentExpr(token token.Token) *IdentExpr {
	return &IdentExpr{Ident: token}
}

type UnaryOp struct {
	Op   token.Kind
	Expr Expression
}

func (e *UnaryOp) ExprStr() string {
	return fmt.Sprintf("(%v %v)", e.Op, e.Expr.ExprStr())
}

func NewUnaryOp(op token.Kind, expr Expression) Expression {
	return &UnaryOp{Op: op, Expr: expr}
}

type BinaryOp struct {
	Lhs Expression
	Op  token.Token
	Rhs Expression
}

func (e *BinaryOp) ExprStr() string {
	return fmt.Sprintf("(%v %v %v)", e.Op, e.Lhs.ExprStr(), e.Rhs.ExprStr())
}

func NewBinaryOp(lhs Expression, op token.Token, rhs Expression) Expression {
	return &BinaryOp{Lhs: lhs, Op: op, Rhs: rhs}
}

type CallExpr struct {
	Callee Expression
	Args   []Expression
}

func (e *CallExpr) ExprStr() string {
	return fmt.Sprintf("(%v() %v)", e.Callee.ExprStr(), exprListStr(e.Args))
}

type DotExpr struct {
	Expr  Expression
	Field *IdentExpr
}

func (e *DotExpr) ExprStr() string {
	return fmt.Sprintf("(%v.%v)", e.Expr.ExprStr(), e.Field.ExprStr())
}

type IndexExpr struct {
	Expr  Expression
	Index Expression
}

func (e *IndexExpr) ExprStr() string {
	return fmt.Sprintf("(%v[%v])", e.Expr.ExprStr(), e.Index.ExprStr())
}

type TypeExpr struct {
	Name string
	List bool
}

func (e *TypeExpr) ExprStr() string {
	if e.List {
		return fmt.Sprint("[]", e.Name)
	}
	return e.Name
}

type LitListExpr struct {
	Elements []Expression
}

func (e *LitListExpr) ExprStr() string {
	return fmt.Sprint("(", exprListStr(e.Elements), ")")
}

type LitMapExpr struct {
	Elements map[Expression]Expression
}

func (e *LitMapExpr) ExprStr() string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for k, v := range e.Elements {
		sb.WriteString(fmt.Sprintf("\t%v: %v", k.ExprStr(), v.ExprStr()))
	}
	sb.WriteString("}")
	return sb.String()
}
