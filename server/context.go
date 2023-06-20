package server

import (
	"errors"
	"time"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/scheduling"
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/service"
	. "github.com/xeronith/diamante/contracts/settings"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/concurrent"
)

type context struct {
	resultType          ID
	server              *baseServer
	operation           IOperation
	actor               IActor
	analyticsProvider   IAnalyticsProvider
	requestId           uint64
	serverVersion       int32
	apiVersion          int32
	clientVersion       int32
	clientLatestVersion int32
	clientName          string
	timestamp           time.Time
}

func acquireContext(timestamp time.Time, server *baseServer, operation IOperation, actor IActor, requestId uint64, serverVersion, apiVersion, clientVersion, clientLatestVersion int32, clientName string, resultType ID) IContext {
	context := new(context)

	context.server = server
	context.operation = operation
	context.actor = actor
	context.analyticsProvider = DefaultProvider
	context.resultType = resultType
	context.requestId = requestId
	context.serverVersion = serverVersion
	context.apiVersion = apiVersion
	context.clientVersion = clientVersion
	context.clientLatestVersion = clientLatestVersion
	context.clientName = clientName
	context.timestamp = timestamp

	return context
}

func (context *context) Configuration() IConfiguration {
	return context.server.configuration
}

func (context *context) SetSecureCookie(key, value string) {
	if context.actor.Writer() != nil {
		context.actor.Writer().SetSecureCookie(key, value)
	}
}

func (context *context) GetSecureCookie(key string) string {
	return context.actor.Writer().GetSecureCookie(key)
}

func (context *context) Token() string {
	return context.actor.Token()
}

func (context *context) Identity() Identity {
	identity := context.actor.Identity()
	if identity.Role() == ADMINISTRATOR {
		identity.SetSystemCallHandler(context.server.systemCall)
	}

	return identity
}

func (context *context) RequestId() uint64 {
	return context.requestId
}

func (context *context) ApiVersion() int32 {
	return context.apiVersion
}

func (context *context) ClientVersion() int32 {
	return context.clientVersion
}

func (context *context) ClientLatestVersion() int32 {
	return context.clientLatestVersion
}

func (context *context) ServerVersion() int32 {
	return context.serverVersion
}

func (context *context) ClientName() string {
	return context.clientName
}

func (context *context) ResultType() ID {
	return context.resultType
}

func (context *context) SetResultType(resultType ID) {
	context.resultType = resultType
}

func (context *context) IncrementActorsCount() {
	context.server.IncrementActorsCount(context.actor)
}

func (context *context) ActorsCount() int {
	return context.server.ActorsCount()
}

func (context *context) Scheduler() IScheduler {
	return context.server.scheduler
}

func (context *context) SetTimeout(callback func(), timeout time.Duration) string {
	return context.Scheduler().SetTimeout(callback, timeout)
}

func (context *context) SetInterval(callback func(), timeout time.Duration) string {
	return context.Scheduler().SetInterval(callback, timeout)
}

func (context *context) CancelSchedule(id string) {
	context.Scheduler().Cancel(id)
}

func (context *context) Logger() ILogger {
	return context.server.logger
}

func (context *context) SecurityHandler() ISecurityHandler {
	return context.server.securityHandler
}

func (context *context) Push(message IPushMessage) error {
	return context.server.Push(context.actor, message)
}

func (context *context) Broadcast(resultType uint64, payload Pointer) error {
	return context.server.Broadcast(resultType, payload)
}

func (context *context) BroadcastSpecific(resultType uint64, payloads map[string]Pointer) error {
	return context.server.BroadcastSpecific(resultType, payloads)
}

func (context *context) SMS(phoneNumber, message string) error {
	provider := context.server.SMSProvider()
	if provider == nil {
		return errors.New("no sms provider")
	}

	if err := provider.Send(phoneNumber, message); err != nil {
		return err
	}

	return nil
}

func (context *context) Timestamp() time.Time {
	return context.timestamp
}

func (context *context) IsStagingEnvironment() bool {
	return context.server.Configuration().IsStagingEnvironment()
}

func (context *context) IsProductionEnvironment() bool {
	return context.server.Configuration().IsProductionEnvironment()
}

func (context *context) ResultContainer() Pointer {
	return context.operation.OutputContainer()
}

func (context *context) SubmitAnalyticsEvent(userId uint64, eventName string, eventData Fields) {
	context.analyticsProvider.SubmitEvent(userId, eventName, eventData)
}

func (context *context) SubmitMeasurement(key string, tags Tags, fields Fields) {
	context.server.measurementsProvider.SubmitMeasurementAsync(key, tags, fields)
}

func (context *context) Async(runnable func()) {
	concurrent.NewAsyncTask(runnable).Run()
}

func (context *context) Lock() {
	context.operation.Lock()
}

func (context *context) Unlock() {
	context.operation.Unlock()
}

func (context *context) SystemCall(args []string) error {
	return context.server.systemCall(context.actor.Identity(), args)
}
