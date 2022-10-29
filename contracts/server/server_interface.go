package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/email"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/settings"
	. "github.com/xeronith/diamante/contracts/sms"
	. "github.com/xeronith/diamante/contracts/system"
)

type IServer interface {
	IBaseServer

	Start()
	Shutdown()

	OnServerStarted(func())
	OnActorConnected(func(string))
	OnActorDisconnected(func(string))

	SetSecurityHandler(ISecurityHandler)

	Version() int32
	RegisterClientVersion(string, int32)
	ResolveClientVersion(string) int32
	Configuration() IConfiguration

	ActiveEndpoint() string
	PassiveEndpoint() string

	TextSerializer() ITextSerializer
	BinarySerializer() IBinarySerializer

	TrafficRecorder() ITrafficRecorder

	MeasurementsProvider() IMeasurementsProvider
	SetMeasurementsProvider(IMeasurementsProvider)

	EmailProvider() IEmailProvider
	SetEmailProvider(IEmailProvider)

	SMSProvider() ISMSProvider
	SetSMSProvider(ISMSProvider)

	Actor(string) (IActor, error)
	Session(string) (ISystemObject, error)
	SetSession(string, ISystemObject) error

	OnActorBinaryData(IActor, []byte) IOperationResult
	OnActorTextData(IActor, string) IOperationResult
	OnActorOperationRequest(IActor, IOperationRequest) IOperationResult

	OnSocketConnected(IActor)
	OnSocketDisconnected(IActor)

	RegisterOperation(IOperation) error
	RegisterOperations(...IOperation) error

	RegisterHttpHandler(IHttpHandler) error
	RegisterHttpHandlers(...IHttpHandler) error

	SetAsciiArt(string)
	SetHUDEnabled(bool)
}
