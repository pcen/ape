#include "../src/lexer.h"

int main(int argc, char* argv[]) {
	std::vector<Token> tokens = Lexer().lex("test/double.ape");

	std::vector<Token> expected{
		Token(TokenType::Def),
		Token(TokenType::Identifier, "twice"),
		Token(TokenType::OpenParen),
		Token(TokenType::Identifier, "a"),
		Token(TokenType::Int),
		Token(TokenType::CloseParen),
		Token(TokenType::Int),
		Token(TokenType::OpenBrace),
		Token(TokenType::Return),
		Token(TokenType::Identifier, "a"),
		Token(TokenType::Star),
		Token(TokenType::Number, "2"),
		Token(TokenType::CloseBrace),
		Token(TokenType::Eof),
	};

	assert(tokens == expected);
}
