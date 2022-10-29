package logging_test

import (
	"testing"

	. "github.com/xeronith/diamante/logging"
)

func TestLogger(test *testing.T) {
	logger := GetDefaultLogger()
	logger.Debug("Lorem ipsum dolor sit amet.")
	logger.Info("Lorem ipsum dolor sit amet.")
	logger.Warning("Lorem ipsum dolor sit amet.")
	logger.Error("Lorem ipsum dolor sit amet.")
	logger.Critical("Lorem ipsum dolor sit amet.")
}
