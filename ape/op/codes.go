package op

type Code uint8

const (
	Assign Code = iota + 1
	Add
	Subtract
	Multiply
	Divide
	Print
)
