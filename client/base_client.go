package client

import (
	. "github.com/xeronith/diamante/contracts/client"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
)

type baseClient struct {
	serializer              ISerializer
	operationResultListener func(IOperationResult)
	connectionEstablished   func(IClient)
	endpoint                string
	token                   string
	version                 int32
	apiVersion              int32
	name                    string
}

func (client *baseClient) Serializer() ISerializer {
	return client.serializer
}

func (client *baseClient) SetOperationResultListener(listener func(IOperationResult)) {
	client.operationResultListener = listener
}

func (client *baseClient) OnConnectionEstablished(callback func(IClient)) {
	client.connectionEstablished = callback
}
