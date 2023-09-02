package server

import (
	"errors"
	"fmt"

	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/securecookie"
	. "github.com/xeronith/diamante/contracts/email"
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/localization"
	. "github.com/xeronith/diamante/contracts/settings"
	. "github.com/xeronith/diamante/contracts/sms"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/scheduling"
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/service"
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/network/http"
	. "github.com/xeronith/diamante/utility/collections"
)

type baseServer struct {
	mutex                   sync.RWMutex
	asciiArt                string
	hudEnabled              bool
	activePort, passivePort int
	diagnosticsPort         int
	running                 bool
	frozen                  bool
	listeners               ISlice
	operations              map[uint64]IOperation
	opcodes                 Opcodes
	configuration           IConfiguration
	clientRegistry          IStringToIntMap
	emailProvider           IEmailProvider
	smsProvider             ISMSProvider
	securityHandler         ISecurityHandler
	serializers             map[string]ISerializer
	scheduler               IScheduler
	actors                  IStringMap
	logger                  ILogger
	localizer               ILocalizer
	cache                   IStringMap
	cacheMiss               int64
	cacheHit                int64
	onStorageUpdated        func(...string)
	measurementsProvider    IMeasurementsProvider
	operationRequestPool    *sync.Pool
	secureCookie            *securecookie.SecureCookie

	// LEGACY
	connectedActors      IPointerMap
	connectedActorsCount int32

	// HANDLERS
	httpGetHandlers  map[string]IHttpHandler
	httpPostHandlers map[string]IHttpHandler

	// EVENTS
	onServerStarted     func()
	onActorConnected    func(string)
	onActorDisconnected func(string)
}

func (server *defaultServer) SetAsciiArt(asciiArtString string) {
	server.asciiArt = asciiArtString
}

func (server *defaultServer) SetHUDEnabled(enabled bool) {
	server.hudEnabled = enabled
}

func (server *baseServer) SetSecurityHandler(handler ISecurityHandler) {
	server.securityHandler = handler
}

func (server *baseServer) getActiveProtocol() string {
	if server.Configuration().GetServerConfiguration().GetTLSConfiguration().IsEnabled() {
		return "wss"
	} else {
		return "ws"
	}
}

func (server *baseServer) getPassiveProtocol() string {
	if server.Configuration().GetServerConfiguration().GetTLSConfiguration().IsEnabled() {
		return "https"
	} else {
		return "http"
	}
}

func (server *baseServer) RegisterClientVersion(clientName string, version int32) {
	server.clientRegistry.Put(clientName, version)
}

func (server *baseServer) ResolveClientVersion(clientName string) int32 {
	if server.clientRegistry.Contains(clientName) {
		return server.clientRegistry.Get(clientName)
	}

	return 0
}

func (server *baseServer) Version() int32 {
	return server.configuration.GetServerConfiguration().GetBuildNumber()
}

func (server *baseServer) ActiveEndpoint() string {
	return fmt.Sprintf("%s://%s:%d", server.getActiveProtocol(), server.Configuration().GetServerConfiguration().GetFQDN(), server.activePort)
}

func (server *baseServer) PassiveEndpoint() string {
	return fmt.Sprintf("%s://%s:%d", server.getPassiveProtocol(), server.Configuration().GetServerConfiguration().GetFQDN(), server.passivePort)
}

func (server *baseServer) OnStorageUpdated() func(...string) {
	return server.onStorageUpdated
}

func (server *baseServer) MeasurementsProvider() IMeasurementsProvider {
	return server.measurementsProvider
}

func (server *baseServer) SetMeasurementsProvider(provider IMeasurementsProvider) {
	server.measurementsProvider = provider
}

func (server *baseServer) EmailProvider() IEmailProvider {
	return server.emailProvider
}

func (server *baseServer) SetEmailProvider(provider IEmailProvider) {
	server.emailProvider = provider
}

func (server *baseServer) SMSProvider() ISMSProvider {
	return server.smsProvider
}

func (server *baseServer) SetSMSProvider(provider ISMSProvider) {
	server.smsProvider = provider
}

func (server *baseServer) getOperations() map[uint64]IOperation {
	return server.operations
}

func (server *baseServer) getSecurityHandler() ISecurityHandler {
	return server.securityHandler
}

func (server *baseServer) Configuration() IConfiguration {
	return server.configuration
}

