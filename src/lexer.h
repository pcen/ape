#pragma once

#include "token.h"

#include <string>
#include <vector>

class ILexer {
public:
	virtual std::vector<Token> lex(const std::string& filename) = 0;
};

class Lexer : public ILexer {
public:
	Lexer();
	std::vector<Token> lex(const std::string& filename) override;

private:
	char next();
	char peek();
	void back();
	void advance();
	bool match (char);
	void skipWspace();

	Token select(char next, TokenType noMatch, TokenType onMatch);

	Token step();
	Token word();
	Token number();
	Token string();
	Token comment();

	int pos;
	std::vector<char> file;
};
