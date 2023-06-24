package io

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
)

type IWriter interface {
	IsClosed() bool
	IsOpen() bool
	ContentType() string
	SetSecureCookie(string, string)
	GetSecureCookie(string) string
	SetAuthCookie(string)
	GetAuthCookie() string
	SetToken(string)
	Write(IOperationResult)
	WriteByte(byte) error
	WriteBytes([]byte)
	End(IOperationResult)
	Serializer() ISerializer
	Close()
}
