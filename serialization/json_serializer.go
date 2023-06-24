package serialization

import (
	"encoding/json"
	"errors"

	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/reflection"
)

type jsonSerializer struct {
}

func NewJsonSerializer() ISerializer {
	return &jsonSerializer{}
}

func (serializer *jsonSerializer) Serialize(object Pointer) ([]byte, error) {
	if !reflection.IsPointer(object) {
		return nil, errors.New("non pointer serialization failed")
	}

	if data, err := json.Marshal(object); err != nil {
		return nil, errors.New("json serialization failed")
	} else {
		return data, err
	}
}

func (serializer *jsonSerializer) Deserialize(data []byte, object Pointer) error {
	if data == nil {
		return nil
	}

	if !reflection.IsPointer(object) {
		return errors.New("non pointer deserialization failed")
	}

	if err := json.Unmarshal(data, object); err != nil {
		return errors.New("json deserialization failed")
	} else {
		return nil
	}
}
