package collections

import . "github.com/xeronith/diamante/contracts/system"

type ISlice interface {
	Append(ISystemObject)
	AppendIfNotExist(ISystemObject)
	Remove(ISystemObject)
	GetSize() int
	Clear()
	ForEach(func(int, ISystemObject))
	ForEachParallel(func(int, ISystemObject))
}
