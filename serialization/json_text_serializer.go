package serialization

import (
	"encoding/json"
	"errors"

	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/reflection"
)

type jsonTextSerializer struct {
}

func NewJsonTextSerializer() ITextSerializer {
	return &jsonTextSerializer{}
}

func (serializer *jsonTextSerializer) Serialize(object Pointer) (string, error) {
	if !reflection.IsPointer(object) {
		return "", errors.New("non pointer serialization failed")
	}

	data, err := json.Marshal(object)
	if err != nil {
		return "", errors.New("json serialization failed")
	} else {
		return string(data), err
	}
}

func (serializer *jsonTextSerializer) Deserialize(data string, object Pointer) error {
	if !reflection.IsPointer(object) {
		return errors.New("non pointer deserialization failed")
	}

	err := json.Unmarshal([]byte(data), object)
	if err != nil {
		return errors.New("json deserialization failed")
	} else {
		return nil
	}
}
