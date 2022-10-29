package operation

import . "github.com/xeronith/diamante/contracts/serialization"

type IBinaryOperationResult interface {
	IOperationResult

	Payload() []byte
	Load(interface{}, IBinarySerializer) error
}
