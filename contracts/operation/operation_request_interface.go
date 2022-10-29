package operation

import . "github.com/xeronith/diamante/contracts/system"

type IOperationRequest interface {
	Id() uint64
	Operation() uint64
	Token() string
	ApiVersion() int32
	ClientVersion() int32
	ClientName() string
	Container() Pointer
}
