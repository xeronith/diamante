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
	var serializer ISerializer
	if contentSerializer, ok := server.Serializers()[contentType]; !ok {
		serializer = server.Serializers()["application/octet-stream"]
	} else {
		serializer = contentSerializer
	}

	return baseWriter{
		logger:      GetDefaultLogger(),
		contentType: contentType,
		serializer:  serializer,
		onClosed:    onClosed,
		closed:      false,
	}
}
