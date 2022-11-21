package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

/*
	This file contains code to produce derivations of a given grammar. The
	syntax used to specify grammars is based on PEG.
*/

// grammar descriptor grammar:
// rule -> NAME "->" or End
// or -> list ("|" list)*
// list -> star*
// star -> ref ?* | primitave ?* | literal ?* | group ?*
// group -> "(" or ")"

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
	ors []*list
}

func (o *or) pick() *list {
	return o.ors[rand.Intn(len(o.ors))]
}

func (o *or) str() string {
	var sb strings.Builder
	for i, l := range o.ors {
		sb.WriteString(l.str())
		if i != len(o.ors)-1 {
			sb.WriteString(" | ")
		}
	}
	return sb.String()
}

// a list of grammar rules that are applied sequentially
type list struct {
	stars []*star
}

func (l *list) str() string {
	var sb strings.Builder
	for i, s := range l.stars {
		sb.WriteString(s.str())
		if i != len(l.stars)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

// a grammar rule that can be applied 0 or more times
type star struct {
	n        node
	repeated bool
}

func (s *star) expand(count int) []node {
	if !s.repeated {
		return []node{s.n}
	}
	repeated := make([]node, count)
	for i := range repeated {
		repeated[i] = s.n
	}
	return repeated
}

func (s *star) str() string {
	if s.repeated {
		return fmt.Sprintf("%v*", s.n.str())
	} else {
		return s.n.str()
	}
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
	for i, line := range p.lines {
		fmt.Printf("parsing rule for line %v: %s\n", i+1, line)
		p.toks = l.Lex(line)
		fmt.Printf("\ttokens: %v\n", p.toks)
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
	o := &or{
		ors: make([]*list, 0),
	}
	o.ors = append(o.ors, p.List())
	for p.match(Pipe) {
		o.ors = append(o.ors, p.List())
	}
	return o
}

func (p *parser) List() *list {
	l := &list{stars: make([]*star, 0)}
	l.stars = append(l.stars, p.Star())
	for p.peekIs(Ref, Primitave, Literal, Lparen) {
		l.stars = append(l.stars, p.Star())
	}
	return l
}

func (p *parser) Star() *star {
	s := &star{}

	switch p.peek().Kind {
	case Ref:
		s.n = p.Ref()
	case Primitave:
		s.n = p.Primitave()
	case Literal:
		s.n = p.Literal()
	case Lparen:
		s.n = p.Group()
	}
	if p.match(Star) {
		s.repeated = true
	}
	return s
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
	Ref Kind = iota + 1
	Primitave
	Literal
	Arrow
	Lparen
	Rparen
	Pipe
	Star
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
	if b == 0 {
		return Token{Kind: End}
	}
	if ws(b) {
		return l.ws()
	}
	if alpha(b) {
		return l.word()
	}
	if b == '"' {
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
  - need to be careful with how stars are expanded or stack will
    overflow
  - might fake tail recursion because go doesn't implement it but
    could speed this up
  - should optimize eval of nodes such as having or return a node that
    is *list if there are multiple rules, but otherwise return the single
    *star rule directly
  - need to figure out how to properly weight * expansion to generate
    good parser tests
*/
func (g *Grammar) eval(n node) {
	// fmt.Printf("evaluating %v:%v\n", reflect.TypeOf(n), n.str())
	// fmt.Printf("current buf: %s\n", g.buf)
	switch n := n.(type) {
	case *or:
		g.eval(n.pick())
	case *list:
		for _, s := range n.stars {
			g.eval(s)
		}
	case *star:
		count, ok := g.Cfg.Get(g.cur)
		if !ok {
			panic("no config for expanding * within " + g.cur)
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
		g.cur = n.name
		g.eval(rule)
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
		g.buf = append(g.buf, n.lexeme)
	}
}

func (g *Grammar) Derive() string {
	g.buf = make([]string, 0, 64)
	start := g.Rules[g.Start]
	g.cur = g.Start
	fmt.Printf("start rule for derivation: %v\n", start.str())
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
				if rand.Intn(100) <= pct {
					return 1
				}
				return 0
			}
		}
	}
	return cfg
}

// func pct(p int) func() int {
// 	return func() int {
// 		if rand.Intn(100) <= p {
// 			return 1
// 		}
// 		return 0
// 	}
// }

func main() {
	var out string
	if len(os.Args) > 1 {
		out = os.Args[1]
	}

	rand.Seed(time.Now().UnixNano())

	lines := ReadGrammarFile("./grammar.txt")
	rules := NewParser(lines).Grammar()

	cfg := ReadConfig("./dist.txt")
	// for _, r := range rules {
	// 	fmt.Println(r.str())
	// }

	g := NewGrammar("program", rules, cfg)

	derivation := g.Derive()

	if out != "" {
		os.WriteFile(out, []byte(derivation), 0664)
	}

	fmt.Println("derivation:")
	fmt.Println(derivation)
}
