package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/pcen/ape/ape"
)

func errWrongTestOutput(t test, got string) error {
	return fmt.Errorf("test %v: expected output %v but got %v", t.file, t.expect, got)
}

type test struct {
	file   string
	expect string
}

func (t test) path() string {
	return filepath.Join("./tests", t.file)
}

var (
	tests = []test{
		{
			"euler1.ape",
			"233168",
		},
		{
			"euler2.ape",
			"4613732",
		},
		{
			"gcd.ape",
			"9",
		},
		{
			"list.ape",
			"210",
		},
		{
			"switches.ape",
			"3\n2\n1",
		},
	}
)

func run(t test) error {
	path := t.path()
	ape.Ape(ape.ApeOpts{Src: path})
	b, err := exec.Command("./bin").Output()
	if err != nil {
		return err
	}
	output := string(bytes.TrimSpace(b))
	if output == t.expect {
		return nil
	}
	return errWrongTestOutput(t, output)
}

func main() {
	failures := 0
	var results []string
	for _, test := range tests {
		err := run(test)
		if err != nil {
			results = append(results, err.Error())
			failures++
		} else {
			results = append(results, fmt.Sprintf("test %v passed", test.file))
		}
	}
	fmt.Println("\n--- test results ---")
	for _, r := range results {
		fmt.Println(r)
	}

	fmt.Println("\n--- test summary ---")
	if failures > 0 {
		fmt.Printf("failed %v / %v tests\n", failures, len(tests))
	} else {
		fmt.Printf("passed all %v tests\n", len(tests))
	}
}
