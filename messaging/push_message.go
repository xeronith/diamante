package messaging

import (
	"github.com/xeronith/diamante/contracts/messaging"
	"github.com/xeronith/diamante/contracts/system"
)

type pushMessage struct {
	_type   uint64
	payload system.Pointer
}

func (broadcast *pushMessage) GetType() uint64 {
	return broadcast._type
}

func (broadcast *pushMessage) GetPayload() system.Pointer {
	return broadcast.payload
}

func NewPushMessage(_type uint64, payload system.Pointer) messaging.IPushMessage {
	return &pushMessage{
		_type:   _type,
		payload: payload,
	}
}
