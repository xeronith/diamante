package http

import (
	"net/http"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/server"
	. "github.com/xeronith/diamante/contracts/system"
)

type dispatcher struct {
	server   IServer
	response http.ResponseWriter
	request  *http.Request
	query    func(string) string
	param    func(string) string
	ip       string
}

func NewDispatcher(
	server IServer,
	response http.ResponseWriter,
	request *http.Request,
	query func(string) string,
	param func(string) string,
	ip string,
) IServerDispatcher {
	return &dispatcher{
		server:   server,
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

func (dispatcher *dispatcher) Serializer() IBinarySerializer {
	return dispatcher.server.BinarySerializer()
}

func (dispatcher *dispatcher) Serialize(object Pointer) ([]byte, error) {
	return dispatcher.server.BinarySerializer().Serialize(object)
}

func (dispatcher *dispatcher) Deserialize(data []byte, object Pointer) error {
	return dispatcher.server.BinarySerializer().Deserialize(data, object)
}

func (dispatcher *dispatcher) OnActorBinaryData(actor IActor, data []byte) IOperationResult {
	return dispatcher.server.OnActorBinaryData(actor, data)
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

func (dispatcher *dispatcher) RemoteAddr() string {
	return dispatcher.ip
}

func (dispatcher *dispatcher) UserAgent() string {
	return dispatcher.request.UserAgent()
}
