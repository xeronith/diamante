package security

import . "github.com/xeronith/diamante/contracts/system"

type Role = uint64

// noinspection GoUnusedConst,GoSnakeCaseUsage
const ENABLE_SECURITY = true

// noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	ANONYMOUS     Role = 0b0000_0000_0000_0000_0000_0000_0000_0000
	USER          Role = 0b0000_0000_0000_0000_0000_0000_0000_0001
	ADMINISTRATOR Role = 0b1111_1111_1111_1111_1111_1111_1111_1111
)

var ROLES = []Role{
	ANONYMOUS,
	USER,
	ADMINISTRATOR,
}

type (
	Identity interface {
		Id() int64
		Username() string
		PhoneNumber() string
		FirstName() string
		LastName() string
		Email() string
		Token() string
		SetToken(token string)
		MultiFactor() bool
		RemoteAddress() string
		SetRemoteAddress(string)
		UserAgent() string
		SetUserAgent(string)
		Hash() string
		Salt() string
		PublicKey() string
		PrivateKey() string
		Permission() uint64
		Restriction() uint32
		Role() Role
		IsInRole(role Role) bool
		IsRestricted() bool
		IsNotRestricted() bool
		Payload() Pointer
		Lock(uint64)
		Unlock(uint64)
		SetSystemCallHandler(func(Identity, []string) error)
		SystemCall(Identity, []string) error
	}

	IdentityCache interface {
		Put(int64, Identity)
		Remove(int64)
		Get(int64) (Identity, bool)
		Size() int
		GetByToken(string) (Identity, bool)
		GetByUsername(string) (Identity, bool)
		GetByPhoneNumber(string) (Identity, bool)
		RefreshToken(Identity, string)
		ForEachValue(func(ISystemObject))
		Load(map[int64]ISystemObject)
		Clear()
		OnChanged(func())

		StoreAuthorizationInfo(string, string, string)
		RetrieveAuthorizationInfo(string) (string, string, error)
	}
)
