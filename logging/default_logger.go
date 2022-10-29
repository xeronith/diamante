package logging

import (
	"fmt"
	"os"

	. "github.com/xeronith/diamante/contracts/logging"
)

// noinspection GoSnakeCaseUsage
const (
	LEVEL_SILENT            Level = 0
	LEVEL_VERBOSE           Level = 1
	LEVEL_SUPPRESS_SYS_COMP Level = 2
)

type logger struct {
	level             Level
	serializationPath string
}

var defaultLogger = &logger{
	level:             LEVEL_VERBOSE,
	serializationPath: "",
}

func GetDefaultLogger() ILogger {
	return defaultLogger
}

func NewLogger(containerized bool) ILogger {
	logger := &logger{
		level: LEVEL_VERBOSE,
	}

	if containerized {
		logger.serializationPath = "/var/log/container-log/"
	}

	return logger
}

func (logger *logger) SetLevel(level Level) {
	logger.level = level
}

// timestamp := fmt.Sprintf("%-25s ─┤", time.Now().Format(time.RFC3339))

func (logger *logger) SysComp(args interface{}) {
	if logger.level == LEVEL_SUPPRESS_SYS_COMP {
		return
	}

	logger.submit("░░░░░", args)
}

func (logger *logger) SysCall(args interface{}) {
	if logger.level == LEVEL_SUPPRESS_SYS_COMP {
		return
	}

	logger.submit("░ S ░", args)
}

func (logger *logger) Debug(args interface{}) {
	logger.submit("░ D ░", args)
}

func (logger *logger) Info(args interface{}) {
	logger.submit("░ I ░", args)
}

func (logger *logger) Warning(args interface{}) {
	logger.submit("▒ W ▒", args)
}

func (logger *logger) Alert(args interface{}) {
	logger.submit("▓ A ▓", args)
}

func (logger *logger) Error(args interface{}) {
	logger.submit("█ E █", args)
}

func (logger *logger) Critical(args interface{}) {
	logger.submit("█ ! █", args)
}

func (logger *logger) Panic(args interface{}) {
	logger.submit("█ P █", args)
}

func (logger *logger) Fatal(args interface{}) {
	logger.submit("█████ ✗", args)
	os.Exit(1)
}

func (logger *logger) submit(args ...interface{}) {
	if logger.level == LEVEL_SILENT {
		return
	}

	fmt.Println(args...)
}

func (logger *logger) SerializationPath() string {
	return logger.serializationPath
}
