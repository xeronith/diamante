package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
)

type IPipeline interface {
	Actor() IActor
	Operation() IOperation
	Serializer() ISerializer
	Opcode() uint64
	RequestId() uint64
	ResultType() uint64
	ContentType() string
	ApiVersion() int32
	ServerVersion() int32
	ClientVersion() int32
	ClientLatestVersion() int32
	ClientName() string
	IsFrozen() bool
	IsSystemCall() bool
	Hash([]byte) string

	ServiceUnavailable(...error) IOperationResult
	InternalServerError(...error) IOperationResult
	NotImplemented(...error) IOperationResult
	Unauthorized(...error) IOperationResult
	BadRequest(...error) IOperationResult
}
