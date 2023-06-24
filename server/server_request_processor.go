package server

import (
	"fmt"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/operation"
	"github.com/xeronith/diamante/operation/binary"
	"github.com/xeronith/diamante/operation/text"
	. "github.com/xeronith/diamante/utility/reflection"
)

func (server *baseServer) OnActorOperationRequest(actor IActor, request IOperationRequest) IOperationResult {
	pipeline := NewPipeline(server, actor, request)

	server.measurement(
		"operations",
		Tags{"type": "r"},
		Fields{
			"operation": int64(pipeline.Opcode()),
			"requestId": int64(pipeline.RequestId()),
		},
	)

	defer func() {
		server.measurement(
			"operations",
			Tags{"type": "f"},
			Fields{
				"operation": int64(pipeline.Opcode()),
				"requestId": int64(pipeline.RequestId()),
			},
		)
	}()

	if server.IsFrozen() && !pipeline.IsSystemCall() {
		return pipeline.ServiceUnavailable()
	}

	operation := pipeline.Operation()
	if operation == nil {
		return pipeline.NotImplemented()
	}

	if err := server.authorize(pipeline); err != nil {
		return pipeline.Unauthorized()
	}

	container := operation.InputContainer()
	if container == nil || !IsPointer(container) {
		return pipeline.InternalServerError(NON_POINTER_PAYLOAD_CONTAINER)
	}

	var err error
	if pipeline.IsBinary() {
		err = server.binarySerializer.
			Deserialize(request.(IBinaryOperationRequest).Payload(), container)
	} else {
		err = server.textSerializer.
			Deserialize(request.(ITextOperationRequest).Payload(), container)
	}

	if err != nil {
		return pipeline.InternalServerError(err)
	}

	context := server.acquireContext(pipeline)
	output, duration, err := server.executeService(context, container, pipeline)

	if err != nil {
		return pipeline.InternalServerError(err)
	}

	if output == nil || !IsPointer(output) {
		return pipeline.InternalServerError(SERVICE_EXECUTION_FAILURE)
	}

	var result IOperationResult
	if pipeline.IsBinary() {
		resultPayload, err := server.binarySerializer.Serialize(output)
		if err != nil {
			return pipeline.InternalServerError(err)
		}

		result = binary.CreateBinaryOperationResult(
			pipeline.RequestId(),
			OK,
			context.ResultType(),
			resultPayload,
			pipeline,
			duration,
			fmt.Sprintf(
				"%x-%x-%x",
				pipeline.Opcode(),
				server.hash(request.(IBinaryOperationRequest).Payload()),
				server.hash(resultPayload),
			),
		)
	} else {
		resultPayload, err := server.textSerializer.Serialize(output)
		if err != nil {
			return pipeline.InternalServerError(err)
		}

		result = text.CreateTextOperationResult(
			pipeline.RequestId(),
			OK,
			context.ResultType(),
			resultPayload,
			pipeline,
			duration,
			fmt.Sprintf(
				"%x-%x-%x",
				pipeline.Opcode(),
				server.hash(request.(ITextOperationRequest).Payload()),
				server.hash(resultPayload),
			),
		)
	}

	return result
}
