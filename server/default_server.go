package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"

	_ "embed"

	"github.com/gorilla/securecookie"
	"github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/settings"
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/localization"
	. "github.com/xeronith/diamante/logging"
	"github.com/xeronith/diamante/operation"
	. "github.com/xeronith/diamante/security"
	. "github.com/xeronith/diamante/serialization"
	. "github.com/xeronith/diamante/utility/collections"
	. "github.com/xeronith/diamante/utility/concurrent"
)

type defaultServer struct {
	baseServer
}

func New(configuration IConfiguration, operationFactory IOperationFactory, handlerFactory IHttpHandlerFactory) (IServer, error) {
	opcodes := make(Opcodes)
	for _, operation := range operationFactory.Operations() {
		requestId, resultId := operation.Id()
		opcodes[requestId] = operation.Tag()
		opcodes[resultId] = operation.Tag()
	}

	activePort, passivePort, diagnosticsPort := configuration.GetPorts()
	hashKey := []byte(configuration.GetServerConfiguration().GetHashKey())
	blockKey := []byte(configuration.GetServerConfiguration().GetBlockKey())

	serializers := map[string]ISerializer{
		"application/octet-stream": NewProtobufSerializer(),
		"application/json":         NewJsonSerializer(),
	}

	server := &defaultServer{
		baseServer{
			opcodes:              opcodes,
			activePort:           activePort,
			passivePort:          passivePort,
			diagnosticsPort:      diagnosticsPort,
			listeners:            NewConcurrentSlice(),
			configuration:        configuration,
			operations:           make(map[uint64]IOperation),
			securityHandler:      NewDefaultSecurityHandler(),
			scheduler:            newScheduler(),
			serializers:          serializers,
			actors:               NewConcurrentStringMap(),
			connectedActors:      NewConcurrentPointerMap(),
			connectedActorsCount: 0,
			logger:               GetDefaultLogger(),
			localizer:            NewLocalizer(),
			clientRegistry:       NewConcurrentStringToIntMap(),
			operationRequestPool: &sync.Pool{New: func() interface{} { return operation.NewOperationRequest() }},
			secureCookie:         securecookie.New(hashKey, blockKey),
			httpGetHandlers:      make(map[string]IHttpHandler),
			httpPostHandlers:     make(map[string]IHttpHandler),
			hudEnabled:           false,
			onServerStarted:      nil,
			onActorConnected:     nil,
			onActorDisconnected:  nil,
		},
	}

	if configuration.IsTestEnvironment() {
		server.activePort = rand.Intn(8999) + 1000
		server.passivePort = rand.Intn(8999) + 1000
		server.diagnosticsPort = rand.Intn(8999) + 1000
	}

	if operationFactory != nil {
		operations := operationFactory.Operations()
		if len(operations) > 0 {
			for _, operation := range operations {
				if err := server.RegisterOperation(operation); err != nil {
					return nil, err
				}
			}
		}
	}

	if handlerFactory != nil {
		handlers := handlerFactory.Handlers()
		if len(handlers) > 0 {
			for _, handler := range handlers {
				if err := server.RegisterHttpHandler(handler); err != nil {
					return nil, err
				}
			}
		}
	}

	return server, nil
}

func (server *defaultServer) Start() {
	if server.measurementsProvider == nil {
		server.Logger().Fatal("Server has no measurements provider.")
	}

	if server.running {
		server.Logger().Warning("Server is already running.")
		return
	}

	if server.asciiArt != "" {
		// https://fsymbols.com/generators/tarty/
		fmt.Println(server.asciiArt)
	}

	for opcode, role := range server.securityHandler.AccessControlHandler().AccessControls() {
		if operation, exists := server.operations[opcode]; exists {
			operation.SetRole(role)
		}
	}

	tasks := CreateAsyncTaskPool(false)

	tasks.Submit(
		func() { server.startActiveServer() },
		func() { server.startPassiveServer() },
		func() { server.startServerScheduler() },
		func() { server.startDiagnosticsServer() },
	)

	server.running = true
	server.measurement("core", analytics.Tags{"type": "i"}, analytics.Fields{"event": "0"})

	if server.onServerStarted != nil {
		server.onServerStarted()
	}

	tasks.Run().Join()
}

func (server *defaultServer) Shutdown() {
	server.listeners.ForEach(func(index int, object ISystemObject) {
		if err := object.(net.Listener).Close(); err != nil {
			log.Println(err)
		}
	})
}
