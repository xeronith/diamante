package server

import (
	"time"

	"github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/scheduling"
	. "github.com/xeronith/diamante/contracts/security"
	"github.com/xeronith/diamante/contracts/service"
	. "github.com/xeronith/diamante/contracts/settings"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/utility/concurrent"
)

type mockContext struct {
	token     string
	timestamp time.Time
}

func CreateMockContext() service.IContext {
	return &mockContext{
		token:     "SAMPLE-TOKEN",
		timestamp: time.Now(),
	}
}

func (context *mockContext) Configuration() IConfiguration {
	return nil
}

func (context *mockContext) SetCookie(string, string) {
}

func (context *mockContext) GetCookie(string) string {
	return ""
}

func (context *mockContext) Token() string {
	return context.token
}

func (context *mockContext) Identity() Identity {
	return nil
}

func (context *mockContext) RequestId() uint64 {
	return 0
}

func (context *mockContext) ApiVersion() int32 {
	return 0
}

func (context *mockContext) ClientVersion() int32 {
	return 0
}

func (context *mockContext) ClientLatestVersion() int32 {
	return 0
}

func (context *mockContext) ServerVersion() int32 {
	return 0
}

func (context *mockContext) ClientName() string {
	return ""
}

func (context *mockContext) ResultType() ID {
	return 0
}

func (context *mockContext) SetResultType(resultType ID) {
}

func (context *mockContext) IncrementActorsCount() {
}

func (context *mockContext) ActorsCount() int {
	return 0
}

func (context *mockContext) Scheduler() IScheduler {
	return nil
}

func (context *mockContext) SetTimeout(callback func(), timeout time.Duration) string {
	return context.Scheduler().SetTimeout(callback, timeout)
}

func (context *mockContext) SetInterval(callback func(), timeout time.Duration) string {
	return context.Scheduler().SetInterval(callback, timeout)
}

func (context *mockContext) CancelSchedule(id string) {
	context.Scheduler().Cancel(id)
}

func (context *mockContext) Logger() ILogger {
	return nil
}

func (context *mockContext) SecurityHandler() ISecurityHandler {
	return nil
}

func (context *mockContext) Push(message IPushMessage) error {
	return nil
}

func (context *mockContext) Broadcast(resultType uint64, payload Pointer) error {
	return nil
}

func (context *mockContext) BroadcastSpecific(resultType uint64, payloads map[string]Pointer) error {
	return nil
}

func (context *mockContext) SMS(phoneNumber, message string) error {
	return nil
}

func (context *mockContext) Timestamp() time.Time {
	return context.timestamp
}

func (context *mockContext) IsStagingEnvironment() bool {
	return true
}

func (context *mockContext) IsProductionEnvironment() bool {
	return !context.IsStagingEnvironment()
}

func (context *mockContext) ResultContainer() Pointer {
	return nil
}

func (context *mockContext) SubmitAnalyticsEvent(uint64, string, analytics.Fields) {
}

func (context *mockContext) SubmitMeasurement(string, analytics.Tags, analytics.Fields) {
}

func (context *mockContext) Async(runnable func()) {
	concurrent.NewAsyncTask(runnable).Run()
}

func (context *mockContext) Lock() {
}

func (context *mockContext) Unlock() {
}

func (context *mockContext) SystemCall(_ []string) error {
	return nil
}
