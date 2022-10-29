package operation

import . "github.com/xeronith/diamante/contracts/serialization"

type IBinaryOperationRequest interface {
	IOperationRequest

	Payload() []byte
	Load(interface{}, IBinarySerializer) error
}
