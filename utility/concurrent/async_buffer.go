package concurrent

import (
	"sync"
	"time"
)

const (
	DefaultBufferSize       = 10000
	DefaultAutoFlushTimeout = 5 * time.Second
)

var autoFlushCallback func()

func StartSingleThreadScheduledExecutor() {
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})

	go func(autoFlushCallback func()) {
		for {
			select {
			case <-ticker.C:
				if autoFlushCallback != nil {
					autoFlushCallback()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(autoFlushCallback)
}

type asyncBuffer struct {
	sync.Mutex

	waitGroup         sync.WaitGroup
	frontBuffer       []interface{}
	bufferSize        int
	autoFlushTimeout  time.Duration
	autoFlushCallback func()
	lastUpdate        int64
	executionCallback func([]interface{})
}

func NewAsyncBuffer(executionCallback func([]interface{}), config ...int) *asyncBuffer {

	size := DefaultBufferSize
	if len(config) > 0 {
		size = config[0]
	}

	timeout := DefaultAutoFlushTimeout
	if len(config) > 1 {
		timeout = time.Duration(config[1]) * time.Millisecond
	}

	instance := &asyncBuffer{
		bufferSize:        size,
		autoFlushTimeout:  timeout,
		waitGroup:         sync.WaitGroup{},
		frontBuffer:       make([]interface{}, 0),
		lastUpdate:        0,
		executionCallback: executionCallback,
	}

	instance.autoFlushCallback = func() {
		if time.Duration(time.Now().UnixNano()-instance.lastUpdate) > instance.autoFlushTimeout {
			instance.Flush()
			autoFlushCallback = nil
		}
	}

	return instance
}

func (asyncBuffer *asyncBuffer) Submit(entity interface{}) {
	var backBuffer []interface{} = nil

	func() {
		asyncBuffer.Lock()
		defer asyncBuffer.Unlock()
		{
			asyncBuffer.frontBuffer = append(asyncBuffer.frontBuffer, entity)
			asyncBuffer.lastUpdate = time.Now().UnixNano()
			autoFlushCallback = asyncBuffer.autoFlushCallback

			if len(asyncBuffer.frontBuffer) == asyncBuffer.bufferSize {
				backBuffer = asyncBuffer.frontBuffer
				asyncBuffer.frontBuffer = make([]interface{}, 0)
				autoFlushCallback = nil
			}
		}
	}()

	if backBuffer != nil {
		asyncBuffer.process(backBuffer)
	}
}

func (asyncBuffer *asyncBuffer) Flush() {
	var backBuffer []interface{} = nil

	func() {
		asyncBuffer.Lock()
		defer asyncBuffer.Unlock()
		{
			backBuffer = asyncBuffer.frontBuffer
			asyncBuffer.frontBuffer = make([]interface{}, 0)
		}
	}()

	if backBuffer != nil && len(backBuffer) > 0 {
		asyncBuffer.process(backBuffer)
	}
}

func (asyncBuffer *asyncBuffer) Wait() {
	asyncBuffer.waitGroup.Wait()
}

func (asyncBuffer *asyncBuffer) process(buffer []interface{}) {
	asyncBuffer.waitGroup.Add(1)
	go func() {
		defer asyncBuffer.waitGroup.Done()
		asyncBuffer.executionCallback(buffer)
	}()
}
