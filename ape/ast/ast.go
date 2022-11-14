package ast

import "fmt"

type Node interface{}

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
