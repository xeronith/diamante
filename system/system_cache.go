package system

import (
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/utility/collections"
)

type cache struct {
	items     IInt64Map
	onChanged func()
}

func NewCache() ICache {
	return &cache{
		items: NewConcurrentOrderedInt64Map(),
	}
}

func (cache *cache) Put(key int64, value ISystemObject) {
	cache.items.Put(key, value)
	cache.notifyChanged()
}

func (cache *cache) Remove(key int64, _ ISystemObject) {
	cache.items.Remove(key)
	cache.notifyChanged()
}

func (cache *cache) Get(key int64) (ISystemObject, bool) {
	return cache.items.Get(key)
}

func (cache *cache) Size() int {
	return cache.items.GetSize()
}

func (cache *cache) ForEachValue(iterator func(ISystemObject)) {
	cache.items.ForEachValue(iterator)
}

func (cache *cache) Load(collection map[int64]ISystemObject) {
	cache.items.Load(collection)
	cache.notifyChanged()
}

func (cache *cache) Clear() {
	cache.items.Clear()
	cache.notifyChanged()
}

func (cache *cache) OnChanged(callback func()) {
	cache.onChanged = callback
}

func (cache *cache) notifyChanged() {
	if cache.onChanged != nil {
		cache.onChanged()
	}
}
