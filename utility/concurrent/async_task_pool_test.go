package concurrent_test

import (
	"sync/atomic"
	"testing"

	. "github.com/xeronith/diamante/utility/concurrent"
)

func TestAsyncTaskPool(test *testing.T) {
	var total uint32 = 1000
	var counter uint32
	pool := NewAsyncTaskPool()

	for i := uint32(1); i <= total; i++ {
		func() {
			// index := i
			pool.Submit(
				func() {
					// Println("Function ", index)
					atomic.AddUint32(&counter, 1)
				})
		}()
	}

	pool.Run().Join()

	count := atomic.LoadUint32(&counter)
	if count != total {
		test.Fail()
	}
}
