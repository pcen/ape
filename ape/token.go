package ape

type TokenType int

const (
	Invalid TokenType = iota + 1
	// keywords
	If   // if
	Elif // elif
	Else // else

	For   // for
	While // while
	Break // break

	Switch // switch
	Case   // case

	And // and
	Or  // or

	Type    // type
	Class   // class
	Def     // def
	Public  // public
	Private // private
	Val     // val
	Var     // var

	Return // return

	True  // true
	False // false

	Module // module
	Import // import

	// built-in types
	// TODO: probably do not need these
	//       since types are identifiers
	Int    // int
	Int8   // int8
	Int16  // int16
	Int32  // int32
	Int64  // int64
	Uint   // uint
	Uint8  // uint8
	Uint16 // uint16
	Uint32 // uint32
	Uint64 // uint64
	Bool   // bool
	Float  // float
	Double // double
	Char   // char
	String // string

	// arithmetic
	Plus     // +
	PlusEq   // +=
	Minus    // -
	MinusEq  // -=
	Divide   // /
	DivideEq // /=
	Star     // *
	StarEq   // *=
	Power    // **
	PowerEq  // **=
	Assign   // =

	// comparison
	Equal     // ==
	NotEqual  // !=
	Less      // <
	LessEq    // <=
	Greater   // >
	GreaterEq // >=

	// unary
	Bang      // !
	Increment // ++
	Decrement // --

	// bitwise
	BitAnd    // &
	BitOr     // |
	BitNegate // ~
	BitXOR    // ^

	Dot        // .
	Comma      // ,
	OpenParen  // (
	CloseParen // )
	OpenBrace  // {
	CloseBrace // }
	OpenBrack  // [
	CloseBrack // ]

	Sep // ; or \n

	Comment
	Number
	Identifier
	Eof
)

func GetKeyword(identifier string) (TokenType, bool) {
	keywords := map[string]TokenType{
		"if":      If,
		"elif":    Elif,
		"else":    Else,
		"for":     For,
		"while":   While,
		"break":   Break,
		"switch":  Switch,
		"case":    Case,
		"and":     And,
		"or":      Or,
		"type":    Type,
		"class":   Class,
		"def":     Def,
		"public":  Public,
		"private": Private,
		"val":     Val,
		"var":     Var,
		"return":  Return,
		"true":    True,
		"false":   False,
		"module":  Module,
		"import":  Import,
	}
	if tt, ok := keywords[identifier]; ok {
		return tt, true
	}
	return Invalid, false
}

var (
	tokenLexemes = []string{
		Invalid: "<INVALID>",

		If:   "if",
		Elif: "elif",
		Else: "else",

		For:   "for",
		While: "while",
		Break: "break",

		Switch: "switch",
		Case:   "case",

		And: "and",
		Or:  "or",

		Type:    "type",
		Class:   "class",
		Def:     "def",
		Public:  "public",
		Private: "private",
		Val:     "val",
		Var:     "var",

		Return: "return",

		True:  "true",
		False: "false",

		Module: "module",
		Import: "import",

		Int:    "int",
		Int8:   "int8",
		Int16:  "int16",
		Int32:  "int32",
		Int64:  "int64",
		Uint:   "uint",
		Uint8:  "uint8",
		Uint16: "uint16",
		Uint32: "uint32",
		Uint64: "uint64",
		Bool:   "bool",
		Float:  "float",
		Double: "double",
		Char:   "char",
		String: "string",

		Plus:     "+",
		PlusEq:   "+=",
		Minus:    "-",
		MinusEq:  "-=",
		Divide:   "/",
		DivideEq: "/=",
		Star:     "*",
		StarEq:   "*=",
		Power:    "**",
		PowerEq:  "**=",
		Assign:   "=",

		Equal:     "==",
		NotEqual:  "!=",
		Less:      "<",
		LessEq:    "<=",
		Greater:   ">",
		GreaterEq: ">=",

		Bang:      "!",
		Increment: "++",
		Decrement: "--",

		BitAnd:    "&",
		BitOr:     "|",
		BitNegate: "~",
		BitXOR:    "^",

		Dot:        ".",
		Comma:      ",",
		OpenParen:  "(",
		CloseParen: ")",
		OpenBrace:  "{",
		CloseBrace: "}",
		OpenBrack:  "[",
		CloseBrack: "]",

		Sep: ";",

		Eof: "<EOF>",
	}
)

func (tt TokenType) String() string {
	if tt == 0 {
		panic("no string for 0 initialized TokenType")
	}
	return tokenLexemes[tt]
}

type Token struct {
	Type   TokenType
	Lexeme string
}

func (t Token) String() string {
	if t.Lexeme != "" {
		return t.Lexeme
	}
	return t.Type.String()
}

func NewToken(tokenType TokenType) Token {
	return Token{Type: tokenType}
}

func NewLexemeToken(tokenType TokenType, lexeme string) Token {
	return Token{Type: tokenType, Lexeme: lexeme}
}
