package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/*
	This file contains code to produce derivations of a given grammar. The
	syntax used to specify grammars is based on PEG. How star expansion is
	handled in derivations is determined by a config file specifying a rule
	name followed by the distribution from which n is picked, where n is the
	number of times to repeat rule*
*/

// grammar descriptor grammar:
// rule -> NAME "->" or End
// or -> list ("|" list)*
// list -> star*
// star -> ref "*"? | primitave "*"? | literal "*"? | group "*"?
// group -> "(" or ")"

// config file is a list of lines with the format:
// rule: distribution
// where distribution is one of the following:
// u[x, y] uniform distribution from x to y
// b[pct]  pct% chance n is 1, otherwise n is 0

type parser struct {
	lines [][]byte
	toks  []Token
	pos   int
}

func NewParser(lines [][]byte) *parser {
	return &parser{
		lines: lines,
	}
}

type node interface {
	str() string
}

type rule struct {
	name string
	or   *or
}

func (r *rule) str() string {
	return fmt.Sprintf("%v: %v", r.name, r.or.str())
}

// sequence of pipe-separated grammar rules
// of which one is used in derivation
type or struct {
	ors [][]*unary
}

func (o *or) pick() []*unary {
	return o.ors[rand.Intn(len(o.ors))]
}

func (o *or) str() string {
	var sb strings.Builder
	for i, l := range o.ors {
		for j, u := range l {
			sb.WriteString(u.str())
			if j != len(l)-1 {
				sb.WriteString(" ")
			}
		}
		if i != len(o.ors)-1 {
			sb.WriteString(" | ")
		}
	}
	return sb.String()
}

// a grammar rule followed by * or ?
type unary struct {
	n  node
	op Kind
}

func (u *unary) expand(count int) []node {
	switch u.op {
	case None:
		return []node{u.n}
	case Star:
		repeated := make([]node, count)
		for i := range repeated {
			repeated[i] = u.n
		}
		return repeated
	case Question:
		// TODO: also configure the chance of evaluating an optional rule
		// but since these shouldn't blow up the way * can 50% should be ok
		if CheckPct(50) {
			return []node{u.n}
		}
		return []node{}
	default:
		panic("invalid operator for unary:" + u.op.String())
	}
}

func (u *unary) str() string {
	switch u.op {
	case Star:
		return fmt.Sprintf("%v*", u.n.str())
	case Question:
		return fmt.Sprintf("%v?", u.n.str())
	case None:
		return u.n.str()
	}
	panic("invalid operator for unary:" + u.op.String())
}

// groups a set of grammar rules in parentheses
type group struct {
	or *or
}

func (g *group) str() string {
	return fmt.Sprintf("(%v)", g.or.str())
}

// refers to a rule in the grammar by name
// rules always start with a lower case character
type ref struct {
	name string
}

func (r *ref) str() string {
	return r.name
}

// refers to a predefined rule such as IDENT or NUMBER
// primitaves are always in upper case
type primitave struct {
	name string
}

func (p *primitave) str() string {
	return p.name
}

// a literal terminal (ie. ";")
// terminals are always enclosed in quotation marks
type literal struct {
	lexeme string
}

func (l *literal) str() string {
	return fmt.Sprintf("\"%v\"", l.lexeme)
}

func (p *parser) Grammar() (rules []*rule) {
	var l Lexer
	for _, line := range p.lines {
		p.toks = l.Lex(line)
		// fmt.Printf("parsing rule for line %v: %s\n", i+1, line)
		// fmt.Printf("\ttokens: %v\n", p.toks)
		p.pos = 0
		rule := p.Rule()
		rules = append(rules, rule)
	}
	return rules
}

func (p *parser) Rule() *rule {
	r := &rule{}
	p.consume(Ref)
	r.name = p.prev().Lex
	p.consume(Arrow)
	r.or = p.Or()
	p.consume(End)
	return r
}

func (p *parser) Or() *or {
	o := &or{ors: make([][]*unary, 0)}
	o.ors = append(o.ors, p.List())
	for p.match(Pipe) {
		o.ors = append(o.ors, p.List())
	}
	return o
}

func (p *parser) List() (list []*unary) {
	list = append(list, p.Unary())
	for p.peekIs(Ref, Primitave, Literal, Lparen) {
		list = append(list, p.Unary())
	}
	return list
}

func (p *parser) Unary() *unary {
	u := &unary{}
	switch p.peek().Kind {
	case Ref:
		u.n = p.Ref()
	case Primitave:
		u.n = p.Primitave()
	case Literal:
		u.n = p.Literal()
	case Lparen:
		u.n = p.Group()
	}
	switch {
	case p.match(Star):
		u.op = Star
	case p.match(Question):
		u.op = Question
	default:
		u.op = None
	}
	return u
}

func (p *parser) Group() *group {
	g := &group{}
	p.consume(Lparen)
	g.or = p.Or()
	p.consume(Rparen)
	return g
}

