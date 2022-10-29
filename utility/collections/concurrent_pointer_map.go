package collections

import (
	"sync"
	"time"

	. "github.com/xeronith/diamante/contracts/system"
)

type concurrentPointerMap struct {
	sync.RWMutex
	collection map[Pointer]ISystemObject
}

func NewConcurrentPointerMap() IPointerMap {
	_map := &concurrentPointerMap{
		collection: make(map[Pointer]ISystemObject),
	}

	return _map
}

func (_map *concurrentPointerMap) Put(key Pointer, value ISystemObject) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentPointerMap) Get(key Pointer) ISystemObject {
	_map.RLock()
	defer _map.RUnlock()
	return _map.collection[key]
}

func (_map *concurrentPointerMap) Remove(key Pointer) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentPointerMap) getFirst(remove bool) (Pointer, ISystemObject) {
	_map.Lock()
	defer _map.Unlock()

	if len(_map.collection) == 0 {
		return nil, nil
	}

	var (
		firstKey   Pointer
		firstValue ISystemObject
	)

	for key, value := range _map.collection {
		firstKey = key
		firstValue = value
		break
	}

	if remove {
		delete(_map.collection, firstKey)
	}

	return firstKey, firstValue
}

func (_map *concurrentPointerMap) RemoveFirst() (Pointer, ISystemObject) {
	return _map.getFirst(true)
}

func (_map *concurrentPointerMap) GetFirst() (Pointer, ISystemObject) {
	return _map.getFirst(false)
}

func (_map *concurrentPointerMap) Contains(key Pointer) bool {
	_map.RLock()
	defer _map.RUnlock()
	_, ok := _map.collection[key]
	return ok
}

func (_map *concurrentPointerMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	return len(_map.collection)
}

func (_map *concurrentPointerMap) GetAll() map[Pointer]ISystemObject {
	_map.RLock()
	defer _map.RUnlock()
	result := make(map[Pointer]ISystemObject)

	if _map.collection != nil {
		for pointer, object := range _map.collection {
			result[pointer] = object
		}
	}

	return result
}

func (_map *concurrentPointerMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[Pointer]ISystemObject)
}

func (_map *concurrentPointerMap) ForEachValue(iterator func(ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentPointerMap) ForEachKey(iterator func(Pointer)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentPointerMap) ForEach(iterator func(Pointer, ISystemObject)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentPointerMap) ForEachParallel(iterator func(Pointer, ISystemObject)) {
	if iterator == nil {
		return
	}

	collection := func() map[Pointer]ISystemObject {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[Pointer]ISystemObject, len(_map.collection))
		for key, value := range _map.collection {
			collection[key] = value
		}

		return collection
	}()

	length := len(collection)
	if length < 1 {
		return
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(length)

	for key, value := range collection {
		go func(key Pointer, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentPointerMap) ForEachParallelWithInitialization(initializer func(int) error, iterator func(Pointer)) error {
	if iterator == nil {
		return nil
	}

	collection := func() map[Pointer]ISystemObject {
		_map.RLock()
		defer _map.RUnlock()

		source := _map.collection
		destination := make(map[Pointer]ISystemObject, len(source))
		for key, value := range source {
			destination[key] = value
		}

		return destination
	}()

	length := len(collection)
	if length < 1 {
		return nil
	}

	if initializer != nil {
		if err := initializer(length); err != nil {
			return err
		}
	}

	go func() {
		for key := range collection {
			go func(key Pointer) {
				iterator(key)
			}(key)

			// runtime.Gosched()
		}
	}()

	return nil
}

func (_map *concurrentPointerMap) ForEachParallelWithInterval(duration time.Duration, iterator func(Pointer, ISystemObject)) {
	if iterator == nil {
		return
	}

	collection := func() map[Pointer]ISystemObject {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[Pointer]ISystemObject, len(_map.collection))
		for key, value := range _map.collection {
			collection[key] = value
		}

		return collection
	}()

	length := len(collection)
	if length < 1 {
		return
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(length)

	for key, value := range collection {
		time.Sleep(duration)
		go func(key Pointer, value ISystemObject) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentPointerMap) Filter(predicate func(ISystemObject) bool) []ISystemObject {
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

func (_map *concurrentPointerMap) Map(predicate func(ISystemObject) ISystemObject) []ISystemObject {
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
