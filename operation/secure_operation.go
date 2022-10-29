package operation

import (
	"sync"
	. "time"

	. "github.com/xeronith/diamante/contracts/security"
)

type SecureOperation struct {
	sync.Mutex
	role Role
}

func (operation *SecureOperation) ExecutionTimeLimits() (Duration, Duration, Duration) {
	return DEFAULT_TIME_LIMIT_WARNING, DEFAULT_TIME_LIMIT_ALERT, DEFAULT_TIME_LIMIT_CRITICAL
}

func (operation *SecureOperation) Role() Role {
	if operation.role == 0 {
		operation.role = USER
	}

	return operation.role
}

func (operation *SecureOperation) SetRole(role Role) {
	operation.role = role
}
