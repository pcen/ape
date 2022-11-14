package types

type TypeKind uint

type Type struct {
	Kind TypeKind
	Name string
}

func NewType(t TypeKind) Type {
	return Type{Kind: t}
}

const (
	Invalid TypeKind = iota + 1
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
	Class
)

var (
	typeNames = []string{
		Invalid: "<INVALID TYPE>",
		Int:     "int",
		Int8:    "int8",
		Int16:   "int16",
		Int32:   "int32",
		Int64:   "int64",
		Uint:    "uint",
		Uint8:   "uint8",
		Uint16:  "uint16",
		Uint32:  "uint32",
		Uint64:  "uint64",
		Bool:    "bool",
		Float:   "float",
		Double:  "double",
		Char:    "char",
		String:  "string",
		Class:   "class",
	}

	typeLookup = map[string]TypeKind{
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

func LookupPrimitive(lexeme string) TypeKind {
	return typeLookup[lexeme]
}

func (tk TypeKind) String() string {
	return typeNames[tk]
}

func (t Type) String() string {
	if t.Kind == Class {
		return t.Name
	}
	return t.Kind.String()
}

func (t Type) Is(kind TypeKind) bool {
	return t.Kind == kind
}
