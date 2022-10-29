package collections

import (
	"sort"
	"sync"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentOrderedInt32Map struct {
	sync.RWMutex
	keys       []int32
	collection map[int32]ISystemObject
}

func NewConcurrentOrderedInt32Map() IOrderedInt32Map {
	return &concurrentOrderedInt32Map{
		keys:       make([]int32, 0),
		collection: make(map[int32]ISystemObject),
	}
}

func (_map *concurrentOrderedInt32Map) Load(collection map[int32]ISystemObject) {
	keys := make([]int32, 0)
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

func (_map *concurrentOrderedInt32Map) Put(key int32, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	if _, exists := _map.collection[key]; !exists {
		_map.keys = append(_map.keys, key)
	}

	_map.collection[key] = value
}

func (_map *concurrentOrderedInt32Map) Get(key int32) (ISystemObject, bool) {
	_map.RLock()
	defer _map.RUnlock()
	value, exists := _map.collection[key]
	return value, exists
}

func (_map *concurrentOrderedInt32Map) Remove(key int32) {
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

func (_map *concurrentOrderedInt32Map) Contains(key int32) bool {
	_, exists := _map.Get(key)
	return exists
}

func (_map *concurrentOrderedInt32Map) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return len(_map.collection)
}

func (_map *concurrentOrderedInt32Map) IndexOf(object ISystemObject) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, key := range _map.keys {
		if _map.collection[key] == object {
			return index
		}
	}

	return -1
}

func (_map *concurrentOrderedInt32Map) IndexOfKey(key int32) int {
	_map.RLock()
	defer _map.RUnlock()

	for index, _key := range _map.keys {
		if key == _key {
			return index
		}
	}

	return -1
}

func (_map *concurrentOrderedInt32Map) GetKeyAt(index int) int32 {
	_map.RLock()
	defer _map.RUnlock()

	return _map.keys[index]
}

func (_map *concurrentOrderedInt32Map) GetValueAt(index int) ISystemObject {
	_map.RLock()
	defer _map.RUnlock()

	return _map.collection[_map.keys[index]]
}

func (_map *concurrentOrderedInt32Map) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int32]ISystemObject)
	_map.keys = make([]int32, 0)
}

func (_map *concurrentOrderedInt32Map) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(_map.collection[key])
	}
}

func (_map *concurrentOrderedInt32Map) ForEachKey(iterator func(int32)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key)
	}
}

func (_map *concurrentOrderedInt32Map) ForEach(iterator func(int32, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, key := range _map.keys {
		iterator(key, _map.collection[key])
	}
}

func (_map *concurrentOrderedInt32Map) ForEachParallel(iterator func(int32, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(_map.collection))

	for _, key := range _map.keys {
		go func(key int32, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, _map.collection[key])
	}

	waitGroup.Wait()
}

func (_map *concurrentOrderedInt32Map) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentOrderedInt32Map) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
