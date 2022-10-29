package collections

import (
	"time"

	. "github.com/xeronith/diamante/contracts/system"
)

type IPointerMap interface {
	Put(Pointer, ISystemObject)
	Get(Pointer) ISystemObject
	Remove(Pointer)
	RemoveFirst() (Pointer, ISystemObject)
	GetFirst() (Pointer, ISystemObject)
	Contains(Pointer) bool
	GetSize() int
	GetAll() map[Pointer]ISystemObject
	Clear()
	ForEachValue(func(ISystemObject))
	ForEachKey(func(Pointer))
	ForEach(func(Pointer, ISystemObject))
	ForEachParallel(func(Pointer, ISystemObject))
	ForEachParallelWithInitialization(func(int) error, func(Pointer)) error
	ForEachParallelWithInterval(time.Duration, func(Pointer, ISystemObject))
	Filter(func(ISystemObject) bool) []ISystemObject
	Map(func(ISystemObject) ISystemObject) []ISystemObject
}
