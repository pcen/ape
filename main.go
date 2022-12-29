package main

import (
	"fmt"
	"os"

	"github.com/pcen/ape/ape"
)

/*
 compiles demo script
*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("supply name of ape script as argument")
		os.Exit(1)
	}
	script := os.Args[1]
	ape.EndToEndC(script)
	// compiled, err := compile(script)
	// if err != nil {
	// 	fmt.Printf("error compiling %v: %v\n", script, err.Error())
	// 	os.Exit(1)
	// }

	// _, err = exec.Command("gcc", compiled, "-o", "bin").CombinedOutput()
	// if err != nil {
	// 	fmt.Printf("error compiling: %v\n", err.Error())
	// }
}
