package concurrent

import (
	"runtime"
	"sync"
)

type asyncTaskPool struct {
	waitGroup    sync.WaitGroup
	lockOSThread bool
	tasks        []func()
}

func NewAsyncTaskPool() IAsyncTaskPool {
	return &asyncTaskPool{}
}

func CreateAsyncTaskPool(lockOSThread bool) IAsyncTaskPool {
	return &asyncTaskPool{
		lockOSThread: lockOSThread,
	}
}

func (pool *asyncTaskPool) Submit(tasks ...func()) IAsyncTaskPool {
	for _, task := range tasks {
		if pool.tasks == nil {
			pool.tasks = make([]func(), 0)
		}

		pool.tasks = append(pool.tasks, task)
	}

	return pool
}

func (pool *asyncTaskPool) Run() IAsyncTaskPool {
	pool.waitGroup = sync.WaitGroup{}
	pool.waitGroup.Add(len(pool.tasks))

	for i := 0; i < len(pool.tasks); i++ {
		if pool.lockOSThread {
			go func(i int) {
				runtime.LockOSThread()
				defer func() {
					pool.waitGroup.Done()
					runtime.UnlockOSThread()
				}()

				pool.tasks[i]()
			}(i)
		} else {
			go func(i int) {
				defer pool.waitGroup.Done()
				pool.tasks[i]()
			}(i)
		}
	}

	return pool
}

func (pool *asyncTaskPool) Join() {
	pool.waitGroup.Wait()
}
