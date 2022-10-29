package collections

import (
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentInt32Map struct {
	sync.RWMutex
	collection map[int32]ISystemObject
}

func NewConcurrentInt32Map() IInt32Map {
	return &concurrentInt32Map{
		collection: make(map[int32]ISystemObject),
	}
}

func (_map *concurrentInt32Map) Load(collection map[int32]ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = collection
}

func (_map *concurrentInt32Map) Put(key int32, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentInt32Map) Get(key int32) (ISystemObject, bool) {
	_map.Lock()
	defer _map.Unlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentInt32Map) Remove(key int32) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentInt32Map) Contains(key int32) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentInt32Map) GetSize() int {
	_map.Lock()
	defer _map.Unlock()
	return len(_map.collection)
}

func (_map *concurrentInt32Map) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int32]ISystemObject)
}

func (_map *concurrentInt32Map) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentInt32Map) ForEachKey(iterator func(int32)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentInt32Map) ForEach(iterator func(int32, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentInt32Map) ForEachParallel(iterator func(int32, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for key, value := range _map.collection {
		go func(key int32, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentInt32Map) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentInt32Map) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
