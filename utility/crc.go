package utility

import (
	"fmt"
	"hash/crc32"
)

func CRC32(value string) string {
	return fmt.Sprintf("0x%.8X", crc32.Checksum([]byte(value), crc32.MakeTable(crc32.IEEE)))
}
