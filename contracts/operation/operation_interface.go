package operation

import (
	"sync"
	. "time"

	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/service"
	. "github.com/xeronith/diamante/contracts/system"
)

type (
	Generator func() interface{}

	IOperation interface {
		sync.Locker
		Tag() string
		Id() (ID, ID)
		Role() Role
		SetRole(Role)
		ActiveRunner() uint
		SetActiveRunner(uint)
		InputContainer() Pointer
		OutputContainer() Pointer
		Execute(IContext, Pointer) (Pointer, error)
		ExecutionTimeLimits() (Duration, Duration, Duration)
		IsCacheable() bool
	}

	IOperationFactory interface {
		Operations() []IOperation
	}
)
