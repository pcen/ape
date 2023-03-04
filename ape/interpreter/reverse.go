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

/** A reverse annotation holds a statement that is hopefully the inverse of some other statement */
type ReverseAnnotation struct {
	stmt *ast.Statement
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
	case ReverseAnnotation:
		twi.executeStmt(*t.stmt)
	}
}
