package collections

import . "github.com/xeronith/diamante/contracts/system"

type IInt64Map interface {
	Load(map[int64]ISystemObject)
	Put(int64, ISystemObject)
	Get(int64) (ISystemObject, bool)
	Remove(int64)
	Contains(int64) bool
	GetSize() int
	Clear()
	ForEachValue(func(ISystemObject))
	ForEachKey(func(int64))
	ForEach(func(int64, ISystemObject))
	ForEachParallel(func(int64, ISystemObject))
	Filter(func(ISystemObject) bool) []ISystemObject
	Map(func(ISystemObject) ISystemObject) []ISystemObject
}

type IOrderedInt64Map interface {
	IInt64Map
	IndexOf(ISystemObject) int
	IndexOfKey(int64) int
	GetKeyAt(int) int64
	GetValueAt(int) ISystemObject
}
