package server

type IServerError interface {
	GetMessage() string
	GetDescription() string
}
