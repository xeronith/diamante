package operation

import "time"

//noinspection GoSnakeCaseUsage
const (
	DEFAULT_TIME_LIMIT_WARNING  time.Duration = 1000_000_000
	DEFAULT_TIME_LIMIT_ALERT    time.Duration = 1500_000_000
	DEFAULT_TIME_LIMIT_CRITICAL time.Duration = 2000_000_000
)
