package federation

type IWebfinger interface {
	Self() string
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
