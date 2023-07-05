package serialization

import (
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
		return nil, ERROR_NON_POINTER_SERIALIZATION_FAILED
	}

	if message, ok := object.(proto.Message); ok {
		return proto.Marshal(message)
	} else {
		return nil, ERROR_PROTO_SERIALIZATION_FAILED
	}
}

func (serializer *protobufSerializer) Deserialize(data []byte, object Pointer) error {
	if data == nil {
		return nil
	}

	if !reflection.IsPointer(object) {
		return ERROR_NON_POINTER_DESERIALIZATION_FAILED
	}

	if message, ok := object.(proto.Message); ok {
		return proto.Unmarshal(data, message)
	} else {
		return ERROR_PROTO_DESERIALIZATION_FAILED
	}
}
