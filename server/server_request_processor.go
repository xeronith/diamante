package server

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/operation"
	. "github.com/xeronith/diamante/utility/reflection"
)

func (server *baseServer) OnOperationRequest(pipeline IPipeline, request IOperationRequest) IOperationResult {
	operation := pipeline.Operation()
	if operation == nil {
		return pipeline.NotImplemented()
	}

	if err := server.authorize(pipeline); err != nil {
		return pipeline.Unauthorized()
	}

	if item, exists := server.cache.Get(pipeline.Actor().RequestHash()); exists {
		result := item.(IOperationResult).ResetDuration()
		if pipeline.IsAcceptable(result) {
			return result
		}
	}

	container := operation.InputContainer()
	if container == nil || !IsPointer(container) {
		return pipeline.InternalServerError(NON_POINTER_PAYLOAD_CONTAINER)
	}

	if err := pipeline.Serializer().
		Deserialize(request.Payload(), container); err != nil {
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

	if payload, err := pipeline.Serializer().Serialize(output); err != nil {
		return pipeline.InternalServerError(err)
	} else {
		result := CreateOperationResult(pipeline.RequestId(), OK, context.ResultType(), payload, pipeline, duration)
		server.cache.Put(result.Hash(), result)
		return result
	}
}
