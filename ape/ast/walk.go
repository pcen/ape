package ast

func WalkExpr(root Expression, fn func(Node) bool) {
	switch e := root.(type) {
	case *GroupExpr:
		if fn(e) {
			WalkExpr(e.Expr, fn)
		}
	case *LiteralExpr:
		fn(e)
	case *IdentExpr:
		fn(e)
	case *UnaryOp:
		if fn(e) {
			WalkExpr(e.Expr, fn)
		}
	case *BinaryOp:
		if fn(e) {
			WalkExpr(e.Lhs, fn)
			WalkExpr(e.Rhs, fn)
		}
	case *CallExpr:
		if fn(e) {
			WalkExpr(e.Callee, fn)
			for _, arg := range e.Args {
				WalkExpr(arg, fn)
			}
		}
	case *DotExpr:
		if fn(e) {
			WalkExpr(e.Expr, fn)
			WalkExpr(e.Field, fn)
		}
	case *IndexExpr:
		if fn(e) {
			WalkExpr(e.Expr, fn)
			WalkExpr(e.Index, fn)
		}
	case *TypeExpr:
		fn(e)
	case *LitListExpr:
		if fn(e) {
			for _, el := range e.Elements {
				WalkExpr(el, fn)
			}
		}
	}
}

func WalkStmt(root Statement, fn func(Node) bool) {
	switch s := root.(type) {
	case *BlockStmt:
		if fn(s) {
			for _, stmt := range s.Content {
				WalkStmt(stmt, fn)
			}
		}
	case *ExprStmt:
		if fn(s) {
			WalkExpr(s.Expr, fn)
		}
	case *ReturnStmt:
		if fn(s) {
			WalkExpr(s.Expr, fn)
		}
	case *TypedDeclStmt:
		if fn(s) {
			WalkDecl(s.Decl, fn)
		}
	case *CondBlockStmt:
		if fn(s) {
			WalkExpr(s.Cond, fn)
			WalkStmt(s.Body, fn)
		}
	case *IfStmt:
		if fn(s) {
			WalkStmt(s.If, fn)
			for _, elif := range s.Elifs {
				WalkStmt(elif, fn)
			}
			WalkStmt(s.Else, fn)
		}
	case *ForStmt:
		if fn(s) {
			WalkDecl(s.Init, fn)
			WalkExpr(s.Cond, fn)
			WalkStmt(s.Incr, fn)
			WalkStmt(s.Body, fn)
		}
	case *IncStmt:
		if fn(s) {
			WalkExpr(s.Expr, fn)
		}
	case *AssignmentStmt:
		if fn(s) {
			WalkExpr(s.Lhs, fn)
			WalkExpr(s.Rhs, fn)
		}
	case *BreakStmt:
		fn(s)
	case *SwitchStmt:
		if fn(s) {
			WalkExpr(s.Expr, fn)
			for _, cas := range s.Cases {
				WalkStmt(cas, fn)
			}
		}
	case *CaseStmt:
		if fn(s) {
			WalkExpr(s.Expr, fn)
			WalkStmt(s.Body, fn)
		}
	}
}

func WalkDecl(root Declaration, fn func(Node) bool) {
	switch d := root.(type) {
	case *FuncDecl:
		if fn(d) {
			for _, param := range d.Params {
				WalkDecl(param, fn)
			}
			WalkExpr(d.ReturnType, fn)
			WalkStmt(d.Body, fn)
		}
	case *ClassDecl:
		if fn(d) {
			for _, decl := range d.Body {
				WalkDecl(decl, fn)
			}
		}
	case *MemberDecl:
		if fn(d) {
			WalkExpr(d.Type, fn)
		}
	}
}

func Walk(root Node, fn func(Node) bool) {
	switch r := root.(type) {
	case Declaration:
		WalkDecl(r, fn)
	case Statement:
		WalkStmt(r, fn)
	case Expression:
		WalkExpr(r, fn)
	}
}
