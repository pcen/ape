package types

import "fmt"

type PrimitaveType uint

type Type interface {
	String() string
}

const (
	Invalid PrimitaveType = iota + 1
	Undefined
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

	Func
)

var (
	primitaveTypes = []PrimitaveType{
		Int,
		Int8,
		Int16,
		Int32,
		Int64,
		Uint,
		Uint8,
		Uint16,
		Uint32,
		Uint64,
		Bool,
		Float,
		Double,
		Char,
		String,

		Func,
	}
)

type NamedType struct {
	name string
}

func NewNamedType(name string) NamedType {
	return NamedType{name: name}
}

func (nt NamedType) String() string {
	return nt.name
}

// assert all types implement Type interface
var (
	_ Type = Invalid
	_ Type = NamedType{}
)

var (
	typeNames = []string{
		Invalid:   "<INVALID TYPE>",
		Undefined: "<UNDEFINED TYPE>",
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
		Func:      "func",
	}

	typeLookup = map[string]PrimitaveType{
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

func LookupPrimitive(lexeme string) PrimitaveType {
	return typeLookup[lexeme]
}

func (p PrimitaveType) String() string {
	return typeNames[p]
}

func Same(t1, t2 Type) bool {
	switch t1 := t1.(type) {
	case PrimitaveType:
		if t2, ok := t2.(PrimitaveType); ok {
			return t1 == t2
		}
	case NamedType:
		fmt.Println("named")
		if t2, ok := t2.(NamedType); ok {
			return t1.name == t2.name
		}
	}
	return false
}
