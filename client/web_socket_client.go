package client

import (
	"crypto/tls"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/fasthttp-contrib/websocket"
	. "github.com/xeronith/diamante/contracts/client"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
	"github.com/xeronith/diamante/operation/binary"
	"github.com/xeronith/diamante/server"
	"github.com/xeronith/diamante/utility/reflection"
)

type webSocketClient struct {
	mutex      sync.Mutex
	base       baseClient
	connection *websocket.Conn
}

func NewWebSocketClient() IClient {
	return &webSocketClient{
		base: baseClient{},
	}
}

func CreateWebSocketClient(listener func(IBinaryOperationResult)) IClient {
	return &webSocketClient{
		base: baseClient{
			binaryOperationResultListener: listener,
		},
	}
}

// noinspection GoUnusedExportedFunction
func CreateDistinctWebSocketClient(version int32, name string, listener func(IBinaryOperationResult)) IClient {
	return &webSocketClient{
		base: baseClient{
			name:                          name,
			version:                       version,
			binaryOperationResultListener: listener,
		},
	}
}

func (client *webSocketClient) IsActive() bool {
	return true
}

func (client *webSocketClient) SetName(name string) {
	client.base.name = name
}

func (client *webSocketClient) SetToken(token string) {
	client.base.token = token
}

func (client *webSocketClient) SetVersion(version int32) {
	client.base.version = version
}

func (client *webSocketClient) SetApiVersion(apiVersion int32) {
	client.base.apiVersion = apiVersion
}

func (client *webSocketClient) BinarySerializer() IBinarySerializer {
	return client.base.binarySerializer
}

func (client *webSocketClient) SetBinaryOperationResultListener(listener func(IBinaryOperationResult)) {
	client.base.binaryOperationResultListener = listener
}

func (client *webSocketClient) OnConnectionEstablished(callback func(IClient)) {
	client.base.connectionEstablished = callback
}

type fakeSessionCache struct{}

func (fakeSessionCache) Get(sessionKey string) (*tls.ClientSessionState, bool) {
	return nil, false
}

func (fakeSessionCache) Put(sessionKey string, cs *tls.ClientSessionState) {
	// no-op
}

func (client *webSocketClient) Connect(endpoint string, token string) error {
	client.base.token = token
	client.base.endpoint = endpoint
	client.base.binarySerializer = server.DefaultBinarySerializer

	var (
		err     error
		headers map[string][]string
	)

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	config.ClientSessionCache = fakeSessionCache{}
	websocket.DefaultDialer.TLSClientConfig = config

	client.connection, _, err = websocket.DefaultDialer.Dial(client.base.endpoint, headers)
	if err != nil {
		return err
	}

	quit := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				rand.Seed(time.Now().UnixNano())
				duration := time.Millisecond * time.Duration(3000+rand.Intn(2000))
				time.Sleep(duration)

				client.mutex.Lock()
				err := client.connection.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second*4))
				client.mutex.Unlock()

				if err != nil {
					if client.base.endpoint != "" {
						log.Println("CLIENT CONTROL ERROR: ", err)
					}

					break
				}
			}
		}
	}()

	go func() {
		for {
			_, message, err := client.connection.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					log.Println("CLIENT SOCKET READ ERROR: ", err)
				}

				quit <- true
				return
			}

			if client.base.binaryOperationResultListener != nil {
				operationResult := binary.NewBinaryOperationResult()
				err = client.base.binarySerializer.Deserialize(message, operationResult.Container())
				if err != nil {
					log.Println("SOCKET DATA DESERIALIZATION ERROR: ", err)
				} else {
					client.base.binaryOperationResultListener(operationResult)
				}
			}
		}
	}()

	if client.base.connectionEstablished != nil {
		go client.base.connectionEstablished(client)
	}

	return nil
}

func (client *webSocketClient) Send(id uint64, operation uint64, payload Pointer) error {
	if !reflection.IsPointer(payload) {
		return errors.New("payload should be a pointer")
	}

	var (
		data []byte
		err  error
	)

	binaryOperationRequest := binary.CreateBinaryOperationRequest(id, operation, client.base.name, client.base.version, client.base.apiVersion, client.base.token, nil)
	if err := binaryOperationRequest.Load(payload, client.base.binarySerializer); err != nil {
		return err
	}

	data, err = client.base.binarySerializer.Serialize(binaryOperationRequest.Container())
	if err != nil {
		return err
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	return client.connection.WriteMessage(websocket.TextMessage, data)
}

func (client *webSocketClient) Disconnect() error {
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")

	client.mutex.Lock()
	defer client.mutex.Unlock()
	result := client.connection.WriteMessage(websocket.CloseMessage, message)
	if result != nil {
		return result
	}

	client.base.endpoint = ""
	return nil
}
