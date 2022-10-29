package operation

import (
	. "github.com/xeronith/diamante/contracts/serialization"
)

type ITextOperationResult interface {
	IOperationResult

	Payload() string
	Load(interface{}, ITextSerializer) error
	SerializeWith(ITextSerializer) (string, error)
}
