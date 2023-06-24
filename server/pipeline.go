package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/localization"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
)

var NO_PIPELINE_INFO = &pipeline{}

type pipeline struct {
	localizer           ILocalizer
	textSerializer      ITextSerializer
	binarySerializer    IBinarySerializer
	actor               IActor
	operation           IOperation
	opcode              uint64
	requestId           uint64
	resultType          uint64
	binary              bool
	apiVersion          int32
	serverVersion       int32
	clientVersion       int32
	clientLatestVersion int32
	clientName          string
}

func NewPipeline(server *baseServer, actor IActor, request IOperationRequest) IPipeline {
	binary := true
	switch request.(type) {
	case ITextOperationRequest:
		binary = false
	}

	operation := server.operations[request.Operation()]

	var resultType uint64
	if operation != nil {
		_, resultType = operation.Id()
	}

	actor.SetToken(request.Token())

	return &pipeline{
		localizer:           server.localizer,
		textSerializer:      server.textSerializer,
		binarySerializer:    server.binarySerializer,
		actor:               actor,
		operation:           operation,
		opcode:              request.Operation(),
		requestId:           request.Id(),
		resultType:          resultType,
		binary:              binary,
		apiVersion:          request.ApiVersion(),
		serverVersion:       server.Version(),
		clientVersion:       request.ClientVersion(),
		clientLatestVersion: server.ResolveClientVersion(request.ClientName()),
		clientName:          request.ClientName(),
	}
}

func (pipeline *pipeline) Actor() IActor {
	return pipeline.actor
}

func (pipeline *pipeline) Operation() IOperation {
	return pipeline.operation
}

func (pipeline *pipeline) IsBinary() bool {
	return pipeline.binary
}

func (pipeline *pipeline) Opcode() uint64 {
	return pipeline.opcode
}

func (pipeline *pipeline) RequestId() uint64 {
	return pipeline.requestId
}

func (pipeline *pipeline) ResultType() uint64 {
	return pipeline.resultType
}

func (pipeline *pipeline) ApiVersion() int32 {
	return pipeline.apiVersion
}

func (pipeline *pipeline) ServerVersion() int32 {
	return pipeline.serverVersion
}

func (pipeline *pipeline) ClientVersion() int32 {
	return pipeline.clientVersion
}

func (pipeline *pipeline) ClientLatestVersion() int32 {
	return pipeline.clientVersion
}

func (pipeline *pipeline) ClientName() string {
	return pipeline.clientName
}

func (pipeline *pipeline) IsSystemCall() bool {
	return pipeline.opcode == SYSTEM_CALL_REQUEST
}
