package server

import (
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/server"
)

func (server *baseServer) OnActorBinaryData(actor IActor, data []byte) IOperationResult {
	server.TrafficRecorder().Record(BINARY_REQUEST, data)

	request := server.binaryOperationRequestPool.Get().(IBinaryOperationRequest)
	if err := server.binarySerializer.Deserialize(data, request.Container()); err != nil {
		return server.serverError(0, BadRequest, err, false, request.ApiVersion(), server.Version(), request.ClientVersion())
	} else {
		actor.SetToken(request.Token())
		return server.OnActorOperationRequest(actor, request)
	}
}

func (server *baseServer) OnActorTextData(actor IActor, data string) IOperationResult {
	server.TrafficRecorder().Record(TEXT_REQUEST, data)

	request := server.textOperationRequestPool.Get().(ITextOperationRequest)
	if err := server.textSerializer.Deserialize(data, request.Container()); err != nil {
		return server.serverError(0, BadRequest, err, true, request.ApiVersion(), server.Version(), request.ClientVersion())
	} else {
		actor.SetToken(request.Token())
		return server.OnActorOperationRequest(actor, request)
	}
}
