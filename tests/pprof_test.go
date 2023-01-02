package tests

import (
	"runtime"
	"testing"

	"github.com/pcen/ape/ape"
)

// go test -bench -v -cpuprofile cpu.prof tests/pprof_test.go
// go tool pprof -http :8080 cpu.prof

func TestProfile(t *testing.T) {
	runtime.SetCPUProfileRate(2000)
	for i := 0; i < 10; i++ {
		ape.EndToEndC("./gcd.ape")
	}
}
