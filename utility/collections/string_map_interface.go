package collections

import . "github.com/xeronith/diamante/contracts/system"

type IStringMap interface {
	Put(string, ISystemObject)
	Get(string) (ISystemObject, bool)
	Remove(string)
	Contains(string) bool
	GetSize() int
	Clear()
	ForEachValue(func(ISystemObject))
	ForEachKey(func(string))
	ForEach(func(string, ISystemObject))
	ForEachParallel(func(string, ISystemObject))
	Filter(func(ISystemObject) bool) []ISystemObject
	Map(func(ISystemObject) ISystemObject) []ISystemObject
}
