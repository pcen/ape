package ape

import (
	"fmt"
	"strings"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
	"github.com/pcen/ape/ape/types"
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
		token.Identifier: true,
		token.Func:       true,
	}

	stmtStart = map[token.Kind]bool{
		token.Identifier: true,
		token.Return:     true,
		token.OpenBrace:  true,
		token.If:         true,
		token.For:        true,
		token.While:      true,
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
	BlockStmt() *ast.BlockStmt
	Errors() ([]ParseError, bool)
}

type parser struct {
	tokens []token.Token
	pos    int
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

func (p *parser) peekn(n int) token.Token {
	return p.tokens[p.pos+n-1]
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
		lhs = ast.NewBinaryOp(lhs, p.prev(), rule())
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
		return ast.NewBinaryOp(base, p.prev(), p.Power())
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
	case token.OpenBrack:
		return p.LitList()
	case token.OpenBrace:
		return p.LitMap()
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

func (p *parser) LitList() ast.Expression {
	p.consume(token.OpenBrack, "start of list literal")
	// need to abstract function for comma separated list of expressions
	var elements []ast.Expression
	for !p.peekIs(token.CloseBrack) {
		elements = append(elements, p.Expression())
		if !p.peekIs(token.CloseBrack) {
			p.consume(token.Comma, "list literal elements must be comma separated")
		}
	}
	p.consume(token.CloseBrack, "end of list literal")
	return &ast.LitListExpr{Elements: elements}
}

func (p *parser) LitMap() ast.Expression {
	p.consume(token.OpenBrace, "start of map literal")
	elements := make(map[ast.Expression]ast.Expression)
	for !p.peekIs(token.CloseBrace) {
		k := p.Expression()
		p.consume(token.Colon, "colon separates map key and value in kvp")
		v := p.Expression()
		if p.peekIs(token.Comma) {
			p.consume(token.Comma, "comma separates map key-value pairs")
		}
		elements[k] = v
	}
	p.consume(token.CloseBrace, "end of map literal")
	return &ast.LitMapExpr{Elements: elements}
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
		// ast.PrettyPrint(p.decls)
		p.skipTo(stmtStart)
	})

	switch p.peek().Kind {

	// the first rule in a simple statement is always an expression
	// - parse simple statement on any of the possible first terminals in an expression
	// - unfortunately, since a simple statement can be an expression, this includes list
	//   literals. we could prevent literals here, which would be simple to parse but would
	//   technically complicate the grammar
	case token.Identifier, token.True, token.False, token.Integer, token.Rational, token.String, token.OpenParen, token.OpenBrack, // atom
		token.Bang, token.Minus, token.Tilde, token.Reverse: // unary operators
		s = p.SimpleStmt(true)
		p.separator("simple stmt")

	case token.Return:
		s = p.ReturnStmt()
		p.separator("return stmt")

	case token.Break:
		p.next()
		s = &ast.BreakStmt{}
		p.separator("break stmt")

	case token.Switch:
		s = p.SwitchStmt()
		p.separator("end of switch statement")

	case token.Fallthrough:
		p.next()
		s = &ast.FallthroughtStmt{}
		p.separator("fallthrough stmt")

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

	case token.Skip:
		s = p.SkipStmt()
		p.separator("skip stmt")

	default:
		// ast.PrettyPrint(p.decls)
		for _, err := range p.errors {
			fmt.Println(err)
		}
		s = &ast.ErrStmt{}
		panic("invalid token for statement start: " + p.peek().Kind.String())
	}

	return s
}

func (p *parser) SimpleStmt(annotateable bool) ast.Statement {
	// reverse
	if p.peekIs(token.Reverse) {
		return p.ReverseStmt()
	}

	// declaration
	if p.peekIs(token.Identifier) && p.peekn(2).Kind == token.Colon {
		return &ast.TypedDeclStmt{Decl: p.VarDecl()}
	}

	// increment / decrement
	lhs := p.Expression()
	if p.match(token.Increment, token.Decrement) {
		return &ast.IncStmt{
			Expr: lhs,
			Op:   p.prev(),
		}
	}

	// assignment
	if p.match(token.Assign, token.PlusEq, token.MinusEq, token.StarEq, token.DivideEq, token.PowerEq, token.ModEq) {
		return ast.NewAssignmentStmt(lhs, p.prev(), p.Expression())
	}

	// expression
	annotations := make(map[string]ast.Statement)
	for p.match(token.At) {
		if !annotateable {
			p.err("@ at end of non-annotatable statement")
		}
		if !p.match(token.Identifier) {
			p.err("@ must be followed by annotation name, but got %v", p.peek())
		}
		annotations[p.prev().Lexeme] = p.SimpleStmt(false)
	}
	return &ast.ExprStmt{Expr: lhs, Annotations: annotations}
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

func (p *parser) SwitchStmt() *ast.SwitchStmt {
	// TODO: make sure that there is only 1 default case in the switch statement
	stmt := &ast.SwitchStmt{Cases: make([]*ast.CaseStmt, 0)}
	p.consume(token.Switch, "switch stmt start")
	stmt.Token = p.prev()
	stmt.Expr = p.Expression()
	p.consume(token.OpenBrace, "switch stmt open brace")
	for p.peekIs(token.Case, token.Default) {
		stmt.Cases = append(stmt.Cases, p.CaseStmt())
	}
	p.consume(token.CloseBrace, "switch stmt end")
	return stmt
}

func (p *parser) CaseStmt() *ast.CaseStmt {
	stmt := &ast.CaseStmt{}
	stmt.Body = &ast.BlockStmt{Content: make([]ast.Statement, 0)}
	if p.peekIs(token.Case) {
		p.consume(token.Case, "start of case statement")
		stmt.Token = p.prev()
		stmt.Expr = p.Expression()
	} else {
		p.consume(token.Default, "start of default case statement")
		stmt.Token = p.prev()
	}
	p.consume(token.Colon, "case expression is followed by colon")
	// The block statement in a switch case is parsed differently than
	// regular block statements because it does not require opening and
	// closing curly braces. A case block statement is done when the next
	// case (including default) begins, or the closing brace of the entire
	// switch statement is next.
	for !p.peekIs(token.Case, token.Default, token.CloseBrace) {
		stmt.Body.Content = append(stmt.Body.Content, p.Statement())
	}
	return stmt
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
		Cond: p.Expression(),
		Body: p.BlockStmt(),
	}
}

