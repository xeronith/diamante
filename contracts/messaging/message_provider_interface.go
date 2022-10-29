package messaging

type IMessagingHandler func(receiver string, message string) error

type IMessagingProvider interface {
	Send(receiver string, message string) error
}
