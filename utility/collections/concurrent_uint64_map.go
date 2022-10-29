package collections

import (
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentUInt64Map struct {
	sync.RWMutex
	collection map[uint64]ISystemObject
}

func NewConcurrentUInt64Map() IUInt64Map {
	return &concurrentUInt64Map{
		collection: make(map[uint64]ISystemObject),
	}
}

func (_map *concurrentUInt64Map) Put(key uint64, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentUInt64Map) Get(key uint64) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentUInt64Map) Remove(key uint64) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentUInt64Map) Contains(key uint64) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentUInt64Map) GetSize() int64 {
	_map.RLock()
	defer _map.RUnlock()
	length := len(_map.collection)
	return int64(length)
}

func (_map *concurrentUInt64Map) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[uint64]ISystemObject)
}

func (_map *concurrentUInt64Map) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentUInt64Map) ForEachKey(iterator func(uint64)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentUInt64Map) ForEach(iterator func(uint64, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentUInt64Map) ForEachParallel(iterator func(uint64, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for key, value := range _map.collection {
		go func(key uint64, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentUInt64Map) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentUInt64Map) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
