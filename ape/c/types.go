package c

import "github.com/pcen/ape/ape/types"

/*
 Code to generate the c code snippet denoting the given types.Type
*/

func (cg *codegen) typstr(typ types.Type) string {
	switch t := typ.(type) {
	case types.Primitive:
		if t.Is(types.String) {
			return "char*"
		}
		return t.String()
	case types.Named:
		panic("cannot generate c type string for named types")
	case types.List:
		switch {
		case t.Is(types.IntList):
			return vectorImplementations[t]
		case t.Is(types.StringList):
			return vectorImplementations[t]
		}
	}
	panic("cannot generate code for unknown type " + typ.String())
}
