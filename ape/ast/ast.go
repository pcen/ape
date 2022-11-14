package ast

import "fmt"

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

func PrintSlice(nodes []Node) {
	for _, node := range nodes {
		PrintTree(node)
	}
}

func PrintTree(root Node) {
	switch n := root.(type) {
	case Declaration:
		fmt.Println(n.DeclStr())
	case Statement:
		fmt.Println(n.StmtStr())
	case Expression:
		fmt.Println(n.ExprStr())
	}
}
