#include "parser.h"

#include <cstring>
#include <iostream> // error reporting

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" | "|" | "^" ) factor )* ;
// factor         → unary ( ( "/" | "*" | "&" ) unary )* ;
// unary          → ( "!" | "-" | "~" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | group ;
// group          → "(" expression ")" ;

ParseError::ParseError(const std::string& message)
	: msg(new char[message.size()])
{
	std::strcpy(msg, message.c_str());
}

char* ParseError::what() {
	return msg;
}

Parser::Parser() {}

Node* Parser::parse(TokenStream* stream) {
	reset(stream);
}

void Parser::reset(TokenStream* stream) {
	delete this->ts;
	this->ts = stream;
}

Expr* Parser::expression() {
	try {
		return comparison();
	} catch (ParseError& e) {
		std::cerr << "error parsing expression: " << e.what() << std::endl;
	}
	return nullptr; // TODO: return bad expression node with error info
}

Expr* Parser::comparison() {
	Expr* lhs = term();
	while (ts->match({TokenType::Equal, TokenType::NotEqual})) {
		lhs = new BinaryOp(lhs, ts->prev().type, term());
	}
	return lhs;
}

Expr* Parser::term() {
	Expr* lhs = factor();
	while (ts->match({TokenType::Plus, TokenType::Minus, TokenType::BitOr, TokenType::BitXOR})) {
		lhs = new BinaryOp(lhs, ts->prev().type, factor());
	}
	return lhs;
}

Expr* Parser::factor() {
	Expr* lhs = unary();
	while (ts->match({TokenType::Divide, TokenType::Star, TokenType::BitAnd})) {
		lhs = new BinaryOp(lhs, ts->prev().type, unary());
	}
	return lhs;
}

Expr* Parser::unary() {
	if (ts->match({TokenType::Bang, TokenType::Minus, TokenType::BitNegate})) {
		return new UnaryOp(ts->prev().type, unary());
	}
	return primary();
}

Expr* Parser::primary() {
	switch (ts->peek().type) {
	case TokenType::Number:
	case TokenType::String:
	case TokenType::True:
	case TokenType::False:
		return new ValueLiteral(ts->next());
	case TokenType::Identifier:
		return new Identifier(ts->next());
	case TokenType::OpenParen:
		return group();
	}
	throw ParseError("invalid token type for primary expression");
}

Expr* Parser::group() {
	if (ts->match(TokenType::OpenParen)) {
		Expr* expr = expression();
		if (ts->match(TokenType::CloseParen)) {
			return new GroupExpr(expr);
		}
		throw ParseError("group expression missing closing parenthesis");
	}
	// this is a programming error in the compiler
	throw ParseError("group must start with opening parenthesis");
}
