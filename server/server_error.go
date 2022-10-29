package server

import (
	"fmt"
	"runtime/debug"
	"strings"

	"errors"

	. "github.com/xeronith/diamante/contracts/operation"
	"github.com/xeronith/diamante/contracts/server"
	"github.com/xeronith/diamante/operation/binary"
	"github.com/xeronith/diamante/operation/text"
	"github.com/xeronith/diamante/protobuf"
)

func (server *baseServer) createServerError(developmentEnvironment bool, status int32, err error) server.IServerError {
	var (
		message     = ""
		description = ""
	)

	if err != nil {
		message = err.Error()
		if developmentEnvironment {
			description = fmt.Sprintf(string(debug.Stack()))
		}
	}

	description = server.localizer.Get(message)

	pattern := "Data too long for column '"
	index := strings.Index(message, pattern)
	if description == "" && index >= 0 {
		description = message[index+len(pattern):]
		index = strings.Index(description, "' at")
		if index >= 0 {
			description = description[:index]
		}

		description = fmt.Sprintf("%s_too_long", description)
		message = fmt.Sprintf("ERROR_MESSAGE_%s", strings.ToUpper(description))
	}

	return &protobuf.ServerError{
		Message:     message,
		Description: description,
	}
}

func (server *baseServer) serverError(id uint64, status int32, _error error, nonBinary bool, apiVersion, serverVersion, clientVersion int32) IOperationResult {
	serverError := server.createServerError(server.Configuration().IsDevelopmentEnvironment(), status, _error)

	var result IOperationResult
	if nonBinary {
		data, err := server.textSerializer.Serialize(serverError)
		var factory = text.CreateTextOperationResult
		if err != nil {
			result = factory(id, InternalServerError, ERROR, "", apiVersion, serverVersion, clientVersion, 0)
		} else {
			result = factory(id, status, ERROR, data, apiVersion, serverVersion, clientVersion, 0)
		}
	} else {
		data, err := server.binarySerializer.Serialize(serverError)
		var factory = binary.CreateBinaryOperationResult
		if err != nil {
			result = factory(id, InternalServerError, ERROR, nil, apiVersion, serverVersion, clientVersion, 0)
		} else {
			result = factory(id, status, ERROR, data, apiVersion, serverVersion, clientVersion, 0)
		}
	}

	return result
}

func (server *baseServer) internalServerError(id uint64, _error error, nonBinary bool, apiVersion, serverVersion, clientVersion int32) IOperationResult {
	if _error == nil {
		return server.serverError(id, InternalServerError, errors.New("internal server error"), nonBinary, apiVersion, serverVersion, clientVersion)
	} else {
		return server.serverError(id, InternalServerError, _error, nonBinary, apiVersion, serverVersion, clientVersion)
	}
}

func (server *baseServer) unauthorized(id uint64, _error error, nonBinary bool, apiVersion, serverVersion, clientVersion int32) IOperationResult {
	if _error == nil {
		return server.serverError(id, Unauthorized, errors.New("unauthorized"), nonBinary, apiVersion, serverVersion, clientVersion)
	} else {
		return server.serverError(id, Unauthorized, _error, nonBinary, apiVersion, serverVersion, clientVersion)
	}
}

func (server *baseServer) notImplemented(id uint64, _error error, nonBinary bool, apiVersion, serverVersion, clientVersion int32) IOperationResult {
	if _error == nil {
		return server.serverError(id, NotImplemented, errors.New("not implemented"), nonBinary, apiVersion, serverVersion, clientVersion)
	} else {
		return server.serverError(id, NotImplemented, _error, nonBinary, apiVersion, serverVersion, clientVersion)
	}
}

func (server *baseServer) serviceUnavailable(id uint64, _error error, nonBinary bool, apiVersion, serverVersion, clientVersion int32) IOperationResult {
	if _error == nil {
		return server.serverError(id, ServiceUnavailable, errors.New("service unavailable due to system maintenance"), nonBinary, apiVersion, serverVersion, clientVersion)
	} else {
		return server.serverError(id, ServiceUnavailable, _error, nonBinary, apiVersion, serverVersion, clientVersion)
	}
}
