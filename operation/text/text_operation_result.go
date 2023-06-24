package text

import (
	"time"

	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/protobuf"
)

type textOperationResult struct {
	container protobuf.TextOperationResult
	duration  time.Duration
}

func NewTextOperationResult() ITextOperationResult {
	return &textOperationResult{}
}

func CreateTextOperationResult(
	id ID,
	status int32,
	resultType uint64,
	payload string,
	pipelineInfo IPipeline,
	duration time.Duration,
	hash string,
) ITextOperationResult {
	return &textOperationResult{
		container: protobuf.TextOperationResult{
			Id:            id,
			Status:        status,
			Type:          resultType,
			Payload:       payload,
			ApiVersion:    pipelineInfo.ApiVersion(),
			ServerVersion: pipelineInfo.ServerVersion(),
			ClientVersion: pipelineInfo.ClientVersion(),
			Hash:          hash,
		},
		duration: duration,
	}
}

func (result *textOperationResult) Id() uint64 {
	return result.container.Id
}

func (result *textOperationResult) Status() int32 {
	return result.container.Status
}

func (result *textOperationResult) Type() uint64 {
	return result.container.Type
}

func (result *textOperationResult) Payload() string {
	return result.container.Payload
}

func (result *textOperationResult) Container() Pointer {
	return &result.container
}

func (result *textOperationResult) ServerVersion() int32 {
	return result.container.ServerVersion
}

func (result *textOperationResult) ExecutionDuration() time.Duration {
	return result.duration
}

func (result *textOperationResult) Hash() string {
	return result.container.Hash
}

func (result *textOperationResult) Load(payload interface{}, serializer ITextSerializer) error {
	data, err := serializer.Serialize(payload)
	result.container.Payload = data
	return err
}

func (result *textOperationResult) SerializeWith(serializer ITextSerializer) (string, error) {
	return serializer.Serialize(result)
}
