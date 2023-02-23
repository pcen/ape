package interpreter

import (
	"fmt"
	"strconv"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

type Interpreter interface {
	Interpret(ast.Node)
}

type TWI struct {
	GlobalScope  *Scope
	CurrentScope *Scope
}

func NewTWI() *TWI {
	scope := &Scope{
		Enclosing: nil,
		Values:    make(map[string]value),
	}
	return &TWI{
		GlobalScope:  scope,
		CurrentScope: scope,
	}
}

// ==== TODO: Temp for testing ====
func (twi *TWI) Interpret(decl ast.Declaration) {
	twi.executeDecl(decl)
}

func (twi *TWI) RunMain() {

	println("STATE OF GLOBALS")
	for k, v := range twi.GlobalScope.Values {
		println(fmt.Sprintf("KEY %s VALUE %+v", k, v.(val_func)))
	}

	call_expr := ast.CallExpr{
		Callee: ast.NewLiteralExpr(token.NewLexeme(token.String, "main", token.Position{1, 1})),
		Args:   []ast.Expression{},
	}
	resp := twi.evaluateExpr(&call_expr)
	println(resp.(val_int).Value)
}

// ====== TESTING =====

/** === Expression Code Begins ==== */
func (twi *TWI) evaluateExpr(expr ast.Expression) value {
	switch t := expr.(type) {
	case *ast.LiteralExpr:
		return twi.visitLiteralExpr(t)
	case *ast.IdentExpr:
		return twi.visitIdentExpr(t)
	case *ast.BinaryOp:
		return twi.visitBinaryExpr(t)
	case *ast.GroupExpr:
		return twi.visitGroupExpr(t)
	case *ast.CallExpr:
		return twi.visitCallExpr(t)
	default:
		panic(fmt.Sprintf("Expression type cannot be evaluated: %+v", t))
	}
}

func (twi *TWI) visitLiteralExpr(literal *ast.LiteralExpr) value {
	switch literal.Kind {
	case token.String:
		return val_str{literal.Lexeme}
	case token.Integer:
		val, _ := strconv.Atoi(literal.Lexeme)
		return val_int{val}
	case token.Rational:
		val, _ := strconv.ParseFloat(literal.Lexeme, 64)
		return val_rational{val}
	case token.True:
		return val_bool{true}
	case token.False:
		return val_bool{false}
	default:
		panic(fmt.Sprintf("Unknown literal expression kind: %s", literal.Kind))
	}
}

func (twi *TWI) visitIdentExpr(ident *ast.IdentExpr) value {
	return val_str{ident.Ident.Lexeme}
}

// TODO: Handle the equal expressions (+=, -=, etc...)
func (twi *TWI) visitBinaryExpr(bin *ast.BinaryOp) value {
	lv := twi.evaluateExpr(bin.Lhs)
	rv := twi.evaluateExpr(bin.Rhs)

	switch bin.Op.Kind {
	case token.Plus:
		switch lv.(type) {
		case val_str:
			return val_str{lv.(val_str).Value + rv.(val_str).Value}
		default:
			return lv.(number).Add(rv.(number)).(value) // Typechecked: we know this is a number
		}
	case token.Minus:
		return lv.(number).Subtract(rv.(number)).(value)
	case token.Star:
		return lv.(number).Multiply(rv.(number)).(value)
	case token.Divide:
		return lv.(number).Divide(rv.(number)).(value)
	case token.Power:
		return lv.(number).Power(rv.(number)).(value)
	case token.Mod:
		return lv.(val_int).Mod(rv.(val_int))
	case token.Less:
		return lv.(number).LessThan(rv.(number))
	case token.LessEq:
		return lv.(number).LessThanEq(rv.(number))
	case token.Greater:
		return lv.(number).GreaterThan(rv.(number))
	case token.GreaterEq:
		return lv.(number).GreaterThanEq(rv.(number))
	case token.Equal:
		return val_bool{lv.Equals(rv)}
	case token.NotEqual:
		return val_bool{!lv.Equals(rv)}
	case token.And:
		return val_bool{lv.(val_bool).Value && rv.(val_bool).Value}
	case token.Or:
		return val_bool{lv.(val_bool).Value || rv.(val_bool).Value}
	}

	panic(fmt.Sprintf("Unknown binary operation: %s", bin.Op.Kind))
}

func (twi *TWI) visitUnaryExpr(unary *ast.UnaryOp) value {
	val := twi.evaluateExpr(unary.Expr)

	switch unary.Op {
	case token.Bang:
		return val_bool{!val.(val_bool).Value}
	default:
		panic("Unknown unary token")
	}
}

func (twi *TWI) visitGroupExpr(group *ast.GroupExpr) value {
	return twi.evaluateExpr(group.Expr)
}

func (twi *TWI) visitCallExpr(expr *ast.CallExpr) (return_val value) {
	// Resolved value we are calling
	// Could be other_fn() or something more convoluted:
	// fn_generator("hello")(" Alex") AKA call a fn returned from a fn
	callee := twi.evaluateExpr(expr.Callee)

	// Evaluate all the arguments
	args := []value{}
	for _, arg := range expr.Args {
		args = append(args, twi.evaluateExpr(arg))
	}

	switch t := callee.(type) {
	case val_str:
		fn_name := t.Value
		fn := twi.CurrentScope.Get(fn_name).(val_func)
		fn_scope := MakeFnScope(twi.GlobalScope, args, fn.Params)

		fmt.Println("CALLING: ", fn_name)

		defer func() {
			ret_val := recover()
			if ret_val != nil {
				switch ret_val.(type) {
				case *value:
				}
				println(ret_val)
				return_val = ret_val.(value) // Use named return to update on panic
			} else {
				return_val = val_void{}
			}
		}()
		twi.visitBlockStmt(&fn_scope, fn.Body)
		return val_void{}

	default:
		panic(fmt.Sprintf("Trying to call a non function: %s", t))
	}
}

/** === Expression Code Ends === */

/** === Statement Code Begins === */
func (twi *TWI) executeStmt(stmt ast.Statement) {
	switch t := stmt.(type) {
	case *ast.BlockStmt:
		twi.visitBlockStmt(&Scope{twi.CurrentScope, make(map[string]value)}, t)
	case *ast.ReturnStmt:
		twi.visitReturnStmt(t)
	case *ast.ExprStmt:
		twi.evaluateExpr(t.Expr)
	}
}

func (twi *TWI) visitBlockStmt(scope *Scope, stmt *ast.BlockStmt) {
	prev_scope := twi.CurrentScope
	twi.CurrentScope = scope

	for _, s := range stmt.Content {
		twi.executeStmt(s)
	}

	twi.CurrentScope = prev_scope
}

/**
This may be a little confusing. We need to stop execution and return to the caller
of the function on a return stmt so we panic with the value. In visitCallExpr you will
see how this value is used.
*/
func (twi *TWI) visitReturnStmt(ret *ast.ReturnStmt) {
	val := twi.evaluateExpr(ret.Expr)
	panic(val)
}

/** === Statement Code Ends === */

/** === Declaration Code Begins === */
func (twi *TWI) executeDecl(decl ast.Declaration) {
	switch t := decl.(type) {
	case *ast.FuncDecl:
		twi.visitFuncDecl(t)
	}
}

func (twi *TWI) visitFuncDecl(fn_decl *ast.FuncDecl) {
	param_names := []string{}

	for _, p := range fn_decl.Params {
		param_names = append(param_names, p.Ident.ExprStr())
	}

	fn := val_func{
		Name:   fn_decl.Name.Lexeme,
		Params: param_names,
		Body:   fn_decl.Body,
	}

	twi.CurrentScope.Values[fn.Name] = fn
}

/** === Declaration Code Ends === */
