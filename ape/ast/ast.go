package ast

import (
	"fmt"
	"strings"
)

type Node interface {
}

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

func pf(level int) string {
	return strings.Repeat("\t", level)
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
	pfx := pf(level)

	switch n := node.(type) {
	case Declaration:
		switch decl := n.(type) {
		case *FuncDecl:
			fmt.Printf("%vfunc %v (%v) {\n", pfx, decl.Name, paramDeclsStr(decl.Params))
			prettyPrint(decl.Body, level+1)
			fmt.Printf("%v}\n\n", pfx)
		case *ClassDecl:
			fmt.Printf("%vclass %v {\n", pfx, decl.Name)
			fmt.Printf("%v}\n\n", pfx)
		default:
			fmt.Printf("%v%v\n", pfx, decl.DeclStr())
		}

	case Statement:
		switch stmt := n.(type) {
		case *BlockStmt:
			prettyPrintStmtList(stmt.Content, level)

		case *IfStmt:
			fmt.Printf("%vif %v {\n", pfx, stmt.If.Cond.ExprStr())
			prettyPrint(stmt.If.Body, level+1)
			if len(stmt.Elifs) > 0 {
				for _, b := range stmt.Elifs {
					fmt.Printf("%v} elif %v {\n", pfx, b.Cond.ExprStr())
					prettyPrint(b.Body, level+1)
				}
			}
			if stmt.Else != nil {
				fmt.Printf("%v} else {\n", pfx)
				prettyPrint(stmt.Else, level+1)
			}
			fmt.Printf("%v}\n", pfx)

		case *ForStmt:
			if stmt.Init == nil {
				fmt.Printf("%vwhile %v {\n", pfx, stmt.Cond.ExprStr())
			} else {
				fmt.Printf("%vfor %v; %v; %v {\n", pfx, stmt.Init.DeclStr(), stmt.Cond.ExprStr(), stmt.Incr.StmtStr())
			}
			prettyPrint(stmt.Body, level+1)
			fmt.Printf("%v}\n", pfx)

		default:
			fmt.Printf("%v%v\n", pfx, stmt.StmtStr())
		}

	case Expression:
		fmt.Printf("%v%v\n", pfx, n.ExprStr())
	}
}
