package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type connectionHeader struct {
	VersionMajor byte
	VersionMinor byte
}

func decodeConnectionHeader(reader io.Reader, header *connectionHeader) error {
	magic := make([]byte, 8, 8)

	remaining := 8
	for remaining > 0 {
		from := 8 - remaining
		n, err := reader.Read(magic[from:])
		if err != nil {
			return err
		}
		remaining -= n
	}

	if string(magic) != "AVIONICA" {
		return fmt.Errorf("invalid connection header: %s", magic)
	}

	return binary.Read(reader, binary.BigEndian, header)
}
