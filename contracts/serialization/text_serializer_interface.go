package serialization

import . "github.com/xeronith/diamante/contracts/system"

type ITextSerializer interface {
	ISerializer

	Serialize(Pointer) (string, error)
	Deserialize(string, Pointer) error
}
