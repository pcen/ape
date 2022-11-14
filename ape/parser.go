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

// blockStmt   -> "{" statement* "}"

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

var (
	declStart = map[token.Kind]bool{
		token.Val:  true,
		token.Var:  true,
		token.Func: true,
	}
)

type ParseError struct {
	Pos  token.Position
	What string
}

func NewParseError(pos token.Position, format string, a ...interface{}) ParseError {
	return ParseError{Pos: pos, What: fmt.Sprintf(format, a...)}
}

func (p ParseError) String() string {
	return fmt.Sprintf("%v: %v", p.Pos, p.What)
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

func (p *parser) errExpected(kind token.Kind, context string) {
	pos := p.prev().Position
	got := p.prev().String()
	err := NewParseError(pos, fmt.Sprintf("%v: expected %v, got %v parsing %v", pos, kind, got, context))
	p.errors = append(p.errors, err)
	panic(err)
}

func (p *parser) err(format string, args ...interface{}) {
	pos := p.prev().Position
	err := NewParseError(pos, fmt.Sprintf(format, args...))
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

func (p *parser) consume(tk token.Kind, context string) {
	if p.tokens[p.pos].Kind == tk {
		p.pos++
		return
	}
	p.errExpected(tk, context)
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
		op := p.prev().Kind
		rhs := rule()
		if err, ok := rhs.(*ast.InvalidExpr); ok {
			fmt.Printf("invalid operand, expected expression for binary operator: %v\n", err.What)
		}
		lhs = ast.NewBinaryOp(lhs, op, rhs)
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
		p.consume(token.CloseParen, "expect \")\" at end of call")
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
	switch p.peek().Kind {
	case token.Number, token.String, token.True, token.False:
		return ast.NewLiteralExpr(p.next())
	case token.Identifier:
		return ast.NewIdentExpr(p.next())
	case token.OpenParen:
		return p.GroupExpr()
	default:
		return &ast.InvalidExpr{}
	}
}

func (p *parser) GroupExpr() (expr ast.Expression) {
	p.consume(token.OpenParen, "group starts with '('")
	expr = p.Expression()
	p.consume(token.CloseParen, "group ends with ')'")
	return &ast.GroupExpr{Expr: expr}
}

// Statements

func (p *parser) Statement() (s ast.Statement) {
	defer func() {
		if err, ok := recover().(ParseError); ok {
			fmt.Println("parser: error parsing statement:")
			fmt.Println(err)
		}
		s = nil
	}()
	switch p.peek().Kind {
	case token.Return:
		s = p.ReturnStmt()
	case token.OpenBrace:
		s = p.BlockStmt()
	default:
		s = p.ExprStmt()
	}
	return s
}

func (p *parser) ReturnStmt() *ast.ReturnStmt {
	p.consume(token.Return, "return statement begins with \"return\"")
	return &ast.ReturnStmt{Expr: p.Expression()}
}

func (p *parser) BlockStmt() *ast.BlockStmt {
	content := make([]ast.Statement, 0)
	p.consume(token.OpenBrace, "block begins with \"{\"")
	for !p.match(token.CloseBrace) {
		content = append(content, p.Statement())
	}
	return &ast.BlockStmt{Content: content}
}

func (p *parser) ExprStmt() ast.Statement {
	return &ast.ExprStmt{Expr: p.Expression()}
}

// Declarations

func syncFunc[N ast.Node](f func()) {
	err := recover()
	if _, ok := err.(ParseError); err != nil && !ok {
		panic(err)
	}
	f()
}

func (p *parser) Declaration() (d ast.Declaration) {
	defer syncFunc[ast.Declaration](func() {
		d = &ast.ErrDecl{}
		fmt.Println("skipping to decl start")
		p.skipTo(declStart)
	})
	// defer func() {
	// 	err := recover()
	// 	if _, ok := err.(ParseError); err != nil && !ok {
	// 		panic(err)
	// 	}
	// 	d = &ast.ErrDecl{}
	// 	p.skipTo(declStart)
	// }()
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
		// empty parameter list
		return decls
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

func (p *parser) TypedDecl() ast.Declaration {
	decl := &ast.TypedDecl{}
	if p.match(token.Val, token.Var) {
		decl.Kind = p.prev().Kind

		if p.match(token.Identifier) {
			decl.Ident = p.prev()
		}

		if p.match(token.Identifier) {
			decl.Type = p.prev().Lexeme
		} else {
			p.err("missing type for variable declaration %v", decl.Ident)
		}

		if p.match(token.Assign) {
			decl.Value = p.Expression()
		}
	} else {
		p.err("missing val or var for typed decl")
	}
	return decl
}

func (p *parser) FuncDecl() ast.Declaration {
	fd := &ast.FuncDecl{}
	p.consume(token.Func, "function declaration expects \"func\"")
	if p.match(token.Identifier) {
		fd.Name = p.prev()
	}
	if p.match(token.OpenParen) {
		fd.Params = p.Parameters()
	}
	p.consume(token.CloseParen, "parameters end with \")\"")
	fd.Body = p.BlockStmt()
	return fd
}
