#pragma once

#include "token.h"

#include <string>
#include <vector>

class ILexer {
public:
	ILexer();
	virtual ~ILexer();
	virtual std::vector<Token> Lex(const std::string& filename) = 0;
};
