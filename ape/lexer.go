package ape

import (
	"os"
	"unicode"

	"github.com/pcen/ape/ape/token"
)

const (
	byteToTokenFactor = 2
)

var (
	stmtEndTokens = map[token.Kind]bool{
		token.Identifier: true,
		token.Integer:    true,
		token.Rational:   true,
		token.String:     true,
		token.True:       true,
		token.False:      true,
		token.Break:      true,
		token.Decrement:  true,
		token.Increment:  true,
		token.Return:     true,
		token.CloseParen: true,
		token.CloseBrace: true,
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
	buf     []byte
	idx     int
	done    bool
	pos     token.Position
	prevPos token.Position
	tokens  []token.Token
}

func (l *lexer) NewToken(kind token.Kind) token.Token {
	return token.New(kind, l.pos)
}

func (l *lexer) NewLexemeToken(kind token.Kind, lexeme string) token.Token {
	return token.NewLexeme(kind, lexeme, l.pos)
}

func (l *lexer) LexFile(file string) []token.Token {
	var err error
	l.buf, err = os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return l.lex()
}

func (l *lexer) LexString(source string) []token.Token {
	l.buf = []byte(source)
	return l.lex()
}

func (l *lexer) lex() []token.Token {
	l.tokens = make([]token.Token, 0, len(l.buf)/byteToTokenFactor)
	l.idx = 0
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
	l.done = l.idx == len(l.buf)
	if l.done {
		return 0, false
	}

	b := l.buf[l.idx]
	l.idx++

	if b == '\n' {
		l.prevPos = l.pos
		l.pos.Line++
		l.pos.Column = 0
	} else {
		l.pos.Column++
	}
	return b, true
}

func (l *lexer) back() {
	if l.done {
		return
	}
	l.idx--
	if l.buf[l.idx] == '\n' {
		l.pos = l.prevPos
	} else {
		l.pos.Column--
	}
}

func (l *lexer) peek() byte {
	if l.idx == len(l.buf) {
		return 0
	}
	return l.buf[l.idx]
}

func (l *lexer) match(b byte) bool {
	isMatch := l.peek() == b
	if isMatch {
		l.next() // consume matching token
	}
	return isMatch
}

func (l *lexer) pick(b byte, match token.Kind, noMatch token.Kind) token.Token {
	if l.match(b) {
		return l.NewToken(match)
	}
	return l.NewToken(noMatch)
}

func (l *lexer) shouldInsertSemi(current byte) bool {
	if current != 0 && current != '\n' {
		return false
	}
	if len(l.tokens) == 0 {
		return false // no previous token to check
	}
	return stmtEndTokens[l.tokens[len(l.tokens)-1].Kind]
}

func (l *lexer) skipWhiteSpace() bool {
	for {
		b, ok := l.next()
		if !ok {
			if l.shouldInsertSemi(b) {
				l.tokens = append(l.tokens, l.NewToken(token.Sep))
			}
			return true
		}
		if !iswspace(b) {
			l.back()
			return false
		}
		if l.shouldInsertSemi(b) {
			l.tokens = append(l.tokens, l.NewToken(token.Sep))
		}
	}
}

func (l *lexer) identifier() token.Token {
	start, end := l.idx, 0
	for {
		b, _ := l.next()
		if !isalpha(b) && !isdigit(b) && b != '_' {
			l.back()
			end = l.idx
			break
		}
	}
	lexeme := string(l.buf[start:end])
	if kind, keyword := token.GetKeyword(lexeme); keyword {
		return l.NewToken(kind)
	}
	return l.NewLexemeToken(token.Identifier, lexeme)
}

func (l *lexer) number() token.Token {
	start, end := l.idx, 0
	negative := false
	kind := token.Integer
	for {
		b, _ := l.next()
		if !(isdigit(b) || b == '.' && kind == token.Integer || b == '-' && !negative) {
			l.back()
			end = l.idx
			break
		}
		if b == '.' {
			kind = token.Rational
		}
		negative = b == '-'
	}
	return l.NewLexemeToken(kind, string(l.buf[start:end]))
}

func (l *lexer) comment() token.Token {
	start, end := l.idx, 0
	for {
		b, _ := l.next()
		if b == '\r' {
			l.next()
			end = l.idx - 2
			break
		} else if b == '\n' {
			end = l.idx - 1
			break
		}
	}
	return l.NewLexemeToken(token.Comment, string(l.buf[start:end]))
}

func (l *lexer) str() token.Token {
	start, end := l.idx, 0
	for {
		b, _ := l.next()
		if b == '"' {
			end = l.idx - 1
			break
		}
	}
	return l.NewLexemeToken(token.String, string(l.buf[start:end]))
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
			return l.NewToken(token.Decrement)
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
		if l.match('<') {
			return l.NewToken(token.ShiftLeft)
		}
		return l.pick('=', token.LessEq, token.Less)

	case '>':
		if l.match('>') {
			return l.NewToken(token.ShiftRight)
		}
		return l.pick('=', token.GreaterEq, token.Greater)

	case '&':
		return l.NewToken(token.Ampersand)

	case '|':
		return l.NewToken(token.Pipe)

	case '~':
		return l.NewToken(token.Tilde)

	case '^':
		return l.NewToken(token.Caret)

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
	panic("invalid byte: " + string(b))
	// return l.NewToken(token.Invalid)
}