func (server *baseServer) Scheduler() IScheduler {
	return server.scheduler
}

func (server *baseServer) Serializers() map[string]ISerializer {
	return server.serializers
}

func (server *baseServer) Serializer(writer IWriter) ISerializer {
	if writer.ContentType() == "" {
		return server.serializers["application/octet-stream"]
	}

	return server.serializers[writer.ContentType()]
}

func (server *baseServer) Logger() ILogger {
	return server.logger
}

func (server *baseServer) Opcodes() Opcodes {
	return server.opcodes
}

func (server *baseServer) Localizer() ILocalizer {
	return server.localizer
}

func (server *baseServer) ActorsCount() int {
	return int(atomic.LoadInt32(&server.connectedActorsCount))
}

func (server *baseServer) IncrementActorsCount(actor IActor) {
	atomic.AddInt32(&server.connectedActorsCount, 1)
	if server.onActorConnected != nil {
		server.onActorConnected(actor.Token())
	}
}

func (server *baseServer) OnSocketConnected(actor IActor) {
	server.connectedActors.Put(actor, "")
	server.measurement("websocket", Tags{"type": "c"}, Fields{"state": 1, "value": server.connectedActors.GetSize()})
}

func (server *baseServer) OnSocketDisconnected(actor IActor) {
	/* TODO: if !server.connectedActors.Contains(actor) {
		return
	}*/

	server.connectedActors.Remove(actor)
	server.measurement("websocket", Tags{"type": "c"}, Fields{"state": 2, "value": server.connectedActors.GetSize()})
	if atomic.LoadInt32(&server.connectedActorsCount) > 0 {
		atomic.AddInt32(&server.connectedActorsCount, -1)
	}

	if server.onActorDisconnected != nil {
		server.onActorDisconnected(actor.Token())
	}
}

func (server *baseServer) Actor(token string) (IActor, error) {
	actor, exists := server.actors.Get(token)
	if !exists {
		return nil, errors.New("actor not found")
	}

	return actor.(IActor), nil
}

func (server *baseServer) Session(token string) (ISystemObject, error) {
	actor, err := server.Actor(token)
	if err != nil {
		return nil, err
	}

	return actor.Session(), err
}

func (server *baseServer) SetSession(token string, session ISystemObject) error {
	actor, err := server.Actor(token)
	if err != nil {
		return err
	}

	actor.SetSession(session)
	return nil
}

func (server *baseServer) RegisterOperation(operation IOperation) error {
	if operation == nil {
		return errors.New("nil operation")
	}

	operationId, _ := operation.Id()
	if operationId < 64 {
		return errors.New("operation ids below 64 are system reserved")
	}

	if server.operations[operationId] != nil {
		return fmt.Errorf("operation id %d already registered", operationId)
	}

	server.getOperations()[operationId] = operation

	return nil
}

func (server *baseServer) RegisterOperations(operations ...IOperation) error {
	for _, operation := range operations {
		if err := server.RegisterOperation(operation); err != nil {
			return err
		}
	}

	return nil
}

func (server *baseServer) RegisterHttpHandler(handler IHttpHandler) error {
	if server.running {
		return errors.New("not allowed to register http handlers when server is running")
	}

	path := handler.Path()
	method := handler.Method()
	handlerFunc := handler.HandlerFunc()

	if path == "/" {
		return errors.New("not_allowed_to_register_root_path")
	}

	if path == "/reports" {
		return errors.New("not_allowed_to_register_reports_path")
	}

	if path == "/mem" {
		return errors.New("not_allowed_to_register_mem_path")
	}

	if path == "/diagnostics" {
		return errors.New("not_allowed_to_diagnostics_mem_path")
	}

	if method != http.MethodGet && method != http.MethodPost {
		return fmt.Errorf("method_%s_not_allowed", method)
	}

	switch handler.Method() {
	case http.MethodGet:
		{
			if _, pathExists := server.httpGetHandlers[path]; pathExists {
				return fmt.Errorf("GET '%s' already_registered", path)
			}

			server.httpGetHandlers[path] = NewHttpHandler(path, method, handlerFunc)
		}

	case http.MethodPost:
		{
			if _, pathExists := server.httpPostHandlers[path]; pathExists {
				return fmt.Errorf("POST '%s' already_registered", path)
			}

			server.httpPostHandlers[path] = NewHttpHandler(path, method, handlerFunc)
		}
	}

	return nil
}

