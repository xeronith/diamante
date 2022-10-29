package logging

type Level int

type ILogger interface {
	SetLevel(Level)
	SysComp(interface{})
	SysCall(interface{})
	Debug(interface{})
	Info(interface{})
	Alert(interface{})
	Warning(interface{})
	Error(interface{})
	Critical(interface{})
	Panic(interface{})
	Fatal(interface{})
	SerializationPath() string
}
