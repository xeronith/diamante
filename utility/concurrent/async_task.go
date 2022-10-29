package concurrent

import . "github.com/xeronith/diamante/logging"

type asyncTask struct {
	runnable func()
}

var logger = GetDefaultLogger()

func NewAsyncTask(runnable func()) IAsyncTask {
	return &asyncTask{
		runnable: runnable,
	}
}

func catch() {
	if reason := recover(); reason != nil {
		logger.Panic(reason)
	}
}

func (task *asyncTask) Run() {
	go func() {
		defer catch()

		task.runnable()
	}()
}
