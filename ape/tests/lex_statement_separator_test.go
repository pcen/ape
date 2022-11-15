package tests

import (
	"testing"

	"github.com/pcen/ape/ape"
	"github.com/pcen/ape/ape/token"
)

func TestInsertStatementSeparators(t *testing.T) {
	source := `
	a += 2
	b++
	return
	break
	return;
	`
	sepIndices := []int{3, 6, 8, 10}

	tokens := ape.NewLexer().LexString(source)

	for _, i := range sepIndices {
		if tokens[i].Kind != token.Sep {
			t.Fatalf("token %v at index %v should be separator", tokens[i], i)
		}
	}
}
