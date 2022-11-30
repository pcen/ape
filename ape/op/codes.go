package op

type Code uint8

const (
	Nil Code = iota
	Set      // set identifier
	Get      // get identifier
	Add
	Subtract
	Multiply
	Divide
	Constant // load a literal value
	Print
)

func (c Code) String() string {
	return "OP_" + []string{
		Nil:      "NIL",
		Set:      "SET",
		Get:      "GET",
		Add:      "ADD",
		Subtract: "SUBTRACT",
		Multiply: "MULTIPLY",
		Divide:   "DIVIDE",
		Constant: "CONSTANT",
		Print:    "PRINT",
	}[c]
}
