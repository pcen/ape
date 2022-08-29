#include <iostream>
#include <vector>

#include "src/token.h"
#include "src/lexer.h"

int main(int argc, char* argv[]) {
	std::vector<Token> toks = Lexer().lex("./test/double.ape");
	for (auto& t : toks) {
		std::cout << "Token: " << t.Lexeme() << std::endl;
	}
	return 0;
}
