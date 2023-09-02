package actor

import (
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
)

type IActor interface {
	Identity() Identity
	SetIdentity(Identity)
	Token() string
	SetToken(string)
	Signature() string
	RemoteAddress() string
	UserAgent() string
	Dispatch(IOperationResult)
	Disconnect(IOperationResult)
	Session() ISystemObject
	SetSession(object ISystemObject)
	Serializer() ISerializer
	Signal(byte)
	UpdateLastActivity()
	LastActivity() int64
	IsActive() bool
	Writer() IWriter
}
