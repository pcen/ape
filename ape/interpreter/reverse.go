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
		switch v_type := t.(type) {
		case val_index_val_pair:
			bc.Scope.Get(bc.Name).(val_map).Data[v_type.Index] = v_type.Value
		default:
			bc.Scope.Set(bc.Name, t)
		}
	case ast.Statement:
		twi.executeStmt(t)
	}
}
