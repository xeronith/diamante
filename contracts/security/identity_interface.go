package security

import . "github.com/xeronith/diamante/contracts/system"

type Role = uint64

//goland:noinspection GoSnakeCaseUsage
const ENABLE_SECURITY = true

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	ANONYMOUS          Role = 0b0000_0000_0000_0000_0000_0000_0000_0000
	USER               Role = 0b0000_0000_0000_0000_0000_0000_0000_0001
	VENDOR_USER        Role = 0b0000_0000_0000_0000_0000_0000_0000_0010
	WAREHOUSE_KEEPER   Role = 0b0000_0000_0000_0000_0000_0000_0000_0100
	WAREHOUSE_MANAGER  Role = 0b0000_0000_0000_0000_0000_0000_0000_1000
	CONTENT_SPECIALIST Role = 0b0000_0000_0000_0000_0000_0000_0001_0000
	ACCOUNT_MANAGER    Role = 0b0000_0000_0000_0000_0000_0000_0010_0000
	COMMERCIAL_MANAGER Role = 0b0000_0000_0000_0000_0000_0000_0100_0000
	SUPPORT_SPECIALIST Role = 0b0000_0000_0000_0000_0000_0000_1000_0000
	SUPPORT            Role = 0b0000_0000_0000_0000_0000_0001_0000_0000
	SUPPORT_MANAGER    Role = 0b0000_0000_0000_0000_0000_0010_0000_0000
	ADMINISTRATOR      Role = 0b1111_1111_1111_1111_1111_1111_1111_1111

	VENDOR_GROUP     = VENDOR_USER
	WAREHOUSE_GROUP  = WAREHOUSE_KEEPER | WAREHOUSE_MANAGER
	COMMERCIAL_GROUP = CONTENT_SPECIALIST | ACCOUNT_MANAGER | COMMERCIAL_MANAGER
	SUPPORT_GROUP    = SUPPORT_SPECIALIST | SUPPORT | SUPPORT_MANAGER
)

var ROLES = []Role{
	VENDOR_USER,
	WAREHOUSE_KEEPER,
	WAREHOUSE_MANAGER,
	CONTENT_SPECIALIST,
	ACCOUNT_MANAGER,
	COMMERCIAL_MANAGER,
	SUPPORT_SPECIALIST,
	SUPPORT,
	SUPPORT_MANAGER,

	VENDOR_GROUP,
	WAREHOUSE_GROUP,
	COMMERCIAL_GROUP,
	SUPPORT_GROUP,
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
