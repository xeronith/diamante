package reflection_test

import (
	"testing"

	"github.com/xeronith/diamante/utility/reflection"
)

func TestReflection(test *testing.T) {
	s := struct{}{}
	p := &struct{}{}

	if reflection.IsPointer(s) {
		test.FailNow()
	}

	if !reflection.IsPointer(p) {
		test.FailNow()
	}
}

func BenchmarkReflection(b *testing.B) {
	sample := &struct{}{}

	result := false
	for i := 0; i < b.N; i++ {
		if reflection.IsPointer(sample) {
			result = true
		}
	}

	if !result {
		b.FailNow()
	}
}
