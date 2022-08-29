#include "lexer.h"

#include <fstream>

std::vector<char> readFile(const std::string& filename) {
	std::ifstream ifs(filename, std::ios::binary | std::ios::ate);
	std::streampos size = ifs.tellg();
	ifs.seekg(0, std::ios::beg);
	std::vector<char> data(size);
	ifs.read(data.data(), size);
	return data;
}

Lexer::Lexer()
	: pos(0) {}

std::vector<Token> Lexer::lex(const std::string& filename) {
	file = readFile(filename);
	std::vector<Token> tokens;
	while (true) {
		Token t = step();
		tokens.push_back(t);
		if (t.type == TokenType::Eof || t.type == TokenType::Invalid) {
			break;
		}
	}
	return tokens;
}

char Lexer::next() {
	return pos < file.size() ? file[pos++] : '\0';
}

void Lexer::back() {
	if (pos < file.size()) {
		pos--;
	}
}

void Lexer::advance() {
	if (pos < file.size()) {
		pos++;
	}
}

char Lexer::peek() {
	return pos+1 < file.size() ? file[pos+1] : '\0';
}

bool Lexer::match(char c) {
	if (peek() == c) {
		advance();
		return true;
	}
	return false;
}

Token Lexer::select(char next, TokenType onMatch, TokenType noMatch) {
	return match(next) ? Token(onMatch) : Token(noMatch);
}

// word lexes identifiers and keywords
Token Lexer::word() {
	int start = pos-1;
	while (true) {
		char c = next();
		if (!std::isalnum(c) && c != '_') {
			back();
			break;
		}
	}
	int end = pos;
	std::string lexeme(file.begin() + start, file.begin() + end);
	return Token(TokenType::Identifier, lexeme);
}

// number lexes number literals
Token Lexer::number() {
	int start = pos-1;
	while (true) {
		char c = next();
		if (!std::isdigit(c)) {
			back();
			break;
		}
	}
	int end = pos;
	std::string lexeme(file.begin() + start, file.begin() + end);
	return Token(TokenType::Number, lexeme);
}

// string lexes string literals
Token Lexer::string() {
	int start = pos; // drop start "
	while (true) {
		char c = next();
		if (c == '"') {
			// consume end "
			break;
		}
	}
	int end = pos-1;
	std::string lexeme(file.begin() + start, file.begin() + end);
	return Token(TokenType::String, lexeme);
}

Token Lexer::comment() {
	int start = pos; // drop the #
	while (next() != '\n') {}
	// drop the newline
	int end = file[pos-2] == '\r' ? pos-3 : pos-2;
	std::string lexeme(file.begin() + start, file.begin() + end);
	return Token(TokenType::Comment, lexeme);
}

void Lexer::skipWspace() {
	while (true) {
		char c = next();
		if (!std::iswspace(c)) {
			back();
			return;
		}
	}
}

Token Lexer::step() {
	skipWspace();

	char c = next();

	if (std::isalpha(c) || c == '_') {
		return word();
	} else if (std::isdigit(c) || (c == '-' && std::isdigit(peek()))) {
		return number();
	}

	switch (c) {
	case '\0':
		if (pos < file.size()) {
			// null character is invalid in source code
			return Token(TokenType::Invalid);
		}
		return Token(TokenType::Eof);

	case '#':
		return comment();

	case '"':
		return string();

	case '+':
		if (match('=')) {
			return Token(TokenType::PlusEq);
		} else if (match('+')) {
			return Token(TokenType::Increment);
		}
		return Token(TokenType::Plus);

	case '-':
		// already checked for negative number literals earlier
		if (match('=')) {
			return Token(TokenType::MinusEq);
		} else if (match('-')) {
			return Token(TokenType::Decrement);
		}
		return Token(TokenType::Minus);

	case '/':
		return select('=', TokenType::DivideEq, TokenType::Divide);

	case '*':
		if (match('*')) {
			return select('=', TokenType::PowerEq, TokenType::Power);
		}
		return select('=', TokenType::StarEq, TokenType::Star);

	case '=':
		return select('=', TokenType::Equal, TokenType::Assign);

	case '!':
		return select('=', TokenType::NotEqual, TokenType::Bang);

	case '<':
		return select('=', TokenType::LessEq, TokenType::Less);

	case '>':
		return select('=', TokenType::GreaterEq, TokenType::Greater);

	case '&':
		return Token(TokenType::BitAnd);

	case '|':
		return Token(TokenType::BitOr);

	case '~':
		return Token(TokenType::BitNegate);

	case '^':
		return Token(TokenType::BitXOR);

	case '.':
		return Token(TokenType::Dot);

	case ',':
		return Token(TokenType::Comma);

	case '(':
		return Token(TokenType::OpenParen);

	case ')':
		return Token(TokenType::CloseParen);

	case '{':
		return Token(TokenType::OpenBrace);

	case '}':
		return Token(TokenType::CloseBrace);

	case '[':
		return Token(TokenType::OpenBrack);

	case ']':
		return Token(TokenType::CloseBrack);

	case ';': // always separates statements
		return Token(TokenType::Sep);
	}

	return Token(TokenType::Invalid);
}