func (p *parser) ForStmt() *ast.ForStmt {
	s := &ast.ForStmt{}
	switch p.next().Kind {

	case token.For:
		s.Init = p.VarDecl()
		p.separator("after for loop init")
		s.Cond = p.Expression()
		p.separator("after for loop condition")
		s.Incr = p.SimpleStmt(false)

	case token.While:
		s.Cond = p.Expression()

	default:
		p.err("%v cannot start a loop statement", p.prev())
	}
	s.Body = p.BlockStmt()
	return s
}

func (p *parser) ReverseStmt() *ast.ReverseStmt {
	s := &ast.ReverseStmt{}
	p.consume(token.Reverse, "reverse stmt")
	if p.peek().Kind != token.Sep && p.peek().Kind != token.OpenBrace {
		s.Expr = p.Expression()
	}
	return s
}

func (p *parser) SkipStmt() *ast.SkipStmt {
	s := &ast.SkipStmt{}
	p.consume(token.Skip, "skip stmt")
	s.Body = p.BlockStmt()
	s.Seizes = make([]*ast.SeizeStmt, 0)
	for p.peekIs(token.Seize) {
		s.Seizes = append(s.Seizes, p.SeizeStmt())
	}
	return s
}

func (p *parser) SeizeStmt() *ast.SeizeStmt {
	s := &ast.SeizeStmt{}
	p.consume(token.Seize, "seize stmt")
	if !p.peekIs(token.OpenBrace) {
		s.Expr = p.Expression()
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
	for !p.peekIs(token.CloseBrace) {
		stmts = append(stmts, p.Statement())
	}
	return stmts
}

// Declarations

func (p *parser) Declaration() (d ast.Declaration) {
	defer sync(func() {
		d = &ast.ErrDecl{}
		// ast.PrettyPrint(p.decls)
		p.skipTo(declStart)
	})

	switch kind := p.peek().Kind; kind {

	case token.Func:
		d = p.FuncDecl()
		p.separator("end of func decl")

	case token.Class:
		d = p.ClassDecl()
		p.separator("end of class decl")

	case token.Identifier:
		d = p.VarDecl()
		p.separator("end of variable decl")

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
		decl.Ident = &ast.IdentExpr{Ident: p.prev()}
	}
	decl.Type = p.Type()
	return decl
}

func (p *parser) VarDecl() *ast.VarDecl {
	decl := &ast.VarDecl{}
	if p.match(token.Identifier) {
		decl.Ident = p.prev()
	} else {
		panic("var decl must have identifier")
	}
	p.consume(token.Colon, "var decl colon after identifier")
	if p.match(token.Assign, token.Colon) {
		// "foo :=" or "foo ::"
		decl.Mutable = p.prev().Kind == token.Assign
		decl.Value = p.Expression()
	} else {
		// "foo : bar"
		decl.Type = p.Type()
		if p.match(token.Assign, token.Colon) {
			// "foo : bar = baz" or "foo : bar : baz"
			decl.Mutable = p.prev().Kind == token.Assign
			decl.Value = p.Expression()
		}
	}
	return decl
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
	} else {
		fd.ReturnType = &ast.TypeExpr{Name: types.Void.String()}
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
	list := false
	if p.match(token.OpenBrack) && p.match(token.CloseBrack) {
		list = true
	}

	p.consume(token.Identifier, "type name")
	lexemes := make([]string, 0, 1)
	lexemes = append(lexemes, p.prev().Lexeme)

	// type is from a module
	for p.match(token.Dot) {
		p.consume(token.Identifier, "imported type name")
		lexemes = append(lexemes, p.prev().Lexeme)
	}
	return &ast.TypeExpr{Name: strings.Join(lexemes, "."), List: list}
}
