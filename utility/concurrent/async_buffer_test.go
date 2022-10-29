package concurrent_test

import (
	. "fmt"
	"sync/atomic"
	"testing"

	. "github.com/xeronith/diamante/utility/concurrent"
)

func TestAsyncBuffer_SubmitAndFlush(test *testing.T) {
	var total uint32 = 64759
	var counter uint32 = 0

	StartSingleThreadScheduledExecutor()

	buffer := NewAsyncBuffer(func(items []interface{}) {
		atomic.AddUint32(&counter, uint32(len(items)))
	})

	for i := uint32(1); i <= total; i++ {
		buffer.Submit(Sprintf("Item %d", i))
	}

	buffer.Flush()
	buffer.Wait()

	if atomic.LoadUint32(&counter) != total {
		test.Fail()
	}
}

func TestAsyncBuffer_AutoFlush(test *testing.T) {
	var total uint32 = 35
	var counter uint32 = 0

	StartSingleThreadScheduledExecutor()

	buffer := NewAsyncBuffer(func(items []interface{}) {
		atomic.AddUint32(&counter, uint32(len(items)))
	})

	for i := uint32(1); i <= total; i++ {
		buffer.Submit(Sprintf("Item %d", i))
	}

	// time.Sleep(DefaultAutoFlushTimeout + time.Second)

	buffer.Flush()
	buffer.Wait()

	if atomic.LoadUint32(&counter) != total {
		test.Fail()
	}
}
