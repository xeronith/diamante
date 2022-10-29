package utility_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xeronith/diamante/utility"
)

func Test_TimedTokenGenerator(test *testing.T) {
	duration := time.Second
	generator := utility.NewTimedTokenGenerator(duration)

	token := generator.Generate()
	if !generator.IsValid(token) {
		test.Fail()
	}

	time.Sleep(duration + time.Millisecond)

	if generator.IsValid(token) {
		test.Fail()
	}

	if generator.Size() != 0 {
		test.Fail()
	}
}

func Test_GenerateHash(test *testing.T) {
	hash := utility.GenerateHash("SAMPLE_VALUE", "SAMPLE_SALT")
	fmt.Println(hash)
}

func Test_GenerateConfirmationCode(test *testing.T) {
	code := utility.GenerateConfirmationCode()
	fmt.Println(code)
}

func Benchmark_TimedTokenGenerator(benchmark *testing.B) {
	duration := 60 * time.Second
	generator := utility.NewTimedTokenGenerator(duration)

	tokens := make([]string, 0)

	for i := 0; i < benchmark.N; i++ {
		tokens = append(tokens, generator.Generate())
	}

	for i := 0; i < benchmark.N; i++ {
		if !generator.IsValid(tokens[i]) {
			benchmark.Fail()
		}
	}
}

type object struct {
	stringField string
	intField    int
	boolField   bool
}

var objectPool = sync.Pool{New: func() interface{} {
	return &object{}
}}

func Benchmark_ObjectCreation(benchmark *testing.B) {
	objects := make([]*object, 0)
	for i := 0; i < benchmark.N; i++ {
		object := &object{
			stringField: "String Field",
			intField:    89234234,
			boolField:   true,
		}

		objects = append(objects, object)
	}

	if len(objects) < 1 {
		benchmark.Fail()
	}
}

func Benchmark_ObjectPool(benchmark *testing.B) {
	objects := make([]*object, 0)
	for i := 0; i < benchmark.N; i++ {
		object := objectPool.Get().(*object)
		object.stringField = "String Field"
		object.intField = 89234234
		object.boolField = true

		objects = append(objects, object)
	}

	if len(objects) < 1 {
		benchmark.Fail()
	}
}
