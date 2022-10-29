package client

import (
	. "github.com/xeronith/diamante/contracts/client"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
)

type baseClient struct {
	binarySerializer              IBinarySerializer
	binaryOperationResultListener func(IBinaryOperationResult)
	connectionEstablished         func(IClient)
	endpoint                      string
	token                         string
	version                       int32
	apiVersion                    int32
	name                          string
}

func (client *baseClient) BinarySerializer() IBinarySerializer {
	return client.binarySerializer
}

func (client *baseClient) SetBinaryOperationResultListener(listener func(IBinaryOperationResult)) {
	client.binaryOperationResultListener = listener
}

func (client *baseClient) OnConnectionEstablished(callback func(IClient)) {
	client.connectionEstablished = callback
}
