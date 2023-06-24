package server

import (
	"errors"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/messaging"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/operation/binary"
	"github.com/xeronith/diamante/utility/reflection"
)

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

			binaryOperationResult := binary.CreateBinaryOperationResult(BROADCAST, OK, resultType, data, NO_PIPELINE_INFO, 0, "")

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

				binaryOperationResult := binary.CreateBinaryOperationResult(BROADCAST, OK, resultType, serializedPayload, NO_PIPELINE_INFO, 0, "")
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

	binaryOperationResult := binary.CreateBinaryOperationResult(BROADCAST, OK, resultType, data, NO_PIPELINE_INFO, 0, "")

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
