package utility

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/xeronith/diamante/contracts/system"
)

type timedTokenGenerator struct {
	sync.RWMutex
	duration time.Duration
	tokens   map[string]time.Time
}

func NewTimedTokenGenerator(duration time.Duration) system.ITimedTokenGenerator {
	return &timedTokenGenerator{
		tokens:   make(map[string]time.Time),
		duration: duration,
	}
}

func (generator *timedTokenGenerator) Generate() string {
	token := uuid.New().String()

	generator.Lock()
	defer generator.Unlock()
	generator.tokens[token] = time.Now()

	return token
}

func (generator *timedTokenGenerator) IsValid(token string) bool {
	generator.RLock()
	timestamp, exists := generator.tokens[token]
	generator.RUnlock()

	if !exists {
		return false
	}

	if time.Since(timestamp) > generator.duration {
		generator.Lock()
		defer generator.Unlock()
		delete(generator.tokens, token)

		return false
	}

	return true
}

func (generator *timedTokenGenerator) Size() int {
	generator.RLock()
	defer generator.RUnlock()

	return len(generator.tokens)
}
