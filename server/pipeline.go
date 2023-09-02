package server

import (
	"fmt"
	"strings"

	"github.com/go-faster/city"
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/localization"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
)

var NO_PIPELINE_INFO = &pipeline{}

type pipeline struct {
	server              *baseServer
	localizer           ILocalizer
	serializer          ISerializer
	actor               IActor
	operation           IOperation
	request             IOperationRequest
	opcode              uint64
	requestId           uint64
	resultType          uint64
	contentType         string
	apiVersion          int32
	serverVersion       int32
	clientVersion       int32
	clientLatestVersion int32
	clientName          string
}

func NewPipeline(server *baseServer, actor IActor, request IOperationRequest) IPipeline {
	contentType := actor.Writer().ContentType()
	operation := server.operations[request.Operation()]

	var serializer ISerializer
	if contentSerializer, ok := server.serializers[contentType]; !ok {
		serializer = server.serializers["application/octet-stream"]
	} else {
		serializer = contentSerializer
	}

	var resultType uint64
	if operation != nil {
		_, resultType = operation.Id()
	}

	actor.SetToken(request.Token())

	return &pipeline{
		server:              server,
		localizer:           server.localizer,
		serializer:          serializer,
		actor:               actor,
		operation:           operation,
		request:             request,
		opcode:              request.Operation(),
		requestId:           request.Id(),
		resultType:          resultType,
		contentType:         contentType,
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

func (pipeline *pipeline) Signature() string {
	return pipeline.actor.Signature()
}

func (pipeline *pipeline) Operation() IOperation {
	return pipeline.operation
}

func (pipeline *pipeline) Serializer() ISerializer {
	return pipeline.serializer
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

func (pipeline *pipeline) ContentType() string {
	return pipeline.contentType
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

func (pipeline *pipeline) IsFrozen() bool {
	return pipeline.server.IsFrozen() && pipeline.opcode != SYSTEM_CALL_REQUEST
}

func (pipeline *pipeline) IsSystemCall() bool {
	return pipeline.opcode == SYSTEM_CALL_REQUEST
}

func (pipeline *pipeline) Sign(payload []byte) string {
	if payload == nil || pipeline.request == nil {
		return ""
	}

	return fmt.Sprintf(
		"%x%x%x%x",
		city.Hash64([]byte(pipeline.actor.Token())),
		city.Hash64(pipeline.request.Payload()),
		pipeline.opcode,
		city.Hash64(payload),
	)
}

func (pipeline *pipeline) IsAcceptable(result IOperationResult) bool {
	if pipeline.request == nil {
		return false
	}

	return strings.HasPrefix(result.Signature(), fmt.Sprintf(
		"%x%x%x",
		city.Hash64([]byte(pipeline.actor.Token())),
		city.Hash64(pipeline.request.Payload()),
		pipeline.opcode,
	))
}
