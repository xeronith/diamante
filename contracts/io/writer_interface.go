package io

import (
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
)

type IWriter interface {
	IsClosed() bool
	IsOpen() bool
	SetCookie(string, string)
	GetCookie(string) string
	SetToken(string)
	Write(IOperationResult)
	WriteByte(byte)
	WriteBytes(int, []byte)
	End(IOperationResult)
	BinarySerializer() IBinarySerializer
	TextSerializer() ITextSerializer
	Close()
}
