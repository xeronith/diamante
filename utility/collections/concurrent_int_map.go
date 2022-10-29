package collections

import (
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentIntMap struct {
	sync.RWMutex
	collection map[int]ISystemObject
}

func NewConcurrentIntMap() IIntMap {
	return &concurrentIntMap{
		collection: make(map[int]ISystemObject),
	}
}

func (_map *concurrentIntMap) Put(key int, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentIntMap) Get(key int) ISystemObject {
	_map.Lock()
	defer _map.Unlock()
	value := _map.collection[key]
	return value
}

func (_map *concurrentIntMap) Remove(key int) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentIntMap) Contains(key int) bool {
	return _map.Get(key) != nil
}

func (_map *concurrentIntMap) GetSize() int {
	_map.Lock()
	defer _map.Unlock()
	length := len(_map.collection)
	return length
}

func (_map *concurrentIntMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int]ISystemObject)
}

func (_map *concurrentIntMap) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.Lock()
	defer _map.Unlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentIntMap) ForEachKey(iterator func(int)) {
	if iterator == nil {
		return
	}

	_map.Lock()
	defer _map.Unlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentIntMap) ForEach(iterator func(int, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.Lock()
	defer _map.Unlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentIntMap) ForEachParallel(iterator func(int, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.Lock()
	defer _map.Unlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for key, value := range _map.collection {
		go func(key int, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentIntMap) Filter(predicate func(ISystemObject) bool) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)
	_map.Lock()
	defer _map.Unlock()
	for _, value := range _map.collection {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentIntMap) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)
	_map.Lock()
	defer _map.Unlock()
	for _, value := range _map.collection {
		result = append(result, predicate(value))
	}

	return result
}
