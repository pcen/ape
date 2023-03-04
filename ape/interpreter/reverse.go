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

/** Empty interface but in reality, only Value and ReverseAnnotation should be used for this */
type Reversible interface{}

type BreadCrumb struct {
	Prev       *BreadCrumb
	SkipMarker *ast.SkipStmt
	Scope      *Scope
	Name       string
	PrevVal    Reversible
}

func (bc BreadCrumb) Reverse(twi *TWI) {
	switch t := bc.PrevVal.(type) {
	case value:
		bc.Scope.Set(bc.Name, t)
	case ast.Statement:
		println("HERE")
		twi.executeStmt(t)
	}
}
