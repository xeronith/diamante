package collections

import (
	"sort"
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentSortedInt64Map struct {
	sync.RWMutex
	keys       []int64
	collection map[int64]ISystemObject
}

func NewConcurrentSortedInt64Map() IOrderedInt64Map {
	return &concurrentSortedInt64Map{
		keys:       make([]int64, 0),
		collection: make(map[int64]ISystemObject),
	}
}

func (_map *concurrentSortedInt64Map) Load(collection map[int64]ISystemObject) {
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

func (_map *concurrentSortedInt64Map) Put(key int64, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	if _, exists := _map.collection[key]; !exists {
		_map.keys = append(_map.keys, key)
		sort.Slice(_map.keys, func(a, b int) bool {
			return _map.keys[a] < _map.keys[b]
		})
	}

	_map.collection[key] = value
}

func (_map *concurrentSortedInt64Map) Get(key int64) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentSortedInt64Map) Remove(key int64) {
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

func (_map *concurrentSortedInt64Map) Contains(key int64) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentSortedInt64Map) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return len(_map.collection)
}

func (_map *concurrentSortedInt64Map) IndexOf(object ISystemObject) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, key := range _map.keys {
		if _map.collection[key] == object {
			return index
		}
	}

	return -1
}

func (_map *concurrentSortedInt64Map) IndexOfKey(key int64) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, _key := range _map.keys {
		if key == _key {
			return index
		}
	}

	return -1
}

func (_map *concurrentSortedInt64Map) GetKeyAt(index int) int64 {
	_map.RLock()
	defer _map.RUnlock()

	return _map.keys[index]
}

func (_map *concurrentSortedInt64Map) GetValueAt(index int) ISystemObject {
	_map.RLock()
	defer _map.RUnlock()

	return _map.collection[_map.keys[index]]
}

func (_map *concurrentSortedInt64Map) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int64]ISystemObject)
	_map.keys = make([]int64, 0)
}

func (_map *concurrentSortedInt64Map) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.Lock()
	defer _map.Unlock()
	for _, key := range _map.keys {
		iterator(_map.collection[key])
	}
}

func (_map *concurrentSortedInt64Map) ForEachKey(iterator func(int64)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key)
	}
}

func (_map *concurrentSortedInt64Map) ForEach(iterator func(int64, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key, _map.collection[key])
	}
}

func (_map *concurrentSortedInt64Map) ForEachParallel(iterator func(int64, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for _, key := range _map.keys {
		go func(key int64, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, _map.collection[key])
	}

	waitGroup.Wait()
}

func (_map *concurrentSortedInt64Map) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentSortedInt64Map) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
