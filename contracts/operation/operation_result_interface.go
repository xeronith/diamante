package operation

import (
	"time"

	. "github.com/xeronith/diamante/contracts/system"
)

type IOperationResult interface {
	Id() uint64
	Status() int32
	Type() uint64
	Container() Pointer
	ServerVersion() int32
	ExecutionDuration() time.Duration
	Hash() string
}
