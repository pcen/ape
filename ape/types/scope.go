package types

import "fmt"

type Scope struct {
	Parent  *Scope
	Types   map[string]Type
	Symbols map[string]Type
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		Types:   make(map[string]Type),
		Symbols: make(map[string]Type),
	}
}

func (s *Scope) LookupType(name string) (Type, bool) {
	typ, ok := s.Types[name]
	if !ok && s.Parent != nil {
		typ, ok = s.Parent.LookupType(name)
	}
	return typ, ok
}

func (s *Scope) DeclareType(name string) error {
	if _, ok := s.LookupType(name); ok {
		return fmt.Errorf("type \"%v\" already declared in this scope", name)
	}
	s.Types[name] = NewNamed(name)
	return nil
}

func (s *Scope) LookupSymbol(name string) (Type, bool) {
	typ, ok := s.Symbols[name]
	if !ok && s.Parent != nil {
		typ, ok = s.Parent.LookupSymbol(name)
	}
	if typ == nil {
		typ = Invalid
	}
	return typ, ok
}

func (s *Scope) DeclareSymbol(name string, typ Type) error {
	if _, ok := s.LookupSymbol(name); ok {
		return fmt.Errorf("cannot redeclare \"%v\"", name)
	}
	s.Symbols[name] = typ
	return nil
}

func (s *Scope) Print() {
	fmt.Println("types:")
	for name := range s.Types {
		fmt.Printf("type: %v\n", name)
	}
	fmt.Println("\nsymbols:")
	for ident, typ := range s.Symbols {
		fmt.Printf("%v: %v\n", ident, typ)
	}
}

func GlobalScope() *Scope {
	scope := NewScope(nil)
	for typ := range primitives {
		scope.Types[typ.String()] = typ
	}
	// TODO: properly define all builtin function signatures somewhere
	scope.Symbols["println"] = NewFunction([]Type{Any}, []Type{Void})
	return scope
}
