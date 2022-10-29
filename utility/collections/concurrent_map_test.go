package collections_test

import (
	"testing"

	. "github.com/xeronith/diamante/utility/collections"
)

func TestConcurrentMap(test *testing.T) {
	key := 1234
	total := 100000

	cm := NewConcurrentMap[int, int]()

	for i := 1; i <= total; i++ {
		cm.Put(key, i)
	}

	if val := cm.Get(key); val != total {
		test.Fail()
	}

	if cm.Contains(321) {
		test.Fail()
	}

	for i := 1; i <= total; i++ {
		cm.Put(key, i)
	}

	cm.Clear()
	for i := 1; i <= total; i++ {
		cm.Put(i, i)
	}

	if cm.GetSize() != total {
		test.Fail()
	}

	for i := 1; i <= total; i++ {
		val := cm.Get(i)

		if val != i {
			test.FailNow()
		}
	}
}
