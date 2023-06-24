package operation

import (
	"sync"
	. "time"

	. "github.com/xeronith/diamante/contracts/security"
)

type Operation struct {
	sync.Mutex
	role         Role
	activeRunner int
}

func (operation *Operation) ExecutionTimeLimits() (Duration, Duration, Duration) {
	return DEFAULT_TIME_LIMIT_WARNING, DEFAULT_TIME_LIMIT_ALERT, DEFAULT_TIME_LIMIT_CRITICAL
}

func (operation *Operation) Role() Role {
	if operation.role == 0 {
		operation.role = ANONYMOUS
	}

	return operation.role
}

func (operation *Operation) SetRole(role Role) {
	operation.role = role
}

func (operation *Operation) ActiveRunner() int {
	if operation.activeRunner < 0 {
		return 0
	}

	return operation.activeRunner
}

func (operation *Operation) SetActiveRunner(value int) {
	operation.activeRunner = value
}
