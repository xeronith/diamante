package io

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo"
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
)

type httpWriter struct {
	sync.RWMutex
	base         baseWriter
	context      echo.Context
	timestamp    time.Time
	opcodes      Opcodes
	secureCookie *securecookie.SecureCookie
}

func CreateHttpWriter(server IServer, context echo.Context, secureCookie *securecookie.SecureCookie) IWriter {
	return &httpWriter{
		base:         createBaseWriter(server, nil),
		context:      context,
		timestamp:    time.Now(),
		opcodes:      server.Opcodes(),
		secureCookie: secureCookie,
	}
}

func (writer *httpWriter) IsClosed() bool {
	return !writer.IsOpen()
}

func (writer *httpWriter) IsOpen() bool {
	writer.RLock()
	defer writer.RUnlock()

	return !writer.base.closed
}

func (writer *httpWriter) SetCookie(key, value string) {
	encoded, err := writer.secureCookie.Encode(key, value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  key,
			Value: encoded,
			Path:  "/",
			// HttpOnly: true,
		}
		http.SetCookie(writer.context.Response().Writer, cookie)
	}
}

func (writer *httpWriter) GetCookie(key string) string {
	var value string
	if cookie, err := writer.context.Request().Cookie(key); err == nil {
		if err = writer.secureCookie.Decode(key, cookie.Value, &value); err == nil {
			return value
			// fmt.Fprintln(context.Response().Writer, value)
		} else {
			// return echo.NewHTTPError(http.StatusUnauthorized, "")
		}
	}

	return ""
}

func (writer *httpWriter) SetToken(token string) {
	writer.base.token = token
}

func (writer *httpWriter) Write(operation IOperationResult) {
	defer writer.catch()

	writer.Lock()
	defer writer.Unlock()

	if writer.base.closed {
		writer.base.logger.Warning("HTTP WRITE ERROR: writer closed")
		return
	}

	switch operation.(type) {
	case IBinaryOperationResult:
		result := operation.(IBinaryOperationResult)
		data, err := writer.base.binarySerializer.Serialize(result.Container())

		action := writer.opcodes[result.Type()]
		serviceDuration := float64(result.ExecutionDuration().Microseconds()) / 1000
		pipelineDuration := float64(time.Since(writer.timestamp).Microseconds()) / 1000
		serverVersion := result.ServerVersion()

		writer.context.Response().Header().Add("Server-Timing", fmt.Sprintf("action;desc=\"%s\",version;desc=\"Build %d\",pipeline;desc=\"Pipeline\";dur=%f,service;desc=\"Service\";dur=%f", action, serverVersion, pipelineDuration, serviceDuration))

		if err == nil {
			if err := writer.context.Blob(int(result.Status()), echo.MIMEOctetStream, data); err == nil {
				writer.base.trafficRecorder.Record(BINARY_RESULT, data)
			} else {
				writer.base.logger.Error(fmt.Sprintf("HTTP/BOR WRITE ERROR: %s", err))
			}
		} else {
			writer.context.Response().Status = http.StatusInternalServerError
			if _, err := fmt.Fprintf(writer.context.Response(), err.Error()); err == nil {
				//TODO: writer.base.trafficRecorder.Record(BINARY_RESULT, data)
			} else {
				writer.base.logger.Error(fmt.Sprintf("HTTP/BERR WRITE ERROR: %s", err))
			}
		}
	case ITextOperationResult:
		result := operation.(ITextOperationResult)
		data, err := writer.base.textSerializer.Serialize(result.Container())
		if err == nil {
			if err := writer.context.String(int(result.Status()), data); err == nil {
				writer.base.trafficRecorder.Record(TEXT_RESULT, data)
			} else {
				writer.base.logger.Error(fmt.Sprintf("HTTP/TOR WRITE ERROR: %s", err))
			}
		} else {
			writer.context.Response().Status = http.StatusInternalServerError
			if _, err := fmt.Fprintf(writer.context.Response(), err.Error()); err == nil {
				//TODO: writer.base.trafficRecorder.Record(TEXT_RESULT, data)
			} else {
				writer.base.logger.Error(fmt.Sprintf("HTTP/TERR WRITE ERROR: %s", err))
			}
		}
	default:
		writer.base.logger.Error("HTTP WRITE ERROR: not supported")
	}
}

// noinspection GoStandardMethods
func (writer *httpWriter) WriteByte(code byte) {
	_ = code
	writer.base.logger.Error("not supported")
}

func (writer *httpWriter) WriteBytes(_type int, data []byte) {
	_ = _type
	_ = data
	writer.base.logger.Warning("HTTP WRITER: WriteBytes not implemented")
}

func (writer *httpWriter) End(operation IOperationResult) {
	writer.Write(operation)
}

func (writer *httpWriter) BinarySerializer() IBinarySerializer {
	return writer.base.binarySerializer
}

func (writer *httpWriter) TextSerializer() ITextSerializer {
	return writer.base.textSerializer
}

func (writer *httpWriter) Close() {
	writer.Lock()
	defer writer.Unlock()

	if writer.base.closed {
		return
	}

	writer.base.closed = true
	writer.context = nil

	if writer.base.onClosed != nil {
		writer.base.onClosed()
	}
}

func (writer *httpWriter) catch() {
	if reason := recover(); reason != nil {
		writer.Close()
		writer.base.logger.Panic(reason)
		// debug.PrintStack()
	}
}
