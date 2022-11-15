package ape

import (
	"fmt"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

// program     -> decl*

// decl        -> typedDecl | funcDecl

// typedDecl   -> (VAL | VAR) IDENT type  "=" expression
// funcDecl    -> "func" IDENT "(" parameters? ")" blockStmt

// parameters  -> (paramDecl ",")*
// paramDecl   -> IDENT type

// blockStmt   -> "{" stmtList "}"
// stmtList    -> (stmt;) *

// expression  -> equality ;
// equality    -> comparison ( ( "!=" | "==" ) comparison )*
// comparison  -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
// term        -> factor ( ( "-" | "+" | "|" | "^" ) factor )*
// factor      -> unary ( ( "/" | "*" | "&" ) unary )*
// unary       -> ( "!" | "-" | "~" ) unary | call
// call        -> primary ( "(" arguments? ")" )*
// primary     -> NUMBER | STRING | IDENT | "true" | "false" | group
// group       -> "(" expression ")"

// arguments   -> expression ( "," expression ) *

// propagates panic errors that are not ParseError
func sync(f func()) {
	err := recover()
	if err == nil {
		return
	}
	if _, ok := err.(ParseError); !ok {
		panic(err)
	}
	f()
}

var (
	declStart = map[token.Kind]bool{
		token.Val:  true,
		token.Var:  true,
		token.Func: true,
	}

	stmtStart = map[token.Kind]bool{
		token.Return:    true,
		token.OpenBrace: true,
	}
)

type ParseError struct {
	Pos    token.Position
	What   string
	Parsed ast.Node
}

func NewParseError(pos token.Position, parsed ast.Node, format string, a ...interface{}) ParseError {
	return ParseError{Pos: pos, Parsed: parsed, What: fmt.Sprintf(format, a...)}
}

func (p ParseError) String() string {
	if p.Parsed == nil {
		return fmt.Sprintf("%v: %v", p.Pos, p.What)
	}
	return fmt.Sprintf("%v: %v, ast: %v", p.Pos, p.What, ast.NodeString(p.Parsed))
}

type Parser interface {
	Program() []ast.Node
	Errors() ([]ParseError, bool)
}

type parser struct {
	tokens []token.Token
	pos    uint
	errors []ParseError
}

func NewParser(tokens []token.Token) Parser {
	return &parser{
		tokens: tokens,
		errors: make([]ParseError, 0),
	}
}

func (p *parser) Errors() ([]ParseError, bool) {
	return p.errors, len(p.errors) > 0
}

func (p *parser) errExpected(kind token.Kind, parsed ast.Node, context string) {
	pos, got := p.prev().Position, p.prev().String()
	err := NewParseError(pos, parsed, fmt.Sprintf("expected %v, got %v parsing %v", kind, got, context))
	p.errors = append(p.errors, err)
	panic(err)
}

func (p *parser) err(parsed ast.Node, format string, args ...interface{}) {
	err := NewParseError(p.prev().Position, parsed, fmt.Sprintf(format, args...))
	p.errors = append(p.errors, err)
	panic(err)
}

func (p *parser) skipTo(tokens map[token.Kind]bool) bool {
	for {
		kind := p.peek().Kind
		if kind == token.Eof {
			return false
		}
		if _, ok := tokens[kind]; ok {
			return true
		}
		p.next()
	}
}

func (p *parser) peek() token.Token {
	return p.tokens[p.pos]
}

func (p *parser) next() token.Token {
	p.pos++
	return p.tokens[p.pos-1]
}

func (p *parser) consume(tk token.Kind, parsed ast.Node, context string) {
	if p.tokens[p.pos].Kind == tk {
		p.pos++
		return
	}
	p.errExpected(tk, parsed, context)
}

func (p *parser) prev() token.Token {
	return p.tokens[p.pos-1]
}

func (p *parser) match(tk ...token.Kind) bool {
	for _, t := range tk {
		if p.tokens[p.pos].Kind == t {
			p.pos++
			return true
		}
	}
	return false
}

func (p *parser) Program() (program []ast.Node) {
	for !p.match(token.Eof) {
		program = append(program, p.Declaration())
	}
	return program
}

// Expressions

func (p *parser) Expression() ast.Expression {
	return p.Equality()
}

func (p *parser) leftAssociativeBinaryOp(rule func() ast.Expression, types ...token.Kind) ast.Expression {
	lhs := rule()
	for p.match(types...) {
		lhs = ast.NewBinaryOp(lhs, p.prev().Kind, rule())
	}
	return lhs
}

func (p *parser) Equality() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Comparison, token.Equal, token.NotEqual)
}

func (p *parser) Comparison() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Term, token.Greater, token.GreaterEq, token.Less, token.LessEq)
}

func (p *parser) Term() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Factor, token.Minus, token.Plus, token.BitOr, token.BitXOR)
}

func (p *parser) Factor() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Unary, token.Divide, token.Star, token.BitAnd)
}

func (p *parser) Unary() ast.Expression {
	switch p.peek().Kind {
	case token.Bang, token.Minus, token.BitNegate:
		return ast.NewUnaryOp(p.next().Kind, p.Unary())
	default:
		return p.CallExpr()
	}
}

