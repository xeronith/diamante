package collections_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/xeronith/diamante/model"
	"github.com/xeronith/diamante/utility"
	. "github.com/xeronith/diamante/utility/collections"
)

func TestConcurrentSlice_Append(test *testing.T) {
	slice := NewConcurrentSlice()

	total := 10000
	utility.Repeat(total, func(index int) {
		slice.Append(model.CreateRandomSample())
	})

	if slice.GetSize() != total {
		test.Fail()
	}
}

func BenchmarkConcurrentSlice_New(benchmark *testing.B) {
	slice := NewConcurrentSlice()
	for n := 0; n < benchmark.N; n++ {
		slice.Append(model.CreateRandomSample())
	}

	if slice.GetSize() != benchmark.N {
		benchmark.Fail()
	}
}

func TestConcurrentSlice_Remove(test *testing.T) {
	slice := NewConcurrentSlice()

	n1 := rand.Intn(50000)
	n2 := rand.Intn(50000)
	n3 := rand.Intn(50000)

	s1 := model.CreateRandomSample()
	s2 := model.CreateRandomSample()
	s3 := model.CreateRandomSample()

	utility.Repeat(n1, func(index int) {
		slice.Append(nil)
	})

	slice.Append(s1)

	utility.Repeat(n2, func(index int) {
		slice.Append(nil)
	})

	slice.Append(s2)

	utility.Repeat(n3, func(index int) {
		slice.Append(nil)
	})

	slice.Append(s3)

	if slice.GetSize() != (n1 + n2 + n3 + 3) {
		test.Fail()
	}
}

func BenchmarkConcurrentSlice_Remove(benchmark *testing.B) {
	total := 1000000
	rand.Seed(time.Now().UnixNano())
	sample := rand.Intn(total)
	slice := NewConcurrentSlice()

	utility.Repeat(total, func(index int) {
		slice.Append(index)
	})

	for n := 0; n < benchmark.N; n++ {
		slice.Remove(sample)
	}

	if slice.GetSize() != (total - 1) {
		benchmark.Fail()
	}
}
