package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/operation"
)

type IPipeline interface {
	Actor() IActor
	Operation() IOperation
	IsBinary() bool
	Opcode() uint64
	RequestId() uint64
	ResultType() uint64
	ApiVersion() int32
	ServerVersion() int32
	ClientVersion() int32
	ClientLatestVersion() int32
	ClientName() string

	IsSystemCall() bool

	ServiceUnavailable(...error) IOperationResult
	InternalServerError(...error) IOperationResult
	NotImplemented(...error) IOperationResult
	Unauthorized(...error) IOperationResult
	BadRequest(...error) IOperationResult
}
