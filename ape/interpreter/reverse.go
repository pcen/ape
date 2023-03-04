package interpreter

import (
	"github.com/pcen/ape/ape/ast"
)

/** Wraps the value of a return statement*/
type ReturnHolder struct {
	Value value
}

/** Wraps the value of a reverse statement*/
type ReverseHolder struct {
	Value value
}

type BreadCrumb struct {
	Prev       *BreadCrumb
	SkipMarker *ast.SkipStmt
	Scope      *Scope
	Name       string
	PrevVal    value
}

func (bc BreadCrumb) Reverse() {
	bc.Scope.Set(bc.Name, bc.PrevVal)
}
