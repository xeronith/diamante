package http

import (
	. "github.com/xeronith/diamante/contracts/network/http"
)

type httpHandler struct {
	path        string
	method      string
	handlerFunc HttpHandlerFunc
}

func NewHttpHandler(path, method string, handlerFunc HttpHandlerFunc) IHttpHandler {
	return &httpHandler{
		path:        path,
		method:      method,
		handlerFunc: handlerFunc,
	}
}

func (handler *httpHandler) Path() string {
	return handler.path
}

func (handler *httpHandler) Method() string {
	return handler.method
}

func (handler *httpHandler) HandlerFunc() HttpHandlerFunc {
	return handler.handlerFunc
}
