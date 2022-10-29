package io

import (
	"github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/logging"
)

type baseWriter struct {
	token            string
	logger           logging.ILogger
	binarySerializer IBinarySerializer
	textSerializer   ITextSerializer
	trafficRecorder  ITrafficRecorder
	onClosed         func()
	closed           bool
}

func createBaseWriter(server IServer, onClosed func()) baseWriter {
	return baseWriter{
		logger:           GetDefaultLogger(),
		binarySerializer: server.BinarySerializer(),
		textSerializer:   server.TextSerializer(),
		trafficRecorder:  server.TrafficRecorder(),
		onClosed:         onClosed,
		closed:           false,
	}
}
