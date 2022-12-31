package ape

import (
	"fmt"
	"strings"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

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
	Demo() []ast.Statement
	File() *ast.File
	Program() []ast.Declaration
	Errors() ([]ParseError, bool)
}

type parser struct {
	tokens []token.Token
	pos    uint
	errors []ParseError
	decls  []ast.Declaration
}

func NewParser(tokens []token.Token) Parser {
	return &parser{
		tokens: tokens,
		errors: make([]ParseError, 0),
	}
}

// errors, hasErrors
func (p *parser) Errors() ([]ParseError, bool) {
	return p.errors, len(p.errors) > 0
}

func (p *parser) errExpected(kind token.Kind, context string) {
	pos, got := p.peek().Position, p.peek().String()
	err := NewParseError(pos, fmt.Sprintf("expected %v, got %v parsing %v", kind, got, context))
	p.errors = append(p.errors, err)
	fmt.Println("parser error:", err)
	panic(err)
}

func (p *parser) err(format string, args ...interface{}) {
	err := NewParseError(p.prev().Position, fmt.Sprintf(format, args...))
	p.errors = append(p.errors, err)
	fmt.Println("parser error:", err)
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
	if p.pos == 0 {
		return token.New(token.Invalid, token.Position{Line: 1, Column: 0})
	}
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

func (p *parser) Demo() []ast.Statement {
	stmts := make([]ast.Statement, 0)
	for !p.match(token.Eof) {
		stmts = append(stmts, p.Statement())
	}
	return stmts
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

func (p *parser) Program() []ast.Declaration {
	p.decls = make([]ast.Declaration, 0)
	for !p.match(token.Eof) {
		p.decls = append(p.decls, p.Declaration())
	}
	return p.decls
}

// Expressions

func (p *parser) Expression() ast.Expression {
	return p.Or()
}

func (p *parser) leftAssociativeBinaryOp(rule func() ast.Expression, types ...token.Kind) ast.Expression {
	lhs := rule()
	for p.match(types...) {
		lhs = ast.NewBinaryOp(lhs, p.prev().Kind, rule())
	}
	return lhs
}

func (p *parser) Or() ast.Expression {
	return p.leftAssociativeBinaryOp(p.And, token.Or)
}

func (p *parser) And() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Equality, token.And)
}

func (p *parser) Equality() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Comparison, token.Equal, token.NotEqual)
}

func (p *parser) Comparison() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Shift, token.Greater, token.GreaterEq, token.Less, token.LessEq)
}

func (p *parser) Shift() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Term, token.ShiftLeft, token.ShiftRight)
}

func (p *parser) Term() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Factor, token.Minus, token.Plus, token.Pipe, token.Caret)
}

func (p *parser) Factor() ast.Expression {
	return p.leftAssociativeBinaryOp(p.Unary, token.Divide, token.Star, token.Mod, token.Ampersand)
}

func (p *parser) Unary() ast.Expression {
	switch p.peek().Kind {
	case token.Bang, token.Minus, token.Tilde:
		return ast.NewUnaryOp(p.next().Kind, p.Unary())
	default:
		return p.Power()
	}
}

