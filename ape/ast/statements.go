package ast

import (
	"fmt"
	"strings"

	"github.com/pcen/ape/ape/token"
)

type Statement interface {
	StmtStr() string
}

type ErrStmt struct{}

func (s *ErrStmt) StmtStr() string {
	return "err"
}

type BlockStmt struct {
	Content []Statement
}

func (s *BlockStmt) StmtStr() string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, content := range s.Content {
		sb.WriteString(fmt.Sprintf("\t%v\n", content.StmtStr()))
	}
	sb.WriteString("}")
	return sb.String()
}

type ExprStmt struct {
	Expr Expression
}

func (s *ExprStmt) StmtStr() string {
	return s.Expr.ExprStr()
}

type ReturnStmt struct {
	Expr Expression
}

func (s *ReturnStmt) StmtStr() string {
	return fmt.Sprintf("(return %v)", s.Expr.ExprStr())
}

type TypedDeclStmt struct {
	Decl *TypedDecl
}

func (s *TypedDeclStmt) StmtStr() string {
	return s.Decl.DeclStr()
}

// Compound Statements

type IfStmt struct {
	Cond Expression
	Body Statement
	Else Statement
}

func (s *IfStmt) StmtStr() string {
	if s.Else != nil {
		return fmt.Sprintf("if %v then %v else %v", s.Cond.ExprStr(), s.Body.StmtStr(), s.Else.StmtStr())
	}
	return fmt.Sprintf("if %v then %v", s.Cond.ExprStr(), s.Body.StmtStr())
}

// Simple Statements

type IncStmt struct {
	Expr Expression
	Op   token.Token
}

func (s *IncStmt) StmtStr() string {
	return fmt.Sprintf("(%v %v)", s.Expr.ExprStr(), s.Op)
}

type AssignmentStmt struct {
	Lhs Expression
	Op  token.Token
	Rhs Expression
}

func (s *AssignmentStmt) StmtStr() string {
	return fmt.Sprintf("(%v %v %v)", s.Op, s.Lhs.ExprStr(), s.Rhs.ExprStr())
}
