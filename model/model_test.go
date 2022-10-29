package model_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/xeronith/diamante/model"
	"github.com/xeronith/diamante/protobuf"
)

func TestModel(test *testing.T) {}

func BenchmarkCreateRandomSample_Pool(benchmark *testing.B) {
	var pool = &sync.Pool{
		New: func() interface{} {
			return new(protobuf.Sample)
		},
	}

	for n := 0; n < benchmark.N; n++ {
		sample := pool.Get().(*protobuf.Sample)
		if sample.Int32Field < 0 {
			fmt.Println("Negative")
		}
	}
}

func BenchmarkCreateRandomSample_Init(benchmark *testing.B) {
	for n := 0; n < benchmark.N; n++ {
		sample := model.CreateRandomSample()
		if sample.Int32Field < 0 {
			fmt.Println("Negative")
		}
	}
}
