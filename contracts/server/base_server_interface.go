package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/localization"
	. "github.com/xeronith/diamante/contracts/logging"
	"github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/scheduling"
	. "github.com/xeronith/diamante/contracts/system"
)

type Opcodes map[uint64]string

type IBaseServer interface {
	Opcodes() Opcodes
	ActorsCount() int
	IncrementActorsCount(IActor)
	Scheduler() IScheduler
	Logger() ILogger
	Localizer() ILocalizer
	Push(IActor, messaging.IPushMessage) error
	Broadcast(uint64, Pointer) error
	BroadcastSpecific(uint64, map[string]Pointer) error
}
