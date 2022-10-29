package concurrent_test

import (
	"testing"

	"github.com/xeronith/diamante/utility/concurrent"
)

func TestAsyncTask(t *testing.T) {
	channel := make(chan bool)

	task := concurrent.NewAsyncTask(func() {

		// panic("What happened")

		channel <- true
	})

	task.Run()

	<-channel
}
