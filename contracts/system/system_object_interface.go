package system

type (
	ISystemObject interface {
	}

	ICache interface {
		Put(int64, ISystemObject)
		Remove(int64, ISystemObject)
		Get(int64) (ISystemObject, bool)
		Size() int
		ForEachValue(func(ISystemObject))
		Load(map[int64]ISystemObject)
		Clear()
		OnChanged(func())
	}
)
