package serialization_test

import (
	"testing"

	"github.com/xeronith/diamante/contracts/serialization"
	"github.com/xeronith/diamante/model"
	. "github.com/xeronith/diamante/serialization"
)

func TestCompoundObjectSerialization(test *testing.T) {

	binarySerializers := []serialization.IBinarySerializer{
		NewProtobufSerializer(),
		NewJsonBinarySerializer(),
	}

	for _, serializer := range binarySerializers {
		input := model.CreateRandomCompoundObject()
		data, err := serializer.Serialize(input)
		if err != nil {
			test.Fail()
		}

		output := model.NewCompoundObject()
		err = serializer.Deserialize(data, output)
		if err != nil {
			test.Fail()
		}

		if input.Sample.StringField != output.Sample.StringField {
			test.Fail()
		}
	}

	textSerializers := []serialization.ITextSerializer{
		NewJsonTextSerializer(),
	}

	for _, serializer := range textSerializers {
		input := model.CreateRandomCompoundObject()
		data, err := serializer.Serialize(input)
		if err != nil {
			test.Fail()
		}

		output := model.NewCompoundObject()
		err = serializer.Deserialize(data, output)
		if err != nil {
			test.Fail()
		}

		if input.Sample.StringField != output.Sample.StringField {
			test.Fail()
		}
	}
}
