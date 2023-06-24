package server

import (
	"errors"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/operation"
	"github.com/xeronith/diamante/utility/reflection"
)

func (server *baseServer) Broadcast(resultType uint64, payload Pointer) error {
	if !reflection.IsPointer(payload) {
		// TODO: What if payload is nil?
		return errors.New("broadcast_failure: non_pointer_payload")
	}

	var (
		err               error
		serializedPayload []byte
	)

	return server.connectedActors.ForEachParallelWithInitialization(
		func(count int) error {
			// fmt.Printf("Broadcast to %d clients\n", count)
			serializer := server.serializers["application/octet-stream"]
			serializedPayload, err = serializer.Serialize(payload)
			if err != nil {
				server.logger.Error(err)
				return err
			}

			operationResult := CreateOperationResult(BROADCAST, OK, resultType, serializedPayload, NO_PIPELINE_INFO, 0)

			if serializedPayload, err = serializer.Serialize(operationResult.Container()); err != nil {
				server.logger.Error(err)
				return err
			}

			return nil
		},
		func(object Pointer) {
			actor := object.(IActor)
			if actor.Writer().IsOpen() {
				actor.Writer().WriteBytes(serializedPayload)
			} else {
				server.OnSocketDisconnected(actor)
			}
		})
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
			serializer := server.serializers["application/octet-stream"]
			for token, payload := range payloads {
				// TODO: What if payload is nil?
				serializedPayload, err := serializer.Serialize(payload)
				if err != nil {
					server.logger.Error(err)
					return err
				}

				operationResult := CreateOperationResult(BROADCAST, OK, resultType, serializedPayload, NO_PIPELINE_INFO, 0)
				serializedOperationResult, err := serializer.Serialize(operationResult.Container())
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
					actor.Writer().WriteBytes(data[actor.Token()])
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
		err               error
		serializedPayload []byte
	)

	serializer := server.serializers["application/octet-stream"]
	serializedPayload, err = serializer.Serialize(payload)
	if err != nil {
		server.logger.Error(err)
		return err
	}

	operationResult := CreateOperationResult(BROADCAST, OK, resultType, serializedPayload, NO_PIPELINE_INFO, 0)
	if serializedPayload, err = serializer.Serialize(operationResult.Container()); err != nil {
		server.logger.Error(err)
		return err
	}

	if actor.Writer().IsOpen() {
		actor.Writer().WriteBytes(serializedPayload)
	} else {
		server.OnSocketDisconnected(actor)
	}

	return nil
}
