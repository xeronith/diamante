package serialization

import . "github.com/xeronith/diamante/contracts/system"

type ISerializer interface {
	Serialize(Pointer) ([]byte, error)
	Deserialize([]byte, Pointer) error
}
