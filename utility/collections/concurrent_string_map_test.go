package collections_test

import (
	"math/rand"
	"strconv"
	"testing"

	. "github.com/xeronith/diamante/utility/collections"
)

func TestConcurrentStringMap(test *testing.T) {

	key := "TEST-KEY"
	total := 100000
	stringMap := NewConcurrentStringMap()

	for i := 1; i <= total; i++ {
		stringMap.Put(key, i)
	}

	if val, _ := stringMap.Get(key); val != total {
		test.Fail()
	}

	if val, _ := stringMap.Get("INVALID-KEY"); val != nil {
		test.Fail()
	}

	for i := 1; i <= total; i++ {
		stringMap.Put(key, i)
	}

	stringMap.Clear()
	for i := 1; i <= total; i++ {
		stringMap.Put("KEY-"+strconv.Itoa(i), i)
	}

	if stringMap.GetSize() != total {
		test.Fail()
	}

	index := rand.Intn(total)
	if !stringMap.Contains("KEY-" + strconv.Itoa(index)) {
		test.Fail()
	}

	for i := 1; i <= total; i++ {
		val, _ := stringMap.Get("KEY-" + strconv.Itoa(i))
		if val.(int) != i {
			test.FailNow()
		}
	}
}
