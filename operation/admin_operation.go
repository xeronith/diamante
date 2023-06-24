package operation

import (
	"sync"
	. "time"

	. "github.com/xeronith/diamante/contracts/security"
)

type AdminOperation struct {
	sync.Mutex
	role         Role
	activeRunner uint
}

func (operation *AdminOperation) ExecutionTimeLimits() (Duration, Duration, Duration) {
	return DEFAULT_TIME_LIMIT_WARNING, DEFAULT_TIME_LIMIT_ALERT, DEFAULT_TIME_LIMIT_CRITICAL
}

func (operation *AdminOperation) Role() Role {
	if operation.role == 0 {
		operation.role = ADMINISTRATOR
	}

	return operation.role
}

func (operation *AdminOperation) SetRole(role Role) {
	operation.role = role
}

func (operation *AdminOperation) ActiveRunner() uint {
	return operation.activeRunner
}

func (operation *AdminOperation) SetActiveRunner(value uint) {
	operation.activeRunner = value
}
