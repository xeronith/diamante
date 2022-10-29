package collections_test

import (
	"testing"
	. "time"

	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/utility/collections"
)

func TestSnapshot(test *testing.T) {
	collection := NewConcurrentOrderedInt64Map()

	total := 100000
	for i := 1; i <= total; i++ {
		collection.Put(int64(i), i)
	}

	temp := make([]ISystemObject, 0)
	collection.ForEachValue(func(systemObject ISystemObject) {
		temp = append(temp, systemObject)
	})

	if len(temp) != total {
		test.Fail()
	}
}

func TestConcurrency(test *testing.T) {
	collection := NewConcurrentOrderedInt64Map()

	total := 10
	for i := 1; i <= total; i++ {
		collection.Put(int64(i), i)
	}

	go func() {
		Sleep(Millisecond * 250)
		collection.Put(0, 0)
	}()

	collection.ForEachValue(func(ISystemObject) {
		Sleep(Millisecond * 100)
		_, _ = collection.Get(1)
	})

	if collection.GetSize() != total+1 {
		test.Fail()
	}
}
