package c

import "fmt"

func declareVector(name, ctype string) string {
	format := `typedef struct %v {
	%v *data;
	int length;
	int capacity;
} %v;
`
	return fmt.Sprintf(format, name, ctype, name)
}

func implementVector(name, ctype string) {

}
