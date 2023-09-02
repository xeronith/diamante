package operation

import (
	"time"

	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/protobuf"
)

type operationResult struct {
	container   protobuf.OperationResult
	contentType string
	duration    time.Duration
}

func NewOperationResult() IOperationResult {
	return &operationResult{}
}

func CreateOperationResult(
	id ID,
	status int32,
	resultType uint64,
	payload []byte,
	pipeline IPipeline,
	duration time.Duration,
) IOperationResult {
	return &operationResult{
		container: protobuf.OperationResult{
			Id:            id,
			Status:        status,
			Type:          resultType,
			Payload:       payload,
			ApiVersion:    pipeline.ApiVersion(),
			ServerVersion: pipeline.ServerVersion(),
			Hash:          pipeline.Sign(payload),
		},
		contentType: pipeline.ContentType(),
		duration:    duration,
	}
}

func (result *operationResult) Id() uint64 {
	return result.container.Id
}

func (result *operationResult) Status() int32 {
	return result.container.Status
}

func (result *operationResult) ContentType() string {
	return result.contentType
}

func (result *operationResult) Type() uint64 {
	return result.container.Type
}

func (result *operationResult) Payload() []byte {
	return result.container.Payload
}

func (result *operationResult) Container() Pointer {
	return &result.container
}

func (result *operationResult) ServerVersion() int32 {
	return result.container.ServerVersion
}

func (result *operationResult) ExecutionDuration() time.Duration {
	return result.duration
}

func (result *operationResult) ResetDuration() IOperationResult {
	result.duration = 0
	return result
}

func (result *operationResult) Signature() string {
	return result.container.Hash
}

func (result *operationResult) Load(payload interface{}, serializer ISerializer) error {
	data, err := serializer.Serialize(payload)
	result.container.Payload = data
	return err
}
