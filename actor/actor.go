package actor

import (
	"time"

	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/io"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/contracts/security"
	. "github.com/xeronith/diamante/contracts/serialization"
	. "github.com/xeronith/diamante/contracts/system"
)

type actor struct {
	token         string
	requestHash   string
	remoteAddress string
	userAgent     string
	identity      Identity
	session       ISystemObject
	lastActivity  int64
	writer        IWriter
	active        bool
}

func CreateActor(writer IWriter, active bool, requestHash string, remoteAddress string, userAgent string) IActor {
	return &actor{
		requestHash:   requestHash,
		remoteAddress: remoteAddress,
		userAgent:     userAgent,
		writer:        writer,
		active:        active,
	}
}

func (actor *actor) Dispatch(result IOperationResult) {
	if actor.writer == nil {
		//TODO: For testing purposes only
		return
	}

	actor.writer.Write(result)
}

func (actor *actor) Disconnect(result IOperationResult) {
	actor.writer.End(result)
}

func (actor *actor) Signal(code byte) {
	_ = actor.writer.WriteByte(code)
}

func (actor *actor) Token() string {
	return actor.token
}

func (actor *actor) SetToken(token string) {
	actor.token = token
	if actor.writer != nil {
		actor.writer.SetToken(token)
	}
}

func (actor *actor) RequestHash() string {
	return actor.requestHash
}

func (actor *actor) RemoteAddress() string {
	return actor.remoteAddress
}

func (actor *actor) UserAgent() string {
	return actor.userAgent
}

func (actor *actor) Identity() Identity {
	return actor.identity
}

func (actor *actor) SetIdentity(identity Identity) {
	actor.identity = identity
}

func (actor *actor) Session() ISystemObject {
	return actor.session
}

func (actor *actor) SetSession(session ISystemObject) {
	actor.session = session
}

func (actor *actor) Serializer() ISerializer {
	return actor.writer.Serializer()
}

func (actor *actor) LastActivity() int64 {
	return actor.lastActivity
}

func (actor *actor) UpdateLastActivity() {
	actor.lastActivity = time.Now().UnixNano()
}

func (actor *actor) IsActive() bool {
	return actor.active
}

func (actor *actor) Writer() IWriter {
	return actor.writer
}
