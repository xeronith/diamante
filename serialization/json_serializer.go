package serialization

import (
	"encoding/json"

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
		return nil, ERROR_NON_POINTER_SERIALIZATION_FAILED
	}

	if data, err := json.Marshal(object); err != nil {
		return nil, ERROR_JSON_SERIALIZATION_FAILED
	} else {
		return data, err
	}
}

func (serializer *jsonSerializer) Deserialize(data []byte, object Pointer) error {
	if data == nil {
		return nil
	}

	if !reflection.IsPointer(object) {
		return ERROR_NON_POINTER_DESERIALIZATION_FAILED
	}

	if err := json.Unmarshal(data, object); err != nil {
		return ERROR_JSON_DESERIALIZATION_FAILED
	} else {
		return nil
	}
}
