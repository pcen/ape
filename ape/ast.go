package ape

import "fmt"

type Expression interface {
	ExprStr() string
}

type Statement interface {
	StmtStr() string
}

// Expressions

type GroupExpr struct {
	Expr Expression
}

type LiteralExpr struct {
	Token
}

func (e LiteralExpr) ExprStr() string {
	return fmt.Sprintf("%v", e.Token)
}

func NewLiteralExpr(token Token) Expression {
	return LiteralExpr{Token: token}
}

type IdentExpr struct {
	Token
}

func (e IdentExpr) ExprStr() string {
	return fmt.Sprintf("(%v)", e.Token)
}

func NewIdentExpr(token Token) Expression {
	return IdentExpr{Token: token}
}

type UnaryOp struct {
	Op   TokenType
	Expr Expression
}

func (e UnaryOp) ExprStr() string {
	return fmt.Sprintf("(%v %v)", e.Op, e.Expr.ExprStr())
}

func NewUnaryOp(op TokenType, expr Expression) Expression {
	return UnaryOp{Op: op, Expr: expr}
}

type BinaryOp struct {
	Lhs Expression
	Op  TokenType
	Rhs Expression
}

func (e BinaryOp) ExprStr() string {
	return fmt.Sprintf("(%v %v %v)", e.Op, e.Lhs.ExprStr(), e.Rhs.ExprStr())
}

func NewBinaryOp(lhs Expression, op TokenType, rhs Expression) Expression {
	return BinaryOp{Lhs: lhs, Op: op, Rhs: rhs}
}
