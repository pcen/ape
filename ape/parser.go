package ape

import (
	"fmt"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

// program      -> module decl*

// module       -> "module" IDENT
// decl         -> typedDecl | funcDecl

// typedDecl    -> (VAL | VAR) IDENT type  "=" expression
// funcDecl     -> "func" IDENT "(" parameters? ")" blockStmt
// classDecl    -> "class" IDENT "{" classBody "}"
// classBody    ->

// parameters   -> (paramDecl ",")*
// paramDecl    -> IDENT type

// blockStmt    -> "{" stmtList "}"
// stmtList     -> (stmt;) *

// stmt         -> simpleStmt | compoundStmt
// simpleStmt   -> incStmt | assignment
// incStmt      -> expression [++ | --]
// assignment   -> expression assign_op expression
// assign_op    -> = | += | *= | -= | /= | **=

// compoundStmt -> ifStmt | forStmt

// expression   -> equality ;
// equality     -> comparison ( ( "!=" | "==" ) comparison )*
// comparison   -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
// term         -> factor ( ( "-" | "+" | "|" | "^" ) factor )*
// factor       -> unary ( ( "/" | "*" | "&" ) unary )*
// unary        -> ( "!" | "-" | "~" ) unary | primary
// primary         -> atom ( ["(" arguments? ")" | "." IDENT] )*
// atom         -> NUMBER | STRING | IDENT | "true" | "false" | group
// selector     -> "." IDENT
// group        -> "(" expression ")"

// arguments    -> expression ( "," expression ) *

// propagates panic errors that are not ParseError
func sync(f func()) {
	if err := recover(); err == nil {
		return
	} else if _, ok := err.(ParseError); !ok {
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
		token.Val:       true,
		token.Var:       true,
		token.Return:    true,
		token.OpenBrace: true,
		token.If:        true,
		token.For:       true,
		token.While:     true,
	}
)

type ParseError struct {
	Pos token.Position
	Msg string
}

func NewParseError(pos token.Position, format string, a ...interface{}) ParseError {
	return ParseError{Pos: pos, Msg: fmt.Sprintf(format, a...)}
}

func (p ParseError) String() string {
	return fmt.Sprint(p.Pos, ": ", p.Msg)
}

type Parser interface {
	File() *ast.File
	Program() []ast.Declaration
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

func (p *parser) errExpected(kind token.Kind, context string) {
	pos, got := p.prev().Position, p.prev().String()
	err := NewParseError(pos, fmt.Sprintf("expected %v, got %v parsing %v", kind, got, context))
	fmt.Println("parser error:", err)
	p.errors = append(p.errors, err)
	panic(err)
}

func (p *parser) err(format string, args ...interface{}) {
	err := NewParseError(p.prev().Position, fmt.Sprintf(format, args...))
	fmt.Println("parser error:", err)
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

func (p *parser) peekIs(kinds ...token.Kind) bool {
	for _, kind := range kinds {
		if p.peek().Kind == kind {
			return true
		}
	}
	return false
}

func (p *parser) next() token.Token {
	p.pos++
	return p.tokens[p.pos-1]
}

func (p *parser) consume(tk token.Kind, context string) {
	if !p.match(tk) {
		p.errExpected(tk, context)
	}
}

func (p *parser) prev() token.Token {
	return p.tokens[p.pos-1]
}

func (p *parser) match(tk ...token.Kind) bool {
	for _, t := range tk {
		if p.peek().Kind == t {
			p.pos++
			return true
		}
	}
	return false
}

func (p *parser) File() (file *ast.File) {
	f := ast.NewFile("")
	p.consume(token.Module, "module declaration")
	p.consume(token.Identifier, "module name")
	f.Module = p.prev().Lexeme
	p.separator("end of module declaration")
	f.Ast = p.Program()
	return f
}

func (p *parser) Program() (decls []ast.Declaration) {
	for !p.match(token.Eof) {
		decls = append(decls, p.Declaration())
	}
	return decls
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
		return p.Primary()
	}
}

// unary and binary operators work on primary expressions
func (p *parser) Primary() ast.Expression {
	expr := p.Atom()
	for p.peekIs(token.OpenParen, token.Dot, token.OpenBrack) {
		// foo(bar)
		if p.match(token.OpenParen) {
			args := p.Arguments()
			p.consume(token.CloseParen, "end of call expr")
			expr = &ast.CallExpr{Callee: expr, Args: args}
		}
		// foo.bar
		if p.match(token.Dot) {
			p.consume(token.Identifier, "field in dot expr")
			expr = &ast.DotExpr{Expr: expr, Field: ast.NewIdentExpr(p.prev())}
		}
		// foo[bar]
		if p.match(token.OpenBrack) {
			index := p.Expression()
			p.consume(token.CloseBrack, "end of index expr")
			expr = &ast.IndexExpr{Expr: expr, Index: index}
		}
	}
	return expr
}

func (p *parser) Arguments() (args []ast.Expression) {
	for !p.peekIs(token.CloseParen) {
		args = append(args, p.Expression())
		if !p.match(token.Comma) {
			break
		}
	}
	return args
}

func (p *parser) Atom() ast.Expression {
	switch p.peek().Kind {
	case token.Number, token.String, token.True, token.False:
		return ast.NewLiteralExpr(p.next())
	case token.Identifier:
		return ast.NewIdentExpr(p.next())
	case token.OpenParen:
		return p.GroupExpr()
	default:
		p.err("invalid token for expression: %v", p.peek())
		return nil // err unwinds stack
	}
}

func (p *parser) GroupExpr() (expr ast.Expression) {
	p.consume(token.OpenParen, "start of group expr")
	expr = p.Expression()
	p.consume(token.CloseParen, "end of group expr")
	return &ast.GroupExpr{Expr: expr}
}

// Statements

func (p *parser) separator(context string) {
	if p.match(token.Sep) || p.peekIs(token.CloseBrace) {
		return
	}
	p.errExpected(token.Sep, fmt.Sprint(context, ": expected statement separator"))
}

func (p *parser) Statement() (s ast.Statement) {
	defer sync(func() {
		s = &ast.ErrStmt{}
		p.skipTo(stmtStart)
	})

	switch p.peek().Kind {

	case token.Identifier:
		s = p.SimpleStmt()
		p.separator("simple stmt")

	case token.Val, token.Var:
		s = &ast.TypedDeclStmt{Decl: p.TypedDecl()}
		p.separator("typed decl stmt")

	case token.Return:
		s = p.ReturnStmt()
		p.separator("return stmt")

	case token.OpenBrace:
		s = p.BlockStmt()

	case token.If:
		s = p.IfStmt()

	case token.For, token.While:
		s = p.ForStmt()

	case token.Eof:
		panic("stmt at eof")

	default:
		s = p.ExprStmt()
	}

	return s
}

func (p *parser) SimpleStmt() ast.Statement {
	lhs := p.Expression()
	if p.match(token.Increment, token.Decrement) {
		return &ast.IncStmt{
			Expr: lhs,
			Op:   p.prev(),
		}
	}
	if p.match(token.Assign, token.PlusEq, token.MinusEq, token.StarEq, token.DivideEq, token.PowerEq) {
		return &ast.AssignmentStmt{
			Lhs: lhs,
			Op:  p.prev(),
			Rhs: p.Expression(),
		}
	}
	p.err("invalid token for simple stmt: %v", p.peek())
	return nil // unreachable
}

func (p *parser) ReturnStmt() *ast.ReturnStmt {
	p.consume(token.Return, "return stmt")
	return &ast.ReturnStmt{Expr: p.Expression()}
}

func (p *parser) BlockStmt() *ast.BlockStmt {
	p.consume(token.OpenBrace, "block stmt start")
	content := p.StmtList()
	p.consume(token.CloseBrace, "block stmt end")
	return &ast.BlockStmt{Content: content}
}

func (p *parser) IfStmt() *ast.IfStmt {
	stmt := &ast.IfStmt{
		Elifs: make([]*ast.CondBlockStmt, 0),
	}
	p.consume(token.If, "if stmt start")
	stmt.If = p.CondBlockStmt()
	for p.match(token.Elif) {
		stmt.Elifs = append(stmt.Elifs, p.CondBlockStmt())
	}
	if p.match(token.Else) {
		stmt.Else = p.BlockStmt()
	}
	return stmt
}

func (p *parser) CondBlockStmt() *ast.CondBlockStmt {
	if p.peek().Kind == token.OpenBrace {
		p.err("missing predicate expression for conditional block")
	}
	return &ast.CondBlockStmt{
		Cond: p.Equality(),
		Body: p.BlockStmt(),
	}
}

func (p *parser) ForStmt() *ast.ForStmt {
	s := &ast.ForStmt{}
	switch p.next().Kind {

	case token.For:
		s.Init = p.TypedDecl()
		p.separator("after for loop init")
		s.Cond = p.Expression()
		p.separator("after for loop condition")
		s.Incr = p.SimpleStmt()

	case token.While:
		s.Cond = p.Expression()

	default:
		p.err("%v cannot start a loop statement", p.prev())
	}
	s.Body = p.BlockStmt()
	return s
}

// TODO: when the last statement in a statement list is invalid,
//
//	Statement() skips the closing curly brace in attempt to
//	find the next statement in the list, so the parser will
//	consume the statements in the outer block. Figure out if
//	handling this edge case is worth the complexity.
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

	case token.Class:
		d = p.ClassDecl()

	default:
		panic(fmt.Sprintf("%v not a declaration start", kind))
	}
	return d
}

