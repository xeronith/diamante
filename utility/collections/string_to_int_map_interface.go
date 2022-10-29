package collections

import (
	"time"
)

type IStringToIntMap interface {
	Put(string, int32)
	Get(string) int32
	Remove(string)
	Contains(string) bool
	GetSize() int
	Clear()
	ForEachValue(func(int32))
	ForEachKey(func(string))
	ForEach(func(string, int32))
	ForEachParallel(func(string, int32))
	ForEachParallelWithInitialization(func(int) error, func(string)) error
	ForEachParallelWithInterval(time.Duration, func(string, int32))
	Filter(func(int32) bool) []int32
	Map(func(int32) int32) []int32
}
