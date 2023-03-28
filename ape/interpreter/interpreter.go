package interpreter

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/token"
)

/** === NATIVE FUNCTION CODE LIVES HERE FOR NOW === */

var NATIVE_FUNCTIONS = []val_native_func{
	{
		Name: "println",
		Fn: func(scope *Scope) {
			var sb strings.Builder
			for i := 0; i < len(scope.Values); i++ {
				index := strconv.FormatInt(int64(i), 10)
				val := scope.Values[index]
				sb.WriteString(val.ToString())
			}
			fmt.Println(sb.String())
		},
		Variadic: true,
	},
	{
		Name:   "read",
		Params: []string{"filename"},
		Fn: func(scope *Scope) {
			filename := scope.Get("filename").(val_str)
			bytes, err := os.ReadFile(filename.Value)
			if err != nil {
				panic(err)
			}
			panic(ReturnHolder{Value: val_str{Value: string(bytes)}})
		},
	},
	{
		Name:   "write",
		Params: []string{"filename", "data"},
		Fn: func(scope *Scope) {
			filename := scope.Get("filename").(val_str)
			content := scope.Get("data").(val_str)
			err := os.WriteFile(filename.Value, []byte(content.Value), os.ModePerm)
			if err != nil {
				panic(err)
			}
		},
	},
	{
		Name:   "touch",
		Params: []string{"filename"},
		Fn: func(scope *Scope) {
			filename := scope.Get("filename").(val_str)
			os.Create(filename.Value)
		},
	},
	{
		Name:   "delete",
		Params: []string{"filename"},
		Fn: func(scope *Scope) {
			filename := scope.Get("filename").(val_str)
			os.Remove(filename.Value)
		},
	},
	{
		Name:   "shell",
		Params: []string{"cmd"},
		Fn: func(scope *Scope) {
			cmdstr := scope.Get("cmd").(val_str).Value
			// split := strings.Split(cmdstr, " ")
			// name := split[0]
			// args := []string{}
			// if len(split) > 1 {
			// 	args = split[1:]
			// }
			// fmt.Println(args)
			cmd := exec.Command("bash", "-c", cmdstr)
			b, _ := cmd.CombinedOutput()
			// if err != nil {
			// 	panic(err)
			// }
			fmt.Printf("%s", b)
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
	reversing      bool
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

func (twi *TWI) AddBreadCrumb(node ast.Node) {
	if twi.reversing {
		// only add bread crumbs when executing forwards
		return
	}

	switch n := node.(type) {
	case *ast.AssignmentStmt:

		switch t := n.Lhs.(type) {
		case *ast.IndexExpr:
			name := t.Expr.ExprStr()
			m := twi.evaluateExpr(t.Expr).(val_map)
			idx := twi.evaluateExpr(t.Index)
			twi.LastBreadCrumb = &BreadCrumb{
				Prev:    twi.LastBreadCrumb,
				Scope:   twi.CurrentScope.GetScope(name),
				Name:    name,
				PrevVal: val_index_val_pair{Index: idx, Value: m.Data[idx]},
			}
		default:
			name := t.ExprStr()
			twi.LastBreadCrumb = &BreadCrumb{
				Prev:    twi.LastBreadCrumb,
				Scope:   twi.CurrentScope.GetScope(name),
				Name:    name,
				PrevVal: twi.CurrentScope.Get(name),
			}
		}

	case *ast.ExprStmt:
		// only add bread crumb when the expression is annotated
		if prevVal, ok := n.Annotations["undo"]; ok {
			twi.LastBreadCrumb = &BreadCrumb{
				Prev:    twi.LastBreadCrumb,
				Scope:   twi.CurrentScope,
				PrevVal: prevVal,
			}
		}

	case *ast.IdentExpr:
		name := n.Ident.Lexeme
		twi.LastBreadCrumb = &BreadCrumb{
			Prev:    twi.LastBreadCrumb,
			Scope:   twi.CurrentScope.GetScope(name),
			Name:    name,
			PrevVal: twi.CurrentScope.Get(name),
		}
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
	twi.evaluateExpr(&call_expr)
	// println(resp.(val_int).Value)
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
	case *ast.LitMapExpr:
		return twi.visitLitMapExpr(t)
	case *ast.IndexExpr:
		return twi.visitIndexExpr(t)
	default:
		print(expr.ExprStr())
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

func (twi *TWI) visitIndexExpr(idxExpr *ast.IndexExpr) value {
	m := twi.evaluateExpr(idxExpr.Expr).(val_map)
	return m.Data[twi.evaluateExpr(idxExpr.Index)]
}

func (twi *TWI) visitLitMapExpr(mapVal *ast.LitMapExpr) value {
	val := val_map{Data: map[value]value{}}
	for k, v := range mapVal.Elements {
		res_k := twi.evaluateExpr(k)
		res_v := twi.evaluateExpr(v)
		val.Data[res_k] = res_v
	}
	return val
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
		var fn_scope Scope
		if fn.Variadic {
			fn_scope = MakeVariadicFnScope(twi.GlobalScope, args)
		} else {
			fn_scope = MakeFnScope(twi.GlobalScope, args, fn.Params)
		}
		fn.Fn(&fn_scope)

	default:
		panic(fmt.Sprintf("Trying to call a non function: %s", fn))
	}

	return val_void{}
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
		twi.visitExprStmt(t)
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

func (twi *TWI) visitExprStmt(stmt *ast.ExprStmt) {
	twi.AddBreadCrumb(stmt)
	twi.evaluateExpr(stmt.Expr)
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
	twi.AddBreadCrumb(stmt)

	switch t := stmt.Lhs.(type) {
	case *ast.IndexExpr:
		m := twi.evaluateExpr(t.Expr).(val_map)
		m.Data[twi.evaluateExpr(t.Index)] = twi.evaluateExpr(stmt.Rhs)
	default:
		name := t.ExprStr()
		twi.CurrentScope.Set(name, twi.evaluateExpr(stmt.Rhs))
	}
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
		twi.AddBreadCrumb(t)
		name := t.Ident.Lexeme
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
				twi.reversing = false
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
	twi.reversing = true
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
