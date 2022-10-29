package text

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/protobuf"
)

type textOperationRequest struct {
	container protobuf.TextOperationRequest
}

func NewTextOperationRequest() ITextOperationRequest {
	return &textOperationRequest{}
}

func CreateTextOperationRequest(id uint64, operation uint64, token string) ITextOperationRequest {
	return &textOperationRequest{
		container: protobuf.TextOperationRequest{
			Id:        id,
			Operation: operation,
			Token:     token,
		},
	}
}

func (request *textOperationRequest) Id() uint64 {
	return request.container.Id
}

func (request *textOperationRequest) Operation() uint64 {
	return request.container.Operation
}

func (request *textOperationRequest) Token() string {
	return request.container.Token
}

func (request *textOperationRequest) Payload() string {
	return request.container.Payload
}

func (request *textOperationRequest) ApiVersion() int32 {
	return request.container.ApiVersion
}

func (request *textOperationRequest) ClientVersion() int32 {
	return request.container.ClientVersion
}

func (request *textOperationRequest) ClientName() string {
	return request.container.ClientName
}

func (request *textOperationRequest) Container() Pointer {
	return &request.container
}

func (request *textOperationRequest) Load(payload interface{}, serializer ITextSerializer) error {
	data, err := serializer.Serialize(payload)
	request.container.Payload = data
	return err
}