func (p *parser) Power() ast.Expression {
	base := p.Primary()
	if p.match(token.Power) {
		return ast.NewBinaryOp(base, token.Power, p.Power())
	}
	return base
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
	case token.Integer, token.Rational, token.String, token.True, token.False:
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
		ast.PrettyPrint(p.decls)
		p.skipTo(stmtStart)
	})

	switch p.peek().Kind {

	// the first rule in a simple statement is always an expression
	// - parse simple statement on any of the possible first terminals in an expression
	case token.Identifier, token.True, token.False, token.Integer, token.Rational, token.String, token.OpenParen, // atom
		token.Bang, token.Minus, token.Tilde: // unary operators
		s = p.SimpleStmt()
		p.separator("simple stmt")

	case token.Val, token.Var:
		s = &ast.TypedDeclStmt{Decl: p.TypedDecl()}
		p.separator("typed decl stmt")

	case token.Return:
		s = p.ReturnStmt()
		p.separator("return stmt")

	case token.Break:
		p.next()
		s = &ast.BreakStmt{}
		p.separator("break stmt")

	case token.OpenBrace:
		s = p.BlockStmt()
		p.separator("end of block stmt")

	case token.If:
		s = p.IfStmt()
		p.separator("end of if stmt")

	case token.For, token.While:
		s = p.ForStmt()
		p.separator("end of loop stmt")

	case token.Eof:
		// TODO: this happens a lot and is really uninformative when the parser breaks
		// find a way to make debugging this case easier
		ast.PrettyPrint(p.decls)
		for _, err := range p.errors {
			fmt.Println(err)
		}
		s = &ast.ErrStmt{}
		panic("stmt at eof")

	default:
		ast.PrettyPrint(p.decls)
		for _, err := range p.errors {
			fmt.Println(err)
		}
		s = &ast.ErrStmt{}
		panic("invalid token for statement start: " + p.peek().Kind.String())
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
	if p.match(token.Assign, token.PlusEq, token.MinusEq, token.StarEq, token.DivideEq, token.PowerEq, token.ModEq) {
		return ast.NewAssignmentStmt(lhs, p.prev().Kind, p.Expression())
	}

	return &ast.ExprStmt{Expr: lhs}
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
// Statement() skips the closing curly brace in attempt to
// find the next statement in the list, so the parser will
// consume the statements in the outer block. Figure out if
// handling this edge case is worth the complexity.
func (p *parser) StmtList() (stmts []ast.Statement) {
	for p.peek().Kind != token.CloseBrace {
		stmts = append(stmts, p.Statement())
	}
	return stmts
}

// Declarations

func (p *parser) Declaration() (d ast.Declaration) {
	defer sync(func() {
		d = &ast.ErrDecl{}
		ast.PrettyPrint(p.decls)
		p.skipTo(declStart)
	})

	switch kind := p.peek().Kind; kind {

	case token.Val, token.Var:
		d = p.TypedDecl()
		p.separator("end of val/var decl")

	case token.Func:
		d = p.FuncDecl()
		p.separator("end of func decl")

	case token.Class:
		d = p.ClassDecl()
		p.separator("end of class decl")

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
	decl.Type = p.Type()
	return decl
}

func (p *parser) TypedDecl() *ast.TypedDecl {
	if p.match(token.Val, token.Var) {
		decl := &ast.TypedDecl{}
		decl.Kind = p.prev().Kind

		p.consume(token.Identifier, "typed decl identifier")
		decl.Ident = p.prev()

		decl.Type = p.Type()

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

	if p.peekIs(token.Identifier) {
		fd.ReturnType = p.Type()
	}

	fd.Body = p.BlockStmt()
	return fd
}

// Class Parsing

func (p *parser) ClassDecl() *ast.ClassDecl {
	cd := &ast.ClassDecl{}
	p.consume(token.Class, "class declaration start")
	p.consume(token.Identifier, "class name")
	cd.Name = p.prev()
	cd.Body = p.ClassBody()
	return cd
}

func (p *parser) ClassBody() (decls []ast.Declaration) {
	p.consume(token.OpenBrace, "begin class body")
	for p.peekIs(token.Identifier, token.Func) {
		switch p.peek().Kind {
		case token.Identifier:
			decls = append(decls, p.MemberDecl())
		case token.Func:
			decls = append(decls, p.FuncDecl())
		}
		p.separator("end of declaration in class body")
	}
	p.consume(token.CloseBrace, "end class body")
	return decls
}

func (p *parser) MemberDecl() *ast.MemberDecl {
	p.consume(token.Identifier, "class member name")
	return &ast.MemberDecl{
		Name: p.prev(),
		Type: p.Type(),
	}
}

// Miscellaneous

func (p *parser) Type() *ast.TypeExpr {
	p.consume(token.Identifier, "type name")
	lexemes := make([]string, 0, 1)
	lexemes = append(lexemes, p.prev().Lexeme)

	// type is from a module
	for p.match(token.Dot) {
		p.consume(token.Identifier, "imported type name")
		lexemes = append(lexemes, p.prev().Lexeme)
	}
	return &ast.TypeExpr{Name: strings.Join(lexemes, ".")}
}
