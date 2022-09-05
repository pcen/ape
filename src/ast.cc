#include "ast.h"

#include <sstream>

// concatenates args into a string using streamstring's << operator
template<typename ...Args> std::string concat(Args&& ...args) {
	std::stringstream ss;
	(ss << ... << args);
	return ss.str();
}

// expression nodes

GroupExpr::GroupExpr(Expr* expr)
	: expr(expr) {}

std::string GroupExpr::str() {
	return concat("(", expr->str(), ")");
}


ValueLiteral::ValueLiteral(Token token)
	: tok(token) {}

std::string ValueLiteral::str() {
	return tok.lexeme;
}


Identifier::Identifier(Token token)
	: tok(token) {}

std::string Identifier::str() {
	return tok.lexeme;
}


UnaryOp::UnaryOp(TokenType op, Expr* expr)
	: op(op), expr(expr) {}

std::string UnaryOp::str() {
	return concat("(", getTokenTypeLexeme(op), " ", expr->str(), ")");
}


BinaryOp::BinaryOp(Expr* lhs, TokenType op, Expr* rhs)
	: lhs(lhs), op(op), rhs(rhs) {}

std::string BinaryOp::str() {
	return concat("(", getTokenTypeLexeme(op), " ", lhs->str(), " ", rhs->str(), ")");
}
