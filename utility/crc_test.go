package utility_test

import (
	"testing"

	"github.com/xeronith/diamante/utility"
)

func Test_CRC32(test *testing.T) {
	value := ""
	result := utility.CRC32(value)
	_ = result
}
