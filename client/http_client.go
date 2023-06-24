package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"

	. "github.com/xeronith/diamante/contracts/client"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/operation"
	. "github.com/xeronith/diamante/serialization"
	. "github.com/xeronith/diamante/utility/reflection"
)

type httpClient struct {
	baseClient
	internalClient *http.Client
}

func NewHttpClient() IClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &httpClient{
		baseClient: baseClient{},
		internalClient: &http.Client{
			Jar: jar,
		},
	}
}

func CreateHttpClient(listener func(IOperationResult)) IClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &httpClient{
		baseClient: baseClient{
			operationResultListener: listener,
		},
		internalClient: &http.Client{
			Jar: jar,
		},
	}
}

// noinspection GoUnusedExportedFunction
func CreateDistinctHttpClient(version int32, name string, listener func(IOperationResult)) IClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &httpClient{
		baseClient: baseClient{
			name:                    name,
			version:                 version,
			operationResultListener: listener,
		},
		internalClient: &http.Client{
			Jar: jar,
		},
	}
}

func (client *httpClient) IsActive() bool {
	return false
}

func (client *httpClient) SetName(name string) {
	client.name = name
}

func (client *httpClient) SetToken(token string) {
	client.token = token
}

func (client *httpClient) SetVersion(version int32) {
	client.version = version
}

func (client *httpClient) SetApiVersion(apiVersion int32) {
	client.apiVersion = apiVersion
}

func (client *httpClient) Connect(endpoint string, token string) error {
	client.endpoint = endpoint
	client.token = token
	client.serializer = NewProtobufSerializer()

	if client.connectionEstablished != nil {
		client.connectionEstablished(client)
	}

	return nil
}

func (client *httpClient) Send(id uint64, operation uint64, payload Pointer) error {
	if !IsPointer(payload) {
		return errors.New("payload should be a pointer")
	}

	var (
		data []byte
		err  error
	)

	operationRequest := CreateOperationRequest(id, operation, client.name, client.version, client.apiVersion, client.token, nil)
	if err := operationRequest.Load(payload, client.serializer); err != nil {
		return err
	}

	data, err = client.serializer.Serialize(operationRequest.Container())
	if err != nil {
		return err
	}

	return client.send(data)
}

func (client *httpClient) send(data []byte) error {
	request, err := http.NewRequest("POST", client.endpoint, bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	response, err := client.internalClient.Do(request)
	if err != nil {
		return err
	}

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	operationResult := NewOperationResult()
	err = client.serializer.Deserialize(buffer, operationResult.Container())
	if err != nil {
		client.operationResultListener(nil)
		return err
	}

	client.operationResultListener(operationResult)

	return nil
}

func (client *httpClient) Disconnect() error {
	return nil
}
