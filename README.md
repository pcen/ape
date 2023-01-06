# ape programming language

## compile a program
- run `go run main.go ./tests/<source file>.ape`
- check the generated c code in `./out/<source file>.i`
- run the binary `./bin`

## layout
[ape/lexer.go](./ape/lexer.go) - tokenizes files into the tokens defined in [ape/token](./ape/token) \
[ape/parser.go](./ape/parser.go) - recursive descent parser that parses the grammar described in [grammar.txt](./grammar.txt), and generates an ast described in [ape/ast](./ape/ast) \
[ape/types](./ape/types) - type checks ast and produces a 'type environment' used by code generator \
[ape/c](./ape/c) - code generation that compiles an ast tree to c code

## pipeline
The directories above are listed in the same order they are used in the compilation pipeline: \
lex file into tokens -> parse tokens into ast -> type check ast -> generate code
<br>
<br>
[ape/util.go](./ape/util.go) implements `func EndToEndC(path string)` which provides an example of how each module is called in order to compile a file. `EndToEndC` compiles the ape source file `path` to the binary `./bin`, and outputs the generated c code in `./out` (this is an empty directory not checked in to git). `EndToEndC` invokes `gcc` directly to compile the generated c code, so this <b>will not work</b> unless `gcc` is installed (but other c compilers should be able to compile the generated c code).

## testing
[cmd/gen/main.go](./cmd/gen/main.go) - Generates input for the parser based on the grammar described in [grammar.txt](./grammar.txt), using [dist.txt](./dist.txt) to control how grammar rules are expanded. These test programs are valid grammar derivations, but are not validly typed. \
[tests](./tests) - Contains handwritten programs to test the language.

## scratch

parse statements as string chunk first
- idea is that expressions that cannot be created within a single c expression are only needed before the statement they are used in
- allows constructing temporary results ahead of statement
- write supplementary code + 'actual' line to source
