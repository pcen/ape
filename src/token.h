#pragma once

#include <string>

enum class TokenType {
	Invalid,
	// keywords
	If,     // if
	Elif,   // elif
	Else,   // else

	For,    // for
	While,  // while
	Break,  // break

	Switch, // switch
	Case,   // case

	And,    // and
	Or,     // or

	Type,    // type
	Class,   // class
	Def,     // def
	Public,  // public
	Private, // private
	Var,     // var
	Let,     // let

	Return,  // return

	True,    // true
	False,   // false

	// built-in types
	Int,    // int
	Int8,   // int8
	Int16,  // int16
	Int32,  // int32
	Int64,  // int64
	Uint,   // uint
	Uint8,  // uint8
	Uint16, // uint16
	Uint32, // uint32
	Uint64, // uint64
	Bool,   // bool
	Float,  // float
	Double, // double
	Char,   // char
	String, // string

	// arithmetic
	Plus,     // +
	PlusEq,   // +=
	Minus,    // -
	MinusEq,  // -=
	Divide,   // /
	DivideEq, // /=
	Star,     // *
	StarEq,   // *=
	Power,    // **
	PowerEq,  // **=
	Assign,   // =

	// comparison
	Equal,     // ==
	NotEqual,  // !=
	Less,      // <
	LessEq,    // <=
	Greater,   // >
	GreaterEq, // >=

	// unary
	Bang,      // !
	Increment, // ++
	Decrement, // --

	// bitwise
	BitAnd,    // &
	BitOr,     // |
	BitNegate, // ~
	BitXOR,    // ^


	Dot,        // .
	Comma,      // ,
	OpenParen,  // (
	CloseParen, // )
	OpenBrace,  // {
	CloseBrace, // }
	OpenBrack,  // [
	CloseBrack, // ]

	Sep, // ; or \n

	Comment,
	Number,
	Identifier,
	Eof,
};

std::string getTokenTypeLexeme(TokenType);

struct Token {
	TokenType type;
	std::string lexeme;

	Token(TokenType type);
	Token(TokenType type, const std::string& lexeme);

	std::string Lexeme() const;
};
