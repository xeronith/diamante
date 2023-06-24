package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/operation"
)

func (server *baseServer) OnData(actor IActor, data []byte) IOperationResult {
	request := server.operationRequestPool.Get().(IOperationRequest)
	if err := actor.Serializer().Deserialize(data, request.Container()); err != nil {
		pipeline := NewPipeline(server, actor, request)
		return pipeline.InternalServerError(INPUT_STREAM_DESERIALIZATION_FAILURE)
	}

	pipeline := NewPipeline(server, actor, request)
	// r: request_initiated, f: request_finalized, op: operation, id: request_id
	fields := Fields{"op": int64(pipeline.Opcode()), "id": int64(pipeline.RequestId())}
	/* //////// */ server.measurement("operations", Tags{"type": "r"}, fields)
	defer func() { server.measurement("operations", Tags{"type": "f"}, fields) }()

	if pipeline.IsFrozen() {
		return pipeline.ServiceUnavailable()
	}

	return server.OnOperationRequest(pipeline, request)
}
