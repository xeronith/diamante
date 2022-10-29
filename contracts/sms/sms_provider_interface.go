package sms

type ISMSProvider interface {
	Send(receiver string, message string) error
}
