package collections

import "time"

type IIntToStringMap interface {
	Put(int32, string)
	Get(int32) string
	Remove(int32)
	Contains(int32) bool
	GetSize() int
	Clear()
	ForEachValue(func(string))
	ForEachKey(func(int32))
	ForEach(func(int32, string))
	ForEachParallel(func(int32, string))
	ForEachParallelWithInitialization(func(int) error, func(int32)) error
	ForEachParallelWithInterval(time.Duration, func(int32, string))
	Filter(func(string) bool) []string
	Map(func(string) string) []string
}
