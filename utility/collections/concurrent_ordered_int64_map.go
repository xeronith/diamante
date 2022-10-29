package collections

import (
	"sort"
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentOrderedInt64Map struct {
	sync.RWMutex
	keys       []int64
	collection map[int64]ISystemObject
}

func NewConcurrentOrderedInt64Map() IOrderedInt64Map {
	return &concurrentOrderedInt64Map{
		keys:       make([]int64, 0),
		collection: make(map[int64]ISystemObject),
	}
}

func (_map *concurrentOrderedInt64Map) Load(collection map[int64]ISystemObject) {
	keys := make([]int64, 0)
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

func (_map *concurrentOrderedInt64Map) Put(key int64, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	if _, exists := _map.collection[key]; !exists {
		_map.keys = append(_map.keys, key)
	}

	_map.collection[key] = value
}

func (_map *concurrentOrderedInt64Map) Get(key int64) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentOrderedInt64Map) Remove(key int64) {
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

func (_map *concurrentOrderedInt64Map) Contains(key int64) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentOrderedInt64Map) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return len(_map.collection)
}

func (_map *concurrentOrderedInt64Map) IndexOf(object ISystemObject) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, key := range _map.keys {
		if _map.collection[key] == object {
			return index
		}
	}

	return -1
}

func (_map *concurrentOrderedInt64Map) IndexOfKey(key int64) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, _key := range _map.keys {
		if key == _key {
			return index
		}
	}

	return -1
}

func (_map *concurrentOrderedInt64Map) GetKeyAt(index int) int64 {
	_map.RLock()
	defer _map.RUnlock()

	return _map.keys[index]
}

func (_map *concurrentOrderedInt64Map) GetValueAt(index int) ISystemObject {
	_map.RLock()
	defer _map.RUnlock()

	return _map.collection[_map.keys[index]]
}

func (_map *concurrentOrderedInt64Map) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int64]ISystemObject)
	_map.keys = make([]int64, 0)
}

func (_map *concurrentOrderedInt64Map) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	keys, collection := _map.snapshot()
	for _, key := range keys {
		iterator(collection[key])
	}
}

func (_map *concurrentOrderedInt64Map) ForEachKey(iterator func(int64)) {
	if iterator == nil {
		return
	}

	keys, _ := _map.snapshot()
	for _, key := range keys {
		iterator(key)
	}
}

func (_map *concurrentOrderedInt64Map) ForEach(iterator func(int64, ISystemObject)) {
	if iterator == nil {
		return
	}

	keys, collection := _map.snapshot()
	for _, key := range keys {
		iterator(key, collection[key])
	}
}

func (_map *concurrentOrderedInt64Map) ForEachParallel(iterator func(int64, ISystemObject)) {
	if iterator == nil {
		return
	}

	keys, collection := _map.snapshot()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(collection))

	for _, key := range keys {
		go func(key int64, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, collection[key])
	}

	waitGroup.Wait()
}

func (_map *concurrentOrderedInt64Map) Filter(predicate func(ISystemObject) bool) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)

	keys, collection := _map.snapshot()
	for _, key := range keys {
		value := collection[key]
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentOrderedInt64Map) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
	if predicate == nil {
		return nil
	}

	result := make([]ISystemObject, 0)

	keys, collection := _map.snapshot()
	for _, key := range keys {
		value := collection[key]
		result = append(result, predicate(value))
	}

	return result
}

func (_map *concurrentOrderedInt64Map) snapshot() ([]int64, map[int64]ISystemObject) {
	_map.RLock()
	defer _map.RUnlock()

	keys := make([]int64, len(_map.keys))
	collection := make(map[int64]ISystemObject)

	copy(keys, _map.keys)
	for _, key := range keys {
		collection[key] = _map.collection[key]
	}

	return keys, collection
}
