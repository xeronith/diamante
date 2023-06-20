package service

import (
	"time"

	"github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/scheduling"
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/settings"
	. "github.com/xeronith/diamante/contracts/system"
)

type IContext interface {
	Configuration() IConfiguration
	SetSecureCookie(string, string)
	GetSecureCookie(string) string
	Token() string
	Identity() Identity
	RequestId() uint64
	ApiVersion() int32
	ClientVersion() int32
	ClientLatestVersion() int32
	ServerVersion() int32
	ClientName() string
	ResultType() ID
	SetResultType(ID)
	IncrementActorsCount()
	ActorsCount() int
	Scheduler() IScheduler
	SetTimeout(func(), time.Duration) string
	SetInterval(func(), time.Duration) string
	CancelSchedule(string)
	Logger() ILogger
	SecurityHandler() ISecurityHandler
	Push(IPushMessage) error
	Broadcast(uint64, Pointer) error
	BroadcastSpecific(uint64, map[string]Pointer) error
	SMS(string, string) error
	Timestamp() time.Time
	IsStagingEnvironment() bool
	IsProductionEnvironment() bool
	ResultContainer() Pointer
	SubmitAnalyticsEvent(uint64, string, analytics.Fields)
	SubmitMeasurement(string, analytics.Tags, analytics.Fields)
	Async(func())
	Lock()
	Unlock()
	SystemCall([]string) error
}
