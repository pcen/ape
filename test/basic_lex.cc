#include "../src/lexer.h"
#include "../external/googletest/googletest/include/gtest/gtest.h"

TEST(Lexer, DoubleDotApe) {
	TokenStream* tokens = Lexer().lex("test/double.ape");
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

	ASSERT_EQ(tokens->readFull(), expected);
}

int main(int argc, char** argv) {
	testing::InitGoogleTest(&argc, argv);
	return RUN_ALL_TESTS();
}
