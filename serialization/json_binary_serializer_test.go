package serialization_test

import (
	"testing"

	. "github.com/xeronith/diamante/model"
	. "github.com/xeronith/diamante/serialization"
)

func TestJsonBinarySerializer(test *testing.T) {
	var (
		data []byte
		err  error
	)

	input := CreateRandomSample()
	serializer := NewJsonBinarySerializer()
	data, err = serializer.Serialize(&input)
	if err != nil {
		test.Fatal(err)
	}

	output := NewSample()
	err = serializer.Deserialize(data, &output)
	if err != nil {
		test.Fatal(err)
	}

	if output.StringField != input.StringField ||
		output.Int32Field != input.Int32Field ||
		output.BoolField != input.BoolField {
		test.Fatal()
	}
}
