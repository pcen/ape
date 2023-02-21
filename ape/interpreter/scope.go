package interpreter

import (
	"fmt"

	"github.com/pcen/ape/ape/types"
)

type Scoped interface {
	Get(name string) types.Type
	Set(name string, value types.Type)
	Define(name string, value types.Type)
}

type Scope struct {
	Enclosing *Scope

	Values map[string]types.Type // Vars declared within this scope
}

/** Travel up scopes looking for given identifier. */
func (s *Scope) Get(name string) types.Type {
	val, exists := s.Values[name]
	if exists {
		return val
	}

	if s.Enclosing != nil {
		return s.Enclosing.Get(name)
	}

	panic(fmt.Sprintf("Failed to find variable with name: %s", name))
}

/** Travel up scopes looking for given identifier to assign to. */
func (s *Scope) Set(name string, value types.Type) {
	if _, exists := s.Values[name]; exists {
		s.Values[name] = value
		return
	}

	if s.Enclosing != nil {
		s.Enclosing.Set(name, value)
	}

	panic(fmt.Sprintf("Failed to find variable with name: %s", name))
}

/** Assign identifier with given value at this scope */
func (s *Scope) Define(name string, value types.Type) {
	s.Values[name] = value
}
