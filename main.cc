#include <iostream>

#include "src/token.h"

int main(int argc, char* argv[]) {
	Token t(TokenType::True);
	std::cout << "TokenType::True lexeme: " << t.Lexeme() << std::endl;
	return 0;
}
