package collections

import (
	"time"
)

type IStringToStringMap interface {
	Put(string, string)
	Get(string) string
	Remove(string)
	Contains(string) bool
	GetSize() int
	Clear()
	ForEachValue(func(string))
	ForEachKey(func(string))
	ForEach(func(string, string))
	ForEachParallel(func(string, string))
	ForEachParallelWithInitialization(func(int) error, func(string)) error
	ForEachParallelWithInterval(time.Duration, func(string, string))
	Filter(func(string) bool) []string
	Map(func(string) string) []string
}