func (p *parser) CallExpr() ast.Expression {
	primary := p.Primary()
	if p.match(token.OpenParen) {
		args := p.Arguments()
		p.consume(token.CloseParen, primary, "end of call expr")
		return &ast.CallExpr{
			Callee: primary,
			Args:   args,
		}
	}
	return primary
}

func (p *parser) Arguments() (args []ast.Expression) {
	if p.peek().Kind == token.CloseParen {
		// empty argument list
		return args
	}
	for {
		args = append(args, p.Expression())
		if !p.match(token.Comma) {
			return args
		}
	}
}

func (p *parser) Primary() ast.Expression {
	switch kind := p.peek().Kind; kind {
	case token.Number, token.String, token.True, token.False:
		return ast.NewLiteralExpr(p.next())
	case token.Identifier:
		return ast.NewIdentExpr(p.next())
	case token.OpenParen:
		return p.GroupExpr()
	default:
		p.err(nil, "invalid token for expression: %v", p.peek())
		return nil // err unwinds stack
	}
}

func (p *parser) GroupExpr() (expr ast.Expression) {
	p.consume(token.OpenParen, nil, "start of group expr")
	expr = p.Expression()
	p.consume(token.CloseParen, expr, "end of group expr")
	return &ast.GroupExpr{Expr: expr}
}

// Statements

func (p *parser) separator() {
	if p.match(token.Sep) || p.peek().Kind == token.CloseBrace {
		return
	}
	p.errExpected(token.Sep, nil, "expected statement separator")
}

func (p *parser) Statement() (s ast.Statement) {
	defer sync(func() {
		s = &ast.ErrStmt{}
		p.skipTo(stmtStart)
	})

	kind := p.peek().Kind
	switch kind {
	case token.Val, token.Var:
		s = &ast.TypedDeclStmt{Decl: p.TypedDecl()}
		p.separator()
	case token.Return:
		s = p.ReturnStmt()
		p.separator()
	case token.OpenBrace:
		s = p.BlockStmt()
	case token.Eof:
		panic("stmt at eof")
	default:
		s = p.ExprStmt()
		p.separator()
	}

	return s
}

func (p *parser) ReturnStmt() *ast.ReturnStmt {
	p.consume(token.Return, nil, "return stmt")
	return &ast.ReturnStmt{Expr: p.Expression()}
}

func (p *parser) BlockStmt() *ast.BlockStmt {
	p.consume(token.OpenBrace, nil, "block stmt start")
	content := p.StmtList()
	p.consume(token.CloseBrace, content, "block stmt end")
	return &ast.BlockStmt{Content: content}
}

func (p *parser) StmtList() (stmts []ast.Statement) {
	for p.peek().Kind != token.CloseBrace {
		stmts = append(stmts, p.Statement())
	}
	return stmts
}

func (p *parser) ExprStmt() ast.Statement {
	return &ast.ExprStmt{Expr: p.Expression()}
}

// Declarations

func (p *parser) Declaration() (d ast.Declaration) {
	defer sync(func() {
		d = &ast.ErrDecl{}
		p.skipTo(declStart)
	})
	switch kind := p.peek().Kind; kind {
	case token.Val, token.Var:
		d = p.TypedDecl()
	case token.Func:
		d = p.FuncDecl()
	default:
		panic(fmt.Sprintf("%v not a declaration start", kind))
	}
	return d
}

func (p *parser) Parameters() (decls []*ast.ParamDecl) {
	if p.peek().Kind == token.CloseParen {
		return decls // empty parameter list
	}
	for {
		decls = append(decls, p.ParamDecl())
		if !p.match(token.Comma) {
			return decls
		}
	}
}

func (p *parser) ParamDecl() *ast.ParamDecl {
	decl := &ast.ParamDecl{}
	if p.match(token.Identifier) {
		decl.Ident = p.prev()
	}
	if p.match(token.Identifier) {
		decl.Type = p.prev().Lexeme
	}
	return decl
}

func (p *parser) TypedDecl() *ast.TypedDecl {
	decl := &ast.TypedDecl{}
	if p.match(token.Val, token.Var) {
		decl.Kind = p.prev().Kind

		if p.match(token.Identifier) {
			decl.Ident = p.prev()
		}

		if p.match(token.Identifier) {
			decl.Type = p.prev().Lexeme
		} else {
			p.err(decl, "missing type for variable declaration %v", decl.Ident)
		}

		if p.match(token.Assign) {
			decl.Value = p.Expression()
		}
	} else {
		p.err(decl, "missing val or var for typed decl")
	}
	return decl
}

func (p *parser) FuncDecl() ast.Declaration {
	fd := &ast.FuncDecl{}
	p.consume(token.Func, nil, "function declaration start")
	if p.match(token.Identifier) {
		fd.Name = p.prev()
	}
	if p.match(token.OpenParen) {
		fd.Params = p.Parameters()
	}
	p.consume(token.CloseParen, fd, "end of function signature parameters")
	fd.Body = p.BlockStmt()
	return fd
}
