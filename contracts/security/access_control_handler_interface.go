package security

type IAccessControlHandler interface {
	AddOrUpdateAccessControl(key uint64, value uint64, editor Identity) error
	AccessControls() map[uint64]uint64
}
