package collections_test

import (
	"math/rand"
	"testing"

	. "github.com/xeronith/diamante/utility/collections"
)

func TestConcurrentIntMap(test *testing.T) {

	key := 1234
	total := 100000
	intMap := NewConcurrentIntMap()

	for i := 1; i <= total; i++ {
		intMap.Put(key, i)
	}

	if val := intMap.Get(key); val != total {
		test.Fail()
	}

	if val := intMap.Get(321); val != nil {
		test.Fail()
	}

	for i := 1; i <= total; i++ {
		intMap.Put(key, i)
	}

	intMap.Clear()
	for i := 1; i <= total; i++ {
		intMap.Put(i, i)
	}

	if intMap.GetSize() != total {
		test.Fail()
	}

	index := rand.Intn(total)
	if !intMap.Contains(index) {
		test.Fail()
	}

	for i := 1; i <= total; i++ {
		val := intMap.Get(i)
		if val.(int) != i {
			test.FailNow()
		}
	}
}
