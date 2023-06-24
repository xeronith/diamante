package io

import (
	"fmt"
	"sync"
	"time"

	. "github.com/gorilla/websocket"
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
)

type webSocketWriter struct {
	sync.RWMutex
	base       baseWriter
	connection *Conn
}

func CreateWebSocketWriter(server IServer, connection *Conn, onClosed func()) IWriter {
	return &webSocketWriter{
		base:       createBaseWriter(server, onClosed, "application/octet-stream"),
		connection: connection,
	}
}

func (writer *webSocketWriter) ContentType() string {
	return writer.base.contentType
}

func (writer *webSocketWriter) IsClosed() bool {
	return !writer.IsOpen()
}

func (writer *webSocketWriter) IsOpen() bool {
	writer.RLock()
	defer writer.RUnlock()

	return !writer.base.closed
}

func (writer *webSocketWriter) SetSecureCookie(_, _ string) {
	writer.base.logger.Error("WEBSOCKET WRITER: SetSecureCookie not supported")
}

func (writer *webSocketWriter) GetSecureCookie(_ string) string {
	writer.base.logger.Error("WEBSOCKET WRITER: GetSecureCookie not supported")
	return ""
}

func (writer *webSocketWriter) SetAuthCookie(_ string) {
	writer.base.logger.Error("WEBSOCKET WRITER: SetAuthCookie not supported")
}

func (writer *webSocketWriter) GetAuthCookie() string {
	writer.base.logger.Error("WEBSOCKET WRITER: GetAuthCookie not supported")
	return ""
}

func (writer *webSocketWriter) SetToken(token string) {
	writer.base.token = token
}

func (writer *webSocketWriter) Write(result IOperationResult) {
	defer writer.catch()

	if data, err := writer.base.serializer.Serialize(result.Container()); err != nil {
		writer.base.logger.Error(fmt.Sprintf("SOCKET/OR SERIALIZATION ERROR {%s}: %s", writer.base.token, err))
	} else {
		writer.WriteBytes(data)
	}
}

func (writer *webSocketWriter) WriteBytes(data []byte) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if duration > time.Millisecond*75 {
			writer.base.logger.Alert(fmt.Sprintf("SOCKET/WRT {%s}: %s", writer.base.token, duration.String()))
		}
	}()

	// TODO: What if data is nil? (Check other places too)
	var closed bool
	func() {
		defer writer.catch()

		writer.Lock()
		defer writer.Unlock()

		if writer.base.closed {
			writer.base.logger.Warning("SOCKET WRITE ERROR: writer closed")
			return
		}

		if err := writer.connection.SetWriteDeadline(time.Now().Add(time.Second * 4)); err != nil {
			//TODO: Probably should close the writer but according to implementation the method above won't return any error
			writer.base.logger.Error(fmt.Sprintf("SOCKET/SWD ERROR {%s}: %s", writer.base.token, err))
			return
		}

		if err := writer.connection.WriteMessage(BinaryMessage, data); err != nil {
			closed = true
			writer.base.closed = true
			writer.base.logger.Error(fmt.Sprintf("SOCKET/OR WRITE ERROR {%s}: %s", writer.base.token, err))
		}
	}()

	if closed {
		writer.finalize()
	}
}

func (writer *webSocketWriter) WriteByte(code byte) error {
	defer writer.catch()

	writer.Lock()
	defer writer.Unlock()

	if writer.base.closed {
		writer.base.logger.Warning("SOCKET WRITE ERROR: writer closed")
		return nil
	}

	if err := writer.connection.WriteMessage(BinaryMessage, []byte{code}); err != nil {
		writer.base.closed = true
		writer.base.logger.Error(fmt.Sprintf("SOCKET/SIG WRITE ERROR: %s", err))
	}

	return nil
}

func (writer *webSocketWriter) End(operation IOperationResult) {
	defer func() {
		defer writer.catch()

		writer.Lock()
		defer writer.Unlock()

		if writer.base.closed {
			writer.base.logger.Warning("SOCKET WRITE ERROR: writer closed")
			return
		}

		if err := writer.connection.Close(); err != nil {
			writer.base.logger.Error(fmt.Sprintf("SOCKET CLOSE ERROR: %s", err))
		}

		writer.base.closed = true
	}()

	writer.Write(operation)
}

func (writer *webSocketWriter) Serializer() ISerializer {
	return writer.base.serializer
}

func (writer *webSocketWriter) Close() {
	writer.Lock()
	if writer.base.closed {
		writer.Unlock()
		return
	}

	writer.base.closed = true
	writer.Unlock()

	writer.finalize()
}

func (writer *webSocketWriter) finalize() {
	if err := writer.connection.Close(); err != nil {
		// writer.base.logger.Error(fmt.Sprintf("SOCKET/DFR CLOSE ERROR: %s", err))
		_ = err
	}

	if writer.base.onClosed != nil {
		writer.base.onClosed()
	}
}

func (writer *webSocketWriter) OnClosed(callback func()) {
	writer.base.onClosed = callback
}

func (writer *webSocketWriter) catch() {
	if reason := recover(); reason != nil {
		writer.Close()
		writer.base.logger.Panic(reason)
		// debug.PrintStack()
	}
}
