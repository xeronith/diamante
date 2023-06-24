package server

import (
	"errors"

	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/operation"
	. "github.com/xeronith/diamante/protobuf"
)

var (
	INVALID_PARAMETERS                            = errors.New("invalid parameters")
	NON_POINTER_PAYLOAD_CONTAINER                 = errors.New("non_pointer_payload_container")
	SERVICE_EXECUTION_FAILURE                     = errors.New("service_execution_failure")
	INPUT_STREAM_DESERIALIZATION_FAILURE          = errors.New("input_stream_deserialization_failure")
	SERVICE_UNAVAILABLE_DUE_TO_SYSTEM_MAINTENANCE = errors.New("system_maintenance")
	NOT_IMPLEMENTED                               = errors.New("not_implemented")
	INTERNAL_SERVER_ERROR                         = errors.New("internal_server_error")
	UNAUTHORIZED                                  = errors.New("unauthorized")
	BAD_REQUEST                                   = errors.New("bad_request")
)

func (pipeline *pipeline) ServiceUnavailable(errors ...error) IOperationResult {
	err := SERVICE_UNAVAILABLE_DUE_TO_SYSTEM_MAINTENANCE
	if len(errors) > 0 {
		err = errors[0]
	}

	return pipeline.serverError(ServiceUnavailable, err)
}

func (pipeline *pipeline) NotImplemented(errors ...error) IOperationResult {
	err := NOT_IMPLEMENTED
	if len(errors) > 0 {
		err = errors[0]
	}

	return pipeline.serverError(NotImplemented, err)
}

func (pipeline *pipeline) InternalServerError(errors ...error) IOperationResult {
	err := INTERNAL_SERVER_ERROR
	if len(errors) > 0 {
		err = errors[0]
	}

	return pipeline.serverError(InternalServerError, err)
}

func (pipeline *pipeline) Unauthorized(errors ...error) IOperationResult {
	err := UNAUTHORIZED
	if len(errors) > 0 {
		err = errors[0]
	}

	return pipeline.serverError(Unauthorized, err)
}

func (pipeline *pipeline) BadRequest(errors ...error) IOperationResult {
	err := BAD_REQUEST
	if len(errors) > 0 {
		err = errors[0]
	}

	return pipeline.serverError(BadRequest, err)
}

func (pipeline *pipeline) serverError(status int32, err error) IOperationResult {
	serverError := &ServerError{}
	if err != nil {
		serverError.Message = err.Error()
		serverError.Description = pipeline.localizer.Get(serverError.Message)
	}

	var (
		payload []byte
		result  IOperationResult
	)

	if data, serializationErr := pipeline.serializer.Serialize(serverError); serializationErr != nil {
		status = InternalServerError
	} else {
		payload = data
	}

	result = CreateOperationResult(
		pipeline.RequestId(), // id
		status,               // status
		ERROR,                // resultType
		payload,              // payload
		pipeline,             // pipelineInfo
		0,                    // duration
	)

	return result
}
