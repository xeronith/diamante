package serialization

import "errors"

var (
	ERROR_NON_POINTER_SERIALIZATION_FAILED   = errors.New("non_pointer_serialization_failed")
	ERROR_NON_POINTER_DESERIALIZATION_FAILED = errors.New("non_pointer_deserialization_failed")
	ERROR_PROTO_SERIALIZATION_FAILED         = errors.New("proto_serialization_failed")
	ERROR_PROTO_DESERIALIZATION_FAILED       = errors.New("proto_deserialization_failed")
	ERROR_JSON_SERIALIZATION_FAILED          = errors.New("json_serialization_failed")
	ERROR_JSON_DESERIALIZATION_FAILED        = errors.New("json_deserialization_failed")
	ERROR_XML_SERIALIZATION_FAILED           = errors.New("xml_serialization_failed")
	ERROR_XML_DESERIALIZATION_FAILED         = errors.New("xml_deserialization_failed")
)
