package serialization

import (
	"errors"

	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/reflection"
	"google.golang.org/protobuf/proto"
)

type protobufSerializer struct {
}

func NewProtobufSerializer() ISerializer {
	return &protobufSerializer{}
}

func (serializer *protobufSerializer) Serialize(object Pointer) ([]byte, error) {
	if !reflection.IsPointer(object) {
		return nil, errors.New("non pointer serialization failed")
	}

	if message, ok := object.(proto.Message); ok {
		return proto.Marshal(message)
	} else {
		return nil, errors.New("invalid proto message")
	}
}

func (serializer *protobufSerializer) Deserialize(data []byte, object Pointer) error {
	if data == nil {
		return nil
	}

	if !reflection.IsPointer(object) {
		return errors.New("non pointer deserialization failed")
	}

	if message, ok := object.(proto.Message); ok {
		return proto.Unmarshal(data, message)
	} else {
		return errors.New("invalid proto message")
	}
}