func (p *parser) Ref() *ref {
	p.consume(Ref)
	return &ref{name: p.prev().Lex}
}

func (p *parser) Primitave() *primitave {
	p.consume(Primitave)
	return &primitave{name: p.prev().Lex}
}

func (p *parser) Literal() *literal {
	p.consume(Literal)
	return &literal{lexeme: p.prev().Lex}
}

func (p *parser) next() Token {
	p.pos++
	return p.toks[p.pos-1]
}

func (p *parser) peek() Token {
	return p.toks[p.pos]
}

func (p *parser) prev() Token {
	return p.toks[p.pos-1]
}

func (p *parser) peekIs(kinds ...Kind) bool {
	for _, k := range kinds {
		if p.peek().Kind == k {
			return true
		}
	}
	return false
}

func (p *parser) match(k Kind) bool {
	if p.peek().Kind == k {
		p.next()
		return true
	}
	return false
}

func (p *parser) consume(k Kind) {
	if !p.match(k) {
		err := fmt.Sprintf("could not consume: %v, got %v: %v", k, p.peek().Kind, p.peek().Lex)
		panic(err)
	}
}

type Kind int

const (
	None Kind = iota + 1
	Ref
	Primitave
	Literal
	Arrow    // ->
	Lparen   // (
	Rparen   // )
	Pipe     // |
	Star     // *
	Question // ?
	Ws
	End
)

func (k Kind) String() string {
	return []string{
		Ref:       "<REF>",
		Primitave: "<PRIMITAVE>",
		Literal:   "<LITERAL>",
		Arrow:     "<ARROW>",
		Lparen:    "<LPAREN>",
		Rparen:    "<RPAREN>",
		Pipe:      "<PIPE>",
		Star:      "<STAR>",
		Question:  "<QUESTION>",
		Ws:        "<WS>",
		End:       "<END>",
	}[k]
}

var (
	KindMap = map[byte]Kind{
		'(': Lparen,
		')': Rparen,
		'|': Pipe,
		'*': Star,
		'?': Question,
	}
)

type Token struct {
	Kind
	Lex string
}

type Lexer struct {
	line []byte
	pos  int
}

func (l *Lexer) next() byte {
	if l.pos == len(l.line) {
		return 0
	}
	b := l.line[l.pos]
	l.pos++
	return b
}

func (l *Lexer) back() {
	if l.pos != len(l.line) {
		l.pos--
	}
}

func (l *Lexer) peek() byte {
	if l.pos != len(l.line) {
		return l.line[l.pos]
	}
	return 0
}

func (l *Lexer) Lex(line []byte) (tokens []Token) {
	l.line = line
	l.pos = 0
	for {
		tok := l.step()
		if tok.Kind == Ws {
			continue
		}
		tokens = append(tokens, tok)
		if tok.Kind == End {
			return tokens
		}
	}
}

func alpha(b byte) bool {
	return unicode.IsLetter(rune(b))
}

func isUpper(b byte) bool {
	return bytes.ToUpper([]byte{b})[0] == b
}

func ws(b byte) bool {
	return unicode.IsSpace(rune(b))
}

func (l *Lexer) word() Token {
	first := l.pos - 1
	kind := Ref
	if isUpper(l.line[first]) {
		kind = Primitave
	}
	for {
		b := l.peek()
		if !alpha(b) {
			return Token{
				Kind: kind,
				Lex:  string(l.line[first:l.pos]),
			}
		}
		l.next()
	}
}

func (l *Lexer) lit() Token {
	first := l.pos // skip "
	for {
		b := l.peek()
		if b == '"' {
			end := l.pos
			l.next()
			return Token{
				Kind: Literal,
				Lex:  string(l.line[first:end]),
			}
		}
		l.next()
	}
}

func (l *Lexer) ws() Token {
	for {
		if !ws(l.next()) {
			l.back()
			return Token{Kind: Ws}
		}
	}
}

func (l *Lexer) step() Token {
	b := l.next()
	switch {
	case b == 0:
		return Token{Kind: End}
	case ws(b):
		return l.ws()
	case alpha(b):
		return l.word()
	case b == '"':
		return l.lit()
	}
	if b == '-' && l.pos < len(l.line) && l.line[l.pos] == '>' {
		l.next()
		return Token{Kind: Arrow, Lex: "->"}
	}
	if kind, ok := KindMap[b]; ok {
		return Token{Kind: kind, Lex: string(b)}
	}
	panic("couldn't lex char " + string(b))
}

func NewScanner(file string) *bufio.Scanner {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	return bufio.NewScanner(f)
}

func ReadLines(file string) (lines [][]byte) {
	scanner := NewScanner(file)
	for scanner.Scan() {
		b := bytes.TrimSpace(scanner.Bytes())
		if len(b) > 0 {
			lines = append(lines, b)
		}
	}
	return lines
}

