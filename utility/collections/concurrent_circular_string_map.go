package collections

import (
	"sort"
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentCircularStringMap struct {
	sync.RWMutex
	current    int
	capacity   int
	keys       []string
	collection map[string]ISystemObject
}

func NewConcurrentCircularStringMap(capacity int) IOrderedStringMap {
	return &concurrentCircularStringMap{
		current:    -1,
		capacity:   capacity,
		keys:       make([]string, capacity),
		collection: make(map[string]ISystemObject, capacity),
	}
}

func (_map *concurrentCircularStringMap) Load(collection map[string]ISystemObject) {
	keys := make([]string, 0)
	for key := range collection {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(a, b int) bool {
		return keys[a] < keys[b]
	})

	_map.Lock()
	defer _map.Unlock()
	_map.keys = keys
	_map.collection = collection
}

func (_map *concurrentCircularStringMap) Put(key string, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	if _, exists := _map.collection[key]; !exists {
		_map.current = (_map.current + 1) % _map.capacity
		if _map.keys[_map.current] != "" {
			delete(_map.collection, _map.keys[_map.current])
		}

		_map.keys[_map.current] = key
	}

	_map.collection[key] = value
}

func (_map *concurrentCircularStringMap) Get(key string) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentCircularStringMap) Remove(key string) {
	panic("not_implemented")
}

func (_map *concurrentCircularStringMap) Contains(key string) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentCircularStringMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return _map.capacity
}

func (_map *concurrentCircularStringMap) IndexOf(object ISystemObject) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, key := range _map.keys {
		if _map.collection[key] == object {
			return index
		}
	}

	return -1
}

func (_map *concurrentCircularStringMap) IndexOfKey(key string) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, _key := range _map.keys {
		if key == _key {
			return index
		}
	}

	return -1
}

func (_map *concurrentCircularStringMap) GetKeyAt(index int) string {
	_map.RLock()
	defer _map.RUnlock()

	return _map.keys[index]
}

func (_map *concurrentCircularStringMap) GetValueAt(index int) ISystemObject {
	_map.RLock()
	defer _map.RUnlock()

	return _map.collection[_map.keys[index]]
}

func (_map *concurrentCircularStringMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.keys = make([]string, _map.capacity)
	_map.collection = make(map[string]ISystemObject, _map.capacity)
}

func (_map *concurrentCircularStringMap) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(_map.collection[key])
	}
}

func (_map *concurrentCircularStringMap) ForEachKey(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key)
	}
}

func (_map *concurrentCircularStringMap) ForEach(iterator func(string, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key, _map.collection[key])
	}
}

func (_map *concurrentCircularStringMap) ForEachParallel(iterator func(string, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for _, key := range _map.keys {
		go func(key string, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, _map.collection[key])
	}

	waitGroup.Wait()
}

func (_map *concurrentCircularStringMap) Filter(predicate func(ISystemObject) bool) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		value := _map.collection[key]
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentCircularStringMap) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		value := _map.collection[key]
		result = append(result, predicate(value))
	}

	return result
}
