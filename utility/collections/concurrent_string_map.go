package collections

import (
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentStringMap struct {
	sync.RWMutex
	collection map[string]ISystemObject
}

func NewConcurrentStringMap() IStringMap {
	return &concurrentStringMap{
		collection: make(map[string]ISystemObject),
	}
}

func (_map *concurrentStringMap) Put(key string, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentStringMap) Get(key string) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentStringMap) Remove(key string) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentStringMap) Contains(key string) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentStringMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return len(_map.collection)
}

func (_map *concurrentStringMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[string]ISystemObject)
}

func (_map *concurrentStringMap) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentStringMap) ForEachKey(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentStringMap) ForEach(iterator func(string, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentStringMap) ForEachParallel(iterator func(string, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for key, value := range _map.collection {
		go func(key string, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentStringMap) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentStringMap) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
