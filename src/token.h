#pragma once

#include <string>
#include <vector>
#include <array>

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

	Module,  // module
	Import,  // import

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

std::string getTokenTypeLexeme(TokenType tt);

struct Token {
	TokenType type;
	std::string lexeme;

	Token(TokenType type);
	Token(TokenType type, const std::string& lexeme);

	std::string Lexeme() const;

	friend bool operator== (const Token&, const Token&);
	friend bool operator!= (const Token&, const Token&);
};

// Wrap token stream behind a class so it's easy to switch
// from lexing the whole file to lexing asynchronously as
// the next token is requested (ie. for repl)
class TokenStream {
public:
	virtual ~TokenStream() {}
	virtual bool done() = 0;
	virtual Token next() = 0;
	virtual Token peek() = 0;
	virtual bool match(TokenType) = 0;
	virtual bool match(std::initializer_list<TokenType>) = 0;
	virtual Token prev() = 0;
	virtual std::vector<Token> readFull() = 0;
};

class VectorTokenStream : public TokenStream {
public:
	VectorTokenStream(std::vector<Token>&& tokens);
	bool done() override;
	Token next() override;
	Token peek() override;
	bool match(TokenType) override;
	bool match(std::initializer_list<TokenType>) override;
	Token prev() override;
	std::vector<Token> readFull() override;

private:
	std::size_t pos;
	std::vector<Token> tokens;
};
