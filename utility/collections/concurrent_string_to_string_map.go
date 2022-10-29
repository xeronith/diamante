package collections

import (
	"sync"
	"time"
)

type concurrentStringToStringMap struct {
	sync.RWMutex
	collection map[string]string
}

func NewConcurrentStringToStringMap() IStringToStringMap {
	return &concurrentStringToStringMap{
		collection: make(map[string]string),
	}
}

func (_map *concurrentStringToStringMap) Put(key string, value string) {
	_map.Lock()
	defer _map.Unlock()
	_map.collection[key] = value
}

func (_map *concurrentStringToStringMap) Get(key string) string {
	_map.RLock()
	defer _map.RUnlock()
	return _map.collection[key]
}

func (_map *concurrentStringToStringMap) Remove(key string) {
	_map.Lock()
	defer _map.Unlock()
	delete(_map.collection, key)
}

func (_map *concurrentStringToStringMap) Contains(key string) bool {
	_map.RLock()
	defer _map.RUnlock()
	_, ok := _map.collection[key]
	return ok
}

func (_map *concurrentStringToStringMap) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()
	length := len(_map.collection)
	return length
}

func (_map *concurrentStringToStringMap) Clear() {
	_map.Lock()
	defer _map.Unlock()
	_map.collection = make(map[string]string)
}

func (_map *concurrentStringToStringMap) ForEachValue(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentStringToStringMap) ForEachKey(iterator func(string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentStringToStringMap) ForEach(iterator func(string, string)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()
	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentStringToStringMap) ForEachParallel(iterator func(string, string)) {
	if iterator == nil {
		return
	}

	collection := func() map[string]string {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[string]string, len(_map.collection))
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
		go func(key string, value string) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentStringToStringMap) ForEachParallelWithInitialization(initializer func(int) error, iterator func(string)) error {
	if iterator == nil {
		return nil
	}

	collection := func() map[string]string {
		_map.RLock()
		defer _map.RUnlock()

		source := _map.collection
		destination := make(map[string]string, len(source))
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

func (_map *concurrentStringToStringMap) ForEachParallelWithInterval(duration time.Duration, iterator func(string, string)) {
	if iterator == nil {
		return
	}

	collection := func() map[string]string {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[string]string, len(_map.collection))
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
		go func(key string, value string) {
			defer waitGroup.Done()
			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentStringToStringMap) Filter(predicate func(string) bool) []string {
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

func (_map *concurrentStringToStringMap) Map(predicate func(string) string) []string {
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
