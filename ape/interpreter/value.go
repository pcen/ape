package interpreter

type value interface{}

type val_bool struct {
	Value bool
}

type val_str struct {
	Value string
}

type val_int struct {
	Value int
}

type val_rational struct {
	Value float64
}
