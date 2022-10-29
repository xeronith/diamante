package messaging

import . "github.com/xeronith/diamante/contracts/system"

type IPushMessage interface {
	GetType() uint64
	GetPayload() Pointer
}
