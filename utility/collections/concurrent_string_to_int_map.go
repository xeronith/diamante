package collections

import (
	"sync"
	"time"
)

type concurrentStringToIntMap struct {
	sync.RWMutex
	collection map[string]int32
}

func NewConcurrentStringToIntMap() IStringToIntMap {
	return &concurrentStringToIntMap{
		collection: make(map[string]int32),
	}
}

func (_map *concurrentStringToIntMap) Put(key string, value int32) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentStringToIntMap) Get(key string) int32 {
	_map.RLock()
	defer _map.RUnlock()
	return _map.collection[key]
}

func (_map *concurrentStringToIntMap) Remove(key string) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentStringToIntMap) Contains(key string) bool {
	_map.RLock()
	defer _map.RUnlock()
	_, ok := _map.collection[key]
	return ok
}

func (_map *concurrentStringToIntMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	length := len(_map.collection)
	return length
}

func (_map *concurrentStringToIntMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[string]int32)
}

func (_map *concurrentStringToIntMap) ForEachValue(iterator func(int32)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentStringToIntMap) ForEachKey(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentStringToIntMap) ForEach(iterator func(string, int32)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentStringToIntMap) ForEachParallel(iterator func(string, int32)) {
	if iterator == nil {
		return
	}

	collection := func() map[string]int32 {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[string]int32, len(_map.collection))
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
		go func(key string, value int32) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentStringToIntMap) ForEachParallelWithInitialization(initializer func(int) error, iterator func(string)) error {
	if iterator == nil {
		return nil
	}

	collection := func() map[string]int32 {
		_map.RLock()
		defer _map.RUnlock()

		source := _map.collection
		destination := make(map[string]int32, len(source))
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
			go func(key string) {
				iterator(key)
			}(key)
		}
	}()

	return nil
}

func (_map *concurrentStringToIntMap) ForEachParallelWithInterval(duration time.Duration, iterator func(string, int32)) {
	if iterator == nil {
		return
	}

	collection := func() map[string]int32 {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[string]int32, len(_map.collection))
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
		go func(key string, value int32) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentStringToIntMap) Filter(predicate func(int32) bool) []int32 {
	if predicate == nil {
		return nil
	}

	result := make([]int32, 0)
	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentStringToIntMap) Map(predicate func(int32) int32) []int32 {
	if predicate == nil {
		return nil
	}

	result := make([]int32, 0)
	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		result = append(result, predicate(value))
	}

	return result
}
