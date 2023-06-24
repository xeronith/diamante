package operation

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/protobuf"
)

type operationRequest struct {
	container protobuf.OperationRequest
}

func NewOperationRequest() IOperationRequest {
	return &operationRequest{}
}

func CreateOperationRequest(
	id uint64,
	operation uint64,
	clientName string,
	clientVersion,
	apiVersion int32,
	token string,
	payload []byte,
) IOperationRequest {
	return &operationRequest{
		container: protobuf.OperationRequest{
			Id:            id,
			ApiVersion:    apiVersion,
			ClientName:    clientName,
			ClientVersion: clientVersion,
			Operation:     operation,
			Token:         token,
			Payload:       payload,
		},
	}
}

func (request *operationRequest) Id() uint64 {
	return request.container.Id
}

func (request *operationRequest) Operation() uint64 {
	return request.container.Operation
}

func (request *operationRequest) Token() string {
	return request.container.Token
}

func (request *operationRequest) Payload() []byte {
	return request.container.Payload
}

func (request *operationRequest) ApiVersion() int32 {
	return request.container.ApiVersion
}

func (request *operationRequest) ClientVersion() int32 {
	return request.container.ClientVersion
}

func (request *operationRequest) ClientName() string {
	return request.container.ClientName
}

func (request *operationRequest) Container() Pointer {
	return &request.container
}

func (request *operationRequest) Load(payload interface{}, serializer ISerializer) error {
	data, err := serializer.Serialize(payload)
	request.container.Payload = data
	return err
}
