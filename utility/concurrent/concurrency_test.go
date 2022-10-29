package concurrent_test

import (
	"sync"
	"testing"

	"github.com/xeronith/diamante/utility/concurrent"
)

func BenchmarkLock(benchmark *testing.B) {
	mutex := sync.Mutex{}
	i := 0
	for n := 0; n < benchmark.N; n++ {
		mutex.Lock()
		i++
		mutex.Unlock()
	}
}

func BenchmarkChannel(benchmark *testing.B) {
	channel := make(chan bool, 1)
	i := 0
	for n := 0; n < benchmark.N; n++ {
		channel <- true
		i++
		<-channel
	}
}

func BenchmarkRaw(benchmark *testing.B) {
	i := 0
	for n := 0; n < benchmark.N; n++ {
		i++
	}
}

func TestFlag(test *testing.T) {
	flag := concurrent.NewFlag()

	flag.Set()

	if flag.IsSet() {
		flag.Clear()
	}

	if flag.IsSet() {
		test.Fail()
	}
}

func BenchmarkFlag(benchmark *testing.B) {
	flag := concurrent.NewFlag()

	for i := 0; i < benchmark.N; i++ {
		flag.Set()
		flag.Clear()
	}

	if flag.IsSet() {
		benchmark.Fail()
	}
}
