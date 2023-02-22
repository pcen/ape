package interpreter

import (
	"fmt"
	"math"
)

/** Base interface of all values */
type value interface {
	Equals(value) bool
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

//TODO: List
// type val_list[T value] struct {
// 	Size int
// 	Data []T
// }

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

/**
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
