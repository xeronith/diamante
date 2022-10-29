package io_test

import (
	"fmt"
	"testing"
)

func TestIO(test *testing.T) {
}

func BenchmarkStd(benchmark *testing.B) {
	for n := 0; n < benchmark.N; n++ {
		fmt.Println("Writing to stdout")
	}
}
