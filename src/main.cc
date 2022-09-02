#include <iostream>

#include "token.h"
#include "lexer.h"

int main(int argc, char* argv[]) {
	auto toks = Lexer().lex("./test/double.ape");
	while (!toks->done()) {
		Token t = toks->next();
		std::cout << "Token: " << t.Lexeme() << std::endl;
	}
	return 0;
}
