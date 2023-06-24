package io

import (
	"github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/logging"
)

type baseWriter struct {
	token       string
	logger      logging.ILogger
	contentType string
	serializer  ISerializer
	onClosed    func()
	closed      bool
}

func createBaseWriter(server IServer, onClosed func(), contentType string) baseWriter {
	return baseWriter{
		logger:      GetDefaultLogger(),
		contentType: contentType,
		serializer:  server.Serializers()[contentType],
		onClosed:    onClosed,
		closed:      false,
	}
}
