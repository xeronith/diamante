package server

//noinspection GoSnakeCaseUsage
const (
	BINARY_REQUEST = 0
	TEXT_REQUEST   = 1

	BINARY_RESULT = 3
	TEXT_RESULT   = 4
)

type ITrafficRecorder interface {
	Record(uint32, ...interface{})
	Start()
	Stop() error
	Load(string) error
	Replay(IServer, float32) error
}
