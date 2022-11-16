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

type CondBlockStmt struct {
	Cond Expression
	Body *BlockStmt
}

type IfStmt struct {
	If    *CondBlockStmt
	Elifs []*CondBlockStmt
	Else  *BlockStmt
}

func (s *IfStmt) StmtStr() string {
	return fmt.Sprintf("(if %v)", s.If.Cond.ExprStr())
}

// ForStmt represents both for and while loops
type ForStmt struct {
	Init Declaration
	Cond Expression
	Incr Statement
	Body *BlockStmt
}

func (s *ForStmt) StmtStr() string {
	if s.Init == nil {
		return fmt.Sprintf("(while %v)", s.Cond.ExprStr())
	}
	return fmt.Sprintf("(for %v)", s.Cond.ExprStr())
}

// Simple Statements

type IncStmt struct {
	Expr Expression
	Op   token.Token
}

func (s *IncStmt) StmtStr() string {
	return fmt.Sprintf("(%v %v)", s.Op, s.Expr.ExprStr())
}

type AssignmentStmt struct {
	Lhs Expression
	Op  token.Token
	Rhs Expression
}

func (s *AssignmentStmt) StmtStr() string {
	return fmt.Sprintf("(%v %v %v)", s.Op, s.Lhs.ExprStr(), s.Rhs.ExprStr())
}
