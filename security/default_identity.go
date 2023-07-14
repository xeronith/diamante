package security

import (
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/system"
)

type defaultIdentity struct {
	token         string
	role          Role
	remoteAddress string
	userAgent     string
}

func CreateDefaultIdentity(token string, role Role, remoteAddress, userAgent string) Identity {
	return &defaultIdentity{
		token:         token,
		role:          role,
		remoteAddress: remoteAddress,
		userAgent:     userAgent,
	}
}

func (identity *defaultIdentity) Id() int64 {
	return 0
}

func (identity *defaultIdentity) Username() string {
	return ""
}

func (identity *defaultIdentity) PhoneNumber() string {
	return ""
}

func (identity *defaultIdentity) FirstName() string {
	return ""
}

func (identity *defaultIdentity) LastName() string {
	return ""
}

func (identity *defaultIdentity) Email() string {
	return ""
}

func (identity *defaultIdentity) Token() string {
	return identity.token
}

func (identity *defaultIdentity) SetToken(token string) {
	identity.token = token
}

func (identity *defaultIdentity) MultiFactor() bool {
	return false
}

func (identity *defaultIdentity) RemoteAddress() string {
	return identity.remoteAddress
}

func (identity *defaultIdentity) SetRemoteAddress(remoteAddress string) {
	identity.remoteAddress = remoteAddress
}

func (identity *defaultIdentity) UserAgent() string {
	return identity.userAgent
}

func (identity *defaultIdentity) SetUserAgent(userAgent string) {
	identity.userAgent = userAgent
}

func (identity *defaultIdentity) Hash() string {
	return ""
}

func (identity *defaultIdentity) Salt() string {
	return ""
}

func (identity *defaultIdentity) Role() Role {
	return identity.role
}

func (identity *defaultIdentity) IsInRole(role Role) bool {
	identityRole := identity.role & role
	return identityRole == role
}

func (identity *defaultIdentity) PublicKey() string {
	return ""
}

func (identity *defaultIdentity) PrivateKey() string {
	return ""
}

func (identity *defaultIdentity) Permission() uint64 {
	return identity.role
}

func (identity *defaultIdentity) Restriction() uint32 {
	return 0
}

func (identity *defaultIdentity) IsRestricted() bool {
	return false
}

func (identity *defaultIdentity) IsNotRestricted() bool {
	return true
}

func (identity *defaultIdentity) Payload() Pointer {
	return nil
}

func (identity *defaultIdentity) Lock(_ uint64) {
}

func (identity *defaultIdentity) Unlock(_ uint64) {
}

func (identity *defaultIdentity) SetSystemCallHandler(_ func(Identity, []string) error) {
}

func (identity *defaultIdentity) SystemCall(_ Identity, _ []string) error {
	return nil
}
