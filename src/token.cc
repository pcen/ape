#include "token.h"

std::string getTokenTypeLexeme(TokenType tt) {
	switch (tt) {
	case TokenType::If:         return "if";
	case TokenType::Elif:       return "elif";
	case TokenType::Else:       return "else";
	case TokenType::For:        return "for";
	case TokenType::While:      return "while";
	case TokenType::Break:      return "break";
	case TokenType::Switch:     return "switch";
	case TokenType::Case:       return "case";
	case TokenType::And:        return "and";
	case TokenType::Or:         return "or";
	case TokenType::Class:      return "class";
	case TokenType::Def:        return "def";
	case TokenType::Public:     return "public";
	case TokenType::Private:    return "private";
	case TokenType::Var:        return "var";
	case TokenType::Let:        return "let";
	case TokenType::Return:     return "return";
	case TokenType::True:       return "true";
	case TokenType::False:      return "false";
	case TokenType::Int:        return "int";
	case TokenType::Int8:       return "int8";
	case TokenType::Int16:      return "int16";
	case TokenType::Int32:      return "int32";
	case TokenType::Int64:      return "int64";
	case TokenType::Uint:       return "uint";
	case TokenType::Uint8:      return "uint8";
	case TokenType::Uint16:     return "uint16";
	case TokenType::Uint32:     return "uint32";
	case TokenType::Uint64:     return "uint64";
	case TokenType::Bool:       return "bool";
	case TokenType::Float:      return "float";
	case TokenType::Double:     return "double";
	case TokenType::Char:       return "char";
	case TokenType::String:     return "string";
	case TokenType::Plus:       return "+";
	case TokenType::PlusEq:     return "+=";
	case TokenType::Minus:      return "-";
	case TokenType::MinusEq:    return "-=";
	case TokenType::Divide:     return "/";
	case TokenType::DivideEq:   return "/=";
	case TokenType::Star:       return "*";
	case TokenType::StarEq:     return "*=";
	case TokenType::Power:      return "**";
	case TokenType::PowerEq:    return "**=";
	case TokenType::Assign:     return "=";
	case TokenType::Equal:      return "==";
	case TokenType::NotEqual:   return "!=";
	case TokenType::Less:       return "<";
	case TokenType::LessEq:     return "<=";
	case TokenType::Greater:    return ">";
	case TokenType::GreaterEq:  return ">=";
	case TokenType::Bang:       return "!";
	case TokenType::Increment:  return "++";
	case TokenType::Decrement:  return "--";
	case TokenType::BitAnd:     return "&";
	case TokenType::BitOr:      return "|";
	case TokenType::BitNegate:  return "~";
	case TokenType::BitXOR:     return "^";
	case TokenType::Dot:        return ".";
	case TokenType::Comma:      return ",";
	case TokenType::OpenParen:  return "(";
	case TokenType::CloseParen: return ")";
	case TokenType::OpenBrace:  return "{";
	case TokenType::CloseBrace: return "}";
	case TokenType::OpenBrack:  return "[";
	case TokenType::CloseBrack: return "]";
	case TokenType::Sep:        return ";";
	default:
		return "";
	}
}

Token::Token(TokenType type)
	: type(type) {}

Token::Token(TokenType type, const std::string& lexeme)
	: type(type), lexeme(lexeme) {}

std::string Token::Lexeme() const {
	std::string lex = getTokenTypeLexeme(this->type);
	if (lex.empty()) {
		lex = this->lexeme;
	}
	return lex;
}
