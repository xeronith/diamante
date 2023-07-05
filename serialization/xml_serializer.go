package serialization

import (
	"encoding/xml"

	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/reflection"
)

type xmlSerializer struct {
}

func NewXmlSerializer() ISerializer {
	return &xmlSerializer{}
}

func (serializer *xmlSerializer) Serialize(object Pointer) ([]byte, error) {
	if !reflection.IsPointer(object) {
		return nil, ERROR_NON_POINTER_SERIALIZATION_FAILED
	}

	if data, err := xml.Marshal(object); err != nil {
		return nil, ERROR_XML_SERIALIZATION_FAILED
	} else {
		return data, err
	}
}

func (serializer *xmlSerializer) Deserialize(data []byte, object Pointer) error {
	if data == nil {
		return nil
	}

	if !reflection.IsPointer(object) {
		return ERROR_NON_POINTER_DESERIALIZATION_FAILED
	}

	if err := xml.Unmarshal(data, object); err != nil {
		return ERROR_XML_DESERIALIZATION_FAILED
	} else {
		return nil
	}
}
