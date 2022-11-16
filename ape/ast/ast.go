package ast

import (
	"fmt"
	"strings"
)

type Node interface{}

func NodeString(n Node) string {
	switch impl := n.(type) {
	case Declaration:
		return impl.DeclStr()
	case Statement:
		return impl.StmtStr()
	case Expression:
		return impl.ExprStr()
	default:
		return "<NODE>"
	}
}

func PrettyPrint(tree []Node) {
	for _, node := range tree {
		prettyPrint(node, 0)
	}
}

func prettyPrintStmtList(stmts []Statement, level int) {
	for _, s := range stmts {
		prettyPrint(s, level)
	}
}

func paramDeclsStr(pds []*ParamDecl) string {
	var sb strings.Builder
	for i, pd := range pds {
		sb.WriteString(fmt.Sprintf("%v %v", pd.Ident, pd.Type))
		if i != len(pds)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func prettyPrint(node Node, level int) {
	prefix := strings.Repeat("\t", level)

	switch n := node.(type) {
	case Declaration:
		switch decl := n.(type) {
		case *FuncDecl:
			fmt.Printf("%vfunc %v (%v) {\n", prefix, decl.Name, paramDeclsStr(decl.Params))
			prettyPrint(decl.Body, level+1)
			fmt.Printf("%v}\n", prefix)
		default:
			fmt.Printf("%v%v\n", prefix, decl.DeclStr())
		}

	case Statement:
		switch stmt := n.(type) {
		case *BlockStmt:
			prettyPrintStmtList(stmt.Content, level)

		case *IfStmt:
			fmt.Printf("%vif %v {\n", prefix, stmt.Cond.ExprStr())
			prettyPrint(stmt.Body, level+1)
			if stmt.Else != nil {
				fmt.Printf("%v} else {\n", prefix)
				prettyPrint(stmt.Else, level+1)
			}
			fmt.Printf("%v}\n", prefix)

		default:
			fmt.Printf("%v%v\n", prefix, stmt.StmtStr())
		}

	case Expression:
		fmt.Printf("%v%v\n", prefix, n.ExprStr())
	}
}
