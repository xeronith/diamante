package security

import (
	. "github.com/xeronith/diamante/contracts/security"
)

type defaultSecurityHandler struct{}

func NewDefaultSecurityHandler() ISecurityHandler {
	return &defaultSecurityHandler{}
}

func (handler *defaultSecurityHandler) AccessControlHandler() IAccessControlHandler {
	return nil
}

func (handler *defaultSecurityHandler) SetAccessControlHandler(_ IAccessControlHandler) {
}

func (handler *defaultSecurityHandler) GetByToken(_ string) (Identity, bool) {
	return nil, false
}

func (handler *defaultSecurityHandler) Validate(_ string, _ string) (string, error) {
	return "", nil
}

func (handler *defaultSecurityHandler) Verify(_ string, _ string) (string, uint64, error) {
	return "", 0, nil
}

func (handler *defaultSecurityHandler) RefreshTokenCache(_ Identity, _ string) error {
	return nil
}

func (handler *defaultSecurityHandler) Authenticate(token string, role Role, remoteAddress string, userAgent string) Identity {
	return CreateDefaultIdentity(token, role, remoteAddress, userAgent)
}

func (handler *defaultSecurityHandler) SignOut(_ Identity) error {
	return nil
}
