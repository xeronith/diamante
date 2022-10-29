package concurrent

import (
	"sync"

	"github.com/xeronith/diamante/contracts/system"
)

type flag struct {
	sync.RWMutex
	value bool
}

func NewFlag() system.Flag {
	return &flag{}
}

func (flag *flag) IsSet() bool {
	flag.RLock()
	defer flag.RUnlock()

	return flag.value
}

func (flag *flag) Set() {
	flag.Lock()
	defer flag.Unlock()

	flag.value = true
}

func (flag *flag) Clear() {
	flag.Lock()
	defer flag.Unlock()

	flag.value = false
}
