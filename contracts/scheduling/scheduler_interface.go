package scheduling

import "time"

type IScheduler interface {
	Start()
	SetTimeout(func(), time.Duration) string
	SetInterval(func(), time.Duration) string
	Cancel(string)
}
