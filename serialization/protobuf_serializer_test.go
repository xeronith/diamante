package serialization_test

import (
	"testing"

	. "github.com/xeronith/diamante/model"
	. "github.com/xeronith/diamante/serialization"
)

func TestProtobufSerializer(test *testing.T) {
	var (
		data []byte
		err  error
	)

	input := CreateRandomSample()
	serializer := NewProtobufSerializer()
	data, err = serializer.Serialize(input)
	if err != nil {
		test.Fatal(err)
	}

	output := NewSample()
	err = serializer.Deserialize(data, output)
	if err != nil {
		test.Fatal(err)
	}

	if output.StringField != input.StringField ||
		output.Int32Field != input.Int32Field ||
		output.BoolField != input.BoolField {
		test.Fail()
	}
}

func BenchmarkProtobufSerializer_Serialize(benchmark *testing.B) {
	serializer := NewProtobufSerializer()
	for n := 0; n < benchmark.N; n++ {
		input := CreateRandomSample()
		_, err := serializer.Serialize(input)
		if err != nil {
			benchmark.Fatal(err)
		}
	}
}

func BenchmarkProtobufSerializer_Deserialize(benchmark *testing.B) {
	serializer := NewProtobufSerializer()
	for n := 0; n < benchmark.N; n++ {

		data := []byte{
			10, 12, 83, 97, 109, 112,
			108, 101, 32, 57, 50, 53, 48,
			56, 16, 220, 210, 5, 24, 1,
		}

		output := NewSample()
		err := serializer.Deserialize(data, output)
		if err != nil {
			benchmark.Fatal(err)
		}
	}
}
