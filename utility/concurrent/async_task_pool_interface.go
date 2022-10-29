package concurrent

type IAsyncTaskPool interface {
	Submit(tasks ...func()) IAsyncTaskPool
	Run() IAsyncTaskPool
	Join()
}