func ReadGrammarFile(file string) (lines [][]byte) {
	all := ReadLines(file)
	for _, line := range all {
		if line[0] == '#' {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

type Grammar struct {
	Start string
	Rules map[string]*or
	buf   []string
	cur   string
	Cfg   *Config
}

func NewGrammar(start string, rules []*rule, cfg *Config) *Grammar {
	m := make(map[string]*or)
	for _, rule := range rules {
		m[rule.name] = rule.or
	}
	return &Grammar{
		Start: start,
		Rules: m,
		Cfg:   cfg,
	}
}

/*
eval is the recursive function that drives generating derivations

  - need to be careful with how stars are expanded or stack will overflow
    since derivations are unbounded in size

  - a config text file controls the behaviour of * expansion for given parent
    rules

  - * expansion handled well, but or selection leads to bad input ie. unary
    repeatedly selects the option ( "!" | "-" | "~" ) unary leading to
    unnecessary long chains of unary operators in test input
*/
func (g *Grammar) eval(n node) {
	// fmt.Printf("evaluating %v:%v\n", reflect.TypeOf(n), n.str())
	// fmt.Printf("current buf: %s\n", g.buf)
	switch n := n.(type) {

	case *or:
		for _, u := range n.pick() {
			g.eval(u)
		}

	case *unary:
		count := 1
		if n.op == Star {
			var ok bool
			count, ok = g.Cfg.Get(g.cur)
			if !ok {
				panic("no config for expanding * within " + g.cur)
			}
		}
		expanded := n.expand(count)
		for _, n := range expanded {
			g.eval(n)
		}
	case *group:
		g.eval(n.or)
	case *ref:
		rule, ok := g.Rules[n.name]
		if !ok {
			panic("no rule in grammar for " + n.name)
		}
		// need to restore current rule name in case a list of rules
		// is being traversed
		prev := g.cur
		g.cur = n.name
		g.eval(rule)
		g.cur = prev

	case *primitave:
		switch n.name {
		case "IDENT":
			if g.cur == "type" {
				g.buf = append(g.buf, "int")
			} else {
				g.buf = append(g.buf, "foo")
			}
		case "NUMBER":
			g.buf = append(g.buf, "123")
		case "STRING":
			g.buf = append(g.buf, `"bar"`)
		default:
			panic("unknown primitave: " + n.name)
		}

	case *literal:
		lex := n.lexeme
		if n.lexeme == ";" {
			lex += "\n"
		}
		g.buf = append(g.buf, lex)
	}
}

func (g *Grammar) Derive() string {
	g.buf = make([]string, 0, 64)
	start := g.Rules[g.Start]
	g.cur = g.Start
	g.eval(start)
	g.buf = append(g.buf, "\n")
	return strings.Join(g.buf, " ")
}

type Config struct {
	Dist map[string]func() int
}

func (c *Config) Get(rule string) (int, bool) {
	if f, ok := c.Dist[rule]; ok {
		return f(), true
	}
	return 0, false
}

func CheckPct(pct int) bool {
	return 1+rand.Intn(100) <= pct
}

func ReadConfig(file string) *Config {
	cfg := &Config{
		Dist: make(map[string]func() int),
	}
	lines := ReadGrammarFile(file)
	for _, line := range lines {
		l := bytes.Split(line, []byte{':'})
		rule := string(bytes.TrimSpace(l[0]))
		distStr := string(bytes.TrimSpace(l[1]))
		switch distStr[0] {
		case 'u':
			var low, high int
			fmt.Sscanf(distStr, "u[%d, %d]", &low, &high)
			cfg.Dist[rule] = func() int {
				return low + rand.Intn(low+high+1)
			}
		case 'b':
			var pct int
			fmt.Sscanf(distStr, "b[%d]", &pct)
			cfg.Dist[rule] = func() int {
				if CheckPct(pct) {
					return 1
				}
				return 0
			}
		}
	}
	return cfg
}

func main() {
	output := len(os.Args) > 2
	var (
		directory string
		n         int64
	)
	if output {
		directory = os.Args[1]
		n, _ = strconv.ParseInt(os.Args[2], 10, 32)
	}
	start := time.Now()

	rand.Seed(start.UnixNano())

	lines := ReadGrammarFile("./grammar.txt")
	rules := NewParser(lines).Grammar()

	cfg := ReadConfig("./dist.txt")

	g := NewGrammar("program", rules, cfg)

	if output {
		for i := int64(0); i < n; i++ {
			derivation := g.Derive()
			os.WriteFile(fmt.Sprintf("%v/fuzz_%v_%v.ape", directory, start.Unix(), i), []byte(derivation), 0664)
		}
		fmt.Println("generated", n, "files in", directory)
	} else {
		derivation := g.Derive()
		fmt.Printf("derivation:\n%v\n", derivation)
	}
}
