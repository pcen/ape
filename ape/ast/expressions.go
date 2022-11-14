package ast

import (
	"fmt"

	"github.com/pcen/ape/ape/token"
)

type Expression interface {
	ExprStr() string
}

type InvalidExpr struct {
	What string
}

func (e *InvalidExpr) ExprStr() string {
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
	token.Token
}

func (e *IdentExpr) ExprStr() string {
	return fmt.Sprintf("%v", e.Token)
}

func NewIdentExpr(token token.Token) Expression {
	return &IdentExpr{Token: token}
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
	Op  token.Kind
	Rhs Expression
}

func (e *BinaryOp) ExprStr() string {
	return fmt.Sprintf("(%v %v %v)", e.Op, e.Lhs.ExprStr(), e.Rhs.ExprStr())
}

func NewBinaryOp(lhs Expression, op token.Kind, rhs Expression) Expression {
	return &BinaryOp{Lhs: lhs, Op: op, Rhs: rhs}
}

type CallExpr struct {
	Callee Expression
	Args   []Expression
}

func (e *CallExpr) ExprStr() string {
	return fmt.Sprintf("(%v() %v)", e.Callee.ExprStr(), e.Args)
}
