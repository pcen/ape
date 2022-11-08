package ape

// statement -> decl
// decl -> (VAL | VAR) ident type  "=" expression
// type -> ident

// expression  -> equality ;
// equality    -> comparison ( ( "!=" | "==" ) comparison )* ;
// comparison  -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term        -> factor ( ( "-" | "+" | "|" | "^" ) factor )* ;
// factor      -> unary ( ( "/" | "*" | "&" ) unary )* ;
// unary       -> ( "!" | "-" | "~" ) unary | primary ;
// primary     -> NUMBER | STRING | "true" | "false" | group ;
// group       -> "(" expression ")" ;

type Parser interface {
	Program() Statement
}

type parser struct {
	tokens []Token
	pos    uint
}

func NewParser(tokens []Token) Parser {
	return &parser{tokens: tokens}
}

func (p *parser) prev() Token {
	return p.tokens[p.pos-1]
}

func (p *parser) match(tt ...TokenType) bool {
	for _, t := range tt {
		if p.tokens[p.pos].Type == t {
			p.pos++
			return true
		}
	}
	return false
}

func (p *parser) Program() Statement {
	stmt := p.Statement()
	if p.match(Eof) {
		return stmt
	}
	panic("program must end with <EOF>")
}

// Expressions

func (p *parser) Expression() Expression {
	return p.Equality()
}

func (p *parser) leftAssociativeBinaryOp(rule func() Expression, types ...TokenType) Expression {
	lhs := rule()
	for p.match(types...) {
		lhs = NewBinaryOp(lhs, p.prev().Type, rule())
	}
	return lhs
}

func (p *parser) Equality() Expression {
	return p.leftAssociativeBinaryOp(p.Comparison, Equal, NotEqual)
}

func (p *parser) Comparison() Expression {
	return p.leftAssociativeBinaryOp(p.Term, Greater, GreaterEq, Less, LessEq)
}

func (p *parser) Term() Expression {
	return p.leftAssociativeBinaryOp(p.Factor, Minus, Plus, BitOr, BitXOR)
}

func (p *parser) Factor() Expression {
	return p.leftAssociativeBinaryOp(p.Unary, Divide, Star, BitAnd)
}

func (p *parser) Unary() Expression {
	if p.match(Bang, Minus, BitNegate) {
		return NewUnaryOp(p.prev().Type, p.Unary())
	}
	return p.Primary()
}

func (p *parser) Primary() Expression {
	if p.match(Number, String, True, False) {
		return NewLiteralExpr(p.prev())
	}
	return p.Group()
}

func (p *parser) Group() (expr Expression) {
	if p.match(OpenParen) {
		expr = p.Expression()
	} else {
		panic("group must start with (")
	}
	if p.match(CloseParen) {
		return expr
	} else {
		panic("group must end with )")
	}
}

// Statements

func (p *parser) Statement() Statement {
	return p.Decl()
}

func (p *parser) Decl() Statement {
	decl := &DeclStmt{}
	if p.match(Val, Var) {
		decl.Kind = p.prev().Type
		if p.match(Identifier) {
			decl.Ident = p.prev().Lexeme
		}
		if p.match(Identifier) {
			decl.Type = p.prev().Lexeme
		}
		if p.match(Assign) {
			decl.Value = p.Expression()
		}
	} else {
		panic("no val or var keyword")
	}
	return decl
}
