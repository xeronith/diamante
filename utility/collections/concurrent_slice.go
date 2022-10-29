package collections

import (
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentSlice struct {
	sync.RWMutex
	collection []ISystemObject
}

func NewConcurrentSlice() ISlice {
	return &concurrentSlice{
		collection: make([]ISystemObject, 0),
	}
}

func (slice *concurrentSlice) Append(object ISystemObject) {
	slice.Lock()
	defer slice.Unlock()

	slice.collection = append(slice.collection, object)
}

func (slice *concurrentSlice) AppendIfNotExist(object ISystemObject) {
	slice.Lock()
	defer slice.Unlock()

	found := false
	for _, item := range slice.collection {
		if object == item {
			found = true
			break
		}
	}

	if !found {
		slice.collection = append(slice.collection, object)
	}
}

func (slice *concurrentSlice) Remove(object ISystemObject) {
	slice.Lock()
	defer slice.Unlock()
	temp := slice.collection[:0]
	for _, item := range slice.collection {
		if item != object {
			temp = append(temp, item)
		}
	}

	slice.collection = temp
}

func (slice *concurrentSlice) GetSize() int {
	slice.RLock()
	defer slice.RUnlock()
	return len(slice.collection)
}

func (slice *concurrentSlice) Clear() {
	slice.Lock()
	defer slice.Unlock()
	slice.collection = make([]ISystemObject, 0)
}

func (slice *concurrentSlice) ForEach(iterator func(int, ISystemObject)) {
	if iterator == nil {
		return
	}

	slice.RLock()
	defer slice.RUnlock()
	for index, value := range slice.collection {
		iterator(index, value)
	}
}

func (slice *concurrentSlice) ForEachParallel(iterator func(int, ISystemObject)) {
	if iterator == nil {
		return
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(slice.GetSize())

	slice.RLock()
	defer slice.RUnlock()
	for index, value := range slice.collection {
		go func(index int, value ISystemObject) {
			defer waitGroup.Done()
			iterator(index, value)
		}(index, value)
	}

	waitGroup.Wait()
}
