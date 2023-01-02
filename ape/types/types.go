package types

import (
	"fmt"

	"github.com/pcen/ape/ape/ast"
)

type Type interface {
	String() string
	Underlying() Type
	Is(Type) bool
}

type Environment struct {
	Expressions map[ast.Expression]Type
}

type Primitive int

const (
	Invalid Primitive = iota + 1
	Undefined
	Void
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Bool
	Float
	Double
	Char
	String
)

func (p Primitive) String() string {
	return primitives[p]
}

func (p Primitive) Underlying() Type {
	return p
}

func (p Primitive) Is(other Type) bool {
	o, ok := other.(Primitive)
	return ok && p == o
}

var (
	primitives = map[Primitive]string{
		Invalid:   "<INVALID TYPE>",
		Undefined: "<UNDEFINED TYPE>",
		Void:      "<VOID>",
		Int:       "int",
		Int8:      "int8",
		Int16:     "int16",
		Int32:     "int32",
		Int64:     "int64",
		Uint:      "uint",
		Uint8:     "uint8",
		Uint16:    "uint16",
		Uint32:    "uint32",
		Uint64:    "uint64",
		Bool:      "bool",
		Float:     "float",
		Double:    "double",
		Char:      "char",
		String:    "string",
	}
)

type Named struct {
	name string
}

func NewNamed(name string) Type {
	return Named{name: name}
}

func (n Named) String() string {
	return n.name
}

func (n Named) Underlying() Type {
	return n
}

func (n Named) Is(other Type) bool {
	if o, ok := other.(Named); ok {
		return n.name == o.name
	}
	return false
}

type Function struct {
	Params  []Type
	Returns []Type
}

func NewFunction(params []Type, returns []Type) Type {
	return Function{Params: params, Returns: returns}
}

func (f Function) String() string {
	return fmt.Sprintf("func %v -> %v", f.Params, f.Returns)
}

func (f Function) Underlying() Type {
	return f
}

func typeSlicesEqual(x, y []Type) bool {
	if len(x) == len(y) {
		for i := range x {
			if !x[i].Is(y[i]) {
				return false
			}
		}
	}
	return true
}

func (f Function) Is(other Type) bool {
	if o, ok := other.(Function); ok {
		return typeSlicesEqual(f.Params, o.Params) && typeSlicesEqual(f.Returns, o.Returns)
	}
	return false
}

type List struct {
	Data Type
}

func NewList(t Type) Type {
	return List{Data: t}
}

func (l List) Is(other Type) bool {
	o, ok := other.(List)
	return ok && l.Data.Is(o.Data)
}

func (l List) String() string {
	return fmt.Sprintf("list %v", l.Data)
}

func (l List) Underlying() Type {
	return l
}

// lists
var (
	IntList    = NewList(Int)
	StringList = NewList(String)
)

// assert all types implement Type interface
var (
	_ Type = Invalid
	_ Type = Named{}
	_ Type = Function{}
	_ Type = List{}
)

var (
	typeLookup = map[string]Primitive{
		"int":    Int,
		"int8":   Int8,
		"int16":  Int16,
		"int32":  Int32,
		"int64":  Int64,
		"uint":   Uint,
		"uint8":  Uint8,
		"uint16": Uint16,
		"uint32": Uint32,
		"uint64": Uint64,
		"bool":   Bool,
		"float":  Float,
		"double": Double,
		"char":   Char,
		"string": String,
	}
)

func LookupPrimitive(lexeme string) Primitive {
	return typeLookup[lexeme]
}
