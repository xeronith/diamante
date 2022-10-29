package security

import (
	"errors"
	"sync"
	"time"

	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/utility/collections"
)

type identityCache struct {
	sync.RWMutex
	idToIdentityMap          IInt64Map
	tokenToIdentityMap       IStringMap
	usernameToIdentityMap    IStringMap
	phoneNumberToIdentityMap IStringMap

	registrationInfoMap IStringMap
	expirationWindow    time.Duration

	onChanged func()
}

func NewIdentityCache() IdentityCache {
	return &identityCache{
		idToIdentityMap:          NewConcurrentInt64Map(),
		tokenToIdentityMap:       NewConcurrentStringMap(),
		usernameToIdentityMap:    NewConcurrentStringMap(),
		phoneNumberToIdentityMap: NewConcurrentStringMap(),
		registrationInfoMap:      NewConcurrentStringMap(),
		expirationWindow:         time.Minute * 10,
	}
}

func (cache *identityCache) Put(_ int64, identity Identity) {
	cache.Lock()
	defer cache.Unlock()

	cache.idToIdentityMap.Put(identity.Id(), identity)
	cache.tokenToIdentityMap.Put(identity.Token(), identity)
	cache.usernameToIdentityMap.Put(identity.Username(), identity)
	cache.phoneNumberToIdentityMap.Put(identity.PhoneNumber(), identity)

	cache.notifyChanged()
}

func (cache *identityCache) Remove(id int64) {
	cache.Lock()
	defer cache.Unlock()

	if _identity, exists := cache.idToIdentityMap.Get(id); exists {
		identity := _identity.(Identity)
		cache.idToIdentityMap.Remove(identity.Id())
		cache.tokenToIdentityMap.Remove(identity.Token())
		cache.usernameToIdentityMap.Remove(identity.Username())
		cache.phoneNumberToIdentityMap.Remove(identity.PhoneNumber())

		cache.notifyChanged()
	}
}

func (cache *identityCache) Get(id int64) (Identity, bool) {
	cache.RLock()
	defer cache.RUnlock()

	if identity, exists := cache.idToIdentityMap.Get(id); exists {
		return identity.(Identity), true
	}

	return nil, false
}

func (cache *identityCache) Size() int {
	return cache.idToIdentityMap.GetSize()
}

func (cache *identityCache) GetByToken(token string) (Identity, bool) {
	cache.RLock()
	defer cache.RUnlock()

	if identity, exists := cache.tokenToIdentityMap.Get(token); exists {
		return identity.(Identity), true
	}

	return nil, false
}

func (cache *identityCache) GetByUsername(username string) (Identity, bool) {
	cache.RLock()
	defer cache.RUnlock()

	if identity, exists := cache.usernameToIdentityMap.Get(username); exists {
		return identity.(Identity), true
	}

	return nil, false
}

func (cache *identityCache) GetByPhoneNumber(phoneNumber string) (Identity, bool) {
	cache.RLock()
	defer cache.RUnlock()

	if identity, exists := cache.phoneNumberToIdentityMap.Get(phoneNumber); exists {
		return identity.(Identity), true
	}

	return nil, false
}

func (cache *identityCache) RefreshToken(identity Identity, token string) {
	cache.Lock()
	defer cache.Unlock()

	if cache.registrationInfoMap.Contains(identity.Token()) {
		cache.registrationInfoMap.Remove(identity.Token())
	}

	cache.tokenToIdentityMap.Remove(identity.Token())
	identity.SetToken(token)
	cache.tokenToIdentityMap.Put(identity.Token(), identity)
}

func (cache *identityCache) ForEachValue(iterator func(ISystemObject)) {
	cache.RLock()
	defer cache.RUnlock()

	cache.idToIdentityMap.ForEachValue(iterator)
}

func (cache *identityCache) Load(collection map[int64]ISystemObject) {
	cache.Lock()
	defer cache.Unlock()

	for _, object := range collection {
		if identity, ok := object.(Identity); ok {
			cache.idToIdentityMap.Put(identity.Id(), identity)
			cache.tokenToIdentityMap.Put(identity.Token(), identity)
			cache.usernameToIdentityMap.Put(identity.Username(), identity)
			cache.phoneNumberToIdentityMap.Put(identity.PhoneNumber(), identity)
		}
	}

	cache.notifyChanged()
}

func (cache *identityCache) Clear() {
	cache.Lock()
	defer cache.Unlock()

	cache.idToIdentityMap.Clear()
	cache.tokenToIdentityMap.Clear()
	cache.usernameToIdentityMap.Clear()
	cache.phoneNumberToIdentityMap.Clear()
	cache.registrationInfoMap.Clear()

	cache.notifyChanged()
}

func (cache *identityCache) OnChanged(callback func()) {
	cache.onChanged = callback
}

func (cache *identityCache) notifyChanged() {
	if cache.onChanged != nil {
		cache.onChanged()
	}
}

func (cache *identityCache) StoreAuthorizationInfo(token, phoneNumber, confirmationCode string) {
	cache.registrationInfoMap.Put(token, &struct {
		phoneNumber      string
		confirmationCode string
		timestamp        time.Time
		attempts         int
	}{
		phoneNumber:      phoneNumber,
		confirmationCode: confirmationCode,
		timestamp:        time.Now(),
		attempts:         0,
	})
}

var invalidToken = errors.New("ERROR_MESSAGE_INVALID_TOKEN")

func (cache *identityCache) RetrieveAuthorizationInfo(token string) (string, string, error) {
	if object, exists := cache.registrationInfoMap.Get(token); !exists {
		return "", "", invalidToken
	} else {
		registrationInfo := object.(*struct {
			phoneNumber      string
			confirmationCode string
			timestamp        time.Time
			attempts         int
		})

		registrationInfo.attempts++
		if registrationInfo.attempts > 3 || time.Since(registrationInfo.timestamp) > cache.expirationWindow {
			cache.registrationInfoMap.Remove(token)
			return "", "", invalidToken
		}

		return registrationInfo.phoneNumber, registrationInfo.confirmationCode, nil
	}
}
