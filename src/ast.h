#pragma once

#include "token.h"

#include <string>

struct Node {
	virtual ~Node() {}
	virtual std::string str() = 0;
};

struct Expr : public Node {
	virtual ~Expr() {}
};

// expression nodes

struct GroupExpr : public Expr {
	GroupExpr(Expr*);
	Expr* expr;
	std::string str() override;
};

struct ValueLiteral : public Expr {
	ValueLiteral(Token);
	Token tok;
	std::string str() override;
};

struct Identifier : public Expr {
	Identifier(Token);
	Token tok;
	std::string str() override;
};

struct UnaryOp : public Expr {
	UnaryOp(TokenType, Expr*);
	TokenType op;
	Expr* expr;
	std::string str() override;
};

struct BinaryOp : public Expr {
	BinaryOp(Expr*, TokenType, Expr*);
	Expr* lhs;
	TokenType op;
	Expr* rhs;
	std::string str() override;
};
