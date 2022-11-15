package token

import "fmt"

type Kind int

const (
	Invalid Kind = iota + 1
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
	Func    // func
	Public  // public
	Private // private
	Val     // val
	Var     // var

	Return // return

	True  // true
	False // false

	Module // module
	Import // import

	String // string literal

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
		Func:    "func",
		Public:  "public",
		Private: "private",
		Val:     "val",
		Var:     "var",

		Return: "return",

		True:  "true",
		False: "false",

		Module: "module",
		Import: "import",

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

	keywords = map[string]Kind{
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
		"func":    Func,
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
)

func GetKeyword(identifier string) (Kind, bool) {
	if tt, ok := keywords[identifier]; ok {
		return tt, true
	}
	return Invalid, false
}

func (tt Kind) String() string {
	if tt == 0 {
		panic("no string for 0 initialized TokenType")
	}
	return tokenLexemes[tt]
}

type Position struct {
	Line   uint
	Column uint
}

func (p Position) String() string {
	return fmt.Sprintf("%v:%v", p.Line, p.Column)
}

type Token struct {
	Kind   Kind
	Lexeme string
	Position
}

func (t Token) String() string {
	if t.Lexeme != "" {
		return t.Lexeme
	}
	return t.Kind.String()
}

func New(tokenType Kind, position Position) Token {
	return Token{
		Kind:     tokenType,
		Position: position,
	}
}

func NewLexeme(tokenType Kind, lexeme string, position Position) Token {
	return Token{
		Kind:     tokenType,
		Lexeme:   lexeme,
		Position: position,
	}
}