package server

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	. "github.com/xeronith/diamante/contracts/scheduling"
	"github.com/xeronith/diamante/logging"
)

type scheduler struct {
	futuresChannel chan *future
	futures        map[string]*future
}

type future struct {
	id        string
	duration  time.Duration
	timeout   time.Time
	callback  func()
	recurring bool
	done      bool
}

func newScheduler() IScheduler {
	return &scheduler{
		futuresChannel: make(chan *future, 1000),
	}
}

func (scheduler *scheduler) Start() {
	scheduler.futures = make(map[string]*future)
	ticker := time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case future := <-scheduler.futuresChannel:
			scheduler.futures[future.id] = future
		case <-ticker.C:
			executed := make([]string, 0)
			for id, _future := range scheduler.futures {
				if _future.done {
					executed = append(executed, id)
				} else if _future.timeout.Before(time.Now()) {
					if _future.recurring {
						_future.timeout = time.Now().Add(_future.duration)
					} else {
						_future.done = true
					}
					go func(_future *future) {
						defer catch()
						_future.callback()
					}(_future)
				}
			}

			for _, id := range executed {
				delete(scheduler.futures, id)
			}
		}
	}
}

func catch() {
	if reason := recover(); reason != nil {
		logging.GetDefaultLogger().Panic(fmt.Sprintf("SCHEDULER: %s", reason))
	}
}

func (scheduler *scheduler) SetTimeout(callback func(), timeout time.Duration) string {
	return scheduler.set(callback, timeout, false)
}

func (scheduler *scheduler) SetInterval(callback func(), timeout time.Duration) string {
	return scheduler.set(callback, timeout, true)
}

func (scheduler *scheduler) Cancel(id string) {
	if _future, exists := scheduler.futures[id]; exists {
		_future.done = true
	}
}

func (scheduler *scheduler) set(callback func(), duration time.Duration, recurring bool) string {
	_future := future{
		id:        uuid.NewV4().String(),
		duration:  duration,
		timeout:   time.Now().Add(duration),
		callback:  callback,
		recurring: recurring,
		done:      false,
	}

	scheduler.futuresChannel <- &_future
	return _future.id
}

func (server *defaultServer) startServerScheduler() {
	server.Scheduler().Start()
}