func (server *baseServer) RegisterHttpHandlers(handlers ...IHttpHandler) error {
	for _, handler := range handlers {
		if err := server.RegisterHttpHandler(handler); err != nil {
			return err
		}
	}

	return nil
}

func (server *baseServer) OnServerStarted(callback func()) {
	server.onServerStarted = callback
}

func (server *baseServer) OnActorConnected(callback func(string)) {
	server.onActorConnected = callback
}

func (server *baseServer) OnActorDisconnected(callback func(string)) {
	server.onActorDisconnected = callback
}

func (server *baseServer) executeService(context IContext, container Pointer, pipeline IPipeline) (Pointer, time.Duration, error) {
	operation := pipeline.Operation()
	operationId := pipeline.Opcode()
	requestId := pipeline.RequestId()

	defer server.catch(operationId, pipeline.RequestId())

	output, err := operation.Execute(context, container)
	duration := server.analyzeOperationPerformance(operation, operationId, context.Timestamp())
	server.measurement(
		"operations",
		Tags{"type": "x"},
		Fields{
			"operation": int64(operationId),
			"requestId": int64(requestId),
			"duration":  int64(duration),
		},
	)

	return output, duration, err
}

func (server *baseServer) analyzeOperationPerformance(operation IOperation, operationId uint64, timestamp time.Time) time.Duration {
	timeLimitWarning, timeLimitAlert, timeLimitCritical := operation.ExecutionTimeLimits()
	delta := time.Since(timestamp)
	if delta > timeLimitWarning {
		message := fmt.Sprintf(
			"SED 0x%.8X %016d %s",
			operationId,
			delta,
			server.opcodes[operationId],
		)

		if delta > timeLimitCritical {
			server.logger.Critical(message)
		} else if delta > timeLimitAlert {
			server.logger.Alert(message)
		} else {
			server.logger.Warning(message)
		}
	}

	return delta
}

func (server *baseServer) authorize(pipeline IPipeline) error {
	operation := pipeline.Operation()
	actor := pipeline.Actor()
	role := operation.Role()

	token := ""
	if actor.Writer() != nil {
		token = actor.Writer().GetAuthCookie()
	}

	if token == "" {
		token = actor.Token()
	}

	identity := server.getSecurityHandler().Authenticate(
		token,
		role,
		actor.RemoteAddress(),
		actor.UserAgent(),
	)

	if identity == nil {
		return UNAUTHORIZED
	}

	actor.SetIdentity(identity)
	actor.UpdateLastActivity()

	return nil
}

func (server *baseServer) systemCall(identity Identity, args []string) error {
	if len(args) < 1 {
		return errors.New("command_required")
	}

	switch args[0] {

	case "freeze", "resume":
		{
			server.mutex.Lock()
			server.frozen = args[0] == "freeze"
			server.mutex.Unlock()
			return nil
		}

	case "acl":
		{
			if len(args) < 3 {
				return INVALID_PARAMETERS
			}

			opcode, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return INVALID_PARAMETERS
			}

			if opcode == SYSTEM_CALL_REQUEST {
				return INVALID_PARAMETERS
			}

			role, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return INVALID_PARAMETERS
			}

			if operation, exists := server.operations[opcode]; !exists {
				return INVALID_PARAMETERS
			} else {
				if err := server.securityHandler.AccessControlHandler().
					AddOrUpdateAccessControl(opcode, role, identity); err != nil {
					return err
				}

				operation.SetRole(role)
				return nil
			}
		}

	default:
		return errors.New("syscall: command_not_found " + args[0])
	}
}

func (server *baseServer) IsFrozen() bool {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	return server.frozen
}

func (server *baseServer) measurement(key string, tags Tags, fields Fields) {
	if server.Configuration().IsDevelopmentEnvironment() {
		return
	}

	server.measurementsProvider.SubmitMeasurement(key, tags, fields)
}

func (server *baseServer) catch(operationId uint64, requestId uint64) {
	if reason := recover(); reason != nil {
		server.logger.Panic(
			fmt.Sprintf(
				"OPR 0x%.8X %s\n%s",
				operationId,
				reason,
				debug.Stack(),
			),
		)

		server.measurement(
			"operations",
			Tags{"type": "p"},
			Fields{
				"operation": int64(operationId),
				"requestId": int64(requestId),
			},
		)
	}
}
