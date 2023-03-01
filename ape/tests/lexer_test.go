package tests

import (
	"fmt"
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/token"
)

func lex(source string) []token.Token {
	return ape.NewLexer().LexString(source)
}

func tokensEqual(t1, t2 []token.Token) bool {
	if t1 == nil {
		return t2 == nil
	}
	if len(t1) != len(t2) {
		return false
	}
	for i := range t1 {
		if t1[i].Kind != t2[i].Kind || t1[i].Lexeme != t2[i].Lexeme {
			return false
		}
	}
	return true
}

func TestInsertStatementSeparators(t *testing.T) {
	source := `
	a += 2
	b++
	return
	break
	return;
	`
	source2 := `
	a += 2;
	b++;
	return;
	break;
	return;
	`

	sepIndices := []int{3, 6, 8, 10}

	tokens := lex(source)
	tokens2 := lex(source2)

	fmt.Println(tokens)

	for _, i := range sepIndices {
		if tokens[i].Kind != token.Sep {
			t.Fatalf("token %v at index %v should be separator", tokens[i], i)
		}
	}

	if len(tokens) != len(tokens2) {
		t.Fatalf("different length token outputs")
	}
	for i := range tokens2 {
		if tokens[i].Kind != tokens2[i].Kind {
			t.Fatalf("missmatching tokens at index %v: expect %v got %v", i, tokens2[i], tokens[i])
		}
	}
}

func TestInsertAfterBrace(t *testing.T) {
	source := `
		class foo {

		}
		func bar() {

		}
	`
	source2 := `
		class foo {

		}
		func bar() {

		}
	`
	tokens := lex(source)
	tokens2 := lex(source2)
	fmt.Println(tokens)
	fmt.Println(tokens2)
	if !tokensEqual(tokens, tokens2) {
		t.Fatal("tokens are not equal")
	}
}

func TestSkipSeizeKeywords(t *testing.T) {
	source := `
		func bar() {
			SKIP {
				x := 1 * 2
				REVERSE x
			} SEIZE (2) {
				println("this is two")
			}
		}
	`
	tokens := lex(source)
	fmt.Println(tokens)
}
