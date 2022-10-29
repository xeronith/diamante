package utility

import (
	"log"
	"syscall"
)

func SetLimits() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatalf("ERROR: Failed to get Rlimit: %s", err)
	}

	rLimit.Max = 15_000_000
	rLimit.Cur = 15_000_000

	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Printf("ERROR: Failed to set Rlimit: %s", err)
	}
}
