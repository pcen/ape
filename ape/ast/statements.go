package ast

import (
	"fmt"
	"strings"
)

type Statement interface {
	StmtStr() string
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
	sb.WriteString("{\n")
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
