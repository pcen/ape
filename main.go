package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	s := `var a int = 5`
	r := strings.NewReader(s)
	br := bufio.NewReader(r)
	for {
		b, err := br.ReadByte()
		fmt.Printf("%c\n", b)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
