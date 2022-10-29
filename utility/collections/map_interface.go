package collections

type IMap[T comparable, V any] interface {
	Put(T, V)
	Get(T) V
	Remove(T)
	Contains(T) bool
	GetSize() int
	Clear()
	ForEachValue(func(V))
	ForEachKey(func(T))
	ForEach(func(T, V))
	ForEachParallel(func(T, V))
	Filter(func(V) bool) []V
	Map(func(V) V) []V
}
