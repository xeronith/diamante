package server

import (
	"errors"

	. "github.com/xeronith/diamante/contracts/operation"
	"github.com/xeronith/diamante/operation/binary"
	"github.com/xeronith/diamante/operation/text"
	"github.com/xeronith/diamante/protobuf"
)

var (
	INVALID_PARAMETERS                            = errors.New("invalid parameters")
	NON_POINTER_PAYLOAD_CONTAINER                 = errors.New("non_pointer_payload_container")
	SERVICE_EXECUTION_FAILURE                     = errors.New("service_execution_failure")
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
	serverError := &protobuf.ServerError{}
	if err != nil {
		serverError.Message = err.Error()
		serverError.Description = pipeline.localizer.Get(serverError.Message)
	}

	var result IOperationResult
	if pipeline.IsBinary() {
		var factory = binary.CreateBinaryOperationResult
		if data, serializationErr := pipeline.binarySerializer.Serialize(serverError); serializationErr != nil {
			result = factory(pipeline.RequestId(), InternalServerError, ERROR, nil, pipeline, 0, "")
		} else {
			result = factory(pipeline.RequestId(), status, ERROR, data, pipeline, 0, "")
		}
	} else {
		var factory = text.CreateTextOperationResult
		if data, serializationErr := pipeline.textSerializer.Serialize(serverError); serializationErr != nil {
			result = factory(pipeline.RequestId(), InternalServerError, ERROR, "", pipeline, 0, "")
		} else {
			result = factory(pipeline.RequestId(), status, ERROR, data, pipeline, 0, "")
		}
	}

	return result
}
