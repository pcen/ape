#pragma once

#include "token.h"

#include <string>
#include <vector>

class Lexer {
public:
	Lexer();
	~Lexer();
	TokenStream* lex(const std::string& filename);
	TokenStream* lexString(const std::string& source);

private:
	void reset();

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
