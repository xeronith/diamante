package operation

import . "github.com/xeronith/diamante/contracts/serialization"

type ITextOperationRequest interface {
	IOperationRequest

	Payload() string
	Load(interface{}, ITextSerializer) error
}
