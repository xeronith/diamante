package collections_test

import (
	"fmt"
	"testing"

	"github.com/xeronith/diamante/utility/collections"
)

func TestConcurrentPointerMap_RemoveFirst(test *testing.T) {
	_map := collections.NewConcurrentPointerMap()
	_map.Put("A", "B")
	key, value := _map.RemoveFirst()

	if key == nil && value == nil {
		test.Fail()
	}

	fmt.Println(key, value)
}
