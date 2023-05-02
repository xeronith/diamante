package http

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	. "github.com/xeronith/diamante/actor"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/operation/binary"
	. "github.com/xeronith/diamante/protobuf"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Handle[T, V protoreflect.ProtoMessage](
	x IServerDispatcher,
	entryPoint string,
	requestType,
	resultType uint64,
	input T,
	output V,
	onInputUnmarshalled func(T),
	onRequestProcessed func(V) (string, []byte),
	redirect bool,
) error {
	x.Logger().Debug(entryPoint)

	body, err := io.ReadAll(x.Request().Body)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		if onRequestProcessed != nil {
			body, err = json.Marshal(struct {
				Body string `json:"body"`
			}{Body: string(body)})

			if err != nil {
				return err
			}
		}

		if err := protojson.Unmarshal(body, input); err != nil {
			return err
		}
	}

	if onInputUnmarshalled != nil {
		onInputUnmarshalled(input)
	}

	request := CreateBinaryOperationRequest(
		uint64(time.Now().UnixNano()),
		requestType,
		"pipeline",
		0, 0, "", nil,
	)

	if err := request.Load(input, x.Serializer()); err != nil {
		return err
	}

	data, err := x.Serialize(request.Container())
	if err != nil {
		return err
	}

	actor := CreateActor(nil, false, x.RemoteAddr(), x.UserAgent())
	result := x.OnActorBinaryData(actor, data)
	if result.Type() != resultType {
		if result.Type() == 0 {
			serverErr := &ServerError{}
			if err = x.Deserialize(result.(IBinaryOperationResult).Payload(), serverErr); err != nil {
				return err
			}

			if serverErr.Description != "" {
				return errors.New(serverErr.Description)
			} else {
				return errors.New(serverErr.Message)
			}
		} else {
			return errors.New("internal_handler_error")
		}
	}

	payload := result.(IBinaryOperationResult).Payload()
	if err = x.Deserialize(payload, output); err != nil {
		return err
	}

	if !redirect {
		if onRequestProcessed != nil {
			contentType, data := onRequestProcessed(output)
			x.Response().Header().Add("Content-Type", contentType)
			_, _ = x.Response().Write(data)
		} else {
			data, err := protojson.Marshal(output)
			if err != nil {
				return err
			}

			x.Response().Header().Add("Content-Type", "application/activity+json; charset=utf-8")
			// x.Response().Header().Add("Content-Type", "application/json; charset=utf-8")
			_, _ = x.Response().Write(data)
		}
	}

	return nil
}
