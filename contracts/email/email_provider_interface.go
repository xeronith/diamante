package email

type IEmailProvider interface {
	Send(string, string) error
}
