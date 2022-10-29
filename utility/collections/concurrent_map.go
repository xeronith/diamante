package collections

import "sync"

type concurrentMap[T comparable, V any] struct {
	sync.RWMutex
	collection map[T]V
}

func NewConcurrentMap[T comparable, V any]() IMap[T, V] {
	return &concurrentMap[T, V]{
		collection: make(map[T]V),
	}
}

func (_map *concurrentMap[T, V]) Put(key T, value V) {
	_map.Lock()
	defer _map.Unlock()

	_map.collection[key] = value
}

func (_map *concurrentMap[T, V]) Get(key T) V {
	_map.RLock()
	defer _map.RUnlock()

	return _map.collection[key]
}

func (_map *concurrentMap[T, V]) Remove(key T) {
	_map.Lock()
	defer _map.Unlock()

	delete(_map.collection, key)
}

func (_map *concurrentMap[T, V]) Contains(key T) bool {
	_map.RLock()
	defer _map.RUnlock()

	_, ok := _map.collection[key]
	return ok
}

func (_map *concurrentMap[T, V]) GetSize() int {
	_map.RLock()
	defer _map.RUnlock()

	return len(_map.collection)
}

func (_map *concurrentMap[T, V]) Clear() {
	_map.Lock()
	defer _map.Unlock()

	_map.collection = make(map[T]V)
}

func (_map *concurrentMap[T, V]) ForEachValue(iterator func(V)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	for _, value := range _map.collection {
		iterator(value)
	}
}

func (_map *concurrentMap[T, V]) ForEachKey(iterator func(T)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	for key := range _map.collection {
		iterator(key)
	}
}

func (_map *concurrentMap[T, V]) ForEach(iterator func(T, V)) {
	if iterator == nil {
		return
	}

	_map.RLock()
	defer _map.RUnlock()

	for key, value := range _map.collection {
		iterator(key, value)
	}
}

func (_map *concurrentMap[T, V]) ForEachParallel(iterator func(T, V)) {
	if iterator == nil {
		return
	}

	collection := func() map[T]V {
		_map.RLock()
		defer _map.RUnlock()

		collection := make(map[T]V, len(_map.collection))
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
		go func(key T, value V) {
			defer waitGroup.Done()

			iterator(key, value)
		}(key, value)
	}

	waitGroup.Wait()
}

func (_map *concurrentMap[T, V]) Filter(predicate func(V) bool) []V {
	if predicate == nil {
		return nil
	}

	result := make([]V, 0)

	_map.RLock()
	defer _map.RUnlock()

	for _, value := range _map.collection {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (_map *concurrentMap[T, V]) Map(predicate func(V) V) []V {
	if predicate == nil {
		return nil
	}

	result := make([]V, 0)

	_map.RLock()
	defer _map.RUnlock()

	for _, value := range _map.collection {
		result = append(result, predicate(value))
	}

	return result
}
