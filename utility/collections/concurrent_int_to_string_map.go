package collections

import (
	"sync"
	"time"
)

type concurrentIntToStringMap struct {
	sync.RWMutex
	collection map[int32]string
}

func NewConcurrentIntToStringMap() IIntToStringMap {
	return &concurrentIntToStringMap{
		collection: make(map[int32]string),
	}
}

func (_map *concurrentIntToStringMap) Put(key int32, value string) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentIntToStringMap) Get(key int32) string {
	_map.RLock()
	defer _map.RUnlock()
	return _map.collection[key]
}

func (_map *concurrentIntToStringMap) Remove(key int32) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentIntToStringMap) Contains(key int32) bool {
	_map.RLock()
	defer _map.RUnlock()
	_, ok := _map.collection[key]
	return ok
}

func (_map *concurrentIntToStringMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	length := len(_map.collection)
	return length
}

func (_map *concurrentIntToStringMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[int32]string)
}

func (_map *concurrentIntToStringMap) ForEachValue(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentIntToStringMap) ForEachKey(iterator func(int32)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentIntToStringMap) ForEach(iterator func(int32, string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentIntToStringMap) ForEachParallel(iterator func(int32, string)) {
	if iterator == nil {
		return
	}

	collection := func() map[int32]string {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[int32]string, len(_map.collection))
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
		go func(key int32, value string) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentIntToStringMap) ForEachParallelWithInitialization(initializer func(int) error, iterator func(int32)) error {
	if iterator == nil {
		return nil
	}

	collection := func() map[int32]string {
		_map.RLock()
		defer _map.RUnlock()

		source := _map.collection
		destination := make(map[int32]string, len(source))
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
			go func(key int32) {
				iterator(key)
			}(key)
		}
	}()

	return nil
}

func (_map *concurrentIntToStringMap) ForEachParallelWithInterval(duration time.Duration, iterator func(int32, string)) {
	if iterator == nil {
		return
	}

	collection := func() map[int32]string {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[int32]string, len(_map.collection))
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
		go func(key int32, value string) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentIntToStringMap) Filter(predicate func(string) bool) []string {
	if predicate == nil {
		return nil
	}

	result := make([]string, 0)
	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentIntToStringMap) Map(predicate func(string) string) []string {
	if predicate == nil {
		return nil
	}

	result := make([]string, 0)
	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		result = append(result, predicate(value))
	}

	return result
}
