package ape

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pcen/ape/ape/ast"
	"github.com/pcen/ape/ape/c"
	"github.com/pcen/ape/ape/types"
)

const OutputDir = "out"

func utilWriteCode(path string, sb *strings.Builder) (string, error) {
	base := filepath.Base(path)
	suffix := filepath.Ext(base)
	base = strings.TrimSuffix(base, suffix)

	if err := os.Mkdir(OutputDir, os.ModePerm); err != nil {
		fmt.Printf("error creating output directory %v: %v\n", OutputDir, err.Error())
		os.Exit(1)
	}

	output := fmt.Sprintf("./%v/%v.i", OutputDir, base)
	return output, os.WriteFile(output, []byte(sb.String()), 0664)
}

func utilCompile(path string) (string, error) {
	fmt.Println("lexing...")
	lexStart := time.Now()
	lexer := NewLexer()
	tokens := lexer.LexFile(path)
	lexDur := time.Since(lexStart)
	fmt.Println(tokens)

	fmt.Println("parsing...")
	parseStart := time.Now()
	parser := NewParser(tokens)
	file := parser.File()
	ast.PrettyPrint(file.Ast)
	parseDur := time.Since(parseStart)
	if errs, hasErrs := parser.Errors(); hasErrs {
		return "", fmt.Errorf("parser error(s): %v", errs)
	}

	checker := types.NewChecker(file)
	env := checker.Check()

	fmt.Println("generating code...")
	genStart := time.Now()
	code := c.GenerateCode(file.Ast, env)
	genDur := time.Since(genStart)

	fmt.Println("\ntime summary:")
	fmt.Printf("lex: %v\nparse: %v\ngen: %v\n", lexDur.Microseconds(), parseDur.Microseconds(), genDur.Microseconds())

	return utilWriteCode(path, code.Code)
}

func EndToEndC(path string) {
	compiled, err := utilCompile(path)
	if err != nil {
		fmt.Printf("error compiling %v: %v\n", path, err.Error())
		os.Exit(1)
	}
	start := time.Now()

	gccStart := time.Now()
	_, err = exec.Command("gcc", compiled, "-o", "bin").CombinedOutput()
	if err != nil {
		fmt.Printf("error compiling: %v\n", err.Error())
	}
	gccDur := time.Since(gccStart)
	dur := time.Since(start)
	fmt.Printf("gcc: %v\n", gccDur.Microseconds())
	fmt.Printf("total (ms): %v\n", dur.Milliseconds())
}

type ApeOpts struct {
	Run bool
	Src string
	Out string
}

// CLI program interface
func Ape(opts ApeOpts) {
	if opts.Out == "" {
		opts.Out = "./bin"
	}

	tokens := NewLexer().LexFile(opts.Src)
	file := NewParser(tokens).File()
	env := types.NewChecker(file).Check()
	code := c.GenerateCode(file.Ast, env)
	compiled, _ := utilWriteCode(opts.Src, code.Code)
	exec.Command("gcc", compiled, "-o", opts.Out).Run()
}
