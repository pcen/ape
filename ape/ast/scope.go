package ast

type Scope struct {
	Module  string
	Outer   *Scope
	Symbols map[string]*Symbol
}

type Symbol struct {
	Name string
}
