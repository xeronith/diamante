package security

type ISecurityHandler interface {
	AccessControlHandler() IAccessControlHandler
	SetAccessControlHandler(IAccessControlHandler)
	Validate(phoneNumber string, password string) (string, error)
	Verify(token string, confirmationCode string) (string, uint64, error)
	RefreshTokenCache(identity Identity, token string) error
	Authenticate(token string, role Role, remoteAddress string, userAgent string) Identity
	SignOut(Identity) error
}
