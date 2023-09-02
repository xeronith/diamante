package collections

import (
	"sort"
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentOrderedStringMap struct {
	sync.RWMutex
	capacity   int
	keys       []string
	collection map[string]ISystemObject
}

func NewConcurrentOrderedStringMap(capacity int) IOrderedStringMap {
	return &concurrentOrderedStringMap{
		capacity:   capacity,
		keys:       make([]string, 0),
		collection: make(map[string]ISystemObject),
	}
}

func (_map *concurrentOrderedStringMap) Load(collection map[string]ISystemObject) {
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

func (_map *concurrentOrderedStringMap) Put(key string, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	if _, exists := _map.collection[key]; !exists {
		_map.keys = append(_map.keys, key)
	}

	_map.collection[key] = value

	if _map.capacity > 1 && len(_map.collection) > _map.capacity {
		key := _map.keys[0]

		_map.keys = _map.keys[1:]
		delete(_map.collection, key)
	}
}

func (_map *concurrentOrderedStringMap) Get(key string) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentOrderedStringMap) Remove(key string) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)

	index := -1
	for i, item := range _map.keys {
		if item == key {
			index = i
			break
		}
	}

	if index >= 0 {
		_map.keys = append(_map.keys[:index], _map.keys[index+1:]...)
	}
}

func (_map *concurrentOrderedStringMap) Contains(key string) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentOrderedStringMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return len(_map.collection)
}

func (_map *concurrentOrderedStringMap) IndexOf(object ISystemObject) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, key := range _map.keys {
		if _map.collection[key] == object {
			return index
		}
	}

	return -1
}

func (_map *concurrentOrderedStringMap) IndexOfKey(key string) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, _key := range _map.keys {
		if key == _key {
			return index
		}
	}

	return -1
}

func (_map *concurrentOrderedStringMap) GetKeyAt(index int) string {
	_map.RLock()
	defer _map.RUnlock()

	return _map.keys[index]
}

func (_map *concurrentOrderedStringMap) GetValueAt(index int) ISystemObject {
	_map.RLock()
	defer _map.RUnlock()

	return _map.collection[_map.keys[index]]
}

func (_map *concurrentOrderedStringMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[string]ISystemObject)
	_map.keys = make([]string, 0)
}

func (_map *concurrentOrderedStringMap) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(_map.collection[key])
	}
}

func (_map *concurrentOrderedStringMap) ForEachKey(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key)
	}
}

func (_map *concurrentOrderedStringMap) ForEach(iterator func(string, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key, _map.collection[key])
	}
}

func (_map *concurrentOrderedStringMap) ForEachParallel(iterator func(string, ISystemObject)) {
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

func (_map *concurrentOrderedStringMap) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentOrderedStringMap) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
