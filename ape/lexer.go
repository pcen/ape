package ape

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/pcen/ape/ape/token"
)

// Since source code files should be pretty small, LexFile should
// probably load the entire file into a byte slice instead of using
// bufio since this will be faster, lexemes can be obtained directly
// from the buffer (no need for small byte slices), and less
// backtracking with UnreadByte.

var (
	stmtEndTokens = map[token.Kind]bool{
		token.Identifier: true,
		token.Number:     true,
		token.String:     true,
		token.True:       true,
		token.False:      true,
		token.Break:      true,
		token.Decrement:  true,
		token.Increment:  true,
		token.Return:     true,
	}
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
	LexFile(string) []token.Token
	LexString(string) []token.Token
}

func NewLexer() Lexer {
	return &lexer{
		pos: token.Position{
			Line: 1,
			// column index of the next char to read
			// real column number of the most recent char
			Column: 0,
		},
	}
}

type lexer struct {
	r       *bufio.Reader
	prev    byte
	pos     token.Position
	prevPos token.Position
	done    bool
	tokens  []token.Token
}

func (l *lexer) NewToken(kind token.Kind) token.Token {
	return token.New(kind, l.pos)
}

func (l *lexer) NewLexemeToken(kind token.Kind, lexeme string) token.Token {
	return token.NewLexeme(kind, lexeme, l.pos)
}

func (l *lexer) LexFile(file string) []token.Token {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	l.r = bufio.NewReader(f)
	return l.lex()

}

func (l *lexer) LexString(source string) []token.Token {
	l.r = bufio.NewReader(strings.NewReader(source))
	return l.lex()
}

func (l *lexer) lex() []token.Token {
	l.tokens = make([]token.Token, 0, 64)
	for {
		tok := l.step()
		l.tokens = append(l.tokens, tok)
		if tok.Kind == token.Eof {
			break
		}
	}
	return l.tokens
}

func (l *lexer) next() (byte, bool) {
	if l.done {
		return 0, false
	}

	b, err := l.r.ReadByte()
	if err == io.EOF {
		l.done = true
		return 0, false
	} else if err != nil {
		panic(err)
	}

	if b == '\n' {
		l.prevPos = l.pos
		l.pos.Line++
		l.pos.Column = 0
	} else {
		l.pos.Column++
	}
	l.prev = b
	return b, true
}

func (l *lexer) back() {
	if l.done {
		return
	}

	// TODO: this fails if back is ever called successively
	//       make sure this cannot occur
	if err := l.r.UnreadByte(); err != nil {
		panic(err)
	}
	if l.prev == '\n' {
		l.pos = l.prevPos
	} else {
		l.pos.Column--
	}
}

func (l *lexer) peek() byte {
	if b, err := l.r.Peek(1); err == nil {
		return b[0]
	}
	return 0
}

func (l *lexer) match(b byte) bool {
	isMatch := l.peek() == b
	if isMatch {
		l.next() // consume matching token
	}
	return isMatch
}

func (l *lexer) pick(b byte, onMatch token.Kind, noMatch token.Kind) token.Token {
	if l.match(b) {
		return l.NewToken(onMatch)
	}
	return l.NewToken(noMatch)
}

func (l *lexer) skipWhiteSpace() bool {
	for {
		b, ok := l.next()
		if !ok {
			return true
		}
		if !iswspace(b) {
			l.back()
			return false
		}
		// insert statement separator automatically
		if b == '\n' && len(l.tokens) > 0 && stmtEndTokens[l.tokens[len(l.tokens)-1].Kind] {
			l.tokens = append(l.tokens, l.NewToken(token.Sep))
		}
	}
}

func (l *lexer) identifier() token.Token {
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
	if kind, keyword := token.GetKeyword(lexeme); keyword {
		return l.NewToken(kind)
	}
	return l.NewLexemeToken(token.Identifier, lexeme)
}

func (l *lexer) number() token.Token {
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
	return l.NewLexemeToken(token.Number, string(buf))
}

func (l *lexer) comment() token.Token {
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
	return l.NewLexemeToken(token.Comment, string(buf))
}

func (l *lexer) str() token.Token {
	buf := make([]byte, 0, 16)
	for {
		b, _ := l.next()
		if b == '"' {
			break
		}
		buf = append(buf, b)
	}
	return l.NewLexemeToken(token.String, string(buf))
}

func (l *lexer) step() token.Token {
	atEof := l.skipWhiteSpace()
	if atEof {
		return l.NewToken(token.Eof)
	}
	b, _ := l.next()
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
			return l.NewToken(token.PlusEq)
		} else if l.match('+') {
			return l.NewToken(token.Increment)
		}
		return l.NewToken(token.Plus)

	case '-':
		if l.match('=') {
			return l.NewToken(token.MinusEq)
		} else if l.match('-') {
			l.NewToken(token.Decrement)
		}
		return l.NewToken(token.Minus)

	case '/':
		return l.pick('=', token.DivideEq, token.Divide)

	case '*':
		if l.match('*') {
			return l.pick('=', token.PowerEq, token.Power)
		}
		return l.pick('=', token.StarEq, token.Star)

	case '=':
		return l.pick('=', token.Equal, token.Assign)

	case '!':
		return l.pick('=', token.NotEqual, token.Bang)

	case '<':
		return l.pick('=', token.LessEq, token.Less)

	case '>':
		return l.pick('=', token.GreaterEq, token.Greater)

	case '&':
		return l.NewToken(token.BitAnd)

	case '|':
		return l.NewToken(token.BitOr)

	case '~':
		return l.NewToken(token.BitNegate)

	case '^':
		return l.NewToken(token.BitXOR)

	case '.':
		return l.NewToken(token.Dot)

	case ',':
		return l.NewToken(token.Comma)

	case '(':
		return l.NewToken(token.OpenParen)

	case ')':
		return l.NewToken(token.CloseParen)

	case '{':
		return l.NewToken(token.OpenBrace)

	case '}':
		return l.NewToken(token.CloseBrace)

	case '[':
		return l.NewToken(token.OpenBrack)

	case ']':
		return l.NewToken(token.CloseBrack)

	case ';':
		return l.NewToken(token.Sep)

	}

	return l.NewToken(token.Invalid)
}
