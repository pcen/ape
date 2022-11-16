package tests

import (
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/token"
)

func lex(source string) []token.Token {
	return ape.NewLexer().LexString(source)
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
