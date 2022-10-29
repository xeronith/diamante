package http

import (
	"net/http"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
)

type (
	IServerDispatcher interface {
		Logger() ILogger
		Serializer() IBinarySerializer
		Serialize(Pointer) ([]byte, error)
		Deserialize([]byte, Pointer) error
		OnActorBinaryData(IActor, []byte) IOperationResult
		Request() *http.Request
		Response() http.ResponseWriter
		Redirect(string)
		Query(string) string
		Param(string) string
		RemoteAddr() string
		UserAgent() string
	}

	HttpHandlerFunc func(IServerDispatcher) error

	IHttpHandler interface {
		Method() string
		Path() string
		HandlerFunc() HttpHandlerFunc
	}

	IHttpHandlerFactory interface {
		Handlers() []IHttpHandler
	}
)
