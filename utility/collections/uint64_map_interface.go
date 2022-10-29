package collections

import . "github.com/xeronith/diamante/contracts/system"

type IUInt64Map interface {
	Put(uint64, ISystemObject)
	Get(uint64) (ISystemObject, bool)
	Remove(uint64)
	Contains(uint64) bool
	GetSize() int64
	Clear()
	ForEachValue(func(ISystemObject))
	ForEachKey(func(uint64))
	ForEach(func(uint64, ISystemObject))
	ForEachParallel(func(uint64, ISystemObject))
	Filter(func(ISystemObject) bool) []ISystemObject
	Map(func(ISystemObject) ISystemObject) []ISystemObject
}
