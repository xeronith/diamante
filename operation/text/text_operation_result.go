package text

import (
	"time"

	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
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

func CreateTextOperationResult(id uint64, status int32, _type uint64, payload string, apiVersion, serverVersion, clientVersion int32, duration time.Duration) ITextOperationResult {
	return &textOperationResult{
		container: protobuf.TextOperationResult{
			Id:            id,
			Status:        status,
			Type:          _type,
			Payload:       payload,
			ApiVersion:    apiVersion,
			ServerVersion: serverVersion,
			ClientVersion: clientVersion,
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

func (result *textOperationResult) Load(payload interface{}, serializer ITextSerializer) error {
	data, err := serializer.Serialize(payload)
	result.container.Payload = data
	return err
}

func (result *textOperationResult) SerializeWith(serializer ITextSerializer) (string, error) {
	return serializer.Serialize(result)
}
