package server

import (
	"errors"
	"time"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/analytics"
	. "github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/operation/binary"
	"github.com/xeronith/diamante/operation/text"
	"github.com/xeronith/diamante/utility/reflection"
)

func (server *baseServer) OnActorOperationRequest(actor IActor, request IOperationRequest) IOperationResult {
	opCode := int64(request.Operation())
	apiVersion := request.ApiVersion()
	serverVersion := server.Version()
	clientVersion := server.ResolveClientVersion(request.ClientName())
	requestId := request.Id()

	server.measurement("operations", Tags{"type": "r"}, Fields{"operation": opCode, "requestId": int64(requestId)})
	defer func() {
		server.measurement("operations", Tags{"type": "f"}, Fields{"operation": opCode, "requestId": int64(requestId)})
	}()

	var (
		nonBinary  bool
		serializer ISerializer
	)

	switch request.(type) {
	case IBinaryOperationRequest:
		nonBinary = false
		serializer = server.binarySerializer
	case ITextOperationRequest:
		nonBinary = true
		serializer = server.textSerializer
	}

	operation := server.getOperations()[request.Operation()]
	if operation == nil {
		return server.notImplemented(requestId, nil, nonBinary, apiVersion, serverVersion, clientVersion)
	}

	if err := server.authorize(actor, operation); err != nil {
		return server.unauthorized(requestId, nil, nonBinary, apiVersion, serverVersion, clientVersion)
	}

	server.mutex.RLock()
	frozen := server.frozen
	server.mutex.RUnlock()

	if frozen && opCode != 0x1000 {
		return server.serviceUnavailable(requestId, nil, nonBinary, apiVersion, serverVersion, clientVersion)
	}

	container := operation.InputContainer()
	if container == nil || !reflection.IsPointer(container) {
		return server.internalServerError(requestId, errors.New("non_pointer_payload_container"), nonBinary, apiVersion, serverVersion, clientVersion)
	}

	var err error
	if nonBinary {
		err = serializer.(ITextSerializer).Deserialize(request.(ITextOperationRequest).Payload(), container)
	} else {
		err = serializer.(IBinarySerializer).Deserialize(request.(IBinaryOperationRequest).Payload(), container)
	}

	if err != nil {
		return server.internalServerError(requestId, err, nonBinary, apiVersion, serverVersion, clientVersion)
	}

	_, resultType := operation.Id()

	timestamp := time.Now()
	context := acquireContext(timestamp, server, operation, actor, request.Id(), serverVersion, request.ApiVersion(), request.ClientVersion(), clientVersion, request.ClientName(), resultType)
	output, duration, err := server.executeService(timestamp, operation, requestId, context, container)
	contextResultType := context.ResultType()

	if err != nil {
		return server.internalServerError(requestId, err, nonBinary, apiVersion, serverVersion, clientVersion)
	}

	if output == nil {
		if //noinspection GoNilness
		err == nil {
			return server.internalServerError(requestId, errors.New("service_execution_failure"), nonBinary, apiVersion, serverVersion, clientVersion)
		}

		return server.internalServerError(requestId, errors.New("service_execution_failure: null"), nonBinary, apiVersion, serverVersion, clientVersion)
	}

	if !reflection.IsPointer(output) {
		return server.internalServerError(requestId, errors.New("non_pointer_service_result"), nonBinary, apiVersion, serverVersion, clientVersion)
	}

	var result IOperationResult
	if nonBinary {
		data, err := serializer.(ITextSerializer).Serialize(output)
		if err != nil {
			return server.internalServerError(requestId, err, nonBinary, apiVersion, serverVersion, clientVersion)
		}
		result = text.CreateTextOperationResult(requestId, OK, contextResultType, data, apiVersion, serverVersion, clientVersion, duration)
	} else {
		data, err := serializer.(IBinarySerializer).Serialize(output)
		if err != nil {
			return server.internalServerError(requestId, err, nonBinary, apiVersion, serverVersion, clientVersion)
		}
		result = binary.CreateBinaryOperationResult(requestId, OK, contextResultType, data, apiVersion, serverVersion, clientVersion, duration)
	}

	return result
}

func (server *baseServer) BroadcastSpecific(resultType uint64, payloads map[string]Pointer) error {
	if payloads == nil {
		return errors.New("broadcast_failure: null_payload")
	}

	if len(payloads) < 1 {
		return nil
	}

	data := make(map[string][]byte)
	return server.connectedActors.ForEachParallelWithInitialization(
		func(count int) error {
			serializer := server.binarySerializer
			for token, payload := range payloads {
				// TODO: What if payload is nil?
				serializedPayload, err := serializer.Serialize(payload)
				if err != nil {
					server.logger.Error(err)
					return err
				}

				binaryOperationResult := binary.CreateBinaryOperationResult(BROADCAST, OK, resultType, serializedPayload, 0, 0, 0, 0)
				serializedOperationResult, err := serializer.Serialize(binaryOperationResult.Container())
				if err != nil {
					server.logger.Error(err)
					return err
				}

				data[token] = serializedOperationResult
			}

			return nil
		}, func(object Pointer) {
			actor := object.(IActor)
			if actor.Writer().IsOpen() {
				if _, exists := data[actor.Token()]; exists {
					actor.Writer().WriteBytes(BINARY_RESULT, data[actor.Token()])
				}
			} else {
				server.OnSocketDisconnected(actor)
			}
		})
}

func (server *baseServer) Push(actor IActor, message IPushMessage) error {
	resultType := message.GetType()
	payload := message.GetPayload()

	if !reflection.IsPointer(payload) {
		// TODO: What if payload is nil?
		return errors.New("broadcast_failure: non_pointer_payload")
	}

	var (
		err  error
		data []byte
	)

	serializer := server.binarySerializer
	data, err = serializer.Serialize(payload)
	if err != nil {
		server.logger.Error(err)
		return err
	}

	binaryOperationResult := binary.CreateBinaryOperationResult(BROADCAST, OK, resultType, data, 0, 0, 0, 0)

	if data, err = serializer.Serialize(binaryOperationResult.Container()); err != nil {
		server.logger.Error(err)
		return err
	}

	if actor.Writer().IsOpen() {
		actor.Writer().WriteBytes(BINARY_RESULT, data)
	} else {
		server.OnSocketDisconnected(actor)
	}

	return nil
}

func (server *baseServer) Broadcast(resultType uint64, payload Pointer) error {
	if !reflection.IsPointer(payload) {
		// TODO: What if payload is nil?
		return errors.New("broadcast_failure: non_pointer_payload")
	}

	var (
		err  error
		data []byte
	)

	return server.connectedActors.ForEachParallelWithInitialization(
		func(count int) error {
			// fmt.Printf("Broadcast to %d clients\n", count)
			serializer := server.binarySerializer
			data, err = serializer.Serialize(payload)
			if err != nil {
				server.logger.Error(err)
				return err
			}

			binaryOperationResult := binary.CreateBinaryOperationResult(BROADCAST, OK, resultType, data, 0, 0, 0, 0)

			if data, err = serializer.Serialize(binaryOperationResult.Container()); err != nil {
				server.logger.Error(err)
				return err
			}

			return nil
		},
		func(object Pointer) {
			actor := object.(IActor)
			if actor.Writer().IsOpen() {
				actor.Writer().WriteBytes(BINARY_RESULT, data)
			} else {
				server.OnSocketDisconnected(actor)
			}
		})
}
