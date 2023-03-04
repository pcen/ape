package interpreter

import (
	"fmt"
	"strconv"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

/** === NATIVE FUNCTION CODE LIVES HERE FOR NOW === */

var NATIVE_FUNCTIONS = []val_native_func{
	{
		Name:   "println",
		Params: []string{"val"},
		Fn: func(scope *Scope) value {
			val := scope.Get("val")
			println(val.ToString())
			return val_void{}
		},
	},
}

type Interpreter interface {
	Interpret(ast.Node)
}

type TWI struct {
	GlobalScope    *Scope
	CurrentScope   *Scope
	LastBreadCrumb *BreadCrumb
}

func NewTWI() *TWI {
	scope := &Scope{
		Enclosing: nil,
		Values:    make(map[string]value),
	}

	// Load in all native functions in global scope
	// This means you could override them in more inner scopes..
	for _, nf := range NATIVE_FUNCTIONS {
		scope.Define(nf.Name, nf)
	}

	return &TWI{
		GlobalScope:    scope,
		CurrentScope:   scope,
		LastBreadCrumb: nil,
	}
}

// ==== TODO: Temp for testing ====
func (twi *TWI) Interpret(decl ast.Declaration) {
	twi.executeDecl(decl)
}

func (twi *TWI) RunMain() {
	call_expr := ast.CallExpr{
		Callee: ast.NewIdentExpr(token.NewLexeme(token.Identifier, "main", token.Position{1, 1})),
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
	// pprintScope(twi.CurrentScope)
	return twi.CurrentScope.Get(ident.Ident.Lexeme)
}

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

	// Handle return values here
	defer func() {
		if panic_val := recover(); panic_val != nil {
			switch holder := panic_val.(type) {
			case ReturnHolder:
				return_val = holder.Value // Use named return to update on panic
				return
			}
			panic(panic_val)
		}
		return_val = val_void{}
	}()

	switch fn := callee.(type) {
	case val_func:
		fn_scope := MakeFnScope(twi.GlobalScope, args, fn.Params)
		twi.visitBlockStmt(&fn_scope, fn.Body)
		return val_void{}

	case val_native_func:
		fn_scope := MakeFnScope(twi.GlobalScope, args, fn.Params)
		return fn.Fn(&fn_scope)

	default:
		panic(fmt.Sprintf("Trying to call a non function: %s", fn))
	}
}

/** === Expression Code Ends === */

/** === Statement Code Begins === */
func (twi *TWI) executeStmt(stmt ast.Statement) {
	switch t := stmt.(type) {
	case *ast.ForStmt:
		twi.visitForStmt(t)
	case *ast.BlockStmt:
		twi.visitBlockStmt(&Scope{twi.CurrentScope, make(map[string]value)}, t)
	case *ast.IfStmt:
		twi.visitIfStmt(t)
	case *ast.ReturnStmt:
		twi.visitReturnStmt(t)
	case *ast.ExprStmt:
		twi.evaluateExpr(t.Expr)
	case *ast.TypedDeclStmt:
		twi.executeDecl(t.Decl)
	case *ast.AssignmentStmt:
		twi.visitAssignmentStmt(t)
	case *ast.IncStmt:
		twi.visitIncStmt(t)
	case *ast.SkipStmt:
		twi.visitSkipStmt(t)
	case *ast.ReverseStmt:
		twi.visitReverseStmt(t)
	}
}

func (twi *TWI) visitForStmt(stmt *ast.ForStmt) {
	twi.executeDecl(stmt.Init) // Initializes far in local scope
	for twi.evaluateExpr(stmt.Cond).(val_bool).Value {
		twi.executeStmt(stmt.Body)
		twi.executeStmt(stmt.Incr)
	}
}

func (twi *TWI) visitBlockStmt(scope *Scope, stmt *ast.BlockStmt) {
	prev_scope := twi.CurrentScope
	twi.CurrentScope = scope

	// Must reset the scope even if we encounter a panic
	defer func() {
		twi.CurrentScope = prev_scope
	}()

	for _, s := range stmt.Content {
		twi.executeStmt(s)
	}

}

func (twi *TWI) visitIfStmt(stmt *ast.IfStmt) {
	result := twi.evaluateExpr(stmt.If.Cond).(val_bool)

	if result.Value {
		twi.executeStmt(stmt.If.Body)
		return
	}

	// Was false. Iterate through elifs now
	for _, elif := range stmt.Elifs {
		result = twi.evaluateExpr(elif.Cond).(val_bool)
		if result.Value {
			twi.executeStmt(elif.Body)
			return
		}
	}

	// Else stmt could be nil
	if stmt.Else != nil {
		twi.executeStmt(stmt.Else)
	}
}

/*
*
This may be a little confusing. We need to stop execution and return to the caller
of the function on a return stmt so we panic with the value. In visitCallExpr you will
see how this value is used.
*/
func (twi *TWI) visitReturnStmt(ret *ast.ReturnStmt) {
	val := twi.evaluateExpr(ret.Expr)
	panic(ReturnHolder{val})
}

func (twi *TWI) visitAssignmentStmt(stmt *ast.AssignmentStmt) {
	// TODO: This only works for simple name assignments
	name := stmt.Lhs.ExprStr()
	twi.LastBreadCrumb = &BreadCrumb{
		Prev:    twi.LastBreadCrumb,
		Scope:   twi.CurrentScope.GetScope(name),
		Name:    name,
		PrevVal: twi.CurrentScope.Get(name)}
	twi.CurrentScope.Set(name, twi.evaluateExpr(stmt.Rhs))
}

func (twi *TWI) visitIncStmt(inc *ast.IncStmt) {
	val := twi.evaluateExpr(inc.Expr).(number)
	switch inc.Op.Kind {
	case token.Increment:
		val = val.Add(val_int{1})
	case token.Decrement:
		val = val.Subtract(val_int{1})
	}

	switch t := inc.Expr.(type) {
	case *ast.IdentExpr:
		name := t.Ident.Lexeme
		twi.LastBreadCrumb = &BreadCrumb{
			Prev:    twi.LastBreadCrumb,
			Scope:   twi.CurrentScope.GetScope(name),
			Name:    name,
			PrevVal: twi.CurrentScope.Get(name)}
		twi.CurrentScope.Set(name, val.(value))
	}
}

func (twi *TWI) visitSkipStmt(stmt *ast.SkipStmt) {
	// Handle return values here
	defer func() {
		if panic_val := recover(); panic_val != nil {
			switch holder := panic_val.(type) {
			case ReturnHolder:
				// Reset the last LastBreadCrumb to point to the bread crumb before this skip, without reverse executing
				// This is necessary to support a return within a skip statement
				for twi.LastBreadCrumb.SkipMarker != stmt {
					twi.LastBreadCrumb = twi.LastBreadCrumb.Prev
				}

				twi.LastBreadCrumb = twi.LastBreadCrumb.Prev // Remove the SkipMarker

			case ReverseHolder:
				// Reverse any assignment statements Before the current SkipMarker
				for twi.LastBreadCrumb.SkipMarker != stmt {
					twi.LastBreadCrumb.Reverse(twi)
					twi.LastBreadCrumb = twi.LastBreadCrumb.Prev
				}

				twi.LastBreadCrumb = twi.LastBreadCrumb.Prev // Remove the SkipMarker

				for _, seize := range stmt.Seizes {
					seize_val := twi.evaluateExpr(seize.Expr)

					if seize_val.Equals(holder.Value) {
						twi.executeStmt(seize.Body)
						return // Exit the Panic Loop
					}

				}
			}
			panic(panic_val) // Propagate panic
		}
	}()

	// Mark the start of the current skip
	twi.LastBreadCrumb = &BreadCrumb{
		Prev:       twi.LastBreadCrumb,
		SkipMarker: stmt,
	}

	twi.executeStmt(stmt.Body)
}

func (twi *TWI) visitReverseStmt(rev *ast.ReverseStmt) {
	// Handle reverse values here
	val := twi.evaluateExpr(rev.Expr)
	panic(ReverseHolder{val})
}

/** === Statement Code Ends === */

/** === Declaration Code Begins === */
func (twi *TWI) executeDecl(decl ast.Declaration) {
	switch t := decl.(type) {
	case *ast.FuncDecl:
		twi.visitFuncDecl(t)
	case *ast.VarDecl:
		twi.visitVarDecl(t)
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

	twi.CurrentScope.Define(fn.Name, fn)
}

func (twi *TWI) visitVarDecl(var_decl *ast.VarDecl) {
	scope := twi.CurrentScope
	scope.Define(var_decl.Ident.Lexeme, twi.evaluateExpr(var_decl.Value))

	// pprintScope(twi.CurrentScope)
}

/** === Declaration Code Ends === */
