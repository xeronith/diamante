package client

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
)

type IClient interface {
	SetName(string)
	SetToken(string)
	SetVersion(int32)
	SetApiVersion(int32)
	Connect(string, string) error
	Disconnect() error
	Send(uint64, uint64, Pointer) error
	OnConnectionEstablished(func(IClient))
	SetOperationResultListener(func(IOperationResult))
	Serializer() ISerializer
	IsActive() bool
}
