package operation

import (
	"time"

	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
)

type IOperationResult interface {
	Id() uint64
	Status() int32
	Type() uint64
	ContentType() string
	Container() Pointer
	ServerVersion() int32
	ExecutionDuration() time.Duration
	UpdateStat(bool, int64, int64) IOperationResult
	Signature() string
	Payload() []byte
	Load(interface{}, ISerializer) error
	Stat() (int64, int64)
}
