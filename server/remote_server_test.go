package server_test

import (
	"errors"
	"math/rand"
	"testing"

	. "github.com/xeronith/diamante/api"
	. "github.com/xeronith/diamante/client"
	. "github.com/xeronith/diamante/contracts/client"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/model"
	. "github.com/xeronith/diamante/server"
	. "github.com/xeronith/diamante/utility/concurrent"
)

const (
	activeUrl  = "wss://api.domain.com:7070"
	passiveUrl = "https://api.domain.com:7080"
)

func TestRemoteServer_Load(test *testing.T) {
	const totalRequests = 100

	tasks := NewAsyncTaskPool()

	for i := 0; i < totalRequests; i++ {
		tasks.Submit(func() {
			channel := make(chan IBinaryOperationResult, 1)

			active := rand.Intn(2) >= 1

			var client IClient
			callback := func(data IBinaryOperationResult) {
				channel <- data
			}

			var url string
			if active {
				client = CreateWebSocketClient(callback)
				url = activeUrl
			} else {
				client = CreateHttpClient(callback)
				url = passiveUrl
			}

			if err := client.Connect(url, "DEFAULT-TOKEN"); err != nil {
				test.Fail()
			}

			input := CreateRandomSample()
			if err := client.Send(1, ECHO, input); err != nil {
				test.Fail()
			}

			result := <-channel

			if err := client.Disconnect(); err != nil {
				test.Fail()
			}

			output := NewSample()
			if err := DefaultBinarySerializer.Deserialize(result.Payload(), output); err != nil {
				test.Fail()
			}

			if output.StringField != input.StringField {
				test.Fail()
			}
		})
	}

	tasks.Run()
	tasks.Join()
}

func TestRemoteServer_Passive(test *testing.T) {
	if err := remoteRoundTrip(NewHttpClient(), false); err != nil {
		test.Fatal(err)
	}
}

func TestRemoteServer_Active(test *testing.T) {
	if err := remoteRoundTrip(NewWebSocketClient(), true); err != nil {
		test.Fatal(err)
	}
}

func remoteRoundTrip(client IClient, active bool) error {
	channel := make(chan IBinaryOperationResult, 1)

	client.SetBinaryOperationResultListener(func(data IBinaryOperationResult) {
		channel <- data
	})

	var url string
	if active {
		url = activeUrl
	} else {
		url = passiveUrl
	}

	if err := client.Connect(url, "DEFAULT-TOKEN"); err != nil {
		return err
	}

	input := CreateRandomSample()
	if err := client.Send(1, ECHO, input); err != nil {
		return err
	}

	result := <-channel

	if err := client.Disconnect(); err != nil {
		return err
	}

	output := NewSample()
	if err := DefaultBinarySerializer.Deserialize(result.Payload(), output); err != nil {
		return err
	}

	if output.StringField != input.StringField {
		return errors.New("comparison failed")
	}

	return nil
}
