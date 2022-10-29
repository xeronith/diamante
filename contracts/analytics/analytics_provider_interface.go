package analytics

import "time"

type IAnalyticsProvider interface {
	SubmitEvent(uint64, string, Fields)
	PushNotification(uint64, string, string, time.Duration)
	String() string
}

var DefaultProvider IAnalyticsProvider
