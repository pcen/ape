package interpreter

import (
	"fmt"
	"math"
	"strconv"

	"github.com/pcen/ape/ape/ast"
)

/** Base interface of all values */
type value interface {
	Equals(value) bool
	ToString() string
}

/** Represents nothing. Ex. no return from a func */
type val_void struct{}

func (v val_void) Equals(other value) bool {
	switch other.(type) {
	case val_void:
		return true
	default:
		return false
	}
}

func (v val_void) ToString() string {
	return "VOID"
}

/** Needed to easily support breadcrumb reversal for maps */
type val_index_val_pair struct {
	Index value
	Value value
}

func (vivp val_index_val_pair) Equals(other value) bool {
	return false
}

func (vivp val_index_val_pair) ToString() string {
	return vivp.Index.ToString() + ": " + vivp.Value.ToString()
}

type val_native_func struct {
	Name     string
	Params   []string
	Fn       func(*Scope)
	Variadic bool
}

// Pointless
func (vnf val_native_func) Equals(other value) bool {
	return false
}

func (vnf val_native_func) ToString() string {
	return "NATIVE: " + vnf.Name
}

/** Internal representation of a function. First class citizen */
type val_func struct {
	Name   string
	Params []string
	Body   *ast.BlockStmt
}

func (fn val_func) Equals(other value) bool {
	switch other.(type) {
	case val_func:
		return (fn.Name == other.(val_func).Name)
	default:
		return false // Impossible with type checking
	}
}

func (v val_func) ToString() string {
	return "FUNC: " + v.Name
}

//TODO: List
// type val_list[T value] struct {
// 	Size int
// 	Data []T
// }

type val_map struct {
	Data map[value]value
}

func (m val_map) Equals(other value) bool {
	switch t := other.(type) {
	case val_map:
		if len(m.Data) != len(t.Data) {
			return false
		}

		for k, v := range m.Data {
			o_v, exists := t.Data[k]
			if !exists {
				return false
			}
			if !v.Equals(o_v) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (m val_map) ToString() string {
	out := ""
	for k, v := range m.Data {
		out += k.ToString() + ": " + v.ToString() + ", "
	}
	return out[:len(out)-2]
}

type val_bool struct {
	Value bool
}

func (b val_bool) Equals(other value) bool {
	switch other.(type) {
	case val_bool:
		return b.Value == other.(val_bool).Value
	default:
		return false // Impossible with type checking
	}
}

func (b val_bool) ToString() string {
	if b.Value {
		return "True"
	} else {
		return "False"
	}
}

type val_str struct {
	Value string
}

func (s val_str) Equals(other value) bool {
	switch other.(type) {
	case val_str:
		return s.Value == other.(val_str).Value
	default:
		return false
	}
}

func (s val_str) ToString() string {
	return s.Value
}

/** Interface for both Integers and Rationals*/
type number interface {
	Add(number) number
	Subtract(number) number
	Multiply(number) number
	Divide(number) number
	Power(number) number
	LessThan(number) val_bool
	LessThanEq(number) val_bool
	GreaterThan(number) val_bool
	GreaterThanEq(number) val_bool
}

/*
*
Maybe all numbers should be contained within one struct to reduce the
duplication of operations seen below. Everything could be in the highest precision
float possible and then just when accessing the value we could round up and down.
I though this may be bad for weird floating point errors when the user defined ints
so I opted out.
*/
type val_int struct {
	Value int
}

func (v val_int) Equals(other value) bool {
	switch other.(type) {
	case val_int:
		return v.Value == other.(val_int).Value
	default:
		return false
	}
}

func (v val_int) ToString() string {
	return strconv.FormatInt(int64(v.Value), 10)
}

func (v val_int) Add(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_int{v.Value + other.(val_int).Value}
	case val_rational:
		return val_rational{float64(v.Value) + other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't add int and %s", t))
	}
}

func (v val_int) Subtract(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_int{v.Value - other.(val_int).Value}
	case val_rational:
		return val_rational{float64(v.Value) + other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't subtract int and %s", t))
	}
}

func (v val_int) Multiply(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_int{v.Value * other.(val_int).Value}
	case val_rational:
		return val_rational{float64(v.Value) + other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't multiply int and %s", t))
	}
}

func (v val_int) Divide(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_int{v.Value / other.(val_int).Value}
	case val_rational:
		return val_rational{float64(v.Value) + other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't divide int and %s", t))
	}
}

func (v val_int) Power(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_int{int(math.Pow(float64(v.Value), float64(other.(val_int).Value)))}
	case val_rational:
		return val_rational{math.Pow(float64(v.Value), other.(val_rational).Value)}
	default:
		panic(fmt.Sprintf("Can't exponentiate int and %s", t))
	}
}

func (v val_int) Mod(other val_int) val_int {
	return val_int{v.Value % other.Value}
}

func (v val_int) LessThan(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value < other.(val_int).Value}
	case val_rational:
		return val_bool{float64(v.Value) < other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}

func (v val_int) LessThanEq(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value <= other.(val_int).Value}
	case val_rational:
		return val_bool{float64(v.Value) <= other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}

func (v val_int) GreaterThan(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value > other.(val_int).Value}
	case val_rational:
		return val_bool{float64(v.Value) > other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}

func (v val_int) GreaterThanEq(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value >= other.(val_int).Value}
	case val_rational:
		return val_bool{float64(v.Value) >= other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}

type val_rational struct {
	Value float64
}

func (v val_rational) Equals(other value) bool {
	switch other.(type) {
	case val_rational:
		return v.Value == other.(val_rational).Value
	default:
		return false
	}
}

func (v val_rational) ToString() string {
	return strconv.FormatInt(int64(v.Value), 10)
}

func (v val_rational) Add(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_rational{v.Value + float64(other.(val_int).Value)}
	case val_rational:
		return val_rational{v.Value + other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't add rational and %s", t))
	}
}

func (v val_rational) Subtract(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_rational{v.Value - float64(other.(val_int).Value)}
	case val_rational:
		return val_rational{v.Value - other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't subtract rational and %s", t))
	}
}

func (v val_rational) Multiply(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_rational{v.Value * float64(other.(val_int).Value)}
	case val_rational:
		return val_rational{v.Value * other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't multiply rational and %s", t))
	}
}

func (v val_rational) Divide(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_rational{v.Value / float64(other.(val_int).Value)}
	case val_rational:
		return val_rational{v.Value / other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't divide rational and %s", t))
	}
}

func (v val_rational) Power(other number) number {
	switch t := other.(type) {
	case val_int:
		return val_rational{math.Pow(v.Value, float64(other.(val_int).Value))}
	case val_rational:
		return val_rational{math.Pow(float64(v.Value), other.(val_rational).Value)}
	default:
		panic(fmt.Sprintf("Can't exponentiate rational and %s", t))
	}
}

func (v val_rational) LessThan(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value < float64(other.(val_int).Value)}
	case val_rational:
		return val_bool{v.Value < other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}

func (v val_rational) LessThanEq(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value <= float64(other.(val_int).Value)}
	case val_rational:
		return val_bool{v.Value <= other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}
func (v val_rational) GreaterThan(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value > float64(other.(val_int).Value)}
	case val_rational:
		return val_bool{v.Value > other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}
func (v val_rational) GreaterThanEq(other number) val_bool {
	switch t := other.(type) {
	case val_int:
		return val_bool{v.Value >= float64(other.(val_int).Value)}
	case val_rational:
		return val_bool{v.Value >= other.(val_rational).Value}
	default:
		panic(fmt.Sprintf("Can't compare rational and %s", t))
	}
}
