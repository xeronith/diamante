package binary

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/protobuf"
)

type binaryOperationRequest struct {
	container protobuf.BinaryOperationRequest
}

func NewBinaryOperationRequest() IBinaryOperationRequest {
	return &binaryOperationRequest{}
}

func CreateBinaryOperationRequest(id uint64, operation uint64, clientName string, clientVersion, apiVersion int32, token string, payload []byte) IBinaryOperationRequest {
	return &binaryOperationRequest{
		container: protobuf.BinaryOperationRequest{
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

func (request *binaryOperationRequest) Id() uint64 {
	return request.container.Id
}

func (request *binaryOperationRequest) Operation() uint64 {
	return request.container.Operation
}

func (request *binaryOperationRequest) Token() string {
	return request.container.Token
}

func (request *binaryOperationRequest) Payload() []byte {
	return request.container.Payload
}

func (request *binaryOperationRequest) ApiVersion() int32 {
	return request.container.ApiVersion
}

func (request *binaryOperationRequest) ClientVersion() int32 {
	return request.container.ClientVersion
}

func (request *binaryOperationRequest) ClientName() string {
	return request.container.ClientName
}

func (request *binaryOperationRequest) Container() Pointer {
	return &request.container
}

func (request *binaryOperationRequest) Load(payload interface{}, serializer IBinarySerializer) error {
	data, err := serializer.Serialize(payload)
	request.container.Payload = data
	return err
}