func (p *parser) ParamList() (decls []*ast.ParamDecl) {
	if p.peekIs(token.CloseParen) {
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
	if p.match(token.Val, token.Var) {
		decl := &ast.TypedDecl{}
		decl.Kind = p.prev().Kind

		p.consume(token.Identifier, "typed decl identifier")
		decl.Ident = p.prev()

		p.consume(token.Identifier, "typed decl type")
		decl.Type = p.prev().Lexeme

		if p.match(token.Assign) {
			decl.Value = p.Expression()
		}
		return decl
	}
	p.err("missing val or var for typed decl")
	return nil
}

func (p *parser) FuncDecl() *ast.FuncDecl {
	fd := &ast.FuncDecl{}
	p.consume(token.Func, "function declaration start")

	p.consume(token.Identifier, "function name")
	fd.Name = p.prev()

	p.consume(token.OpenParen, "function signature parameters")
	fd.Params = p.ParamList()
	p.consume(token.CloseParen, "end of function signature parameters")

	fd.Body = p.BlockStmt()
	return fd
}

func (p *parser) ClassDecl() *ast.ClassDecl {
	cd := &ast.ClassDecl{}
	p.consume(token.Class, "class declaration start")
	p.consume(token.Identifier, "class name")
	cd.Name = p.prev()
	p.BlockStmt()

	return cd
}
