package system

type ITimedTokenGenerator interface {
	Generate() string
	IsValid(string) bool
	Size() int
}
