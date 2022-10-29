package collections

import (
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentInt64Map struct {
	sync.RWMutex
	collection map[int64]ISystemObject
}

func NewConcurrentInt64Map() IInt64Map {
	return &concurrentInt64Map{
		collection: make(map[int64]ISystemObject),
	}
}

func (_map *concurrentInt64Map) Load(collection map[int64]ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = collection
}

func (_map *concurrentInt64Map) Put(key int64, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentInt64Map) Get(key int64) (ISystemObject, bool) {
	_map.Lock()
	defer _map.Unlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentInt64Map) Remove(key int64) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentInt64Map) Contains(key int64) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentInt64Map) GetSize() int {
	_map.Lock()
	defer _map.Unlock()
	return len(_map.collection)
}

func (_map *concurrentInt64Map) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int64]ISystemObject)
}

func (_map *concurrentInt64Map) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentInt64Map) ForEachKey(iterator func(int64)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentInt64Map) ForEach(iterator func(int64, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentInt64Map) ForEachParallel(iterator func(int64, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for key, value := range _map.collection {
		go func(key int64, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentInt64Map) Filter(predicate func(ISystemObject) bool) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)
	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentInt64Map) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)
	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		result = append(result, predicate(value))
	}

	return result
}
