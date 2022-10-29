package serialization

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/reflection"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsonBinarySerializer struct {
}

func NewJsonBinarySerializer() IBinarySerializer {
	return &jsonBinarySerializer{}
}

func (serializer *jsonBinarySerializer) Serialize(object Pointer) ([]byte, error) {
	if !reflection.IsPointer(object) {
		return nil, errors.New("non pointer serialization failed")
	}

	data, err := json.Marshal(object)
	if err != nil {
		return nil, errors.New("json serialization failed")
	} else {
		return data, err
	}
}

func (serializer *jsonBinarySerializer) Deserialize(data []byte, object Pointer) error {
	if data == nil {
		return nil
	}

	if !reflection.IsPointer(object) {
		return errors.New("non pointer deserialization failed")
	}

	err := json.Unmarshal(data, object)
	if err != nil {
		return errors.New("json deserialization failed")
	} else {
		return nil
	}
}
