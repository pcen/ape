package ape

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

func isalpha(b byte) bool {
	return unicode.IsLetter(rune(b))
}

func isdigit(b byte) bool {
	return unicode.IsDigit(rune(b))
}

func iswspace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

type Lexer interface {
	LexFile(string) []Token
}

func NewLexer() Lexer {
	return &lexer{}
}

type lexer struct {
	r   bufio.Reader
	pos uint
}

func (l *lexer) LexFile(file string) []Token {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	l.r = *bufio.NewReader(f)
	var tokens []Token
	for {
		token := l.step()
		tokens = append(tokens, token)
		if token.Type == Eof {
			break
		}
	}
	return tokens
}

func (l *lexer) next() (byte, bool) {
	b, err := l.r.ReadByte()
	if err == io.EOF {
		return 0, false
	}
	l.pos++
	return b, true
}

func (l *lexer) back() {
	if err := l.r.UnreadByte(); err != nil {
		panic(err)
	}
	l.pos--
}

func (l *lexer) peek() byte {
	if b, err := l.r.Peek(1); err == nil {
		return b[0]
	}
	return 0
}

func (l *lexer) match(b byte) bool {
	return l.peek() == b
}

func (l *lexer) pick(b byte, onMatch TokenType, noMatch TokenType) Token {
	if l.match(b) {
		return NewToken(onMatch)
	}
	return NewToken(noMatch)
}

func (l *lexer) skipWhiteSpace() {
	for {
		b, ok := l.next()
		if !ok {
			return
		}
		if !iswspace(b) {
			l.back()
			return
		}
	}
}

func (l *lexer) identifier() Token {
	buf := make([]byte, 0, 16)
	for {
		b, _ := l.next()
		if !isalpha(b) && !isdigit(b) && b != '_' {
			l.back()
			break
		}
		buf = append(buf, b)
	}
	lexeme := string(buf)
	if tt, isKeyword := GetKeyword(lexeme); isKeyword {
		return NewToken(tt)
	}
	return NewLexemeToken(Identifier, lexeme)
}

func (l *lexer) number() Token {
	buf := make([]byte, 0, 16)
	dot := false
	for {
		b, _ := l.next()
		if !isdigit(b) || (b == '.' && dot) {
			l.back()
			break
		}
		if b == '.' {
			dot = true
		}
		buf = append(buf, b)
	}
	return NewLexemeToken(Number, string(buf))
}

func (l *lexer) comment() Token {
	buf := make([]byte, 0, 16)
	for {
		b, _ := l.next()
		if b == '\r' {
			l.next()
			break
		} else if b == '\n' {
			break
		}
		buf = append(buf, b)
	}
	return NewLexemeToken(Comment, string(buf))
}

func (l *lexer) str() Token {
	buf := make([]byte, 0, 16)
	for {
		b, _ := l.next()
		if b == '"' {
			break
		}
		buf = append(buf, b)
	}
	return NewLexemeToken(String, string(buf))
}

func (l *lexer) step() Token {
	l.skipWhiteSpace()

	b, err := l.r.ReadByte()
	switch err {
	case io.EOF:
		return NewToken(Eof)
	case nil:
		break
	default:
		panic(err)
	}

	if isalpha(b) || b == '_' {
		// variable or keyword
		l.back()
		return l.identifier()
	}
	if isdigit(b) || (b == '-' && isdigit(l.peek())) {
		// number
		// TODO: parse 0x and b prefixed numbers
		l.back()
		return l.number()
	}

	switch b {
	case '#':
		return l.comment()

	case '"':
		return l.str()

	case '+':
		if l.match('=') {
			return NewToken(PlusEq)
		} else if l.match('+') {
			return NewToken(Increment)
		}
		return NewToken(Plus)

	case '-':
		if l.match('=') {
			return NewToken(MinusEq)
		} else if l.match('-') {
			NewToken(Decrement)
		}
		return NewToken(Minus)

	case '/':
		return l.pick('=', DivideEq, Divide)

	case '*':
		if l.match('*') {
			return l.pick('=', PowerEq, Power)
		}
		return l.pick('=', StarEq, Star)

	case '=':
		return l.pick('=', Equal, Assign)

	case '!':
		return l.pick('=', NotEqual, Bang)

	case '<':
		return l.pick('=', LessEq, Less)

	case '>':
		return l.pick('=', GreaterEq, Greater)

	case '&':
		return NewToken(BitAnd)

	case '|':
		return NewToken(BitOr)

	case '~':
		return NewToken(BitNegate)

	case '^':
		return NewToken(BitXOR)

	case '.':
		return NewToken(Dot)

	case ',':
		return NewToken(Comma)

	case '(':
		return NewToken(OpenParen)

	case ')':
		return NewToken(CloseParen)

	case '{':
		return NewToken(OpenBrace)

	case '}':
		return NewToken(CloseBrace)

	case '[':
		return NewToken(OpenBrack)

	case ']':
		return NewToken(CloseBrack)

	case ';': // always separates statements
		return NewToken(Sep)

	}

	return NewToken(Invalid)
}
