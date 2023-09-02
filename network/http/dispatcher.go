package http

import (
	"net/http"

	. "github.com/xeronith/diamante/actor"
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/serialization"
)

type dispatcher struct {
	server   IServer
	actor    IActor
	response http.ResponseWriter
	request  *http.Request
	query    func(string) string
	param    func(string) string
	ip       string
}

func NewDispatcher(
	server IServer,
	writer IWriter,
	response http.ResponseWriter,
	request *http.Request,
	query func(string) string,
	param func(string) string,
	ip string,
) IServerDispatcher {
	return &dispatcher{
		server: server,
		actor: CreateActor(
			writer,
			false,
			request.Header.Get("X-Request-Signature"),
			ip,
			request.UserAgent(),
		),
		response: response,
		request:  request,
		query:    query,
		param:    param,
		ip:       ip,
	}
}

func (dispatcher *dispatcher) Logger() ILogger {
	return dispatcher.server.Logger()
}

func (dispatcher *dispatcher) Actor() IActor {
	return dispatcher.actor
}

func (dispatcher *dispatcher) Serializer() ISerializer {
	return NewProtobufSerializer()
}

func (dispatcher *dispatcher) Serialize(object Pointer) ([]byte, error) {
	return dispatcher.Serializer().Serialize(object)
}

func (dispatcher *dispatcher) Deserialize(data []byte, object Pointer) error {
	return dispatcher.Serializer().Deserialize(data, object)
}

func (dispatcher *dispatcher) OnData(actor IActor, data []byte) IOperationResult {
	return dispatcher.server.OnData(actor, data)
}

func (dispatcher *dispatcher) Request() *http.Request {
	return dispatcher.request
}

func (dispatcher *dispatcher) Response() http.ResponseWriter {
	return dispatcher.response
}

func (dispatcher *dispatcher) Redirect(url string) {
	http.Redirect(dispatcher.response, dispatcher.request, url, http.StatusMovedPermanently)
}

func (dispatcher *dispatcher) Query(key string) string {
	return dispatcher.query(key)
}

func (dispatcher *dispatcher) Param(key string) string {
	return dispatcher.param(key)
}
