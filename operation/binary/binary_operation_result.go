package binary

import (
	"time"

	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/protobuf"
)

type binaryOperationResult struct {
	container protobuf.BinaryOperationResult
	duration  time.Duration
}

func NewBinaryOperationResult() IBinaryOperationResult {
	return &binaryOperationResult{}
}

func CreateBinaryOperationResult(id uint64, status int32, _type uint64, payload []byte, apiVersion, serverVersion, clientVersion int32, duration time.Duration, hash string) IBinaryOperationResult {
	return &binaryOperationResult{
		container: protobuf.BinaryOperationResult{
			Id:            id,
			Status:        status,
			Type:          _type,
			Payload:       payload,
			ApiVersion:    apiVersion,
			ServerVersion: serverVersion,
			ClientVersion: clientVersion,
			Hash:          hash,
		},
		duration: duration,
	}
}

func (result *binaryOperationResult) Id() uint64 {
	return result.container.Id
}

func (result *binaryOperationResult) Status() int32 {
	return result.container.Status
}

func (result *binaryOperationResult) Type() uint64 {
	return result.container.Type
}

func (result *binaryOperationResult) Payload() []byte {
	return result.container.Payload
}

func (result *binaryOperationResult) Container() Pointer {
	return &result.container
}

func (result *binaryOperationResult) ServerVersion() int32 {
	return result.container.ServerVersion
}

func (result *binaryOperationResult) ExecutionDuration() time.Duration {
	return result.duration
}

func (result *binaryOperationResult) Hash() string {
	return result.container.Hash
}

func (result *binaryOperationResult) Load(payload interface{}, serializer IBinarySerializer) error {
	data, err := serializer.Serialize(payload)
	result.container.Payload = data
	return err
}
