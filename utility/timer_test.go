package utility_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xeronith/diamante/utility"
)

func Test_Time(test *testing.T) {
	start := time.Now()
	total := 1000000
	wg := sync.WaitGroup{}
	wg.Add(total)

	utility.Repeat(total, func(i int) {
		go func() {
			defer wg.Done()
			if time.Since(start) < 0 {
				test.Fail()
			}
		}()
	})

	wg.Wait()
}

func Benchmark_Time(benchmark *testing.B) {
	start := time.Now()
	for i := 0; i < benchmark.N; i++ {
		if time.Since(start) < 0 {
			benchmark.Fail()
		}
	}
}

func Test_UTC(test *testing.T) {
	t := time.Now()
	fmt.Println(t.Unix())
	fmt.Println(t.UTC().Unix())
}
