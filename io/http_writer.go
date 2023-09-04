package io

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"

	_ "embed"
)

var (
	//go:embed quotes.json
	quotesData []byte
	quotes     = &[]*struct{ Author, Quote string }{}
)

func init() {
	_ = json.Unmarshal(quotesData, quotes)
}

type httpWriter struct {
	sync.RWMutex
	base         baseWriter
	context      echo.Context
	timestamp    time.Time
	opcodes      Opcodes
	secureCookie *securecookie.SecureCookie
}

func CreateHttpWriter(
	server IServer,
	context echo.Context,
	secureCookie *securecookie.SecureCookie,
) IWriter {
	return &httpWriter{
		base: createBaseWriter(
			server, nil,
			context.Request().Header.Get("Content-Type"),
		),
		context:      context,
		timestamp:    time.Now(),
		opcodes:      server.Opcodes(),
		secureCookie: secureCookie,
	}
}

func (writer *httpWriter) ContentType() string {
	return writer.base.contentType
}

func (writer *httpWriter) IsClosed() bool {
	return !writer.IsOpen()
}

func (writer *httpWriter) IsOpen() bool {
	writer.RLock()
	defer writer.RUnlock()

	return !writer.base.closed
}

func (writer *httpWriter) SetSecureCookie(key, value string) {
	encoded, err := writer.secureCookie.Encode(key, value)
	if err == nil {
		cookie := &http.Cookie{
			Name:     key,
			Value:    encoded,
			Path:     "/",
			MaxAge:   7 * 24 * 60 * 60,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(writer.context.Response().Writer, cookie)
	}
}

func (writer *httpWriter) GetSecureCookie(key string) string {
	var value string
	if cookie, err := writer.context.Request().Cookie(key); err == nil {
		if err := writer.secureCookie.Decode(key, cookie.Value, &value); err == nil {
			return value
		}
	}

	return ""
}

func (writer *httpWriter) SetAuthCookie(token string) {
	writer.SetSecureCookie("Diamante", token)
}

func (writer *httpWriter) GetAuthCookie() string {
	return writer.GetSecureCookie("Diamante")
}

func (writer *httpWriter) SetToken(token string) {
	writer.base.token = token
}

func (writer *httpWriter) Write(result IOperationResult) {
	defer writer.catch()

	writer.Lock()
	defer writer.Unlock()

	if writer.base.closed {
		writer.base.logger.Warning("HTTP WRITE ERROR: writer closed")
		return
	}

	serviceDuration := float64(result.ExecutionDuration().Microseconds()) / 1000
	pipelineDuration := float64(time.Since(writer.timestamp).Microseconds()) / 1000

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	alloc := memStats.Alloc / 1024 / 1024
	sys := memStats.Sys / 1024 / 1024
	miss, hit := result.Stat()

	writer.context.Response().Header().Add("X-Powered-By", "Magic")
	writer.context.Response().Header().Add("X-Request-Timestamp", fmt.Sprintf("%d", result.Id()))
	writer.context.Response().Header().Add("X-Response-Signature", result.Signature())
	writer.context.Response().Header().Add("Server-Timing", fmt.Sprintf("id;desc=\"Opcode: 0x%X, Alloc: %d MiB, Sys: %d MiB, GC: %d, Miss: %d, Hit: %d\",pipeline;desc=\"Pipeline\";dur=%f,service;desc=\"Service\";dur=%f", result.Type(), alloc, sys, memStats.NumGC, miss, hit, pipelineDuration, serviceDuration))

	if serviceDuration == 0 {
		writer.context.Response().Header().Add("X-Turbo", "On")
	}

	quote := (*quotes)[rand.Intn(len(*quotes))]
	writer.context.Response().Header().Add("X-Quote", fmt.Sprintf("%s: %s", quote.Author, quote.Quote))

	signature := writer.context.Request().Header.Get("X-Request-Signature")
	if signature != "" && signature == result.Signature() {
		writer.context.String(http.StatusNotModified, "")
		return
	}

	data, err := writer.base.serializer.Serialize(result.Container())
	if err == nil {
		if result.ContentType() == "application/json" {
			var response, responsePayload map[string]interface{}
			_ = json.Unmarshal(data, &response)
			_ = json.Unmarshal(result.Payload(), &responsePayload)
			response["payload"] = responsePayload
			data, _ = json.Marshal(response)
		}

		if err := writer.context.Blob(int(result.Status()), result.ContentType(), data); err != nil {
			writer.base.logger.Error(fmt.Sprintf("HTTP/OR WRITE ERROR: %s", err))
		}
	} else {
		writer.context.Response().Status = http.StatusInternalServerError
		if _, err := fmt.Fprint(writer.context.Response(), err.Error()); err != nil {
			writer.base.logger.Error(fmt.Sprintf("HTTP/ERR WRITE ERROR: %s", err))
		}
	}
}

func (writer *httpWriter) WriteByte(_ byte) error {
	writer.base.logger.Error("HTTP WRITER: WriteByte not supported")
	return nil
}

func (writer *httpWriter) WriteBytes(_ []byte) {
	writer.base.logger.Error("HTTP WRITER: WriteBytes not supported")
}

func (writer *httpWriter) End(operation IOperationResult) {
	writer.Write(operation)
}

func (writer *httpWriter) Serializer() ISerializer {
	return writer.base.serializer
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
