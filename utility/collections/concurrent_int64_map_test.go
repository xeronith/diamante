package collections_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/utility/collections"
)

func TestConcurrentInt64Map(test *testing.T) {
	key := int64(1234)
	total := int64(100000)
	intMap := NewConcurrentInt64Map()

	for i := int64(1); i <= total; i++ {
		intMap.Put(key, i)
	}

	if val, _ := intMap.Get(key); val != total {
		test.Fail()
	}

	if val, _ := intMap.Get(321); val != nil {
		test.Fail()
	}

	for i := int64(1); i <= total; i++ {
		intMap.Put(key, i)
	}

	intMap.Clear()
	for i := int64(1); i <= total; i++ {
		intMap.Put(i, i)
	}

	if int64(intMap.GetSize()) != total {
		test.Fail()
	}

	index := rand.Int63n(total)
	if !intMap.Contains(index) {
		test.Fail()
	}

	for i := int64(1); i <= total; i++ {
		val, _ := intMap.Get(i)
		if val.(int64) != i {
			test.FailNow()
		}
	}
}

func TestConcurrentOrderedInt64Map(test *testing.T) {
	_map := NewConcurrentOrderedInt64Map()

	_map.Put(8, "Item 1")
	_map.Put(2, "Item 2")
	_map.Put(3, "Item 3")
	_map.Put(6, "Item 4")
	_map.Put(3, "Item 5")

	_map.ForEach(func(key int64, value system.ISystemObject) {
		fmt.Println(key, value)
	})

	fmt.Println(_map.IndexOf("Item 5"))
	fmt.Println(_map.IndexOfKey(6))

	fmt.Println(_map.GetKeyAt(2))
}

func TestConcurrentSortedInt64Map(test *testing.T) {
	_map := NewConcurrentSortedInt64Map()

	_map.Put(8, "Item 1")
	_map.Put(2, "Item 2")
	_map.Put(3, "Item 3")
	_map.Put(6, "Item 4")
	_map.Put(3, "Item 5")

	_map.ForEach(func(key int64, value system.ISystemObject) {
		fmt.Println(key, value)
	})

	fmt.Println(_map.IndexOf("Item 5"))
	fmt.Println(_map.IndexOfKey(6))

	fmt.Println(_map.GetKeyAt(2))
}
