package collections

import . "github.com/xeronith/diamante/contracts/system"

type IInt32Map interface {
	Load(map[int32]ISystemObject)
	Put(int32, ISystemObject)
	Get(int32) (ISystemObject, bool)
	Remove(int32)
	Contains(int32) bool
	GetSize() int
	Clear()
	ForEachValue(func(ISystemObject))
	ForEachKey(func(int32))
	ForEach(func(int32, ISystemObject))
	ForEachParallel(func(int32, ISystemObject))
	Filter(func(ISystemObject) bool) []ISystemObject
	Map(func(ISystemObject) ISystemObject) []ISystemObject
}

type IOrderedInt32Map interface {
	IInt32Map
	IndexOf(ISystemObject) int
	IndexOfKey(int32) int
	GetKeyAt(int) int32
	GetValueAt(int) ISystemObject
}
