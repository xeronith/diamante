package utility

import (
	"crypto/rand"
	"math/big"
)

func GenerateEntityId() string {
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(19), nil).Sub(max, big.NewInt(1))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return ""
	}

	return n.Text(10)
}
