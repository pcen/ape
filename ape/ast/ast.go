package ast

import (
	"fmt"
	"strings"
)

type Node interface{}

type File struct {
	Path   string
	Module string
	Ast    []Declaration
}

func NewFile(path string) *File {
	return &File{
		Path: path,
		Ast:  make([]Declaration, 0),
	}
}

func PrettyPrint(tree []Declaration) {
	p := &printer{}
	for _, node := range tree {
		p.level = -1
		p.prettyPrint(node)
	}
}

type printer struct {
	level int
}

func (p *printer) printf(format string, a ...interface{}) {
	pf := strings.Repeat("\t", p.level)
	fmt.Print(pf)
	fmt.Printf(format, a...)
}

func (p *printer) printStmtList(stmts []Statement) {
	for _, s := range stmts {
		p.prettyPrint(s)
	}
}

func (p *printer) prettyPrint(node Node) {
	p.level++
	switch n := node.(type) {
	case Declaration:
		switch decl := n.(type) {
		case *FuncDecl:
			p.printf("func %v (%v) {\n", decl.Name, paramDeclsStr(decl.Params))
			p.prettyPrint(decl.Body)
			p.printf("}\n\n")

		case *ClassDecl:
			p.printf("class %v {\n", decl.Name)
			p.printf("}\n\n")

		default:
			p.printf("%v\n", decl.DeclStr())
		}

	case Statement:
		switch stmt := n.(type) {
		case *BlockStmt:
			p.printStmtList(stmt.Content)

		case *IfStmt:
			p.printf("if %v {\n", stmt.If.Cond.ExprStr())
			p.prettyPrint(stmt.If.Body)
			if len(stmt.Elifs) > 0 {
				for _, b := range stmt.Elifs {
					p.printf("} elif %v {\n", b.Cond.ExprStr())
					p.prettyPrint(b.Body)
				}
			}
			if stmt.Else != nil {
				p.printf("} else {\n")
				p.prettyPrint(stmt.Else)
			}
			p.printf("}\n")

		case *ForStmt:
			if stmt.Init == nil {
				p.printf("while %v {\n", stmt.Cond.ExprStr())
			} else {
				p.printf("for %v; %v; %v {\n", stmt.Init.DeclStr(), stmt.Cond.ExprStr(), stmt.Incr.StmtStr())
			}
			p.prettyPrint(stmt.Body)
			p.printf("}\n")

		default:
			p.printf("%v\n", stmt.StmtStr())
		}

	case Expression:
		p.printf("%v\n", n.ExprStr())
	}
	p.level--
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
