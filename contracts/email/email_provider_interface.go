package email

type IEmailProvider interface {
	Send(string, string, ...interface{}) error
}
