#pragma once

#include <string>

#include "ast.h"
#include "token.h"

class ParseError : std::exception {
public:
	ParseError(const std::string&);
	char* what();

private:
	char* msg;
};

class Parser {
public:
	Parser();
	Node* parse(TokenStream*);

	Expr* expression();
	Expr* comparison();
	Expr* term();
	Expr* factor();
	Expr* unary();
	Expr* primary();

	Expr* group();

private:
	void reset(TokenStream*);
	TokenStream* ts;
};
