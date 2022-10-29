package system

type Flag interface {
	IsSet() bool
	Set()
	Clear()
}
