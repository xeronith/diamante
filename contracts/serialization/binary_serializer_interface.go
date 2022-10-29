package serialization

import . "github.com/xeronith/diamante/contracts/system"

type IBinarySerializer interface {
	ISerializer

	Serialize(Pointer) ([]byte, error)
	Deserialize([]byte, Pointer) error
}
