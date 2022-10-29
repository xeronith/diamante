package collections

import . "github.com/xeronith/diamante/contracts/system"

type IIntMap interface {
	Put(int, ISystemObject)
	Get(int) ISystemObject
	Remove(int)
	Contains(int) bool
	GetSize() int
	Clear()
	ForEachValue(func(ISystemObject))
	ForEachKey(func(int))
	ForEach(func(int, ISystemObject))
	ForEachParallel(func(int, ISystemObject))
	Filter(func(ISystemObject) bool) []ISystemObject
	Map(func(ISystemObject) ISystemObject) []ISystemObject
}
