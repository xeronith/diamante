package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func GenerateConfirmationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", 100000+rand.Intn(899999))
}

func GenerateUsername(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s%d", prefix, 1000000+rand.Intn(8999999))
}

func GenerateHash(value, salt string) string {
	content := value + salt
	hash := sha256.New()
	hash.Write([]byte(content))
	return hex.EncodeToString(hash.Sum(nil))
}
