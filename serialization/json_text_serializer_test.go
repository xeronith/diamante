package serialization_test

import (
	"testing"

	. "github.com/xeronith/diamante/model"
	. "github.com/xeronith/diamante/serialization"
)

func TestJsonTextSerializer(test *testing.T) {
	var (
		data string
		err  error
	)

	input := CreateRandomSample()
	serializer := NewJsonTextSerializer()
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
