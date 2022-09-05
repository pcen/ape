#include <iostream>

#include "token.h"
#include "lexer.h"
#include "ast.h"

int main(int argc, char* argv[]) {
	Token t(TokenType::Identifier, "foo");
	Identifier* id = new Identifier(t);
	ValueLiteral* val = new ValueLiteral(Token(TokenType::Number, "2"));
	BinaryOp bo(id, TokenType::Star, val);

	std::cout << bo.str() << std::endl;
	delete id;
	delete val;
	// auto toks = Lexer().lex("./test/double.ape");
	// while (!toks->done()) {
	// 	Token t = toks->next();
	// 	std::cout << "Token: " << t.Lexeme() << std::endl;
	// }
	return 0;
}
